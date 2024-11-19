package repository

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/gorm"
)

type LogRepository struct {
	DB *gorm.DB
}

func NewLogRepository(db *gorm.DB) *LogRepository {
	return &LogRepository{DB: db}
}

func (l LogRepository) GetLogByID(id uint) (*domain.Log, error) {
	var log domain.Log

	if err := l.DB.First(&log, id).Error; err != nil {
		return nil, err
	}

	return &log, nil
}

func (l LogRepository) GetLogs() ([]*domain.Log, error) {
	var logs []*domain.Log

	tx := l.DB.Begin()

	if err := tx.Preload("User").Find(&logs).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return logs, nil
}

func (l LogRepository) CreateLog(log *domain.Log) error {
	tx := l.DB.Begin()

	if err := tx.Create(log).Error; err != nil {
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

func (l LogRepository) UpdateLog(log *domain.Log) error {
	if err := l.DB.Save(log).Error; err != nil {
		return err
	}

	return nil
}

func (l LogRepository) DeleteLog(id uint) error {
	if err := l.DB.Delete(&domain.Log{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (l LogRepository) GetLogsByUserID(id uint) ([]*domain.Log, error) {
	var logs []*domain.Log

	if err := l.DB.Where("user_id = ?", id).Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}
