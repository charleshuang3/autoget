package mteam

import (
	"encoding/json"
	"net/http"

	_ "embed"

	"github.com/charleshuang3/autoget/backend/indexers"
	"github.com/charleshuang3/autoget/backend/internal/errors"
	"github.com/charleshuang3/autoget/backend/internal/scraper"
	"github.com/rs/zerolog/log"
)

const (
	name           = "m-team"
	defaultBaseURL = "https://api.m-team.cc"

	categoryAdult   = "adult"
	categoryNormal  = "normal"
	categoryGayPorn = "440"
)

var (
	//go:embed category/categories.json
	categoriesJSON []byte
)

type Config struct {
	BaseURL           string `yaml:"base_url"`
	APIKey            string `yaml:"api_key"`
	ExcludeGayContent bool   `yaml:"exclude_gay_content"`
}

func (c *Config) GetBaseURL() string {
	if c.BaseURL == "" {
		return defaultBaseURL
	}
	return c.BaseURL
}

type Categories struct {
	CategoryTree   []indexers.Category `json:"category_tree"`
	CategoryToMode map[string]string   `json:"modes"`
	Categories     map[string]string   `json:"categories"`
}

type MTeam struct {
	indexers.IndexerBasicInfo
	scraper.Scraper

	config *Config

	categories Categories
}

func NewMTeam(config *Config) *MTeam {
	if config.APIKey == "" {
		return nil
	}
	m := &MTeam{
		IndexerBasicInfo: *indexers.NewIndexerBasicInfo(name, true),
		Scraper:          *scraper.NewScraper(),
		config:           config,
	}

	if err := json.Unmarshal(categoriesJSON, &m.categories); err != nil {
		log.Fatal().Err(err).Msgf("Failed to unmarshal categories: %v", err)
	}

	return m
}

func (m *MTeam) authHeader() http.Header {
	h := http.Header{}
	h.Add("x-api-key", m.config.APIKey)
	return h
}

func (m *MTeam) Categories() ([]indexers.Category, *errors.HTTPStatusError) {
	return m.categories.CategoryTree, nil
}

func (m *MTeam) Detail(id string) (indexers.Resource, *errors.HTTPStatusError) {
	return indexers.Resource{}, nil
}
