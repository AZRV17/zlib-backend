package delivery

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/service"
	"github.com/AZRV17/zlib-backend/pkg/auth"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initUserRoutes(r *gin.Engine) {
	users := r.Group("/users")
	{
		// Маршруты для авторизации
		users.POST("/sign-in-by-login", h.signInByLogin)
		users.POST("/sign-in-by-email", h.signInByEmail)
		users.POST("/sign-up", h.signUp)
		users.POST("/refresh", h.refreshTokens)
		users.POST("/logout", h.logout)

		// Маршруты для восстановления пароля
		users.POST("/forgot-password", h.requestPasswordReset)
		users.GET("/reset-password", h.validateResetToken)
		users.POST("/reset-password", h.resetPassword)

		// Маршруты с авторизацией
		authorized := users.Group("/")
		authorized.Use(h.AuthMiddleware)
		{
			authorized.GET("/me", h.getUserMe)
			authorized.PUT("/me", h.updateUserMe)
			authorized.PATCH("/me/email", h.updateUserMeEmail)
			authorized.DELETE("/me", h.deleteUserMe)
		}

		// Админские маршруты
		admin := users.Group("/")
		admin.Use(h.AuthMiddleware, h.AdminMiddleware)
		{
			admin.GET("/", h.getAllUsers)
			admin.GET("/:id", h.getUserByID)
			admin.PUT("/:id", h.updateUser)
			admin.DELETE("/:id", h.deleteUser)
			admin.PATCH("/change-role", h.updateUserRole)
		}
	}
}

// Вспомогательные функции для авторизации
func getTokenFromHeader(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")
	if header == "" {
		return "", Unauthorized("пустой заголовок Authorization")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", Unauthorized("неверный формат токена авторизации")
	}

	return headerParts[1], nil
}

