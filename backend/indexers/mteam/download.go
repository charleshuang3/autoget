package mteam

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/charleshuang3/autoget/backend/indexers"
	"github.com/charleshuang3/autoget/backend/internal/errors"
	"github.com/charleshuang3/autoget/backend/internal/helpers"
	"github.com/rs/zerolog/log"
)

type genDownloadLinkResponse struct {
	Code    interface{} `json:"code"` // maybe string or int
	Message string      `json:"message"`
	Data    string      `json:"data"`
}

func (m *MTeam) Download(id, dir string) (*indexers.DownloadResult, *errors.HTTPStatusError) {
	_, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusBadRequest, "invalid id")
	}

	resp := &genDownloadLinkResponse{}
	er := makeMultipartAPICall(m.config.getBaseURL(), "/api/torrent/genDlToken", m.config.APIKey, map[string]string{
		"id": id,
	}, resp)
	if er != nil {
		return nil, er
	}

	if resp.Code != "0" {
		log.Error().Any("code", resp.Code).Str("message", resp.Message).Str("API", "/api/torrent/genDlToken").Msg("API error")
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, resp.Message)
	}

	destFilePath := filepath.Join(dir, name+"."+id+".torrent")

	err = helpers.DownloadFileFromURL(http.DefaultClient, resp.Data, destFilePath)
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, err.Error())
	}

	return &indexers.DownloadResult{
		TorrentFilePath: destFilePath,
	}, nil
}
