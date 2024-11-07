package delivery

import (
	"github.com/AZRV17/zlib-backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) initGenreRoutes(r *gin.Engine) {
	genres := r.Group("/genres")
	{
		genres.GET("/", h.getAllGenres)
		genres.GET("/:id", h.getGenreByID)
		genres.Use(h.AuthMiddleware, h.LibrarianMiddleware).POST("/", h.createGenre)
		genres.Use(h.AuthMiddleware, h.LibrarianMiddleware).PUT("/:id", h.updateGenre)
		genres.Use(h.AuthMiddleware, h.LibrarianMiddleware).DELETE("/:id", h.deleteGenre)
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

	c.JSON(http.StatusOK, gin.H{"message": "genre deleted"})
}
