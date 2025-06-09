package nyaa

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/charleshuang3/autoget/backend/indexers"
	"github.com/charleshuang3/autoget/backend/indexers/nyaa/prefetcheddata"
	"github.com/charleshuang3/autoget/backend/internal/errors"
	"github.com/rs/zerolog/log"
)

const (
	defaultBaseURL = "https://nyaa.si/"
)

type Config struct {
	BaseURL string `yaml:"base_url"`
}

type Client struct {
	config *Config
}

func (c *Config) GetBaseURL() string {
	if c.BaseURL == "" {
		return defaultBaseURL
	}
	return c.BaseURL
}

func NewClient(config *Config) *Client {
	return &Client{
		config: config,
	}
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
	return nil, nil

}

// Detail of a resource.
func (c *Client) Detail(id string, fileList bool) (*indexers.ResourceDetail, *errors.HTTPStatusError) {
	url, err := url.JoinPath(c.config.GetBaseURL(), "view", id)
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, fmt.Sprintf("failed to join path: %v", err))
	}

	resp, err := http.Get(url)
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
	return nil, nil

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
