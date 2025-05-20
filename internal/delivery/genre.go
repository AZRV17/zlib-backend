package delivery

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AZRV17/zlib-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initGenreRoutes(r *gin.Engine) {
	genres := r.Group("/genres")
	{
		genres.GET("/", h.getAllGenres)
		genres.GET("/:id", h.getGenreByID)
		genres.Use(h.AuthMiddleware, h.LibrarianMiddleware).POST("/", h.createGenre)
		genres.Use(h.AuthMiddleware, h.LibrarianMiddleware).PUT("/:id", h.updateGenre)
		genres.Use(h.AuthMiddleware, h.LibrarianMiddleware).DELETE("/:id", h.deleteGenre)
		genres.Use(h.AuthMiddleware, h.LibrarianMiddleware).GET("/export", h.exportGenresToCSV)
		genres.Use(h.AuthMiddleware, h.LibrarianMiddleware).POST("/import", h.importGenresFromCSV)
	}
}

func (h *Handler) getAllGenres(c *gin.Context) {
	genres, err := h.service.GenreServ.GetGenres()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, genres)
}

func (h *Handler) getGenreByID(c *gin.Context) {
	genreID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	genre, err := h.service.GenreServ.GetGenreByID(uint(genreID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, genre)
}

func (h *Handler) createGenre(c *gin.Context) {
	var input service.CreateGenreInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.GenreServ.CreateGenre(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createLogInput := &service.CreateLogInput{
		UserID:  userID,
		Action:  "Создание жанра",
		Date:    time.Now(),
		Details: "Создание жанра: " + input.Name,
	}

	err = h.service.LogServ.CreateLog(createLogInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "genre created"})
}

func (h *Handler) updateGenre(c *gin.Context) {
	var input service.UpdateGenreInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.GenreServ.UpdateGenre(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createLogInput := &service.CreateLogInput{
		UserID:  userID,
		Action:  "Изменение жанра",
		Date:    time.Now(),
		Details: "Изменение жанра ID: " + fmt.Sprintf("%d", input.ID),
	}

	err = h.service.LogServ.CreateLog(createLogInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "genre updated"})
}

func (h *Handler) deleteGenre(c *gin.Context) {
	genreID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.GenreServ.DeleteGenre(uint(genreID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createLogInput := &service.CreateLogInput{
		UserID:  userID,
		Action:  "Удаление жанра",
		Date:    time.Now(),
		Details: "Удаление жанра ID: " + c.Param("id"),
	}

	err = h.service.LogServ.CreateLog(createLogInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "genre deleted"})
}

func (h *Handler) exportGenresToCSV(c *gin.Context) {
	genreData, err := h.service.GenreServ.ExportGenresToCSV()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := "genres.csv"

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")
	c.Header("Content-Length", fmt.Sprint(len(genreData)))

	c.Data(http.StatusOK, "text/csv", genreData)

	userID, err := getUserIDFromContext(c)
	if err != nil {
		log.Printf("Error getting user ID for logging: %v", err)
		return
	}

	createLogInput := &service.CreateLogInput{
		UserID:  userID,
		Action:  "Экспорт жанров в CSV",
		Date:    time.Now(),
		Details: "Экспорт жанров в CSV файл",
	}

	err = h.service.LogServ.CreateLog(createLogInput)
	if err != nil {
		log.Printf("Error creating log: %v", err)
	}
}

func (h *Handler) importGenresFromCSV(c *gin.Context) {
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

	// Импортируем жанры
	count, err := h.service.GenreServ.ImportGenresFromCSV(fileData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при импорте жанров: " + err.Error()})
		return
	}

	// Логируем действие
	userID, err := getUserIDFromContext(c)
	if err != nil {
		log.Printf("Error getting user ID for logging: %v", err)
	} else {
		createLogInput := &service.CreateLogInput{
			UserID:  userID,
			Action:  "Импорт жанров из CSV",
			Date:    time.Now(),
			Details: fmt.Sprintf("Успешно импортировано %d жанров", count),
		}

		err = h.service.LogServ.CreateLog(createLogInput)
		if err != nil {
			log.Printf("Error creating log: %v", err)
		}
	}

	c.JSON(
		http.StatusOK, gin.H{
			"message": fmt.Sprintf("Успешно импортировано %d жанров", count),
			"count":   count,
		},
	)
}
