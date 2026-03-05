package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server   ServerConfig   `toml:"server"`
	Storage  StorageConfig  `toml:"storage"`
	Media    MediaConfig    `toml:"media"`
	FFmpeg   FFmpegConfig   `toml:"ffmpeg"`
	Security SecurityConfig `toml:"security"`
}

type ServerConfig struct {
	Host                  string `toml:"host"`
	Port                  int    `toml:"port"`
	SessionLifetimeHours  int    `toml:"session_lifetime_hours"`
}

type StorageConfig struct {
	MediaRoot string `toml:"media_root"`
	ThumbRoot string `toml:"thumb_root"`
	DBPath    string `toml:"db_path"`
	TempDir   string `toml:"temp_dir"`
}

type MediaConfig struct {
	ThumbMaxWidth      int `toml:"thumb_max_width"`
	ThumbAVIFCRF       int `toml:"thumb_avif_crf"`
	ThumbWebPQuality   int `toml:"thumb_webp_quality"`
	WorkerCount        int `toml:"worker_count"`
	MaxUploadSizeMB    int `toml:"max_upload_size_mb"`
}

type FFmpegConfig struct {
	Bin   string `toml:"bin"`
	UseQSV bool   `toml:"use_qsv"`
}

type SecurityConfig struct {
	MaxLoginAttempts          int `toml:"max_login_attempts"`
	LockoutDurationMinutes    int `toml:"lockout_duration_minutes"`
	DefaultShareExpiryDays    int `toml:"default_share_expiry_days"`
	MaxShareExpiryDays        int `toml:"max_share_expiry_days"`
}

// Load reads a TOML configuration file and returns the Config struct
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: failed to read file: %w", err)
	}

	cfg := &Config{}
	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("config: failed to unmarshal TOML: %w", err)
	}

	// Set defaults if not provided
	cfg.setDefaults()

	return cfg, nil
}

// setDefaults sets default values for configuration
func (c *Config) setDefaults() {
	if c.Server.Host == "" {
		c.Server.Host = "127.0.0.1"
	}
	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	if c.Server.SessionLifetimeHours == 0 {
		c.Server.SessionLifetimeHours = 24
	}

	if c.Storage.TempDir == "" {
		c.Storage.TempDir = "/tmp/simple-gallery"
	}

	if c.Media.ThumbMaxWidth == 0 {
		c.Media.ThumbMaxWidth = 400
	}
	if c.Media.ThumbAVIFCRF == 0 {
		c.Media.ThumbAVIFCRF = 30
	}
	if c.Media.ThumbWebPQuality == 0 {
		c.Media.ThumbWebPQuality = 80
	}
	if c.Media.WorkerCount == 0 {
		c.Media.WorkerCount = 2
	}
	if c.Media.MaxUploadSizeMB == 0 {
		c.Media.MaxUploadSizeMB = 4096
	}

	if c.FFmpeg.Bin == "" {
		c.FFmpeg.Bin = "/usr/bin/ffmpeg"
	}

	if c.Security.MaxLoginAttempts == 0 {
		c.Security.MaxLoginAttempts = 5
	}
	if c.Security.LockoutDurationMinutes == 0 {
		c.Security.LockoutDurationMinutes = 15
	}
	if c.Security.DefaultShareExpiryDays == 0 {
		c.Security.DefaultShareExpiryDays = 7
	}
	if c.Security.MaxShareExpiryDays == 0 {
		c.Security.MaxShareExpiryDays = 30
	}
}
