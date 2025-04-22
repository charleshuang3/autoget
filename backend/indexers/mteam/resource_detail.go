package mteam

import (
	"github.com/charleshuang3/autoget/backend/indexers"
	"github.com/charleshuang3/autoget/backend/internal/errors"
)

type resourceDetail struct {
	searchResponseItem

	OriginFileName string `json:"originFileName"`
	Descr          string `json:"descr"`
	Mediainfo      string `json:"mediainfo"`
}

type detailResponse struct {
	Code    interface{}    `json:"code"` // maybe string or int
	Message string         `json:"message"`
	Data    resourceDetail `json:"data"`
}

func (m *MTeam) Detail(id string) (*indexers.ResourceDetail, *errors.HTTPStatusError) {

	return nil, nil
}
