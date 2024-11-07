package delivery

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) initUserRoutes(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.POST("/sign-up", h.signUp)
		users.POST("/sign-in-by-login", h.signInByLogin)
		users.POST("/sign-in-by-email", h.signInByEmail)
		users.GET("/:id", h.getUserByID)
		users.GET("/", h.getAllUsers)
		users.PUT("/:id", h.updateUser)
		users.GET("/cookie", h.getUserByCookie)
		users.PATCH("/cookie", h.updateUserByCookie)
		users.POST("/logout", h.logout)
		users.DELETE("/:id", h.deleteUser)
		users.DELETE("/cookie", h.deleteUserByCookie)
		users.Use(h.AuthMiddleware, h.AdminMiddleware).PATCH("/change-role", h.updateUserRole)
	}
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

	user, err := h.service.UserServ.SignInByLogin(input.Login, input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(
		c.Writer, &http.Cookie{
			Name:     "id",
			Value:    strconv.Itoa(int(user.ID)), //nolint:gosec
			Expires:  time.Now().Add(24 * time.Hour),
			Path:     "/",
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		},
	)

	c.JSON(http.StatusOK, user)
}

type signInByEmailInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signInByEmail(c *gin.Context) {
	var input signInByEmailInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.UserServ.SignInByEmail(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
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

func (h *Handler) getUserByCookie(c *gin.Context) {
	cookie, err := c.Request.Cookie("id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if cookie.Value == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.service.UserServ.GetUserByID(uint(userID)) //nolint:gosec
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) logout(c *gin.Context) {
	http.SetCookie(
		c.Writer, &http.Cookie{
			Name:     "id",
			Value:    "0", //nolint:gosec
			Path:     "/", // Убедитесь, что Path совпадает
			MaxAge:   -1,  // MaxAge=-1 означает немедленное удаление
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		},
	)

	c.JSON(http.StatusOK, gin.H{"message": "user logged out"})
}

func (h *Handler) updateUserByCookie(c *gin.Context) {
	cookie, err := c.Request.Cookie("id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if cookie.Value == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input updateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.UserServ.GetUserByID(uint(userID)) //nolint:gosec
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func (h *Handler) deleteUserByCookie(c *gin.Context) {
	cookie, err := c.Request.Cookie("id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if cookie.Value == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.UserServ.DeleteUser(uint(userID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
