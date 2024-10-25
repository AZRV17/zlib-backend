package delivery

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) AuthMiddleware(c *gin.Context) {
	cookie, err := c.Request.Cookie("id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if cookie.Value == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if userID == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.service.UserServ.GetUserByID(uint(userID)) //nolint:gosec
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.Set("user", user)
	c.Next()
}

func (h *Handler) AdminMiddleware(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if user.(*domain.User).Role != domain.RoleAdmin {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.Next()
}

func (h *Handler) LibrarianMiddleware(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if user.(*domain.User).Role != domain.RoleLibrarian {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.Next()
}
