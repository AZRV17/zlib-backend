package delivery

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) initBookRoutes(r *gin.Engine) {
	books := r.Group("/books")
	{
		books.GET("/", h.getAllBooks)
		books.GET("/:id", h.getBookByID)
		books.GET("/grouped", h.getGroupedBooksByTitle)
		books.GET("/pagination", h.getBooksWithPagination)
		books.GET("/:id/audio", h.getBookAudiobookFiles)
		books.GET("/audio/:id", h.serveAudioFile)
		books.Use(h.AuthMiddleware).GET("/:id/codes", h.getBookUniqueCodes)
		books.Use(h.AuthMiddleware).POST("/:id", h.reserveBook)
		books.Use(h.AuthMiddleware).GET("/:id/download", h.downloadEpubFile)
		books.Use(h.AuthMiddleware, h.LibrarianMiddleware).PUT("/:id", h.updateBook)
		books.Use(h.AuthMiddleware, h.LibrarianMiddleware).DELETE("/:id", h.deleteBook)
		books.Use(h.AuthMiddleware, h.LibrarianMiddleware).POST("/", h.createBook)
		books.Use(h.AuthMiddleware, h.LibrarianMiddleware).GET("/export", h.exportBooksToCSV)
		books.Use(h.AuthMiddleware, h.LibrarianMiddleware).POST("/import", h.importBooksFromCSV)
		books.Use(h.AuthMiddleware, h.LibrarianMiddleware).POST("/:id/audio/upload", h.uploadAudiobookFiles)
		books.Use(h.AuthMiddleware, h.LibrarianMiddleware).PUT("/:id/audio/:audio_id/reorder", h.reorderAudioBookFile)
	}
}

func (h *Handler) getAllBooks(c *gin.Context) {
	books, err := h.service.BookServ.GetBooks()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}

