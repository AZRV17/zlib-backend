package repository

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/gorm"
	"time"
)

type chatRepo struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepo {
	return &chatRepo{db: db}
}

// Сохранение сообщения
func (r *chatRepo) SaveMessage(message *domain.Message) error {
	tx := r.db.Begin()

	if err := tx.Create(message).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Обновляем время последней активности в чате
	if err := tx.Model(&domain.Chat{}).Where("id = ?", message.ChatID).
		Update("last_activity", message.CreatedAt).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// Создание нового чата
func (r *chatRepo) CreateChat(chat *domain.Chat) error {
	return r.db.Create(chat).Error
}

// Получение чата по ID
func (r *chatRepo) GetChatByID(chatID uint) (*domain.Chat, error) {
	var chat domain.Chat
	if err := r.db.First(&chat, chatID).Error; err != nil {
		return nil, err
	}
	return &chat, nil
}

// Получение сообщений чата
func (r *chatRepo) GetMessagesByChatID(chatID uint) ([]domain.Message, error) {
	var messages []domain.Message
	if err := r.db.Where("chat_id = ?", chatID).Order("created_at").
		Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

// Получение активных чатов для библиотекаря
func (r *chatRepo) GetActiveChatsForLibrarian() ([]domain.Chat, error) {
	var chats []domain.Chat
	if err := r.db.Where("status IN (?, ?)", "active", "waiting").
		Order("last_activity desc").Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}

// Получение чатов пользователя
func (r *chatRepo) GetChatsByUserID(userID uint) ([]domain.Chat, error) {
	var chats []domain.Chat
	if err := r.db.Where("user_id = ?", userID).
		Order("last_activity desc").Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}

// Назначение библиотекаря на чат
func (r *chatRepo) AssignLibrarianToChat(chatID, librarianID uint) error {
	return r.db.Model(&domain.Chat{}).Where("id = ?", chatID).
		Updates(
			map[string]interface{}{
				"librarian_id": librarianID,
				"status":       "active",
			},
		).Error
}

// Закрытие чата
func (r *chatRepo) CloseChat(chatID uint) error {
	return r.db.Model(&domain.Chat{}).Where("id = ?", chatID).
		Update("status", "closed").Error
}

// Отметить сообщения как прочитанные
func (r *chatRepo) MarkMessagesAsRead(chatID, userID uint) error {
	now := time.Now()
	return r.db.Model(&domain.Message{}).
		Where("chat_id = ? AND sender_id != ? AND read_at IS NULL", chatID, userID).
		Update("read_at", now).Error
}

func (r *chatRepo) GetLibrarianChats(librarianID uint) ([]domain.Chat, error) {
	var chats []domain.Chat
	if err := r.db.Where("librarian_id = ? AND status != ?", librarianID, "closed").
		Order("last_activity desc").Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}

func (r *chatRepo) GetUnassignedChats() ([]domain.Chat, error) {
	var chats []domain.Chat
	if err := r.db.Where("librarian_id IS NULL AND status = ?", "waiting").
		Order("created_at desc").Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}
