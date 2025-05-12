package service

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
)

type ChatService struct {
	repo repository.ChatRepo
}

func NewChatService(repo repository.ChatRepo) *ChatService {
	return &ChatService{repo: repo}
}

func (c ChatService) CreateChat(chat *domain.Chat) error {
	return c.repo.CreateChat(chat)
}

func (c ChatService) GetChatByID(chatID uint) (*domain.Chat, error) {
	return c.repo.GetChatByID(chatID)
}

func (c ChatService) GetMessagesByChatID(chatID uint) ([]domain.Message, error) {
	return c.repo.GetMessagesByChatID(chatID)
}

func (c ChatService) GetActiveChatsForLibrarian() ([]domain.Chat, error) {
	return c.repo.GetActiveChatsForLibrarian()
}

func (c ChatService) GetChatsByUserID(userID uint) ([]domain.Chat, error) {
	return c.repo.GetChatsByUserID(userID)
}

func (c ChatService) AssignLibrarianToChat(chatID, librarianID uint) error {
	return c.repo.AssignLibrarianToChat(chatID, librarianID)
}

func (c ChatService) CloseChat(chatID uint) error {
	return c.repo.CloseChat(chatID)
}

func (c ChatService) MarkMessagesAsRead(chatID, userID uint) error {
	return c.repo.MarkMessagesAsRead(chatID, userID)
}

func (c ChatService) GetLibrarianChats(librarianID uint) ([]domain.Chat, error) {
	return c.repo.GetLibrarianChats(librarianID)
}

func (c ChatService) GetUnassignedChats() ([]domain.Chat, error) {
	return c.repo.GetUnassignedChats()
}
