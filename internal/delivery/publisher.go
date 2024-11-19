package delivery

import (
	"github.com/AZRV17/zlib-backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) initPublisherRoutes(r *gin.Engine) {
	publishers := r.Group("/publishers")
	{
		publishers.GET("/", h.getPublishers)
		publishers.GET("/:id", h.getPublisherByID)
		publishers.Use(h.AuthMiddleware, h.LibrarianMiddleware).DELETE("/:id", h.deletePublisher)
		publishers.Use(h.AuthMiddleware, h.LibrarianMiddleware).POST("/", h.createPublisher)
		publishers.Use(h.AuthMiddleware, h.LibrarianMiddleware).PUT("/:id", h.updatePublisher)
	}
}

func (h *Handler) getPublishers(c *gin.Context) {
	publishers, err := h.service.PublisherServ.GetPublishers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, publishers)
}

func (h *Handler) getPublisherByID(c *gin.Context) {
	publisherID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	publisher, err := h.service.PublisherServ.GetPublisherByID(uint(publisherID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, publisher)
}

func (h *Handler) deletePublisher(c *gin.Context) {
	publisherID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.PublisherServ.DeletePublisher(uint(publisherID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cookie, err := c.Request.Cookie("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.LogServ.CreateLogWithCookie(cookie, "Удаление издателя")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "publisher deleted"})
}

func (h *Handler) updatePublisher(c *gin.Context) {
	var input service.UpdatePublisherInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.PublisherServ.UpdatePublisher(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cookie, err := c.Request.Cookie("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.LogServ.CreateLogWithCookie(cookie, "Изменение издателя")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "publisher updated"})
}

func (h *Handler) createPublisher(c *gin.Context) {
	var input service.CreatePublisherInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.PublisherServ.CreatePublisher(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cookie, err := c.Request.Cookie("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.LogServ.CreateLogWithCookie(cookie, "Создание издателя")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "publisher created"})
}
