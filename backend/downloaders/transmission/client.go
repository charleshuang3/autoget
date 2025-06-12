package transmission

import (
	"context"
	"errors"
	"net/http"
	"net/url"
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
		id := c.name + "/" + hash
		uploaded := *t.UploadedEver

		ss := &db.SeedingStatus{
			ID: id,
		}
		err := c.db.First(ss).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ss.UploadHistories = make(map[string]int64)
			ss.AddToday(uploaded)
			c.db.Create(ss)

			continue
		}
		ss.CleanupHistory()
		ss.AddToday(uploaded)
		c.db.Save(ss)

		before, ok := ss.GetXDayBefore(int(c.cfg.SeedingPolicy.IntervalInDays))
		if !ok {
			continue
		}

		if (uploaded - before) > c.cfg.SeedingPolicy.UploadAtLeastInMB*1024*1024 {
			continue
		}

		// stop this torrent
		stopIDs = append(stopIDs, *t.ID)
	}

	if err := c.db.Where("updated_at < ?", time.Now().AddDate(0, 0, -db.StoreMaxDays)).Delete(&db.SeedingStatus{}).Error; err != nil {
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
