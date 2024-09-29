package delivery

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	}
}

type signUpInput struct {
	Login          string `json:"login" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Email          string `json:"email" binding:"required"`
	Role           string `json:"role" binding:"required"`
	PhoneNumber    string `json:"phoneNumber" binding:"required"`
	PassportNumber int    `json:"passportNumber" binding:"required"`
}

func (h *Handler) signUp(c *gin.Context) {
	var input signUpInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := service.SignUpUserInput{
		Login:          input.Login,
		Password:       input.Password,
		Email:          input.Email,
		Role:           domain.Role(input.Role),
		PhoneNumber:    input.PhoneNumber,
		PassportNumber: input.PassportNumber,
	}

	if err := h.service.UserServ.SignUp(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user signed up"})
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
	Password       string      `json:"password"`
	Role           domain.Role `json:"role" gorm:"type:role;default:'user'"`
	Email          string      `json:"email" gorm:"unique"`
	PhoneNumber    string      `json:"phoneNumber" gorm:"unique"`
	PassportNumber int         `json:"passportNumber" gorm:"unique"`
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

	user, err := h.service.UserServ.GetUserByID(uint(userID))
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
