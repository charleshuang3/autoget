package nyaa

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/charleshuang3/autoget/backend/indexers"
	"github.com/charleshuang3/autoget/backend/indexers/nyaa/prefetcheddata"
	"github.com/charleshuang3/autoget/backend/internal/errors"
	"github.com/charleshuang3/autoget/backend/internal/helpers"
	"github.com/rs/zerolog/log"
)

var (
	httpClient = http.DefaultClient
)

const (
	defaultBaseURL  = "https://nyaa.si/"
	defaultPageSize = 75
)

type Config struct {
	BaseURL  string `yaml:"base_url"`
	UseProxy bool   `yaml:"use_proxy"`
	ProxyURL string
}

type Client struct {
	config *Config
}

func (c *Config) getBaseURL() string {
	if c.BaseURL == "" {
		return defaultBaseURL
	}
	return c.BaseURL
}

func (c *Config) proxyURL() string {
	if c.ProxyURL != "" {
		return c.ProxyURL
	}
	return os.Getenv("HTTP_PROXY")
}

func NewClient(config *Config) *Client {
	c := &Client{
		config: config,
	}

	if config.UseProxy {
		proxyURL := config.proxyURL()
		if proxyURL != "" {
			proxy, err := url.Parse(proxyURL)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to parse proxy URL")
			}
			httpClient.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxy),
			}
		}
	}

	return c
}

// Name of the indexer.
func (c *Client) Name() string {
	return "nyaa"
}

// Categories returns indexer's resource categories.
func (c *Client) Categories() ([]indexers.Category, *errors.HTTPStatusError) {
	return prefetcheddata.CategoriesList, nil
}

// List resources in given category and keyword (optional).
func (c *Client) List(req *indexers.ListRequest) (*indexers.ListResult, *errors.HTTPStatusError) {
	// Nyaa only support following query params
	q := url.Values{}
	if req.Category != "" {
		q.Set("c", req.Category)
	}
	if req.Keyword != "" {
		q.Set("q", req.Keyword)
	}
	if req.Page > 0 {
		q.Set("p", strconv.Itoa(int(req.Page)))
	}

	u, err := url.Parse(c.config.getBaseURL())
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, fmt.Sprintf("failed to join path: %v", err))
	}

	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, fmt.Sprintf("failed to fetch list page: %v", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, fmt.Sprintf("failed to read response body: %v", err))
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, fmt.Sprintf("failed to parse HTML: %v", err))
	}

	var resources []indexers.ListResourceItem
	doc.Find("table.torrent-list tbody tr").Each(func(i int, s *goquery.Selection) {
		item := indexers.ListResourceItem{}

		// Column 1: Category
		categoryLink := s.Find("td:nth-child(1) a").AttrOr("href", "")
		if categoryLink != "" {
			item.Category = prefetcheddata.Categories[strings.TrimPrefix(categoryLink, "/?c=")].Name
		}

		// Column 2: Title and ID
		titleLink := s.Find("td:nth-child(2) a").Last()
		item.Title = strings.TrimSpace(titleLink.Text())
		idLink, exists := titleLink.Attr("href")
		if exists {
			item.ID = strings.TrimPrefix(idLink, "/view/")
		}

		// Column 4: Size
		item.Size, _ = humanSizeToBytes(s.Find("td:nth-child(4)").Text())

		// Column 5: CreatedDate
		timestampStr, exists := s.Find("td:nth-child(5)").Attr("data-timestamp")
		if exists {
			timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
			if err == nil {
				item.CreatedDate = timestamp
			}
		}

		// Column 6: Seeders
		seedersStr := s.Find("td:nth-child(6)").Text()
		seeders, err := strconv.ParseUint(seedersStr, 10, 32)
		if err == nil {
			item.Seeders = uint32(seeders)
		}

		// Column 7: Leechers
		leechersStr := s.Find("td:nth-child(7)").Text()
		leechers, err := strconv.ParseUint(leechersStr, 10, 32)
		if err == nil {
			item.Leechers = uint32(leechers)
		}

		resources = append(resources, item)
	})

	// Pagination
	page := doc.Find("ul.pagination li.active").First().Text()
	page = strings.Replace(page, "(current)", "", -1) // Remove "(current)"
	page = strings.TrimSpace(page)                    // Trim leading/trailing whitespace
	currentPage, _ := strconv.Atoi(page)

	listResult := &indexers.ListResult{
		Pagination: indexers.Pagination{
			Page: uint32(currentPage),
			// Nyaa does not provide following information
			TotalPages: 0,
			PageSize:   defaultPageSize,
			Total:      0,
		},
		Resources: resources,
	}

	return listResult, nil
}

