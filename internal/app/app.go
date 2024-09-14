package app

import (
	"github.com/AZRV17/zlib-backend/internal/config"
	"log"
	"log/slog"
)

func Run() {
	cfg, err := config.NewConfig("internal/config/config.yaml")
	if err != nil {
		log.Fatal("error loading config: ", err)
	}

	slog.Info("config", slog.Any("config", cfg))
}
