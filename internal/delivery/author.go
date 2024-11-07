package delivery

import (
	"github.com/AZRV17/zlib-backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) initAuthorRoutes(r *gin.Engine) {
	authors := r.Group("/authors")
	{
		authors.GET("/", h.getAuthors)
		authors.GET("/:id", h.getAuthorByID)
		authors.Use(h.AuthMiddleware, h.LibrarianMiddleware).POST("/", h.createAuthor)
		authors.Use(h.AuthMiddleware, h.LibrarianMiddleware).PUT("/:id", h.updateAuthor)
		authors.Use(h.AuthMiddleware, h.LibrarianMiddleware).DELETE("/:id", h.deleteAuthor)
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

	c.JSON(http.StatusOK, gin.H{"message": "author deleted"})
}
