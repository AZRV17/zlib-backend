package delivery

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initGenreRoutes(r *gin.Engine) {
	genres := r.Group("/genres")
	{
		genres.GET("/", h.getAllGenres)
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
