package nyaa

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/charleshuang3/autoget/backend/indexers"
	"github.com/charleshuang3/autoget/backend/internal/db"
	"github.com/charleshuang3/autoget/backend/internal/errors"
	"github.com/mmcdole/gofeed"
	"github.com/robfig/cron/v3"
)

func (c *Client) RegisterSearchForRSS(s *indexers.RSSSearch) *errors.HTTPStatusError {
	search := &db.RSSSearch{
		Indexer: c.Name(),
		Text:    s.Text,
		Action:  s.Action,
	}
	err := db.AddSearch(c.db, search)
	return errors.NewHTTPStatusError(http.StatusInternalServerError, err.Error())
}

func (c *Client) RegisterRSSCronjob(cron *cron.Cron) {
	cron.AddFunc("@every 5m", c.PullRSSAndSearch)
}

func (c *Client) PullRSSAndSearch() {
	searchs, err := db.GetSearchsByIndexer(c.db, c.Name())
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get searchs from database")
		return
	}

	u, _ := url.Parse(c.getBaseURL())
	query := u.Query()
	query.Set("page", "rss")
	u.RawQuery = query.Encode()

	fp := gofeed.NewParser()
	fp.Client.Transport = c.httpClient.Transport
	feed, err := fp.ParseURL(u.String())
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get RSS feed")
		return
	}
	for _, item := range feed.Items {
		for _, search := range searchs {
			if search.ResID != "" {
				continue
			}
			if strings.Contains(strings.ToLower(item.Title), search.Text) {
				search.Title = item.Title
				search.URL = item.Link
				search.ResID = getResourceIDFromRSSGUID(item.GUID)
				search.Catergory = c.getCategoryFromRSSCategory(
					item.Custom["nyaa:categoryId"], item.Custom["nyaa:category"])

				err = db.UpdateSearch(c.db, search)
				if err != nil {
					logger.Error().Err(err).Msg("Failed to update search")
					continue
				}

				// TODO: download
				// TODO: notification
			}
		}
	}
}

func (c *Client) getCategoryFromRSSCategory(categoryID, category string) string {
	if cat, ok := c.CategoriesMap[categoryID]; ok {
		return cat.Name
	} else {
		return category
	}
}

func getResourceIDFromRSSGUID(guid string) string {
	parts := strings.Split(guid, "/")
	return parts[len(parts)-1]
}
