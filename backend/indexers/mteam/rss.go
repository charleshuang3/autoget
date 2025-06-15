package mteam

import (
	"net/url"

	"github.com/charleshuang3/autoget/backend/indexers"
	"github.com/charleshuang3/autoget/backend/indexers/rsshelper"
	"github.com/charleshuang3/autoget/backend/internal/errors"
	"github.com/mmcdole/gofeed"
	"github.com/robfig/cron/v3"
)

func (m *MTeam) RegisterSearchForRSS(s *indexers.RSSSearch) *errors.HTTPStatusError {
	return nil
}

func (m *MTeam) RegisterRSSCronjob(cron *cron.Cron) {
	if m.config.RSS == "" {
		return
	}

	cron.AddFunc("@every 5m", func() {
		items, err := m.pullRSS()
		if err != nil {
			logger.Error().Err(err).Msg("Failed to pull RSS feed")
			return
		}

		rsshelper.SearchRSS(m, m.db, m.notify, items)
	})
}

func (m *MTeam) pullRSS() ([]*indexers.RSSItem, error) {
	u, _ := url.Parse(m.config.RSS)

	fp := gofeed.NewParser()
	f, err := fp.ParseURL(u.String())
	if err != nil {
		return nil, err
	}

	items := []*indexers.RSSItem{}
	for _, item := range f.Items {
		parsed := m.ParseRSSItem(item)
		if parsed == nil {
			continue
		}
		items = append(items, parsed)
	}

	return items, nil
}

func (m *MTeam) ParseRSSItem(item *gofeed.Item) *indexers.RSSItem {
	category := ""
	url := ""
	if len(item.Categories) > 0 {
		category = item.Categories[0]
	}
	if len(item.Enclosures) > 0 {
		url = item.Enclosures[0].URL
	}

	if url == "" {
		return nil
	}

	return &indexers.RSSItem{
		ResID:     item.GUID,
		Title:     item.Title,
		Catergory: category,
		URL:       url,
	}
}
