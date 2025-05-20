package delivery

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AZRV17/zlib-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initReviewRoutes(r *gin.Engine) {
	reviews := r.Group("/reviews")
	{
		reviews.GET("/:id", h.getReviewsByBookID)
		reviews.Use(h.AuthMiddleware).POST("/", h.createReview)
	}
}

type createReviewInput struct {
	BookID  uint    `json:"book_id"`
	Rating  float32 `json:"rating"`
	Message string  `json:"message"`
}

func (h *Handler) createReview(c *gin.Context) {
	var input createReviewInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	serviceCreateReviewInput := &service.CreateReviewInput{
		UserID:  userID,
		BookID:  input.BookID,
		Rating:  input.Rating,
		Message: input.Message,
	}

	err = h.service.ReviewServ.CreateReview(serviceCreateReviewInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createLogInput := &service.CreateLogInput{
		UserID:  userID,
		Action:  "Добавление отзыва",
		Date:    time.Now(),
		Details: fmt.Sprintf("Добавление отзыва к книге ID: %d", input.BookID),
	}

	err = h.service.LogServ.CreateLog(createLogInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) getReviewsByBookID(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reviews, err := h.service.ReviewServ.GetReviewsByBookID(uint(bookID)) //nolint:gosec
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}
