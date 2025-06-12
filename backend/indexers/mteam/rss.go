package mteam

import (
	"github.com/charleshuang3/autoget/backend/indexers"
	"github.com/charleshuang3/autoget/backend/internal/errors"
	"github.com/robfig/cron/v3"
)

func (m *MTeam) RegisterSearchForRSS(s *indexers.RSSSearch) *errors.HTTPStatusError {
	return nil
}

func (m *MTeam) RegisterRSSCronjob(cron *cron.Cron) {

}
