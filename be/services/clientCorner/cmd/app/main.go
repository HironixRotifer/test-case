package main

import (
	"clientCorner/internal/config"
	"clientCorner/internal/lib/loki"
)

func main() {
	cfg := config.MustLoadConfig(".")
	log := loki.SetupLogger(cfg)
	log.Info()
}
