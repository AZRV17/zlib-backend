package app

import (
	"fmt"
	"github.com/AZRV17/zlib-backend/internal/config"
	"github.com/AZRV17/zlib-backend/pkg/db/psql"
	"log"
	"log/slog"
)

func Run() {
	cfg, err := config.NewConfig("internal/config/config.yaml")
	if err != nil {
		log.Fatal("error loading config: ", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DB,
	)

	if err = psql.Connect(dsn); err != nil {
		log.Fatal("error connecting to db: ", err)
	}

	slog.Info("config", slog.Any("config", cfg))
	slog.Info("DB connected", slog.Any("DB", psql.DB.Name()))
}
