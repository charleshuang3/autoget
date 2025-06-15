package mteam

import (
	"time"

	_ "embed"

	"github.com/charleshuang3/autoget/backend/indexers"
	"github.com/charleshuang3/autoget/backend/indexers/mteam/prefetcheddata"
	"github.com/charleshuang3/autoget/backend/internal/errors"
	"github.com/charleshuang3/autoget/backend/internal/notify"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	_ indexers.IIndexer = (*MTeam)(nil)

	logger = log.With().Str("indexer", name).Logger()
)

const (
	name           = "m-team"
	defaultBaseURL = "https://api.m-team.cc"

	categoryAdult   = "adult"
	categoryNormal  = "normal"
	categoryGayPorn = "440"

	httpTimeout = 5 * time.Second
)

var (
	//go:embed prefetcheddata/data.json
	categoriesJSON []byte
)

type Config struct {
	BaseURL           string `yaml:"base_url"`
	APIKey            string `yaml:"api_key"`
	ExcludeGayContent bool   `yaml:"exclude_gay_content"`
	RSS               string `yaml:"rss"`

	Downloader string `yaml:"downloader"`
}

func (c *Config) getBaseURL() string {
	if c.BaseURL == "" {
		return defaultBaseURL
	}
	return c.BaseURL
}

type MTeam struct {
	indexers.IndexerBasicInfo

	config *Config
	db     *gorm.DB
	notify notify.INotifier

	prefetched *prefetcheddata.Data
	standards  map[string]string

	torrentsDir string
}

func NewMTeam(config *Config, torrentsDir string, db *gorm.DB, notify notify.INotifier) *MTeam {
	if config.APIKey == "" {
		return nil
	}
	m := &MTeam{
		IndexerBasicInfo: *indexers.NewIndexerBasicInfo(name, true),
		config:           config,
		db:               db,
		standards:        map[string]string{},
		torrentsDir:      torrentsDir,
		notify:           notify,
	}

	var err error
	m.prefetched, err = prefetcheddata.Read()
	if err != nil {
		logger.Fatal().Err(err).Msgf("Failed to read prefetched data: %v", err)
	}

	for k, v := range m.prefetched.Standards {
		m.standards[v] = k
	}

	return m
}

func (m *MTeam) Categories() ([]indexers.Category, *errors.HTTPStatusError) {
	return m.prefetched.Categories.Tree, nil
}
