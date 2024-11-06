package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"` // Порт базы данных
	DBPassword string `env:"DB_PASSWORD"`
	DBUser     string `env:"DB_USER"`
	DBName     string `env:"DB_NAME"`
	SSLMode    string `env:"DB_SSLMODE"`
}

func MustLoadConfig(configPath string) *Config {

	// check if file exists
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config
	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}
