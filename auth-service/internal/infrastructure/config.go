package infrastructure

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port    int    `yaml:"port"`
		Mode    string `yaml:"mode"`
		BaseURL string `yaml:"base_url"`
	} `yaml:"server"`
	Mongo struct {
		URI    string `yaml:"uri"`
		DBName string `yaml:"dbname"`
	} `yaml:"mongo"`
	Auth struct {
		AccessSecret   string `yaml:"access_secret"`
		RefreshSecret  string `yaml:"refresh_secret"`
		AccessTTLMin   int    `yaml:"access_ttl_min"`
		RefreshTTLHour int    `yaml:"refresh_ttl_hour"`
	} `yaml:"auth"`
	Email struct {
		From string `yaml:"from"`
		SMTP struct {
			Host string `yaml:"host"`
			Port string `yaml:"port"`
			User string `yaml:"user"`
			Pass string `yaml:"pass"`
		} `yaml:"smtp"`
	} `yaml:"email"`
	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
}

func Load(path string) (*Config, error) {
	cfg := &Config{}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}