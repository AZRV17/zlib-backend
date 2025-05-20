package delivery

import (
	"net/http"
	"strings"

	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) AuthMiddleware(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	// Получаем токен из куки
	token, err := c.Cookie("access_token")

	// Если в куки нет токена, пробуем получить из заголовка Authorization (для API-клиентов)
	if err != nil {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}
	}

	// Если токен не найден ни в куках, ни в заголовке
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "токен авторизации не найден"})
		return
	}

	// Парсим и проверяем JWT токен
	claims, err := h.service.UserServ.ParseToken(token)
	if err != nil {
		// Если токен недействителен, пробуем обновить через refresh token
		refreshToken, refreshErr := c.Cookie("refresh_token")
		if refreshErr == nil {
			newTokens, refreshErr := h.service.UserServ.RefreshTokens(refreshToken)
			if refreshErr == nil {
				// Успешно обновили токены, устанавливаем новые куки
				setAuthCookies(c, newTokens)

				// Повторная проверка с новым токеном
				claims, err = h.service.UserServ.ParseToken(newTokens.AccessToken)
				if err != nil {
					c.AbortWithStatusJSON(
						http.StatusUnauthorized,
						gin.H{"error": "невалидный токен после обновления: " + err.Error()},
					)
					return
				}
			} else {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					gin.H{"error": "не удалось обновить токен: " + refreshErr.Error()},
				)
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "невалидный токен: " + err.Error()})
			return
		}
	}

	// Получаем пользователя по ID из токена
	user, err := h.service.UserServ.GetUserByID(claims.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "пользователь не найден"})
		return
	}

	// Записываем пользователя в контекст запроса
	c.Set("user", user)
	c.Set("userID", claims.UserID)
	c.Next()
}

func (h *Handler) AdminMiddleware(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if user.(*domain.User).Role != domain.RoleAdmin {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "доступ запрещен"})
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

	if user.(*domain.User).Role != "librarian" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "доступ запрещен"})
		return
	}

	c.Next()
}

// Вспомогательная функция для получения ID пользователя из JWT токена
func getUserIDFromContext(c *gin.Context) (uint, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, Unauthorized("ID пользователя не найден в контексте")
	}

	return userID.(uint), nil
}
