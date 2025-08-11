package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Port int `yaml:"port"`
	} `yaml:"app"`

	Database struct {
		Driver   string `yaml:"driver"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`

	dsn string

	Redis struct {
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		Password     string `yaml:"password"`
		DB           int    `yaml:"db"`
		RepoCacheTTL int    `yaml:"repo_cache_ttl"`
	} `yaml:"redis"`

	JWT struct {
		Secret          string `yaml:"secret"`
		Issuer          string `yaml:"issuer"`
		AccessTokenTTL  int    `yaml:"access_token_ttl"`
		RefreshTokenTTL int    `yaml:"refresh_token_ttl"`
	} `yaml:"jwt"`

	Log struct {
		Filename string `yaml:"file"`
	} `yaml:"log"`

	Storage struct {
		Provider string `yaml:"provider"`

		Local struct {
			BasePath  string `yaml:"base_path"`
			MaxSizeMB int    `yaml:"max_size_mb"`
		} `yaml:"local"`

		S3 struct {
			Bucket    string `yaml:"bucket"`
			Region    string `yaml:"region"`
			AccessKey string `yaml:"access_key"`
			SecretKey string `yaml:"secret_key"`
			Endpoint  string `yaml:"endpoint"`
			MaxSizeMB int    `yaml:"max_size_mb"`
		} `yaml:"s3"`

		Cache struct {
			Path               string `yaml:"path"`
			MaxAgeMin          int    `yaml:"max_age_min"`
			CleanupIntervalMin int    `yaml:"cleanup_interval_min"`
		} `yaml:"cache"`
	} `yaml:"storage"`
}

func Load(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(fmt.Errorf("failed to unmarshal yaml: %w", err))
	}

	cfg.loadDsn()
	return &cfg
}

func (cf *Config) loadDsn() {
	if cf.Database.User == "" || cf.Database.Password == "" || cf.Database.Host == "" || cf.Database.Port == 0 || cf.Database.Name == "" {
		fmt.Println("Warning: one or more database config fields are missing")
	}
	cf.dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cf.Database.User,
		cf.Database.Password,
		cf.Database.Host,
		cf.Database.Port,
		cf.Database.Name,
	)
}

func (cf *Config) GetDsn() string {
	return cf.dsn
}
