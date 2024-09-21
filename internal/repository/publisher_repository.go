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

	if err := p.DB.First(&publisher, id).Error; err != nil {
		return nil, err
	}

	return &publisher, nil
}

func (p PublisherRepository) GetPublishers() ([]*domain.Publisher, error) {
	var publishers []*domain.Publisher

	if err := p.DB.Find(&publishers).Error; err != nil {
		return nil, err
	}

	return publishers, nil
}

func (p PublisherRepository) CreatePublisher(publisher *domain.Publisher) error {
	if err := p.DB.Create(publisher).Error; err != nil {
		return err
	}

	return nil
}

func (p PublisherRepository) UpdatePublisher(publisher *domain.Publisher) error {
	if err := p.DB.Save(publisher).Error; err != nil {
		return err
	}

	return nil
}

func (p PublisherRepository) DeletePublisher(id uint) error {
	if err := p.DB.Delete(&domain.Publisher{}, id).Error; err != nil {
		return err
	}

	return nil
}
