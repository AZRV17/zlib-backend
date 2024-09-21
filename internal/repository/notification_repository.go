package repository

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	DB *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{
		DB: db,
	}
}

func (n NotificationRepository) GetNotificationByID(id uint) (*domain.Notification, error) {
	var notification domain.Notification

	if err := n.DB.First(&notification, id).Error; err != nil {
		return nil, err
	}

	return &notification, nil
}

func (n NotificationRepository) GetNotifications() ([]*domain.Notification, error) {
	var notifications []*domain.Notification

	if err := n.DB.Find(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}

func (n NotificationRepository) CreateNotification(notification *domain.Notification) error {
	if err := n.DB.Create(notification).Error; err != nil {
		return err
	}

	return nil
}

func (n NotificationRepository) UpdateNotification(notification *domain.Notification) error {
	if err := n.DB.Save(notification).Error; err != nil {
		return err
	}

	return nil
}

func (n NotificationRepository) DeleteNotification(id uint) error {
	if err := n.DB.Delete(&domain.Notification{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (n NotificationRepository) GetNotificationsByUserID(id uint) ([]*domain.Notification, error) {
	var notifications []*domain.Notification

	if err := n.DB.Where("user_id = ?", id).Find(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}
