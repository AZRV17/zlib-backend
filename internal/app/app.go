package app

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/AZRV17/zlib-backend/internal/config"
	"github.com/AZRV17/zlib-backend/internal/delivery"
	"github.com/AZRV17/zlib-backend/internal/repository"
	httpserver "github.com/AZRV17/zlib-backend/internal/server/http"
	ws "github.com/AZRV17/zlib-backend/internal/server/websocket"
	serv "github.com/AZRV17/zlib-backend/internal/service"
	"github.com/AZRV17/zlib-backend/pkg/db/psql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//nolint:funlen
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

	service := serv.NewService(&repo, psql.DB, cfg)

	r := gin.Default()

	r.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     []string{"http://localhost:3000"},
				AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Set-Cookie", "Accept"},
				AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization"},
				AllowWildcard:    true,
				AllowCredentials: true,
				MaxAge:           12 * time.Hour,
			},
		),
	)

	if _, err := os.Stat("uploads/books"); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll("uploads/books", 0755); err != nil {
			log.Fatalf("Failed to create uploads directory: %v", err)
		}
	}

	r.Static("/uploads", "./uploads")

	chatHub := ws.NewChatHub(repo.ChatRepo, service.UserServ)

	go chatHub.HandleMessages()

	handler := delivery.NewHandler(*service, cfg, chatHub)

	handler.Init(r)

	server := httpserver.NewHTTPServer(cfg, r)

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
