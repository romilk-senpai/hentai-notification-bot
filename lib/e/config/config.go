package config

import (
	"errors"
	"fmt"
	parseflag "hentai-notification-bot-re/flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string     `yaml:"env" env-default:"local"`
	HTMLPath    string     `yaml:"html_path" env-required:"true"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func Load() (*Config, error) {
	if parseflag.ConfigPath == "" {
		return nil, errors.New("config path is empty")
	}

	if _, err := os.Stat(parseflag.ConfigPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist at given path %s", parseflag.ConfigPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(parseflag.ConfigPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %s", err.Error())
	}

	return &cfg, nil
}
