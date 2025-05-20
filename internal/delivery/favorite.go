package delivery

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AZRV17/zlib-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initFavoriteRoutes(r *gin.Engine) {
	favorites := r.Group("/favorites")
	{

		favorites.Use(h.AuthMiddleware)
		{
			favorites.GET("/", h.getFavoritesByUserID)
			favorites.DELETE("/:id", h.deleteFavoriteByBookID)
			favorites.POST("/", h.createFavorite)
		}
	}
}

func (h *Handler) getFavoritesByUserID(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "не удалось получить ID пользователя: " + err.Error()})
		return
	}

	favorites, err := h.service.FavoriteServ.GetFavoriteByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения избранного: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, favorites)
}

func (h *Handler) deleteFavoriteByBookID(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	favorite, err := h.service.FavoriteServ.DeleteFavoriteByUserIDAndBookID(userID, uint(bookID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createLogInput := &service.CreateLogInput{
		UserID:  userID,
		Action:  "Удаление из избранного",
		Date:    time.Now(),
		Details: fmt.Sprintf("Удаление книги %d из избранного", bookID),
	}

	err = h.service.LogServ.CreateLog(createLogInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, favorite)
}

type createFavoriteInput struct {
	BookID uint `json:"book_id" binding:"required"`
}

func (h *Handler) createFavorite(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "не удалось получить ID пользователя: " + err.Error()})
		return
	}

	var input createFavoriteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных: " + err.Error()})
		return
	}

	err = h.service.FavoriteServ.CreateFavorite(
		&service.CreateFavoriteInput{
			BookID: input.BookID,
			UserID: userID,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка добавления в избранное: " + err.Error()})
		return
	}

	createLogInput := &service.CreateLogInput{
		UserID:  userID,
		Action:  "Добавление книги в избранное",
		Date:    time.Now(),
		Details: fmt.Sprintf("Добавление книги %d в избранное", input.BookID),
	}

	err = h.service.LogServ.CreateLog(createLogInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Важно вернуть правильный статус-код (не 307)
	c.JSON(http.StatusOK, gin.H{"message": "book added to favorites"})
}
