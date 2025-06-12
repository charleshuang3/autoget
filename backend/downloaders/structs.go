package downloaders

import (
	"fmt"

	"github.com/charleshuang3/autoget/backend/downloaders/config"
	"github.com/charleshuang3/autoget/backend/downloaders/transmission"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type IDownloader interface {
	RegisterDailySeedingChecker(cron *cron.Cron)
	TorrentsDir() string
	DownloadDir() string
}

func New(name string, cfg *config.DownloaderConfig, db *gorm.DB) (IDownloader, error) {
	if cfg.Transmission == nil {
		return nil, fmt.Errorf("Unknown downloader %s", name)
	}

	return transmission.New(name, cfg, db)
}
