package transmission

import (
	"net/url"

	"github.com/charleshuang3/autoget/backend/downloaders"
	"github.com/hekmon/transmissionrpc/v3"
)

func New(cfg *downloaders.DownloaderConfig) (*transmissionrpc.Client, error) {
	u, err := url.Parse(cfg.Transmission.URL)
	if err != nil {
		return nil, err
	}

	if cfg.Transmission.Username != "" && cfg.Transmission.Password != "" {
		u.User = url.UserPassword(cfg.Transmission.Username, cfg.Transmission.Password)
	}

	client, err := transmissionrpc.New(u, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}
