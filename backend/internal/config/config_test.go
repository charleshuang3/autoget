package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	// Test case 1: Config with Sukebei
	t.Run("Config with Sukebei", func(t *testing.T) {
		configContent := `
port: "8080"
proxy_url: "http://localhost:8888"
mteam:
  base_url: "http://mteam.example.com"
  api_key: "mteam_key"
nyaa:
  base_url: "http://nyaa.example.com"
  use_proxy: true
sukebei:
  base_url: "http://sukebei.example.com"
  use_proxy: true
`
		tmpFile, err := os.CreateTemp("", "config_with_sukebei_*.yaml")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name())
		_, err = tmpFile.WriteString(configContent)
		assert.NoError(t, err)
		tmpFile.Close()

		cfg, err := ReadConfig(tmpFile.Name())
		assert.NoError(t, err)
		assert.NotNil(t, cfg)

		assert.Equal(t, "8080", cfg.Port)
		assert.Equal(t, "http://localhost:8888", cfg.ProxyURL)
		assert.NotNil(t, cfg.MTeam)
		assert.Equal(t, "http://mteam.example.com", cfg.MTeam.BaseURL)
		assert.Equal(t, "mteam_key", cfg.MTeam.APIKey)
		assert.NotNil(t, cfg.Nyaa)
		assert.Equal(t, "http://nyaa.example.com", cfg.Nyaa.BaseURL)
		assert.True(t, cfg.Nyaa.UseProxy)
		assert.NotNil(t, cfg.Sukebei)
		assert.Equal(t, "http://sukebei.example.com", cfg.Sukebei.BaseURL)
		assert.True(t, cfg.Sukebei.UseProxy)
	})

	// Test case 2: Config without Sukebei
	t.Run("Config without Sukebei", func(t *testing.T) {
		configContent := `
port: "8081"
proxy_url: "http://localhost:9999"
mteam:
  base_url: "http://mteam.example.org"
  api_key: "mteam_key_2"
nyaa:
  base_url: "http://nyaa.example.org"
  use_proxy: false
`
		tmpFile, err := os.CreateTemp("", "config_without_sukebei_*.yaml")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name())
		_, err = tmpFile.WriteString(configContent)
		assert.NoError(t, err)
		tmpFile.Close()

		cfg, err := ReadConfig(tmpFile.Name())
		assert.NoError(t, err)
		assert.NotNil(t, cfg)

		assert.Equal(t, "8081", cfg.Port)
		assert.Equal(t, "http://localhost:9999", cfg.ProxyURL)
		assert.NotNil(t, cfg.MTeam)
		assert.Equal(t, "http://mteam.example.org", cfg.MTeam.BaseURL)
		assert.Equal(t, "mteam_key_2", cfg.MTeam.APIKey)
		assert.NotNil(t, cfg.Nyaa)
		assert.Equal(t, "http://nyaa.example.org", cfg.Nyaa.BaseURL)
		assert.False(t, cfg.Nyaa.UseProxy)
		assert.Nil(t, cfg.Sukebei) // Sukebei should be nil
	})
}
