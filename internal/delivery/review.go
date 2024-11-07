package delivery

import (
	"github.com/AZRV17/zlib-backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

	userID, err := c.Cookie("id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if userIDInt == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	serviceCreateReviewInput := &service.CreateReviewInput{
		UserID:  uint(userIDInt), //nolint:gosec
		BookID:  input.BookID,
		Rating:  input.Rating,
		Message: input.Message,
	}

	err = h.service.ReviewServ.CreateReview(serviceCreateReviewInput)
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
