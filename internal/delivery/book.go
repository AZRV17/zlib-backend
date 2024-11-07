package delivery

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) initBookRoutes(r *gin.Engine) {
	books := r.Group("/books")
	{
		books.GET("/", h.getAllBooks)
		books.GET("/:id", h.getBookByID)
		books.GET("/grouped", h.getGroupedBooksByTitle)
		books.Use(h.AuthMiddleware).GET("/:id/codes", h.getBookUniqueCodes)
		books.Use(h.AuthMiddleware).POST("/:id", h.reserveBook)
		books.Use(h.AuthMiddleware, h.LibrarianMiddleware).PUT("/:id", h.updateBook)
		books.Use(h.AuthMiddleware, h.LibrarianMiddleware).DELETE("/:id", h.deleteBook)
		books.Use(h.AuthMiddleware, h.AdminMiddleware).POST("/", h.createBook)
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
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	AuthorID          uint      `json:"author_id"`
	GenreID           uint      `json:"genre_id"`
	PublisherID       uint      `json:"publisher_id"`
	ISBN              int       `json:"isbn"`
	YearOfPublication time.Time `json:"year_of_publication"`
	Picture           string    `json:"picture"`
	Rating            float32   `json:"rating"`
}

func (h *Handler) createBook(c *gin.Context) {
	userID, err := c.Cookie("id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	st, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.UserServ.GetUserByID(uint(st)) //nolint:gosec
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Role != domain.RoleAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input createBookInput

	if err := c.ShouldBindJSON(&input); err != nil {
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
		Picture:           input.Picture,
		Rating:            input.Rating,
	}

	if err := h.service.BookServ.CreateBook(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
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
	ID                uint      `json:"id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	AuthorID          uint      `json:"author"`
	GenreID           uint      `json:"genre"`
	PublisherID       uint      `json:"publisher"`
	ISBN              int       `json:"isbn"`
	YearOfPublication time.Time `json:"year_of_publication"`
	Picture           string    `json:"picture"`
	Rating            float32   `json:"rating"`
}

func (h *Handler) updateBook(c *gin.Context) {
	var input updateBookInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := service.UpdateBookInput{
		ID:                input.ID,
		Title:             input.Title,
		AuthorID:          input.AuthorID,
		GenreID:           input.GenreID,
		PublisherID:       input.PublisherID,
		ISBN:              input.ISBN,
		YearOfPublication: input.YearOfPublication,
		Picture:           input.Picture,
		Rating:            input.Rating,
	}

	if err := h.service.BookServ.UpdateBook(&book); err != nil {
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

	c.JSON(http.StatusOK, gin.H{"code": code})
}
