package downloaders

import (
	"fmt"
)

type TransmissionConfig struct {
	URL         string `yaml:"url"`
	TorrentsDir string `yaml:"torrents_dir"`
	DownloadDir string `yaml:"download_dir"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
}

func (c *TransmissionConfig) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("transmission RPC URL is required")
	}
	if c.TorrentsDir == "" {
		return fmt.Errorf("torrents directory is required")
	}
	if c.DownloadDir == "" {
		return fmt.Errorf("download directory is required")
	}
	return nil
}

// SeedingPolicy we use at least X MB uploaded in last Y days as
// a condition to continue seeding.
type SeedingPolicy struct {
	IntervalInDays    uint64 `yaml:"interval_in_days"`
	UploadAtLeastInMB uint64 `yaml:"upload_at_least_in_mb"`
}

func (p *SeedingPolicy) Validate() error {
	if p.IntervalInDays == 0 {
		return fmt.Errorf("interval in days is required")
	}
	if p.UploadAtLeastInMB == 0 {
		return fmt.Errorf("upload at least in MB is required")
	}
	return nil
}

type DownloaderConfig struct {
	Transmission  *TransmissionConfig `yaml:"transmission"`
	SeedingPolicy *SeedingPolicy      `yaml:"seeding_policy"`
}

func (c *DownloaderConfig) Validate() error {
	if c.Transmission == nil {
		return fmt.Errorf("transmission config is required")
	}
	if err := c.Transmission.Validate(); err != nil {
		return err
	}
	if c.SeedingPolicy != nil {
		if err := c.SeedingPolicy.Validate(); err != nil {
			return err
		}
	}
	return nil
}
