package mteam

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/charleshuang3/autoget/backend/indexers"
	"github.com/charleshuang3/autoget/backend/internal/errors"
	"github.com/rs/zerolog/log"
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
	_, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusBadRequest, "invalid id")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("id", id)
	err = writer.Close()
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, "close multipart")
	}

	req, err := http.NewRequest(http.MethodPost, m.config.GetBaseURL()+"/api/torrent/detail", body)
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, "failed to new request")
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-api-key", m.config.APIKey)

	client := http.Client{
		Timeout: httpTimeout,
	}
	r, err := client.Do(req)
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, "failed to request")
	}

	if r.StatusCode != http.StatusOK {
		log.Error().Err(err).Str("indexer", name).Int("status_code", r.StatusCode).Msg("API error")
		return nil, errors.NewHTTPStatusError(r.StatusCode, "search request failed")
	}

	resp := &detailResponse{}
	err = json.NewDecoder(r.Body).Decode(resp)
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, "failed to unmarshal response")
	}

	if resp.Code != "0" {
		log.Error().Any("code", resp.Code).Str("message", resp.Message).Msg("API error")
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, resp.Message)
	}

	seeders, _ := strconv.Atoi(resp.Data.Status.Seeders)
	leechers, _ := strconv.Atoi(resp.Data.Status.Leechers)
	size, _ := strconv.ParseUint(resp.Data.Size, 10, 64)

	res := &indexers.ResourceDetail{
		ListResourceItem: indexers.ListResourceItem{
			ID:         resp.Data.ID,
			Title:      resp.Data.Name,
			Title2:     resp.Data.SmallDescr,
			Category:   m.prefetched.Categories.Infos[resp.Data.Category].Name,
			Size:       size,
			Resolution: m.prefetched.Standards[resp.Data.Standard],
			Seeders:    uint32(seeders),
			Leechers:   uint32(leechers),
			DBs:        resp.Data.extractDBInfo(),
			Images:     resp.Data.ImageList,
		},
		Mediainfo:   resp.Data.Mediainfo,
		Description: resp.Data.Descr,
	}

	return res, nil
}
