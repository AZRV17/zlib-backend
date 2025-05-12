package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
	"io"
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

func (p PublisherService) ExportPublishersToCSV() ([]byte, error) {
	return p.repository.ExportPublishersToCSV()
}

func (p PublisherService) ImportPublishersFromCSV(data []byte) (int, error) {
	reader := csv.NewReader(bytes.NewReader(data))

	// Пропускаем заголовок
	if _, err := reader.Read(); err != nil {
		return 0, fmt.Errorf("ошибка при чтении заголовков CSV: %w", err)
	}

	importedCount := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return importedCount, fmt.Errorf("ошибка при чтении записи CSV: %w", err)
		}

		// Проверяем, что у нас достаточно полей
		if len(record) < 1 {
			continue // Пропускаем некорректные записи
		}

		// Обработка полей CSV
		// [Name]
		name := record[0]

		// Создаем нового издателя
		publisherInput := &CreatePublisherInput{
			Name: name,
		}

		// Проверяем существование издателя перед созданием
		existingPublishers, _ := p.repository.GetPublishers()
		exists := false

		for _, existing := range existingPublishers {
			if existing.Name == name {
				exists = true
				break
			}
		}

		if exists {
			continue
		}

		// Создаем издателя
		err = p.CreatePublisher(publisherInput)
		if err != nil {
			return importedCount, fmt.Errorf("ошибка при создании издателя: %w", err)
		}

		importedCount++
	}

	return importedCount, nil
}
