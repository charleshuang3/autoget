package mteam

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/charleshuang3/autoget/backend/indexers"
	"github.com/charleshuang3/autoget/backend/internal/errors"
	"github.com/rs/zerolog/log"
)

type searchRequest struct {
	Mode       string   `json:"mode"` // "normal" or "adult"
	Categories []string `json:"categories"`
	Visible    int      `json:"visible"` // fixed 1
	PageNumber uint32   `json:"pageNumber"`
	PageSize   uint32   `json:"pageSize"`

	// Optional
	Keyword   string   `json:"keyword,omitempty"`
	Discount  string   `json:"discount,omitempty"` // "FREE" or ""
	Standards []string `json:"standards,omitempty"`
}

type searchResponseItem struct {
	ID               string   `json:"id"`
	CreatedDate      string   `json:"createdDate"`
	LastModifiedDate string   `json:"lastModifiedDate"`
	Name             string   `json:"name"`
	SmallDescr       string   `json:"smallDescr"`
	Imdb             string   `json:"imdb"`
	ImdbRating       string   `json:"imdbRating"`
	Douban           string   `json:"douban"`
	DoubanRating     string   `json:"doubanRating"`
	DmmCode          string   `json:"dmmCode"` // a url to dmm.
	Author           string   `json:"author"`
	Category         string   `json:"category"`
	Source           string   `json:"source"`
	Medium           string   `json:"medium"`
	Standard         string   `json:"standard"`
	VideoCodec       string   `json:"videoCodec"`
	AudioCodec       string   `json:"audioCodec"`
	Team             string   `json:"team"`
	Processing       string   `json:"processing"`
	Countries        []string `json:"countries"`
	Numfiles         string   `json:"numfiles"`
	Size             string   `json:"size"`
	Labels           string   `json:"labels"`
	LabelsNew        []string `json:"labelsNew"`
	MsUp             string   `json:"msUp"`
	Anonymous        bool     `json:"anonymous"`
	InfoHash         string   `json:"infoHash"`
	Status           struct {
		ID               string      `json:"id"`
		CreatedDate      string      `json:"createdDate"`
		LastModifiedDate string      `json:"lastModifiedDate"`
		PickType         string      `json:"pickType"`
		ToppingLevel     string      `json:"toppingLevel"`
		ToppingEndTime   string      `json:"toppingEndTime"`
		Discount         string      `json:"discount"`
		DiscountEndTime  string      `json:"discountEndTime"`
		TimesCompleted   string      `json:"timesCompleted"`
		Comments         string      `json:"comments"`
		LastAction       string      `json:"lastAction"`
		LastSeederAction string      `json:"lastSeederAction"`
		Views            string      `json:"views"`
		Hits             string      `json:"hits"`
		Support          string      `json:"support"`
		Oppose           string      `json:"oppose"`
		Status           string      `json:"status"`
		Seeders          string      `json:"seeders"`
		Leechers         string      `json:"leechers"`
		Banned           bool        `json:"banned"`
		Visible          bool        `json:"visible"`
		PromotionRule    interface{} `json:"promotionRule"`  // never seen
		MallSingleFree   interface{} `json:"mallSingleFree"` // unused
	} `json:"status"`
	DmmInfo struct {
		CreatedDate      string   `json:"createdDate"`
		LastModifiedDate string   `json:"lastModifiedDate"`
		ID               string   `json:"id"`
		ProductNumber    string   `json:"productNumber"`
		Director         string   `json:"director"`
		Series           string   `json:"series"`
		Maker            string   `json:"maker"`
		Label            string   `json:"label"`
		KeywordList      []string `json:"keywordList"`
		ActressList      []string `json:"actressList"`
	} `json:"dmmInfo"`
	EditedBy   interface{} `json:"editedBy"` // never seen
	EditDate   string      `json:"editDate"`
	Collection bool        `json:"collection"`
	InRss      bool        `json:"inRss"`
	CanVote    bool        `json:"canVote"`
	ImageList  []string    `json:"imageList"`
	ResetBox   string      `json:"resetBox"`
}

