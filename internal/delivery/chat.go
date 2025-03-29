package delivery

import (
	"fmt"
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) initChatRoutes(r *gin.Engine) {
	chats := r.Group("/chats")
	{
		chats.Use(h.AuthMiddleware)
		chats.GET("/", h.getUserChats)
		chats.POST("/", h.createChat)
		chats.GET("/:id", h.getChatByID)
		chats.GET("/:id/messages", h.getChatMessages)
	}

	librarian := r.Group("/librarian/chats")
	{
		librarian.Use(h.AuthMiddleware, h.LibrarianMiddleware)
		librarian.GET("/", h.getLibrarianChats)            // Новый маршрут
		librarian.GET("/unassigned", h.getUnassignedChats) // Новый маршрут
		librarian.POST("/:id/assign", h.assignChatToLibrarian)
		librarian.POST("/:id/close", h.closeChat)
	}
}

func (h *Handler) getUserChats(c *gin.Context) {
	userID, err := h.getUserIDFromCookie(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	chats, err := h.service.ChatServ.GetChatsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chats)
}

func (h *Handler) createChat(c *gin.Context) {
	userID, err := h.getUserIDFromCookie(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		Title string `json:"title" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	chat := &domain.Chat{
		UserID: userID,
		Title:  input.Title,
		Status: "waiting",
	}

	if err := h.service.ChatServ.CreateChat(chat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, chat)
}

func (h *Handler) getChatMessages(c *gin.Context) {
	userID, err := h.getUserIDFromCookie(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.UserServ.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных пользователя"})
		return
	}

	chatID, err := ParseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chat, err := h.service.ChatServ.GetChatByID(chatID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Чат не найден"})
		return
	}

	isLibrarian := user.Role == "librarian"
	if chat.UserID != userID && !isLibrarian {
		c.JSON(http.StatusForbidden, gin.H{"error": "Нет доступа к чату"})
		return
	}

	messages, err := h.service.ChatServ.GetMessagesByChatID(chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = h.service.ChatServ.MarkMessagesAsRead(chatID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *Handler) closeChat(c *gin.Context) {
	chatID, err := ParseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.service.ChatServ.CloseChat(chatID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Чат закрыт"})
}

func (h *Handler) getChatByID(c *gin.Context) {
	chatID, err := ParseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chat, err := h.service.ChatServ.GetChatByID(chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chat)
}

func (h *Handler) getActiveChats(c *gin.Context) {
	chats, err := h.service.ChatServ.GetActiveChatsForLibrarian()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chats)
}

func (h *Handler) assignChatToLibrarian(c *gin.Context) {
	chatID, err := ParseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	librarianID, err := h.getUserIDFromCookie(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ChatServ.AssignLibrarianToChat(chatID, librarianID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Чат успешно назначен библиотекарю"})
}

func ParseUintParam(c *gin.Context, param string) (uint, error) {
	val, err := strconv.ParseUint(c.Param(param), 10, 32)
	if err != nil {
		return 0, fmt.Errorf("неверный параметр %s: %w", param, err)
	}
	return uint(val), nil
}

func (h *Handler) getUserIDFromCookie(c *gin.Context) (uint, error) {
	idCookie, err := c.Cookie("id")
	if err != nil {
		return 0, fmt.Errorf("cookie id не найден: %w", err)
	}

	userID, err := strconv.ParseUint(idCookie, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("некорректный ID в cookie: %w", err)
	}

	return uint(userID), nil
}

// Получение чатов, назначенных библиотекарю
func (h *Handler) getLibrarianChats(c *gin.Context) {
	librarianID, err := h.getUserIDFromCookie(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	chats, err := h.service.ChatServ.GetLibrarianChats(librarianID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chats)
}

// Получение непринятых чатов
func (h *Handler) getUnassignedChats(c *gin.Context) {
	chats, err := h.service.ChatServ.GetUnassignedChats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chats)
}
