package httpserver

import (
	"context"
	"github.com/AZRV17/zlib-backend/internal/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HTTPServer struct {
	httpServer *http.Server
}

func NewHTTPServer(cfg *config.Config, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		httpServer: &http.Server{
			Addr:    cfg.HTTP.Host + ":" + cfg.HTTP.Port,
			Handler: handler,

			ReadHeaderTimeout: 5 * time.Second,
		},
	}
}

// Run Функция для запуска сервера
func (s *HTTPServer) Run() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown Функция для остановки сервера
func (s *HTTPServer) Shutdown(stopped chan struct{}) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sigint
	log.Println("got interruption signal")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP Server Shutdown Error: %v", err)
	}
	close(stopped)
}