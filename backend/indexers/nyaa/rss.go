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
	cron.AddFunc("@every 5m", func() {
		feed, err := c.pullRSS()
		if err != nil {
			logger.Error().Err(err).Msg("Failed to pull RSS feed")
			return
		}

		c.SearchRSS(feed)
	})
}

func (c *Client) pullRSS() (*gofeed.Feed, error) {
	u, _ := url.Parse(c.getBaseURL())
	query := u.Query()
	query.Set("page", "rss")
	u.RawQuery = query.Encode()

	fp := gofeed.NewParser()
	fp.Client = c.httpClient
	return fp.ParseURL(u.String())
}

func (c *Client) parseRSSItem(item *gofeed.Item, target *db.RSSSearch) {
	target.Title = item.Title
	target.URL = item.Link
	target.ResID = getResourceIDFromRSSGUID(item.GUID)
	target.Catergory = c.getCategoryFromRSSCategory(
		item.Extensions["nyaa"]["categoryId"][0].Value,
		item.Extensions["nyaa"]["category"][0].Value)
}

func (c *Client) SearchRSS(feed *gofeed.Feed) {
	searchs, err := db.GetSearchsByIndexer(c.db, c.Name())
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get searchs from database")
		return
	}

	downloadStarted := []string{}
	downloadPendingToStart := []string{}

	for _, item := range feed.Items {
		for _, search := range searchs {
			if search.ResID != "" {
				continue
			}
			if strings.Contains(strings.ToLower(item.Title), search.Text) {
				c.parseRSSItem(item, search)

				err = db.UpdateSearch(c.db, search)
				if err != nil {
					logger.Error().Err(err).Msg("Failed to update search")
					continue
				}

				if search.Action == "download" {
					_, err := c.Download(search.ResID)
					if err != nil {
						logger.Error().Err(err).Msg("Failed to download torrent")
						continue
					}

					if err := db.DeleteSearch(c.db, search.ID); err != nil {
						logger.Error().Err(err).Msg("Failed to delete search")
						continue
					}

					downloadStarted = append(downloadStarted, search.Title)
				} else if search.Action == "notification" {
					downloadPendingToStart = append(downloadPendingToStart, search.Title)
				}
			}
		}
	}

	if len(downloadStarted) > 0 || len(downloadPendingToStart) > 0 {
		msg, err := indexers.RenderRSSResult(c.Name(), downloadStarted, downloadPendingToStart)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to render RSS result")
			return
		}
		if err := c.notify.SendMarkdownMessage(msg); err != nil {
			logger.Error().Err(err).Msg("Failed to send RSS notification")
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
