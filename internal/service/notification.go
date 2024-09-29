package service

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
)

type NotificationService struct {
	repository repository.NotificationRepo
}

func NewNotificationService(repo repository.NotificationRepo) *NotificationService {
	return &NotificationService{repository: repo}
}

func (n NotificationService) GetNotificationByID(id uint) (*domain.Notification, error) {
	return n.repository.GetNotificationByID(id)
}

func (n NotificationService) GetNotifications() ([]*domain.Notification, error) {
	return n.repository.GetNotifications()
}

func (n NotificationService) CreateNotification(notificationInput *CreateNotificationInput) error {
	notification := domain.Notification{
		UserID:  notificationInput.UserID,
		Message: notificationInput.Message,
	}
	return n.repository.CreateNotification(&notification)
}

func (n NotificationService) DeleteNotification(id uint) error {
	return n.repository.DeleteNotification(id)
}
