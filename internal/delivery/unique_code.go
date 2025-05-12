package delivery

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) initUniqueCodeRoutes(r *gin.Engine) {
	uniqueCodes := r.Group("/unique-codes")
	{
		uniqueCodes.GET("/", h.getUniqueCodes)
		uniqueCodes.Use(h.AuthMiddleware, h.LibrarianMiddleware).GET("/:id", h.getUniqueCodeByID)
		uniqueCodes.Use(h.AuthMiddleware, h.LibrarianMiddleware).POST("/", h.createUniqueCode)
		uniqueCodes.Use(h.AuthMiddleware, h.LibrarianMiddleware).DELETE("/:id", h.deleteUniqueCode)
		uniqueCodes.Use(h.AuthMiddleware, h.LibrarianMiddleware).PUT("/:id", h.updateUniqueCode)
	}
}

func (h *Handler) getUniqueCodes(c *gin.Context) {
	codes, err := h.service.BookServ.GetUniqueCodes() //nolint:gosec
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, codes)
}

func (h *Handler) getUniqueCodeByID(c *gin.Context) {
	codeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code, err := h.service.BookServ.GetUniqueCodeByID(uint(codeID)) //nolint:gosec
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, code)
}

type createUniqueCodeInput struct {
	Code   int  `json:"code" binding:"required"`
	BookID uint `json:"book" binding:"required"`
}

func (h *Handler) createUniqueCode(c *gin.Context) {
	var input createUniqueCodeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.BookServ.CreateUniqueCode(
		&domain.UniqueCode{
			Code:        input.Code,
			BookID:      input.BookID,
			IsAvailable: true,
		},
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cookie, err := c.Request.Cookie("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.LogServ.CreateLogWithCookie(cookie, "Создание уникального кода")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "unique code created"})
}

func (h *Handler) deleteUniqueCode(c *gin.Context) {
	codeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.BookServ.DeleteUniqueCode(uint(codeID)) //nolint:gosec
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cookie, err := c.Request.Cookie("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.LogServ.CreateLogWithCookie(cookie, "Удаление уникального кода")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "unique code deleted"})
}

type updateUniqueCodeInput struct {
	ID          uint `json:"id" binding:"required"`
	Code        int  `json:"code" binding:"required"`
	BookID      uint `json:"book_id" binding:"required"`
	IsAvailable bool `json:"is_available"`
}

func (h *Handler) updateUniqueCode(c *gin.Context) {
	var input updateUniqueCodeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.BookServ.UpdateUniqueCode(
		&domain.UniqueCode{
			ID:          input.ID,
			Code:        input.Code,
			BookID:      input.BookID,
			IsAvailable: input.IsAvailable,
		},
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cookie, err := c.Request.Cookie("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.LogServ.CreateLogWithCookie(cookie, "Изменение уникального кода")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "unique code updated"})
}
