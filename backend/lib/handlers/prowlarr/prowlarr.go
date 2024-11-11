package prowlarr

import (
	"github.com/charleshuang3/autoget/backend/lib/config"
	"golift.io/starr"
	prowlarr_api "golift.io/starr/prowlarr"
)

type Prowlarr struct {
	client *prowlarr_api.Prowlarr
}

func New(config *config.Config) *Prowlarr {
	c := starr.New(config.Prowlarr.APIURL, config.Prowlarr.APIKey, 0)
	client := prowlarr_api.New(c)
	return &Prowlarr{
		client: client,
	}
}
