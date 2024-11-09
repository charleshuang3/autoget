package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BasePath string `yaml:"base_path"`
	Port     int    `yaml:"port"`
	Prowlarr struct {
		APIURL string `yaml:"api_url"`
		APIKey string `yaml:"api_key"`
	} `yaml:"prowlarr"`
	Downloaders []struct {
		Name    string `yaml:"name"`
		SeedDir string `yaml:"seed_dir"`
	} `yaml:"downloaders"`
	Database struct {
		SqlitePath string `yaml:"sqlite.path"`
	} `yaml:"database"`
	Auth struct {
		Type       string `yaml:"type"`
		APIKey     string `yaml:"api_key,omitempty"`
		GoogleAuth struct {
			Project               string   `yaml:"project"`
			ClientID              string   `yaml:"client_id"`
			ClientSecret          string   `yaml:"client_secret"`
			AllowedEmailAddresses []string `yaml:"allowed_email_addresses"`
		} `yaml:"google_auth,omitempty"`
	} `yaml:"auth"`
	TMDB struct {
		APIKey string `yaml:"api_key"`
	} `yaml:"tmdb"`
	CompleteDir string            `yaml:"complete_dir"`
	LibraryDirs map[string]string `yaml:"library_dirs"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("could not decode config file: %w", err)
	}

	return &cfg, nil
}

func (cfg *Config) Validate() error {
	if cfg.BasePath == "" || cfg.Port == 0 || cfg.CompleteDir == "" {
		return fmt.Errorf("BasePath, Port, and CompleteDir must be set")
	}

	if cfg.TMDB.APIKey == "" {
		return fmt.Errorf("TMDB API Key must be set")
	}

	for _, downloader := range cfg.Downloaders {
		if downloader.Name == "" || downloader.SeedDir == "" {
			return fmt.Errorf("Downloader Name and SeedDir must be set")
		}
	}

	if cfg.Database.SqlitePath == "" {
		return fmt.Errorf("Database SqlitePath must be set")
	}

	if cfg.Prowlarr.APIURL == "" || cfg.Prowlarr.APIKey == "" {
		return fmt.Errorf("Prowlarr APIURL and APIKey must be set")
	}

	switch cfg.Auth.Type {
	case "api_key":
		if cfg.Auth.APIKey == "" {
			return fmt.Errorf("API Key must be set for api_key auth type")
		}
	case "google_auth":
		if cfg.Auth.GoogleAuth.Project == "" || cfg.Auth.GoogleAuth.ClientID == "" ||
			cfg.Auth.GoogleAuth.ClientSecret == "" || len(cfg.Auth.GoogleAuth.AllowedEmailAddresses) == 0 {
			return fmt.Errorf("All google_auth fields must be set for google_auth auth type")
		}
	case "none":
	default:
		return fmt.Errorf("Invalid auth type: %s", cfg.Auth.Type)
	}
	return nil
}