func (h *Handler) getBookByID(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := h.service.BookServ.GetBookByID(uint(bookID)) //nolint:gosec
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

type createBookInput struct {
	Title             string    `form:"title"`
	Description       string    `form:"description"`
	AuthorID          uint      `form:"author"`
	GenreID           uint      `form:"genre"`
	PublisherID       uint      `form:"publisher"`
	ISBN              int       `form:"isbn"`
	YearOfPublication time.Time `form:"year_of_publication"`
	Rating            float32   `form:"rating"`
}

func (h *Handler) createBook(c *gin.Context) {
	// Get the file from form data
	file, err := c.FormFile("picture")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}

	// Check file type
	if !isAllowedFileType(file.Header.Get("Content-Type")) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file type. Only images are allowed"})
		return
	}

	// Generate unique filename
	filename := generateUniqueFilename(file.Filename)

	// Create uploads directory if it doesn't exist
	uploadsDir := "uploads/books"

	// Save file to server
	filePath := filepath.Join(uploadsDir, filename+".jpg")
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Get Epub file of book
	epubFile, err := c.FormFile("epub")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no epub file uploaded"})
		return
	}

	// Save epub file to server
	epubFilePath := filepath.Join(uploadsDir, filename+".epub")
	if err := c.SaveUploadedFile(epubFile, epubFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save epub file"})
		return
	}

	var input createBookInput

	// Bind other form data
	if err := c.ShouldBind(&input); err != nil {
		// Clean up uploaded file if binding fails
		os.Remove(filePath)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := service.CreateBookInput{
		Title:             input.Title,
		Description:       input.Description,
		AuthorID:          input.AuthorID,
		GenreID:           input.GenreID,
		PublisherID:       input.PublisherID,
		ISBN:              input.ISBN,
		YearOfPublication: input.YearOfPublication,
		Picture:           filePath, // Save the file path to database
		Rating:            input.Rating,
		EpubFile:          epubFilePath,
	}

	if err := h.service.BookServ.CreateBook(&book); err != nil {
		// Clean up uploaded file if database operation fails
		os.Remove(filePath)
		os.Remove(epubFilePath)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cookie, err := c.Request.Cookie("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.LogServ.CreateLogWithCookie(cookie, "Создание книги")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

// Helper functions
func isAllowedFileType(contentType string) bool {
	allowedTypes := map[string]bool{
		"image/jpeg":           true,
		"image/png":            true,
		"image/gif":            true,
		"image/webp":           true,
		"application/epub+zip": true,
		"application/zip":      true,
	}
	return allowedTypes[contentType]
}

func generateUniqueFilename(originalFilename string) string {
	return fmt.Sprintf("%s", uuid.New().String())
}

func (h *Handler) getGroupedBooksByTitle(c *gin.Context) {
	books, err := h.service.BookServ.GetGroupedBooksByTitle()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}

func (h *Handler) getBookUniqueCodes(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	codes, err := h.service.BookServ.GetBookUniqueCodes(uint(bookID)) // nolint:gosec
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, codes)
}

type updateBookInput struct {
	ID                uint      `form:"id"`
	Title             string    `form:"title"`
	Description       string    `form:"description"`
	AuthorID          uint      `form:"author"`
	GenreID           uint      `form:"genre"`
	PublisherID       uint      `form:"publisher"`
	ISBN              int       `form:"isbn"`
	YearOfPublication time.Time `form:"year_of_publication"`
	Rating            float32   `form:"rating"`
}

func (h *Handler) updateBook(c *gin.Context) {
	var input updateBookInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingBook, err := h.service.BookServ.GetBookByID(input.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	picturePath := existingBook.Picture
	epubPath := existingBook.EpubFile

	// Обработка новой картинки
	if file, err := c.FormFile("picture"); err == nil {
		if !isAllowedFileType(file.Header.Get("Content-Type")) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file type"})
			return
		}

		filename := generateUniqueFilename(file.Filename)
		uploadsDir := "uploads/books"
		if err := os.MkdirAll(uploadsDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create directory"})
			return
		}

		newFilePath := filepath.Join(uploadsDir, filename+".jpg")
		if err := c.SaveUploadedFile(file, newFilePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
			return
		}

		if existingBook.Picture != "" {
			os.Remove(existingBook.Picture)
		}
		picturePath = newFilePath
	}

	// Обработка нового EPUB файла
	if epubFile, err := c.FormFile("epub"); err == nil {
		filename := generateUniqueFilename(epubFile.Filename)
		uploadsDir := "uploads/books"
		newEpubPath := filepath.Join(uploadsDir, filename+".epub")

		if err := c.SaveUploadedFile(epubFile, newEpubPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save epub"})
			return
		}

		if existingBook.EpubFile != "" {
			os.Remove(existingBook.EpubFile)
		}
		epubPath = newEpubPath
	}

	book := service.UpdateBookInput{
		ID:                input.ID,
		Title:             input.Title,
		AuthorID:          input.AuthorID,
		GenreID:           input.GenreID,
		Description:       input.Description,
		PublisherID:       input.PublisherID,
		ISBN:              input.ISBN,
		YearOfPublication: input.YearOfPublication,
		Picture:           picturePath,
		Rating:            input.Rating,
		EpubFile:          epubPath,
	}

	if err := h.service.BookServ.UpdateBook(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cookie, err := c.Request.Cookie("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.LogServ.CreateLogWithCookie(cookie, "Изменение книги"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *Handler) deleteBook(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.BookServ.DeleteBook(uint(bookID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cookie, err := c.Request.Cookie("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.LogServ.CreateLogWithCookie(cookie, "Удаление книги")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book deleted"})
}

func (h *Handler) reserveBook(c *gin.Context) {
	userIDCookie, err := c.Cookie("id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID, err := strconv.Atoi(userIDCookie)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code, err := h.service.BookServ.ReserveBook(uint(bookID), uint(userID)) //nolint:gosec
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cookie, err := c.Request.Cookie("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.LogServ.CreateLogWithCookie(cookie, "Бронирование книги")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": code})
}

type PaginatedResponse struct {
	Books       []*domain.Book `json:"books"`
	TotalPages  int            `json:"totalPages"`
	CurrentPage int            `json:"currentPage"`
}

func (h *Handler) getBooksWithPagination(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit := 10
	offset := (page - 1) * limit

	booksTotal, err := h.service.BookServ.GetBooks()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	totalCount := len(booksTotal)
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	query := c.Query("query")
	if query != "" {
		books, err := h.service.BookServ.FindBooks(limit, offset, query)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response := PaginatedResponse{
			Books:       books,
			TotalPages:  totalPages,
			CurrentPage: page,
		}

		c.JSON(http.StatusOK, response)
		return
	}

	title := c.Query("title")
	if title != "" {
		books, err := h.service.BookServ.FindBookByTitle(limit, offset, title)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response := PaginatedResponse{
			Books:       books,
			TotalPages:  totalPages,
			CurrentPage: page,
		}

		c.JSON(http.StatusOK, response)
		return
	}

	books, err := h.service.BookServ.GetBooksWithPagination(limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := PaginatedResponse{
		Books:       books,
		TotalPages:  totalPages,
		CurrentPage: page,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) downloadEpubFile(c *gin.Context) {
	// Получаем ID книги из параметров запроса
	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	// Получаем данные книги по ID
	book, err := h.service.BookServ.GetBookByID(uint(bookID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	// Проверяем, есть ли файл EPUB
	if book.EpubFile == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "epub file not found"})
		return
	}

	// Проверяем, существует ли файл
	if _, err := os.Stat(book.EpubFile); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found on server"})
		return
	}

	// Открываем файл
	file, err := os.Open(book.EpubFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer file.Close()

	// Получаем информацию о файле
	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get file info"})
		return
	}

	// Устанавливаем заголовки вручную
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(book.EpubFile))
	c.Header("Content-Type", "application/epub+zip")
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// Передаём файл клиенту вручную через io.Copy
	if _, err := io.Copy(c.Writer, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send file"})
		return
	}
}

func (h *Handler) uploadAudiobookFiles(c *gin.Context) {
	// Получаем ID книги из параметров запроса
	bookIDStr := c.Param("id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	// Проверяем, существует ли книга
	book, err := h.service.BookServ.GetBookByID(uint(bookID))
	if err != nil || book == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	// Получаем файл из запроса
	audioFile, err := c.FormFile("audioFiles")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get audio file: " + err.Error()})
		return
	}

	// Получаем название главы
	chapterTitle := c.PostForm("chapterTitle")
	if chapterTitle == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chapter title is required"})
		return
	}

	// Проверяем тип файла (опционально, можно расширить для поддержки разных аудиоформатов)
	contentType := audioFile.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "audio/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file must be an audio file"})
		return
	}

	// Открываем файл
	src, err := audioFile.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file: " + err.Error()})
		return
	}
	defer src.Close()

	// Читаем содержимое файла в []byte
	audioData, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file: " + err.Error()})
		return
	}

	// Создаём объект аудиофайла
	audiobookFile := &domain.AudiobookFile{
		BookID:       uint(bookID),
		FilePath:     audioFile.Filename,
		ChapterTitle: chapterTitle,
	}

	// Передаём в сервис для сохранения
	err = h.service.BookServ.CreateAudiobookFile(audiobookFile, audioData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save audio file: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "audiobook file uploaded successfully"})
}