// Detail of a resource.
func (c *Client) Detail(id string, fileList bool) (*indexers.ResourceDetail, *errors.HTTPStatusError) {
	url, err := url.JoinPath(c.config.getBaseURL(), "view", id)
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, fmt.Sprintf("failed to join path: %v", err))
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, fmt.Sprintf("failed to fetch detail page: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewHTTPStatusError(resp.StatusCode, fmt.Sprintf("failed to fetch detail page, status code: %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, fmt.Sprintf("failed to read response body: %v", err))
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, fmt.Sprintf("failed to parse HTML: %v", err))
	}

	detail := &indexers.ResourceDetail{
		ListResourceItem: indexers.ListResourceItem{
			ID: id,
		},
	}

	// First panel include Title, Category, Size, CreatedDate
	firstPanel := doc.Find(".panel-default").First()

	// Extract Title
	detail.Title = strings.TrimSpace(firstPanel.Find(".panel-title").First().Text())

	firstPanelBody := firstPanel.Find(".panel-body").First()

	// firstPanelBody looks like a map:
	// key in .col-md-1 and value in .col-md-5
	keys := []string{}
	firstPanelBody.Find(".col-md-1").Each(func(i int, s *goquery.Selection) {
		key := s.Text()
		keys = append(keys, key)
	})

	firstPanelBody.Find(".col-md-5").Each(func(i int, s *goquery.Selection) {
		if i >= len(keys) {
			return
		}
		key := keys[i]
		switch key {
		case "Category:":
			{
				links := s.Find("a")
				href, exists := links.Last().Attr("href")
				if exists {
					detail.Category = prefetcheddata.Categories[strings.TrimPrefix(href, "/?c=")].Name
				}
			}
		case "File size:":
			{
				sizeStr := s.Text()
				detail.Size, err = humanSizeToBytes(sizeStr)
				if err != nil {
					log.Info().Err(err).Msgf("humanSizeToBytes() invalid size string: %s", sizeStr)
				}
			}
		case "Seeders:":
			{
				seedersStr := s.Text()
				seeders, err := strconv.ParseUint(seedersStr, 10, 32)
				if err == nil {
					detail.Seeders = uint32(seeders)
				}
			}
		case "Leechers:":
			{
				leechersStr := s.Text()
				leechers, err := strconv.ParseUint(leechersStr, 10, 32)
				if err == nil {
					detail.Leechers = uint32(leechers)
				}
			}
		}
	})

	// Extract CreatedDate
	timestampStr, exists := doc.Find("div[data-timestamp]").First().Attr("data-timestamp")
	if exists {
		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err == nil {
			detail.CreatedDate = timestamp
		}
	}

	// Extract Description
	detail.Description = doc.Find("#torrent-description").Text()

	// Extract Files
	if fileList {
		doc.Find(".torrent-file-list ul li").Each(func(i int, s *goquery.Selection) {
			fileName := strings.TrimSpace(s.Contents().Not("span").Text())
			fileSizeStr := s.Find(".file-size").Text()
			fileSizeStr = strings.TrimPrefix(fileSizeStr, "(")
			fileSizeStr = strings.TrimSuffix(fileSizeStr, ")")
			fileSizeStr = strings.TrimSpace(fileSizeStr)

			fileSize, err := humanSizeToBytes(fileSizeStr)
			if err != nil {
				log.Info().Err(err).Msgf("humanSizeToBytes() invalid size string: %s", fileSizeStr)
			}

			detail.Files = append(detail.Files, indexers.File{
				Name: fileName,
				Size: fileSize,
			})
		})
	}

	return detail, nil
}

// Download the torrent file to given dir or return the magnet link.
func (c *Client) Download(id, dir string) (*indexers.DownloadResult, *errors.HTTPStatusError) {
	fileName := fmt.Sprintf("%s.torrent", id)

	url, err := url.JoinPath(c.config.getBaseURL(), "download", fileName)
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, fmt.Sprintf("failed to join path: %v", err))
	}

	err = helpers.DownloadFileFromURL(httpClient, url, filepath.Join(dir, fileName))
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, err.Error())
	}

	return &indexers.DownloadResult{
		TorrentFilePath: filepath.Join(dir, fileName),
	}, nil
}

func humanSizeToBytes(sizeStr string) (uint64, error) {
	if sizeStr == "" {
		return 0, nil
	}
	parts := strings.Split(sizeStr, " ")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid size format: %s", sizeStr)
	}
	value, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size value: %s", parts[0])
	}
	unit := parts[1]
	var multiplier float64
	switch unit {
	case "B":
		multiplier = 1
	case "KiB":
		multiplier = 1024
	case "MiB":
		multiplier = 1024 * 1024
	case "GiB":
		multiplier = 1024 * 1024 * 1024
	case "TiB":
		multiplier = 1024 * 1024 * 1024 * 1024
	default:
		multiplier = 1 // fallback
	}
	return uint64(value * multiplier), nil
}
