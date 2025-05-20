package websocket

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
	"github.com/AZRV17/zlib-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type ChatHub struct {
	clients   map[uint][]*websocket.Conn
	broadcast chan *domain.Message
	repo      repository.ChatRepo
	userSvc   service.UserServ
	mu        sync.RWMutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:3000"
	},
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 10 * time.Second,
}

func NewChatHub(repo repository.ChatRepo, userSvc service.UserServ) *ChatHub {
	return &ChatHub{
		clients:   make(map[uint][]*websocket.Conn),
		broadcast: make(chan *domain.Message, 100), // Буферизованный канал
		repo:      repo,
		userSvc:   userSvc,
	}
}

func (h *ChatHub) HandleConnections(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		log.Printf("Ошибка получения cookie id: %v", err)
		return
	}

	// Получаем информацию о пользователе сразу
	user, err := h.userSvc.GetUserByID(uint(userID))
	if err != nil {
		log.Printf("Ошибка получения пользователя: %v", err)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Ошибка WebSocket upgrade: %v", err)
		return
	}
	defer conn.Close()

	// Настройки соединения
	conn.SetReadLimit(512 * 1024)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(
		func(string) error {
			return conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		},
	)

	// Регистрация соединения
	h.mu.Lock()
	h.clients[uint(userID)] = append(h.clients[uint(userID)], conn)
	h.mu.Unlock()

	// Отложенная очистка соединения
	defer func() {
		h.mu.Lock()
		connections := h.clients[uint(userID)]
		for i, c := range connections {
			if c == conn {
				h.clients[uint(userID)] = append(connections[:i], connections[i+1:]...)
				break
			}
		}
		if len(h.clients[uint(userID)]) == 0 {
			delete(h.clients, uint(userID))
		}
		h.mu.Unlock()
	}()

	// Основной цикл
	var chatID uint
	for {
		var msg struct {
			ChatID  uint   `json:"chat_id"`
			Content string `json:"content"`
		}

		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("Ошибка чтения: %v", err)
			}
			break
		}

		// Пропускаем пустые сообщения (инициализация)
		if msg.Content == "" {
			chatID = msg.ChatID
			continue
		}

		// Проверяем доступ к чату
		chat, err := h.repo.GetChatByID(msg.ChatID)
		if err != nil || chat.ID != chatID {
			log.Printf("Ошибка доступа к чату: %v", err)
			continue
		}

		message := &domain.Message{
			ChatID:     chatID,
			SenderID:   uint(userID),
			SenderRole: string(user.Role),
			SenderName: user.Login,
			Content:    msg.Content,
		}

		if err := h.repo.SaveMessage(message); err != nil {
			log.Printf("Ошибка сохранения сообщения: %v", err)
			continue
		}

		// Безопасная отправка в канал
		select {
		case h.broadcast <- message:
			// Сообщение успешно отправлено
		default:
			log.Printf("Канал broadcast переполнен или закрыт")
		}
	}
}

func (h *ChatHub) HandleMessages() {
	for message := range h.broadcast {
		h.mu.RLock()
		// Находим всех получателей сообщения
		chat, err := h.repo.GetChatByID(message.ChatID)
		if err != nil {
			log.Printf("Ошибка получения чата: %v", err)
			h.mu.RUnlock()
			continue
		}

		// Отправляем сообщение всем подключенным клиентам чата
		recipients := []uint{chat.UserID}
		if chat.LibrarianID != nil {
			recipients = append(recipients, *chat.LibrarianID)
		}

		for _, userID := range recipients {
			if connections, ok := h.clients[userID]; ok {
				for _, conn := range connections {
					if err := conn.WriteJSON(message); err != nil {
						log.Printf("Ошибка отправки сообщения: %v", err)
						conn.Close()
					}
				}
			}
		}
		h.mu.RUnlock()
	}
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

func (h *ChatHub) getUserIDFromContext(c *gin.Context) (uint, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, Unauthorized("ID пользователя не найден в контексте")
	}

	return userID.(uint), nil
}
