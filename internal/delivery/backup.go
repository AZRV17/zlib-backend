package delivery

import (
	"fmt"
	"github.com/AZRV17/zlib-backend/pkg/db/psql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
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

	// Получаем данные бэкапа
	backupData, err := psql.BackupDatabase(dsn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем текущее время для имени файла
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("backup_%s.sql", timestamp)

	// Устанавливаем заголовки для скачивания файла
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")
	c.Header("Content-Length", fmt.Sprint(len(backupData)))

	// Отправляем файл клиенту
	c.Data(http.StatusOK, "application/octet-stream", backupData)

	// Логируем действие
	cookie, err := c.Request.Cookie("id")
	if err != nil {
		log.Printf("Error getting cookie for logging: %v", err)
		return
	}

	err = h.service.LogServ.CreateLogWithCookie(cookie, "Создание бэкапа")
	if err != nil {
		log.Printf("Error creating log: %v", err)
	}
}

func (h *Handler) restoreBackup(c *gin.Context) {
	// Получаем файл из формы
	file, err := c.FormFile("backup")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "backup file is required"})
		return
	}

	// Открываем загруженный файл
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot open uploaded file"})
		return
	}
	defer src.Close()

	// Восстанавливаем базу данных из файла
	err = psql.RestoreDatabase(
		h.config.Postgres.Host,
		h.config.Postgres.Port,
		h.config.Postgres.User,
		h.config.Postgres.Password,
		"testDB",
		src,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Логируем действие
	cookie, err := c.Request.Cookie("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.LogServ.CreateLogWithCookie(cookie, "Восстановление бд")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "backup restored successfully"})
}
