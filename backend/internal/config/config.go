package config

import (
	"os"

	"github.com/charleshuang3/autoget/backend/indexers/mteam"
	"github.com/charleshuang3/autoget/backend/indexers/nyaa"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Port     string `yaml:"port"`
	ProxyURL string `yaml:"proxy_url"`

	MTeam   *mteam.Config `yaml:"mteam"`
	Nyaa    *nyaa.Config  `yaml:"nyaa"`
	Sukebei *nyaa.Config  `yaml:"sukebei"`
}

func ReadConfig(path string) (*Config, error) {
	config := &Config{}

	b, err := os.ReadFile(path)
	yaml.Unmarshal(b, config)
	if err != nil {
		return nil, err
	}

	if config.Nyaa != nil {
		config.Nyaa.SetProxyURL(config.ProxyURL)
	}
	if config.Sukebei != nil {
		config.Sukebei.SetProxyURL(config.ProxyURL)
	}

	return config, nil
}
