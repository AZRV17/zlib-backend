package delivery

import (
	"github.com/AZRV17/zlib-backend/internal/config"
	"github.com/AZRV17/zlib-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
	config  *config.Config
}

func NewHandler(service service.Service, cfg *config.Config) *Handler {
	return &Handler{
		service: service,
		config:  cfg,
	}
}

func (h *Handler) Init(r *gin.Engine) {
	h.initUserRoutes(r)
}
