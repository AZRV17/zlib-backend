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
		// books.GET("/aggregated", h.getAggregatedBooks)
		books.GET("/grouped", h.getGroupedBooksByTitle)
		books.Use(h.AuthMiddleware).GET("/:id/codes", h.getBookUniqueCodes)
		books.Use(h.AdminMiddleware).POST("/reserve", h.reserveBook)
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

type reserveBookInput struct {
	BookID uint `json:"book_id"`
}

func (h *Handler) reserveBook(c *gin.Context) {
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

	var input reserveBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code, err := h.service.BookServ.ReserveBook(input.BookID, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": code})
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

//func (h *Handler) getAggregatedBooks(c *gin.Context) {
//	books, err := h.service.BookServ.GetAggregatedBooks()
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, books)
//}