func (h *Handler) getBookAudiobookFiles(c *gin.Context) {
	// Получаем ID книги из параметров запроса
	bookIDStr := c.Param("id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	// Проверяем, существует ли книга
	book, err := h.service.BookServ.GetBookByID(uint(bookID))
	if err != nil || book == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	// Получаем список аудиофайлов книги
	audiobookFiles, err := h.service.BookServ.GetAudiobookFilesByBookID(uint(bookID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get audiobook files"})
		return
	}

	c.JSON(http.StatusOK, audiobookFiles)
}

func (h *Handler) serveAudioFile(c *gin.Context) {
	fileIDStr := c.Param("id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}

	audioFile, err := h.service.BookServ.GetAudiobookFileByID(uint(fileID))
	if err != nil || audioFile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "audio file not found"})
		return
	}

	fileExt := filepath.Ext(audioFile.FilePath)
	contentType := "audio/mpeg"
	switch strings.ToLower(fileExt) {
	case ".wav":
		contentType = "audio/wav"
	case ".ogg":
		contentType = "audio/ogg"
	case ".aac":
		contentType = "audio/aac"
	}

	fullPath := audioFile.FilePath

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		log.Println("File not found on disk:", fullPath)
		c.JSON(http.StatusNotFound, gin.H{"error": "audio file not found on disk"})
		return
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filepath.Base(audioFile.FilePath)))

	c.File(fullPath)
}

type orderChangeInput struct {
	NewOrder int `json:"new_order"`
}

func (h *Handler) reorderAudioBookFile(c *gin.Context) {
	fileIDStr := c.Param("audio_id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}

	var input orderChangeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	newOrder := input.NewOrder

	audioFile, err := h.service.BookServ.GetAudiobookFileByID(uint(fileID))
	if err != nil || audioFile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "audio file not found"})
		return
	}

	err = h.service.BookServ.UpdateAudiobookFileOrder(uint(fileID), newOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "audiobook file order updated successfully"})
}

func (h *Handler) exportBooksToCSV(c *gin.Context) {
	bookData, err := h.service.BookServ.ExportBooksToCSV()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := "books.csv"

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")
	c.Header("Content-Length", fmt.Sprint(len(bookData)))

	c.Data(http.StatusOK, "text/csv", bookData)

	cookie, err := c.Request.Cookie("id")
	if err != nil {
		log.Printf("Error getting cookie for logging: %v", err)
		return
	}

	err = h.service.LogServ.CreateLogWithCookie(cookie, "Экспорт книг в CSV")
	if err != nil {
		log.Printf("Error creating log: %v", err)
	}
}

func (h *Handler) importBooksFromCSV(c *gin.Context) {
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

	// Импортируем книги
	count, err := h.service.BookServ.ImportBooksFromCSV(fileData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при импорте книг: " + err.Error()})
		return
	}

	// Логируем действие
	cookie, err := c.Request.Cookie("id")
	if err != nil {
		log.Printf("Error getting cookie for logging: %v", err)
	} else {
		err = h.service.LogServ.CreateLogWithCookie(cookie, "Импорт книг из CSV")
		if err != nil {
			log.Printf("Error creating log: %v", err)
		}
	}

	c.JSON(
		http.StatusOK, gin.H{
			"message": fmt.Sprintf("Успешно импортировано %d книг", count),
			"count":   count,
		},
	)
}
