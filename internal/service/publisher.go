package service

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
)

type PublisherService struct {
	repository repository.PublisherRepo
}

func NewPublisherService(repo repository.PublisherRepo) *PublisherService {
	return &PublisherService{repository: repo}
}

func (p PublisherService) GetPublisherByID(id uint) (*domain.Publisher, error) {
	return p.repository.GetPublisherByID(id)
}

func (p PublisherService) GetPublishers() ([]*domain.Publisher, error) {
	return p.repository.GetPublishers()
}

func (p PublisherService) CreatePublisher(publisherInput *CreatePublisherInput) error {
	publisher := &domain.Publisher{
		Name: publisherInput.Name,
	}

	return p.repository.CreatePublisher(publisher)
}

func (p PublisherService) UpdatePublisher(publisherInput *UpdatePublisherInput) error {
	publisher, err := p.repository.GetPublisherByID(publisherInput.ID)
	if err != nil {
		return err
	}

	publisher.Name = publisherInput.Name

	return p.repository.UpdatePublisher(publisher)
}

func (p PublisherService) DeletePublisher(id uint) error {
	return p.repository.DeletePublisher(id)
}
