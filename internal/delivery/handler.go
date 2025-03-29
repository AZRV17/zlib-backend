package delivery

import (
	"github.com/AZRV17/zlib-backend/internal/config"
	"github.com/AZRV17/zlib-backend/internal/server/websocket"
	"github.com/AZRV17/zlib-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
	config  *config.Config
	chatHub *websocket.ChatHub
}

func NewHandler(service service.Service, cfg *config.Config, chatHub *websocket.ChatHub) *Handler {
	return &Handler{
		service: service,
		config:  cfg,
		chatHub: chatHub,
	}
}

func (h *Handler) Init(r *gin.Engine) {
	// WebSocket маршрут должен быть инициализирован до остальных маршрутов
	r.GET("/ws", h.chatHub.HandleConnections)

	h.initUserRoutes(r)
	h.initBookRoutes(r)
	h.initPublisherRoutes(r)
	h.initAuthorRoutes(r)
	h.initGenreRoutes(r)
	h.initFavoriteRoutes(r)
	h.initReservationRoutes(r)
	h.initReviewRoutes(r)
	h.initUniqueCodeRoutes(r)
	h.initBackupRoutes(r)
	h.initLogRoutes(r)
	h.initChatRoutes(r)
}
