package repository

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/gorm"
)

type PublisherRepository struct {
	DB *gorm.DB
}

func NewPublisherRepository(db *gorm.DB) *PublisherRepository {
	return &PublisherRepository{DB: db}
}

func (p PublisherRepository) GetPublisherByID(id uint) (*domain.Publisher, error) {
	var publisher domain.Publisher

	tx := p.DB.Begin()

	if err := tx.First(&publisher, id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &publisher, nil
}

func (p PublisherRepository) GetPublishers() ([]*domain.Publisher, error) {
	var publishers []*domain.Publisher

	tx := p.DB.Begin()

	if err := tx.Find(&publishers).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return publishers, nil
}

func (p PublisherRepository) CreatePublisher(publisher *domain.Publisher) error {
	tx := p.DB.Begin()

	if err := tx.Create(publisher).Error; err != nil {
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

func (p PublisherRepository) UpdatePublisher(publisher *domain.Publisher) error {
	tx := p.DB.Begin()

	if err := tx.Where("id = ?", publisher.ID).Save(publisher).Error; err != nil {
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

func (p PublisherRepository) DeletePublisher(id uint) error {
	tx := p.DB.Begin()

	if err := tx.Delete(&domain.Publisher{}, id).Error; err != nil {
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
