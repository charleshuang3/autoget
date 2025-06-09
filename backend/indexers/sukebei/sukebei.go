package sukebei

import (
	"github.com/charleshuang3/autoget/backend/indexers/nyaa"
	"github.com/charleshuang3/autoget/backend/indexers/sukebei/prefetcheddata"
)

const (
	defaultBaseURL = "https://sukebei.nyaa.si/"
)

type Client struct {
	nyaa.Client
}

func NewClient(config *nyaa.Config) *Client {
	c := &Client{}
	c.Client = *nyaa.NewClient(config)
	c.Client.DefaultBaseURL = defaultBaseURL
	c.Client.CategoriesMap = prefetcheddata.Categories
	c.Client.CategoriesList = prefetcheddata.CategoriesList

	return c
}