func Unauthorized(message string) error {
	return &AuthError{
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

type AuthError struct {
	Message string
	Status  int
}

func (e *AuthError) Error() string {
	return e.Message
}

// Методы авторизации
type signInByLoginInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signInByLogin(c *gin.Context) {
	var input signInByLoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, tokens, err := h.service.UserServ.SignInByLogin(input.Login, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем JWT в куки
	setAuthCookies(c, tokens)

	c.JSON(
		http.StatusOK, gin.H{
			"user": user,
		},
	)
}

type signInByEmailInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signInByEmail(c *gin.Context) {
	var input signInByEmailInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, tokens, err := h.service.UserServ.SignInByEmail(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем JWT в куки
	setAuthCookies(c, tokens)

	c.JSON(
		http.StatusOK, gin.H{
			"user": user,
		},
	)
}

func (h *Handler) refreshTokens(c *gin.Context) {
	// Получаем refresh token из куков
	refreshTokenCookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
		return
	}

	tokens, err := h.service.UserServ.RefreshTokens(refreshTokenCookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем новые JWT в куки
	setAuthCookies(c, tokens)

	c.JSON(http.StatusOK, gin.H{"message": "tokens refreshed successfully"})
}

func (h *Handler) logout(c *gin.Context) {
	// Удаляем куки, устанавливая время жизни в прошлое
	clearAuthCookies(c)

	c.JSON(http.StatusOK, gin.H{"message": "Вы успешно вышли из системы"})
}

// Вспомогательная функция для установки JWT токенов в куки
func setAuthCookies(c *gin.Context, tokens *auth.Tokens) {
	http.SetCookie(
		c.Writer, &http.Cookie{
			Name:     "access_token",
			Value:    tokens.AccessToken,
			Expires:  time.Now().Add(24 * time.Hour * 30),
			Path:     "/",
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		},
	)

	http.SetCookie(
		c.Writer, &http.Cookie{
			Name:     "refresh_token",
			Value:    tokens.RefreshToken,
			Expires:  time.Now().Add(24 * time.Hour * 30),
			Path:     "/",
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		},
	)
}

// Вспомогательная функция для удаления JWT токенов из куков
func clearAuthCookies(c *gin.Context) {
	http.SetCookie(
		c.Writer, &http.Cookie{
			Name:     "access_token",
			Value:    "0",
			Path:     "/",
			MaxAge:   -1,
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		},
	)

	http.SetCookie(
		c.Writer, &http.Cookie{
			Name:     "refresh_token",
			Value:    "0",
			Path:     "/",
			MaxAge:   -1,
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		},
	)
}

type signUpInput struct {
	Login          string `json:"login" binding:"required"`
	FullName       string `json:"full_name" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Email          string `json:"email" binding:"required"`
	Role           string `json:"role" binding:"required"`
	PhoneNumber    string `json:"phone_number" binding:"required"`
	PassportNumber int    `json:"passport_number" binding:"required"`
}

func (h *Handler) signUp(c *gin.Context) {
	var input signUpInput

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := service.SignUpUserInput{
		Login:          input.Login,
		Password:       input.Password,
		Email:          input.Email,
		Role:           domain.Role(input.Role),
		FullName:       input.FullName,
		PhoneNumber:    input.PhoneNumber,
		PassportNumber: input.PassportNumber,
	}

	if err := h.service.UserServ.SignUp(&user); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser, err := h.service.UserServ.GetUserByLogin(input.Login)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newUser)
}

type updateUserInput struct {
	ID             uint        `json:"id" gorm:"primaryKey,autoIncrement"`
	Login          string      `json:"login" gore:"unique"`
	FullName       string      `json:"full_name"`
	Password       string      `json:"password"`
	Role           domain.Role `json:"role" gorm:"type:role;default:'user'"`
	Email          string      `json:"email" gorm:"unique"`
	PhoneNumber    string      `json:"phone_number" gorm:"unique"`
	PassportNumber int         `json:"passport_number" gorm:"unique"`
}

func (h *Handler) updateUser(c *gin.Context) {
	var input updateUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.UserServ.GetUserByID(input.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	servInput := service.UpdateUserInput{
		ID:             user.ID,
		Login:          input.Login,
		Password:       input.Password,
		Role:           input.Role,
		Email:          input.Email,
		PhoneNumber:    input.PhoneNumber,
		PassportNumber: input.PassportNumber,
	}

	if err := h.service.UserServ.UpdateUser(&servInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

func (h *Handler) getUserByID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.UserServ.GetUserByID(uint(userID)) //nolint:gosec
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) getAllUsers(c *gin.Context) {
	users, err := h.service.UserServ.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// getUserMe возвращает информацию о текущем пользователе по JWT токену
func (h *Handler) getUserMe(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.UserServ.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// updateUserMe обновляет данные текущего пользователя
func (h *Handler) updateUserMe(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var input updateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.UserServ.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Не позволяем менять роль через этот метод
	servInput := service.UpdateUserInput{
		ID:             user.ID,
		Login:          input.Login,
		FullName:       input.FullName,
		Password:       user.Password,
		Role:           user.Role,
		Email:          user.Email,
		PhoneNumber:    input.PhoneNumber,
		PassportNumber: input.PassportNumber,
	}

	if err := h.service.UserServ.UpdateUser(&servInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

type updateUserEmailInput struct {
	UserId int    `json:"user_id"`
	Email  string `json:"email"`
}

func (h *Handler) updateUserMeEmail(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var input updateUserEmailInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.UserServ.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Не позволяем менять роль через этот метод
	servInput := service.UpdateUserInput{
		ID:             user.ID,
		Login:          user.Login,
		FullName:       user.FullName,
		Password:       user.Password,
		Role:           user.Role,
		Email:          input.Email,
		PhoneNumber:    user.PhoneNumber,
		PassportNumber: user.PassportNumber,
	}

	if err := h.service.UserServ.UpdateUser(&servInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

type updateUserRoleInput struct {
	ID   uint   `json:"id"`
	Role string `json:"role"`
}

func (h *Handler) updateUserRole(c *gin.Context) {
	var input updateUserRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UserServ.UpdateUserRole(input.ID, domain.Role(input.Role))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createLogInput := &service.CreateLogInput{
		UserID:  userID,
		Action:  "Изменение роли пользователя",
		Date:    time.Now(),
		Details: "Изменение роли пользователя",
	}

	err = h.service.LogServ.CreateLog(createLogInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user role updated"})
}

func (h *Handler) deleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UserServ.DeleteUser(uint(userID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createLogInput := &service.CreateLogInput{
		UserID:  adminID,
		Action:  "Удаление пользователя",
		Date:    time.Now(),
		Details: "Удаление пользователя с ID: " + c.Param("id"),
	}

	err = h.service.LogServ.CreateLog(createLogInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

// deleteUserMe удаляет текущего пользователя
func (h *Handler) deleteUserMe(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UserServ.DeleteUser(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

// requestPasswordReset обрабатывает запрос на восстановление пароля
type forgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *Handler) requestPasswordReset(c *gin.Context) {
	var input forgotPasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UserServ.RequestPasswordReset(input.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "На указанный email отправлена инструкция по восстановлению пароля"})
}

// validateResetToken проверяет валидность токена для сброса пароля
func (h *Handler) validateResetToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Отсутствует токен сброса пароля"})
		return
	}

	user, err := h.service.UserServ.ValidateResetToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK, gin.H{
			"message": "Токен действителен",
			"email":   user.Email,
		},
	)
}

// resetPassword сбрасывает пароль пользователя
type resetPasswordInput struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (h *Handler) resetPassword(c *gin.Context) {
	var input resetPasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UserServ.ResetPassword(input.Token, input.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пароль успешно изменен"})
}
