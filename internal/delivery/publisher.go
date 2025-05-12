package delivery

import (
	"fmt"
	"github.com/AZRV17/zlib-backend/internal/service"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) initPublisherRoutes(r *gin.Engine) {
	publishers := r.Group("/publishers")
	{
		publishers.GET("/", h.getPublishers)
		publishers.GET("/:id", h.getPublisherByID)
		publishers.Use(h.AuthMiddleware, h.LibrarianMiddleware).DELETE("/:id", h.deletePublisher)
		publishers.Use(h.AuthMiddleware, h.LibrarianMiddleware).POST("/", h.createPublisher)
		publishers.Use(h.AuthMiddleware, h.LibrarianMiddleware).PUT("/:id", h.updatePublisher)
		publishers.Use(h.AuthMiddleware, h.LibrarianMiddleware).GET("/export", h.exportPublishersToCSV)
		publishers.Use(h.AuthMiddleware, h.LibrarianMiddleware).POST("/import", h.importPublishersFromCSV)
	}
}

func (h *Handler) getPublishers(c *gin.Context) {
	publishers, err := h.service.PublisherServ.GetPublishers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, publishers)
}

func (h *Handler) getPublisherByID(c *gin.Context) {
	publisherID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	publisher, err := h.service.PublisherServ.GetPublisherByID(uint(publisherID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, publisher)
}

func (h *Handler) deletePublisher(c *gin.Context) {
	publisherID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.PublisherServ.DeletePublisher(uint(publisherID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cookie, err := c.Request.Cookie("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.LogServ.CreateLogWithCookie(cookie, "Удаление издателя")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "publisher deleted"})
}

func (h *Handler) updatePublisher(c *gin.Context) {
	var input service.UpdatePublisherInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.PublisherServ.UpdatePublisher(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cookie, err := c.Request.Cookie("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.LogServ.CreateLogWithCookie(cookie, "Изменение издателя")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "publisher updated"})
}

func (h *Handler) createPublisher(c *gin.Context) {
	var input service.CreatePublisherInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.PublisherServ.CreatePublisher(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cookie, err := c.Request.Cookie("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.LogServ.CreateLogWithCookie(cookie, "Создание издателя")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "publisher created"})
}

func (h *Handler) exportPublishersToCSV(c *gin.Context) {
	publisherData, err := h.service.PublisherServ.ExportPublishersToCSV()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := "publishers.csv"

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")
	c.Header("Content-Length", fmt.Sprint(len(publisherData)))

	c.Data(http.StatusOK, "text/csv", publisherData)

	cookie, err := c.Request.Cookie("id")
	if err != nil {
		log.Printf("Error getting cookie for logging: %v", err)
		return
	}

	err = h.service.LogServ.CreateLogWithCookie(cookie, "Экспорт издательств в CSV")
	if err != nil {
		log.Printf("Error creating log: %v", err)
	}
}

func (h *Handler) importPublishersFromCSV(c *gin.Context) {
	// Получаем файл из запроса
	file, err := c.FormFile("csv")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось получить CSV файл: " + err.Error()})
		return
	}

	// Проверяем расширение файла
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".csv") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "файл должен быть в формате CSV"})
		return
	}

	// Открываем файл
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось открыть файл: " + err.Error()})
		return
	}
	defer src.Close()

	// Читаем содержимое файла
	fileData, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось прочитать файл: " + err.Error()})
		return
	}

	// Импортируем издательства
	count, err := h.service.PublisherServ.ImportPublishersFromCSV(fileData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при импорте издательств: " + err.Error()})
		return
	}

	// Логируем действие
	cookie, err := c.Request.Cookie("id")
	if err != nil {
		log.Printf("Error getting cookie for logging: %v", err)
	} else {
		err = h.service.LogServ.CreateLogWithCookie(cookie, "Импорт издательств из CSV")
		if err != nil {
			log.Printf("Error creating log: %v", err)
		}
	}

	c.JSON(
		http.StatusOK, gin.H{
			"message": fmt.Sprintf("Успешно импортировано %d издательств", count),
			"count":   count,
		},
	)
}
