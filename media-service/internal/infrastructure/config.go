package infrastructure

import (
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Storage  StorageConfig  `yaml:"storage"`
	Media    MediaConfig    `yaml:"media"`
	Log      LogConfig      `yaml:"log"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type DatabaseConfig struct {
	URI      string `yaml:"uri"`
	Database string `yaml:"database"`
}

type StorageConfig struct {
	Type    string `yaml:"type"`
	BasePath string `yaml:"base_path"`
	BaseURL  string `yaml:"base_url"`
}

type MediaConfig struct {
	MaxFileSize   int64    `yaml:"max_file_size"`
	AllowedTypes  []string `yaml:"allowed_types"`
	ThumbnailSize int      `yaml:"thumbnail_size"`
}

type LogConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Override with environment variables
	overrideFromEnv(&config)

	return &config, nil
}

func overrideFromEnv(config *Config) {
	if port := os.Getenv("SERVER_PORT"); port != "" {
		config.Server.Port = port
	}
	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Server.Host = host
	}
	if uri := os.Getenv("DATABASE_URI"); uri != "" {
		config.Database.URI = uri
	}
	if db := os.Getenv("DATABASE_NAME"); db != "" {
		config.Database.Database = db
	}
	if basePath := os.Getenv("STORAGE_BASE_PATH"); basePath != "" {
		config.Storage.BasePath = basePath
	}
	if baseURL := os.Getenv("STORAGE_BASE_URL"); baseURL != "" {
		config.Storage.BaseURL = baseURL
	}
	if maxSize := os.Getenv("MEDIA_MAX_FILE_SIZE"); maxSize != "" {
		if size, err := strconv.ParseInt(maxSize, 10, 64); err == nil {
			config.Media.MaxFileSize = size
		}
	}
	if allowedTypes := os.Getenv("MEDIA_ALLOWED_TYPES"); allowedTypes != "" {
		config.Media.AllowedTypes = strings.Split(allowedTypes, ",")
	}
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.Log.Level = level
	}
}
