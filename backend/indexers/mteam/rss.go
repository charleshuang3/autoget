package mteam

import (
	"github.com/charleshuang3/autoget/backend/indexers"
	"github.com/charleshuang3/autoget/backend/internal/db"
	"github.com/charleshuang3/autoget/backend/internal/errors"
	"github.com/mmcdole/gofeed"
	"github.com/robfig/cron/v3"
)

func (m *MTeam) RegisterSearchForRSS(s *indexers.RSSSearch) *errors.HTTPStatusError {
	return nil
}

func (m *MTeam) RegisterRSSCronjob(cron *cron.Cron) {

}

func (m *MTeam) parseRSSItem(item *gofeed.Item, target *db.RSSSearch) {
	target.ResID = item.GUID
	target.Title = item.Title
	if len(item.Categories) > 0 {
		target.Catergory = item.Categories[0]
	}
	if len(item.Enclosures) > 0 {
		target.URL = item.Enclosures[0].URL
	}
}
