package config

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLoadConfig(t *testing.T) {
	// Simulating a basic correct configuration content
	content := `
base_path: "/autoget"
port: 8081
prowlarr:
  api_url: "http://localhost:9696"
  api_key: "valid_key"
downloaders:
  - name: "qBittorrent"
    seed_dir: "/seed_dir"
database:
  sqlite.path: "/my_database.db"
auth:
  type: "api_key"
  api_key: "your_api_key"
tmdb:
  api_key: "tmdb_key"
complete_dir: "/downloads"
library_dirs:
  Movies: "/movies"
`

	cfg, err := LoadConfigFromContent(content)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validation failed: %v", err)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	// Simulate creating a temp file for correct config
	content := `
base_path: "/autoget"
port: 8081
prowlarr:
  api_url: "http://localhost:9696"
  api_key: "valid_key"
downloaders:
  - name: "qBittorrent"
    seed_dir: "/seed_dir"
database:
  sqlite.path: "/my_database.db"
auth:
  type: "api_key"
  api_key: "your_api_key"
tmdb:
  api_key: "tmdb_key"
complete_dir: "/downloads"
library_dirs:
  Movies: "/movies"
`
	file, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	if _, err := file.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	file.Close()

	cfg, err := LoadConfig(file.Name())
	if err != nil {
		t.Fatalf("Failed to load config from file: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	// Test for non-existing file
	_, err = LoadConfig("/non/existing/path/config.yaml")
	if err == nil || os.IsNotExist(err) {
		t.Fatalf("Expected file not found error, got: %v", err)
	}
}

func TestValidateFailures(t *testing.T) {
	tests := []struct {
		name    string
		content string
		expect  string
	}{
		{
			name:    "base_path missing",
			content: ``,
			expect:  "BasePath, Port, and CompleteDir must be set",
		},
		{
			name:    "port missing",
			content: `base_path: "/"`,
			expect:  "BasePath, Port, and CompleteDir must be set",
		},
		{
			name: "CompleteDir missing",
			content: `
base_path: "/"
port: 1`,
			expect: "BasePath, Port, and CompleteDir must be set",
		},
		{
			name: "tmdb api_key missing",
			content: `
base_path: "/"
port: 1
complete_dir: "/downloads"`,
			expect: "TMDB API Key must be set",
		},
		{
			name: "Downloader Name missing",
			content: `
base_path: "/"
port: 1
complete_dir: "/downloads"
tmdb:
  api_key: "tmdb_key"
downloaders:
  - seed_dir: "/seed_dir"`,
			expect: "Downloader Name and SeedDir must be set",
		},
		{
			name: "Downloader SeedDir missing",
			content: `
base_path: "/"
port: 1
complete_dir: "/downloads"
tmdb:
  api_key: "tmdb_key"
downloaders:
  - name: "qBittorrent"`,
			expect: "Downloader Name and SeedDir must be set",
		},
		{
			name: "SqlitePath missing",
			content: `
base_path: "/"
port: 1
complete_dir: "/downloads"
tmdb:
  api_key: "tmdb_key"
downloaders:
  - name: "qBittorrent"
    seed_dir: "/seed_dir"
`,
			expect: "Database SqlitePath must be set",
		},
		{
			name: "Prowlarr missing",
			content: `
base_path: "/"
port: 1
complete_dir: "/downloads"
tmdb:
  api_key: "tmdb_key"
downloaders:
  - name: "qBittorrent"
    seed_dir: "/seed_dir"
database:
  sqlite.path: "/my_database.db"
`,
			expect: "Prowlarr APIURL and APIKey must be set",
		},
		{
			name: "Auth API key missing",
			content: `
base_path: "/autoget"
port: 8081
complete_dir: "/downloads"
tmdb:
  api_key: "tmdb_key"
database:
  sqlite.path: "/my_database.db"
prowlarr:
  api_url: "http://localhost:9696"
  api_key: "valid_key"
auth:
  type: "api_key"
`,
			expect: "API Key must be set for api_key auth type",
		},
		{
			name: "Auth GoogleAuth project missing",
			content: `
base_path: "/autoget"
port: 8081
complete_dir: "/downloads"
tmdb:
  api_key: "tmdb_key"
database:
  sqlite.path: "/my_database.db"
prowlarr:
  api_url: "http://localhost:9696"
  api_key: "valid_key"
auth:
  type: "google_auth"
  google_auth:
    client_id: "id"
    client_secret: "secret"
    allowed_email_addresses: ["example.com"]
`,
			expect: "All google_auth fields must be set for google_auth auth type",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := LoadConfigFromContent(tc.content)
			if err != nil {
				t.Fatalf("Failed to load config: %v", err)
			}

			if err := cfg.Validate(); err == nil || err.Error() != tc.expect {
				t.Fatalf("Expected error: %v, got: %v", tc.expect, err)
			}
		})
	}
}

// Helper function to load config from string content
func LoadConfigFromContent(content string) (*Config, error) {
	c := &Config{}
	err := yaml.Unmarshal([]byte(content), c)
	return c, err
}
