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

func (h *Handler) initAuthorRoutes(r *gin.Engine) {
	authors := r.Group("/authors")
	{
		authors.GET("/", h.getAuthors)
		authors.GET("/:id", h.getAuthorByID)
		authors.Use(h.AuthMiddleware, h.LibrarianMiddleware).POST("/", h.createAuthor)
		authors.Use(h.AuthMiddleware, h.LibrarianMiddleware).PUT("/:id", h.updateAuthor)
		authors.Use(h.AuthMiddleware, h.LibrarianMiddleware).DELETE("/:id", h.deleteAuthor)
		authors.Use(h.AuthMiddleware, h.LibrarianMiddleware).GET("/export", h.exportAuthorsToCSV)
		authors.Use(h.AuthMiddleware, h.LibrarianMiddleware).POST("/import", h.importAuthorsFromCSV)
	}
}

func (h *Handler) getAuthors(c *gin.Context) {
	authors, err := h.service.AuthorServ.GetAuthors()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, authors)
}

func (h *Handler) getAuthorByID(c *gin.Context) {
	authorID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	author, err := h.service.AuthorServ.GetAuthorByID(uint(authorID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, author)
}

type createAuthorInput struct {
	Name      string    `json:"name"`
	Lastname  string    `json:"lastname"`
	Biography string    `json:"biography"`
	Birthdate time.Time `json:"birthdate"`
}

func (h *Handler) createAuthor(c *gin.Context) {
	var input createAuthorInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.AuthorServ.CreateAuthor(
		&service.CreateAuthorInput{
			Name:      input.Name,
			Lastname:  input.Lastname,
			Biography: input.Biography,
			Birthdate: input.Birthdate,
		},
	)
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
		Action:  "Создание автора",
		Date:    time.Now(),
		Details: fmt.Sprintf("Создание автора: %s %s", input.Name, input.Lastname),
	}

	err = h.service.LogServ.CreateLog(createLogInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "author created"})
}

type updateAuthorInput struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Lastname  string    `json:"lastname"`
	Biography string    `json:"biography"`
	Birthdate time.Time `json:"birthdate"`
}

func (h *Handler) updateAuthor(c *gin.Context) {
	var input updateAuthorInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.AuthorServ.UpdateAuthor(
		&service.UpdateAuthorInput{
			ID:        input.ID,
			Name:      input.Name,
			Lastname:  input.Lastname,
			Biography: input.Biography,
			Birthdate: input.Birthdate,
		},
	)
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
		Action:  "Обновление автора",
		Date:    time.Now(),
		Details: fmt.Sprintf("Обновление автора ID: %d", input.ID),
	}

	err = h.service.LogServ.CreateLog(createLogInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "author updated"})
}

func (h *Handler) deleteAuthor(c *gin.Context) {
	authorID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.AuthorServ.DeleteAuthor(uint(authorID))
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
		Action:  "Удаление автора",
		Date:    time.Now(),
		Details: fmt.Sprintf("Удаление автора ID: %d", authorID),
	}

	err = h.service.LogServ.CreateLog(createLogInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "author deleted"})
}

func (h *Handler) exportAuthorsToCSV(c *gin.Context) {
	authorData, err := h.service.AuthorServ.ExportAuthorsToCSV()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := "authors.csv"
	// Устанавливаем заголовки для скачивания файла
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")
	c.Header("Content-Length", fmt.Sprint(len(authorData)))

	c.Data(http.StatusOK, "text/csv", authorData)

	// Логируем действие
	userID, err := getUserIDFromContext(c)
	if err != nil {
		log.Printf("Error getting user ID for logging: %v", err)
		return
	}

	createLogInput := &service.CreateLogInput{
		UserID:  userID,
		Action:  "Экспорт авторов в CSV",
		Date:    time.Now(),
		Details: "Экспорт авторов в CSV файл",
	}

	err = h.service.LogServ.CreateLog(createLogInput)
	if err != nil {
		log.Printf("Error creating log: %v", err)
	}
}

func (h *Handler) importAuthorsFromCSV(c *gin.Context) {
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

	// Импортируем авторов
	count, err := h.service.AuthorServ.ImportAuthorsFromCSV(fileData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при импорте авторов: " + err.Error()})
		return
	}

	// Логируем действие
	cookie, err := c.Request.Cookie("id")
	if err != nil {
		log.Printf("Error getting cookie for logging: %v", err)
	} else {
		err = h.service.LogServ.CreateLogWithCookie(cookie, "Импорт авторов из CSV")
		if err != nil {
			log.Printf("Error creating log: %v", err)
		}
	}

	c.JSON(
		http.StatusOK, gin.H{
			"message": fmt.Sprintf("Успешно импортировано %d авторов", count),
			"count":   count,
		},
	)
}