func (it *searchResponseItem) extractDBInfo() []indexers.VideoDB {
	var res []indexers.VideoDB

	if it.Douban != "" {
		res = append(res, indexers.VideoDB{
			DB:     "douban",
			Link:   it.Douban,
			Rating: it.DoubanRating,
		})
	}

	if it.Imdb != "" {
		res = append(res, indexers.VideoDB{
			DB:     "imdb",
			Link:   it.Imdb,
			Rating: it.ImdbRating,
		})
	}

	if it.DmmCode != "" {
		res = append(res, indexers.VideoDB{
			DB:   "dmm",
			Link: it.DmmCode,
		})
	}

	return res
}

type searchResponse struct {
	Code    interface{} `json:"code"` // maybe string or int
	Message string      `json:"message"`
	Data    struct {
		PageNumber string               `json:"pageNumber"`
		PageSize   string               `json:"pageSize"`
		Total      string               `json:"total"`
		TotalPages string               `json:"totalPages"`
		Data       []searchResponseItem `json:"data"`
	} `json:"data"`
}

func (m *MTeam) List(listReq *indexers.ListRequest) (*indexers.ListResult, *errors.HTTPStatusError) {
	// check category is known.
	cat, ok := m.prefetched.Categories.Infos[listReq.Category]
	if !ok {
		return nil, errors.NewHTTPStatusError(http.StatusBadRequest, "invalid category")
	}

	req := &searchRequest{
		Mode:       cat.Mode,
		Categories: cat.Categories,
		Visible:    1,
		PageNumber: listReq.Page,
		PageSize:   listReq.PageSize,
		Keyword:    listReq.Keyword,
	}

	if listReq.Free {
		req.Discount = "FREE"
	}

	for _, standard := range listReq.Standards {
		if st, ok := m.standards[standard]; ok {
			req.Standards = append(req.Standards, st)
		}
	}

	if listReq.Category == categoryAdult || listReq.Category == categoryNormal {
		// root category use empty categories list
		req.Categories = []string{}
	}

	reqData, err := json.Marshal(req)
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, "failed to marshal request")
	}

	request, err := http.NewRequest(http.MethodPost, m.config.GetBaseURL()+"/api/torrent/search", bytes.NewReader(reqData))
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, "failed to new request")
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-key", m.config.APIKey)

	client := http.Client{
		Timeout: httpTimeout,
	}
	r, err := client.Do(request)
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, "failed to request")
	}

	if r.StatusCode != http.StatusOK {
		log.Error().Err(err).Str("indexer", name).Int("status_code", r.StatusCode).Msg("API error")
		return nil, errors.NewHTTPStatusError(r.StatusCode, "search request failed")
	}

	var resp searchResponse
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, "failed to unmarshal response")
	}

	if resp.Code != "0" {
		log.Error().Any("code", resp.Code).Str("message", resp.Message).Msg("API error")
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, resp.Message)
	}

	page_, _ := strconv.Atoi(resp.Data.PageNumber)
	pageSize_, _ := strconv.Atoi(resp.Data.PageSize)
	total, _ := strconv.Atoi(resp.Data.Total)
	totalPages, _ := strconv.Atoi(resp.Data.TotalPages)

	ListResult := &indexers.ListResult{
		Pagination: indexers.Pagination{
			Page:       uint32(page_),
			PageSize:   uint32(pageSize_),
			Total:      uint32(total),
			TotalPages: uint32(totalPages),
		},
	}

	for _, item := range resp.Data.Data {
		if m.config.ExcludeGayContent && item.Category == categoryGayPorn {
			continue
		}

		seeders, _ := strconv.Atoi(item.Status.Seeders)
		leechers, _ := strconv.Atoi(item.Status.Leechers)
		size, _ := strconv.ParseUint(item.Size, 10, 64)

		images := []string{}
		if len(item.ImageList) > 0 {
			images = append(images, item.ImageList[0])
		}

		ListResult.Resources = append(ListResult.Resources, indexers.ListResourceItem{
			ID:         item.ID,
			Title:      item.Name,
			Title2:     item.SmallDescr,
			Category:   m.prefetched.Categories.Infos[item.Category].Name,
			Size:       size,
			Resolution: m.prefetched.Standards[item.Standard],
			Seeders:    uint32(seeders),
			Leechers:   uint32(leechers),
			DBs:        item.extractDBInfo(),
			Images:     images,
		})
	}

	return ListResult, nil
}
