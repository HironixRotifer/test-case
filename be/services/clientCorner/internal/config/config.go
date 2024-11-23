package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBPassword string `env:"DB_PASSWORD"`
	DBUser     string `env:"DB_USER"`
	DBName     string `env:"DB_NAME"`
	SSLMode    string `env:"DB_SSLMODE"`

	KafkaAddrs []string `env:"kafka_addrs"`

	LokiURL string `env:"loki_URL"`
	Job     string `env:"job"`
	// TODO:
}

func MustLoadConfig(configPath string) *Config {
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	err = cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		panic("cannot read config: " + configPath)
	}

	return &cfg
}
