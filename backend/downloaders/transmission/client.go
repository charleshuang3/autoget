package transmission

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/charleshuang3/autoget/backend/downloaders/config"
	"github.com/charleshuang3/autoget/backend/internal/db"
	"github.com/hekmon/transmissionrpc/v3"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	logger = log.With().Str("component", "transmission").Logger()

	httpClient = http.DefaultClient
)

type Client struct {
	client *transmissionrpc.Client
	name   string
	db     *gorm.DB
	cfg    *config.DownloaderConfig
}

func New(name string, cfg *config.DownloaderConfig, db *gorm.DB) (*Client, error) {
	u, err := url.Parse(cfg.Transmission.URL)
	if err != nil {
		return nil, err
	}

	if cfg.Transmission.Username != "" && cfg.Transmission.Password != "" {
		u.User = url.UserPassword(cfg.Transmission.Username, cfg.Transmission.Password)
	}

	client, err := transmissionrpc.New(u, &transmissionrpc.Config{
		CustomClient: httpClient,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
		name:   name,
		db:     db,
		cfg:    cfg,
	}, nil
}

func (c *Client) RegisterCronjobs(cron *cron.Cron) {
	c.RegisterDailySeedingChecker(cron)

	go func() {
		time.Sleep(time.Minute)
		c.ProgressChecker()
	}()
}

func (c *Client) ProgressChecker() {
	torrents, err := c.client.TorrentGetAll(context.Background())
	if err != nil {
		logger.Error().Err(err).Str("name", c.name).Msg("failed to get all torrents")
		return
	}

	torrentsByHash := make(map[string]*transmissionrpc.Torrent)
	for _, t := range torrents {
		torrentsByHash[*t.HashString] = &t
	}

	statuses, err := db.GetDownloadStatusByDownloaderAndState(c.db, c.name, db.DownloadStarted)
	if err != nil {
		logger.Error().Err(err).Str("name", c.name).Msg("failed to get download status")
		return
	}

	for _, s := range statuses {
		t, ok := torrentsByHash[s.ID]
		if !ok {
			continue
		}

		s.DownloadProgress = int32(*t.PercentDone * 1000)
		if *t.Status == transmissionrpc.TorrentStatusSeed {
			s.State = db.DownloadSeeding
		}
		db.SaveDownloadStatus(c.db, &s)
	}

	// check if transmission is actively downloading.
	stats, err := c.client.SessionStats(context.Background())
	if err != nil {
		logger.Err(err).Str("name", c.name).Msg("failed to get session stats")
	}

	// if downloadSpeed > 2M/s, consider transimission is still busy
	if stats.DownloadSpeed > 2*1000*1000 {
		return
	}

	// start copys
	statuses, err = db.GetDownloadStatusByDownloaderStateAndMoveState(c.db, c.name, db.DownloadSeeding, db.UnMoved)
	if err != nil {
		logger.Error().Err(err).Str("name", c.name).Msg("failed to get seeding download status")
		return
	}

	for _, s := range statuses {
		t, ok := torrentsByHash[s.ID]
		if !ok {
			continue
		}

		os.MkdirAll(filepath.Join(c.cfg.Transmission.FinishedDir, s.ID), 0755)

		success := true
		for _, f := range t.Files {
			from := filepath.Join(*t.DownloadDir, f.Name)
			target := filepath.Join(c.cfg.Transmission.FinishedDir, s.ID, f.Name)
			fromFile, err := os.Open(from)
			if err != nil {
				success = false
				logger.Error().Err(err).Str("name", c.name).Msg("failed to open file")
				break
			}
			defer fromFile.Close()
			targetFile, err := os.Open(target)
			if err != nil {
				success = false
				logger.Error().Err(err).Str("name", c.name).Msg("failed to open file")
				break
			}
			defer targetFile.Close()

			_, err = io.Copy(targetFile, fromFile)
			if err != nil {
				success = false
				logger.Error().Err(err).Str("name", c.name).Msg("failed to copy file")
				break
			}
		}

		if success {
			s.MoveState = db.Moved
			db.SaveDownloadStatus(c.db, &s)
		}
	}
}

func (c *Client) RegisterDailySeedingChecker(cron *cron.Cron) {
	if c.cfg.SeedingPolicy == nil {
		return
	}

	cron.AddFunc("0 0 8 * * *", func() {
		c.checkDailySeeding()
	})
}

func (c *Client) checkDailySeeding() {
	torrents, err := c.client.TorrentGetAll(context.Background())
	if err != nil {
		logger.Error().Err(err).Str("name", c.name).Msg("failed to get all torrents")
		return
	}

	stopIDs := []int64{}

	for _, t := range torrents {
		// only check seeding torrents
		if *t.Status != transmissionrpc.TorrentStatusSeed {
			continue
		}

		hash := (*t.HashString)
		uploaded := *t.UploadedEver

		ss, err := db.GetDownloadStatus(c.db, c.name, hash)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ss.ID = hash
			ss.Downloader = c.name
			ss.State = db.DownloadSeeding
			ss.UploadHistories = make(map[string]int64)
			ss.AddToday(uploaded)
			db.SaveDownloadStatus(c.db, ss)

			continue
		}
		ss.CleanupHistory()
		ss.AddToday(uploaded)

		db.SaveDownloadStatus(c.db, ss)

		before, ok := ss.GetXDayBefore(int(c.cfg.SeedingPolicy.IntervalInDays))
		if !ok {
			continue
		}

		if (uploaded - before) > c.cfg.SeedingPolicy.UploadAtLeastInMB*1024*1024 {
			continue
		}

		// stop this torrent
		stopIDs = append(stopIDs, *t.ID)

		ss.State = db.DownloadStopped
		db.SaveDownloadStatus(c.db, ss)
	}

	if err := c.db.Where("updated_at < ?", time.Now().AddDate(0, 0, -db.StoreMaxDays)).Delete(&db.DownloadStatus{}).Error; err != nil {
		logger.Error().Err(err).Str("name", c.name).Msg("failed to cleanup seeding status")
	}

	// stop torrents
	if err := c.client.TorrentStopIDs(context.Background(), stopIDs); err != nil {
		logger.Error().Err(err).Str("name", c.name).Msg("failed to stop torrents")
	}
}

func (c *Client) TorrentsDir() string {
	return c.cfg.Transmission.TorrentsDir
}

func (c *Client) DownloadDir() string {
	return c.cfg.Transmission.DownloadDir
}
