package delivery

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initLogRoutes(r *gin.Engine) {
	logs := r.Group("/logs")
	{
		logs.Use(h.AuthMiddleware, h.AdminMiddleware).GET("/", h.getAllLogs)
	}
}

func (h *Handler) getAllLogs(c *gin.Context) {
	logs, err := h.service.LogServ.GetLogs()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}
