package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	KafkaAddrs []string `env:"kafka_addrs"`

	LokiURL string `env:"loki_url"`
	Job     string `env:"job"`
	// TODO:
}

func MustLoadConfig(configPath string) *Config {
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		panic("cannot read config: " + configPath)
	}

	return &cfg
}
