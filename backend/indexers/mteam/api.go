package mteam

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/charleshuang3/autoget/backend/indexers"
	"github.com/charleshuang3/autoget/backend/internal/scraper"
	"github.com/gocolly/colly"

	"github.com/rs/zerolog/log"
)

const (
	CategoriesCacheDuration = 24 * time.Hour
	name                    = "m-team"
)

var (
	fatalOnUnknownRootCategory = false
)

type Config struct {
	BaseURL      string `yaml:"base_url"`
	APIKey       string `yaml:"api_key"`
	BlockGayPorn bool   `yaml:"block_gay_porn"`
}

const defaultBaseURL = "https://api.m-team.cc"

func (c *Config) GetBaseURL() string {
	if c.BaseURL == "" {
		return defaultBaseURL
	}
	return c.BaseURL
}

type MTeam struct {
	scraper.Scraper

	config *Config

	lock                 sync.RWMutex
	categories           []*indexers.Category
	categoriesExpiryTime time.Time
}

func NewMTeam(config *Config) *MTeam {
	if config.APIKey == "" {
		return nil
	}
	m := &MTeam{
		Scraper: *scraper.NewScraper(),
		config:  config,
	}

	m.Categories()
	return m
}

func (m *MTeam) Name() string {
	return name
}

type ListCategories struct {
	Data struct {
		List []struct {
			CreatedDate      string `json:"createdDate"`
			LastModifiedDate string `json:"lastModifiedDate"`
			ID               string `json:"id"`
			Order            string `json:"order"`
			NameChs          string `json:"nameChs"`
			NameCht          string `json:"nameCht"`
			NameEng          string `json:"nameEng"`
			Image            string `json:"image"`
			Parent           string `json:"parent"`
		} `json:"list"`

		// We don't use following fields because they don't contains
		// all subcategories. For example the parent of tvshow(105).
		Adult  []string `json:"adult"`
		Movie  []string `json:"movie"`
		Music  []string `json:"music"`
		Tvshow []string `json:"tvshow"`

		// We don't use following fields
		Waterfall []string `json:"waterfall"`
	} `json:"data"`

	// We don't use following fields
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (l *ListCategories) toCategories(blockGayPorn bool) []*indexers.Category {
	adultRoot := &indexers.Category{
		ID:   "root", // Temporary ID before prefixing
		Name: CategoryAdult,
	}
	normalRoot := &indexers.Category{
		ID:   "root", // Temporary ID before prefixing
		Name: CategoryNormal,
	}
	res := []*indexers.Category{
		adultRoot,
		normalRoot,
	}

	categories := map[string]*indexers.Category{
		CategoryAdult:  adultRoot,
		CategoryNormal: normalRoot,
	}

	for _, cat := range l.Data.List {
		if blockGayPorn && cat.ID == CategoryGayPorn {
			continue
		}
		categories[cat.ID] = &indexers.Category{
			ID:   cat.ID,
			Name: cat.NameChs,
		}
	}

	for _, cat := range l.Data.List {
		parent := cat.Parent
		if parent == "" {
			var ok bool
			parent, ok = rootCategories[cat.ID]
			if !ok {
				if fatalOnUnknownRootCategory {
					log.Fatal().Msgf("Got unknown root category: %s %s", cat.ID, cat.NameChs)
				} else {
					log.Warn().Msgf("Got unknown root category: %s %s", cat.ID, cat.NameChs)
					// Maybe some new root category, mark it in normal.
					parent = CategoryNormal
				}
			}
		}

		p := categories[parent]
		p.SubCategories = append(p.SubCategories, categories[cat.ID])
	}

	// visit all sub-categories apply normal / adult prefix
	buildIDForSubCategories(adultRoot, CategoryAdult)
	buildIDForSubCategories(normalRoot, CategoryNormal)

	return res
}

func buildIDForSubCategories(c *indexers.Category, prefix string) {
	c.ID = fmt.Sprintf(indexers.CategoryIdFormat, name, prefix+"-"+c.ID)
	for _, sub := range c.SubCategories {
		buildIDForSubCategories(sub, prefix)
	}
}

const (
	CategoryAdult   = "adult"
	CategoryNormal  = "normal"
	CategoryGayPorn = "440"
)

var (
	rootCategories = map[string]string{
		"100": CategoryNormal, // Movie
		"105": CategoryNormal, // TV Series
		"444": CategoryNormal, // Documentary
		"110": CategoryNormal, // Music
		"443": CategoryNormal, // edu
		"447": CategoryNormal, // Game
		"449": CategoryNormal, // Anime
		"450": CategoryNormal, // Others
		"115": CategoryAdult,  // AV Censored
		"120": CategoryAdult,  // AV Uncensored
		"445": CategoryAdult,  // IV
		"446": CategoryAdult,  // HCG
	}
)

func (m *MTeam) authHeader() http.Header {
	h := http.Header{}
	h.Add("x-api-key", m.config.APIKey)
	return h
}

func (m *MTeam) Categories() ([]*indexers.Category, error) {
	var err error
	now := time.Now()
	m.lock.RLock()
	// catched categories not exist or expired
	if len(m.categories) == 0 || m.categoriesExpiryTime.Before(now) {
		m.lock.RUnlock()
		m.lock.Lock()
		defer m.lock.Unlock()

		// maybe cache updated
		if len(m.categories) == 0 || m.categoriesExpiryTime.Before(now) {
			c := m.C.Clone()
			c.OnResponse(func(r *colly.Response) {
				resp := &ListCategories{}
				if err = json.Unmarshal(r.Body, resp); err != nil {
					return
				}
				m.categories = resp.toCategories(m.config.BlockGayPorn)
			})

			err = c.Request(http.MethodPost, m.config.GetBaseURL()+"/api/torrent/categoryList", nil, nil, m.authHeader())
		}
	} else {
		m.lock.RUnlock()
	}

	return m.categories, err
}

func (m *MTeam) List(categories []string, keyword string, page, pageSize uint32) ([]indexers.Resource, error) {
	return nil, nil
}

func (m *MTeam) Detail(id string) (indexers.Resource, error) {
	return indexers.Resource{}, nil
}
