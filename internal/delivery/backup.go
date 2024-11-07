package delivery

import (
	"fmt"
	"github.com/AZRV17/zlib-backend/pkg/db/psql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initBackupRoutes(r *gin.Engine) {
	backup := r.Group("/backup")
	{
		backup.Use(h.AuthMiddleware, h.AdminMiddleware).POST("/backup", h.createBackup)
		backup.Use(h.AuthMiddleware, h.AdminMiddleware).POST("/restore", h.restoreBackup)
	}
}

func (h *Handler) createBackup(c *gin.Context) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		h.config.Postgres.Host,
		h.config.Postgres.Port,
		h.config.Postgres.User,
		h.config.Postgres.Password,
		h.config.Postgres.DB,
	)

	err := psql.BackupDatabase(dsn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "backup created"})
}

func (h *Handler) restoreBackup(c *gin.Context) {

	err := psql.RestoreDatabase(
		h.config.Postgres.Host,
		h.config.Postgres.Port,
		h.config.Postgres.User,
		h.config.Postgres.Password,
		h.config.Postgres.DB,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "backup restored"})
}
