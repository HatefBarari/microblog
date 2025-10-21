package infrastructure

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`
	Mongo struct {
		URI    string `yaml:"uri"`
		DBName string `yaml:"dbname"`
	} `yaml:"mongo"`
	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
	Auth struct {
		AccessToken string `yaml:"access_token"`
		RefreshToken string `yaml:"refresh_token"`
	} `yaml:"auth"`
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