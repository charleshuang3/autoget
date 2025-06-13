package sukebei

import (
	"github.com/charleshuang3/autoget/backend/indexers/nyaa"
	"github.com/charleshuang3/autoget/backend/indexers/sukebei/prefetcheddata"
	"gorm.io/gorm"
)

const (
	defaultBaseURL = "https://sukebei.nyaa.si/"
)

type Client struct {
	nyaa.Client
}

func NewClient(config *nyaa.Config, torrentsDir string, db *gorm.DB) *Client {
	c := &Client{}
	c.Client = *nyaa.NewClient(config, torrentsDir, db)
	c.Client.DefaultBaseURL = defaultBaseURL
	c.Client.CategoriesMap = prefetcheddata.Categories
	c.Client.CategoriesList = prefetcheddata.CategoriesList

	return c
}

// Name of the indexer.
func (c *Client) Name() string {
	return "sukebei"
}
