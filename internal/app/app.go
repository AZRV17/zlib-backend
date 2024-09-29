package app

import (
	"errors"
	"fmt"
	"github.com/AZRV17/zlib-backend/internal/config"
	"github.com/AZRV17/zlib-backend/internal/delivery"
	"github.com/AZRV17/zlib-backend/internal/repository"
	httpserver "github.com/AZRV17/zlib-backend/internal/server/http"
	serv "github.com/AZRV17/zlib-backend/internal/service"
	"github.com/AZRV17/zlib-backend/pkg/db/psql"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"net/http"
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

	repo := repository.NewRepository(psql.DB)

	service := serv.NewService(&repo)

	r := gin.Default()

	handler := delivery.NewHandler(*service, cfg)

	handler.Init(r)

	server := httpserver.NewHttpServer(cfg, r)

	stoppedHTTP := make(chan struct{})

	go server.Shutdown(stoppedHTTP)

	slog.Info("starting HTTP server", slog.Any("host", cfg.HTTP.Host), slog.Any("port", cfg.HTTP.Port))

	go func() {
		if err := server.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server ListenAndServe Error: %v", err)
		}
	}()

	<-stoppedHTTP

	slog.Info("HTTP server stopped")
}
