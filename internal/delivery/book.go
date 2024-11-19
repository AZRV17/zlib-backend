package delivery

import (
	"fmt"
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func (h *Handler) initBookRoutes(r *gin.Engine) {
	books := r.Group("/books")
	{
		books.GET("/", h.getAllBooks)
		books.GET("/:id", h.getBookByID)
		books.GET("/grouped", h.getGroupedBooksByTitle)
		//books.GET("/pagination", h.findBookByTitle)
		books.GET("/pagination", h.getBooksWithPagination)
		books.Use(h.AuthMiddleware).GET("/:id/codes", h.getBookUniqueCodes)
		books.Use(h.AuthMiddleware).POST("/:id", h.reserveBook)
		books.Use(h.AuthMiddleware, h.LibrarianMiddleware).PUT("/:id", h.updateBook)
		books.Use(h.AuthMiddleware, h.LibrarianMiddleware).DELETE("/:id", h.deleteBook)
		books.Use(h.AuthMiddleware, h.LibrarianMiddleware).POST("/", h.createBook)
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
	AuthorID          uint      `form:"author_id"`
	GenreID           uint      `form:"genre_id"`
	PublisherID       uint      `form:"publisher_id"`
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
	filePath := filepath.Join(uploadsDir, filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
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
	}

	if err := h.service.BookServ.CreateBook(&book); err != nil {
		// Clean up uploaded file if database operation fails
		os.Remove(filePath)
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
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	return allowedTypes[contentType]
}

func generateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	return fmt.Sprintf("%s%s", uuid.New().String(), ext)
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

	// Bind form data except file
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing book to check old picture path
	existingBook, err := h.service.BookServ.GetBookByID(input.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	picturePath := existingBook.Picture // Keep old picture path by default

	// Check if new file is uploaded
	file, err := c.FormFile("picture")
	if err == nil { // New file was uploaded
		// Check file type
		if !isAllowedFileType(file.Header.Get("Content-Type")) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file type. Only images are allowed"})
			return
		}

		// Generate unique filename
		filename := generateUniqueFilename(file.Filename)

		// Create uploads directory if it doesn't exist
		uploadsDir := "uploads/books"
		if err := os.MkdirAll(uploadsDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create uploads directory"})
			return
		}

		// Save new file
		newFilePath := filepath.Join(uploadsDir, filename)
		if err := c.SaveUploadedFile(file, newFilePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
			return
		}

		// Delete old file if it exists
		if existingBook.Picture != "" {
			if err := os.Remove(existingBook.Picture); err != nil {
				// Log error but continue, as this is not critical
				log.Printf("Failed to delete old file: %v", err)
			}
		}

		picturePath = newFilePath
	}

	book := service.UpdateBookInput{
		ID:                input.ID,
		Title:             input.Title,
		AuthorID:          input.AuthorID,
		GenreID:           input.GenreID,
		PublisherID:       input.PublisherID,
		ISBN:              input.ISBN,
		YearOfPublication: input.YearOfPublication,
		Picture:           picturePath,
		Rating:            input.Rating,
	}

	if err := h.service.BookServ.UpdateBook(&book); err != nil {
		// If update fails and we uploaded a new file, clean it up
		if picturePath != existingBook.Picture {
			os.Remove(picturePath)
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cookie, err := c.Request.Cookie("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.LogServ.CreateLogWithCookie(cookie, "Изменение книги")
	if err != nil {
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

//func (h *Handler) findBookByTitle(c *gin.Context) {
//	title := c.Query("title")
//
//	books, err := h.service.BookServ.FindBookByTitle(title)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, books)
//}
