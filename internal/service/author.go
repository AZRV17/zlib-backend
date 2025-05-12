package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
	"io"
	"time"
)

type AuthorService struct {
	repository repository.AuthorRepo
}

func NewAuthorService(repository repository.AuthorRepo) *AuthorService {
	return &AuthorService{repository: repository}
}

func (a AuthorService) GetAuthorByID(id uint) (*domain.Author, error) {
	return a.repository.GetAuthorByID(id)
}

func (a AuthorService) GetAuthors() ([]*domain.Author, error) {
	return a.repository.GetAuthors()
}

func (a AuthorService) CreateAuthor(authorInput *CreateAuthorInput) error {
	author := domain.Author{
		Name:      authorInput.Name,
		Lastname:  authorInput.Lastname,
		Biography: authorInput.Biography,
		Birthdate: authorInput.Birthdate,
	}

	if err := a.repository.CreateAuthor(&author); err != nil {
		return err
	}

	return nil
}

func (a AuthorService) UpdateAuthor(authorInput *UpdateAuthorInput) error {
	author := domain.Author{
		ID:        authorInput.ID,
		Name:      authorInput.Name,
		Lastname:  authorInput.Lastname,
		Biography: authorInput.Biography,
		Birthdate: authorInput.Birthdate,
	}

	if err := a.repository.UpdateAuthor(&author); err != nil {
		return err
	}

	return nil
}

func (a AuthorService) DeleteAuthor(id uint) error {
	if err := a.repository.DeleteAuthor(id); err != nil {
		return err
	}

	return nil
}

func (a AuthorService) GetAuthorBooks(id uint) ([]*domain.Book, error) {
	return a.repository.GetAuthorBooks(id)
}

func (a AuthorService) CreateAuthorBook(authorBook *domain.AuthorBook) error {

	if err := a.repository.CreateAuthorBook(authorBook); err != nil {
		return err
	}

	return nil
}

func (a AuthorService) DeleteAuthorBook(id uint) error {
	if err := a.repository.DeleteAuthorBook(id); err != nil {
		return err
	}

	return nil
}

func (a AuthorService) ExportAuthorsToCSV() ([]byte, error) {
	return a.repository.ExportAuthorsToCSV()
}

func (a AuthorService) ImportAuthorsFromCSV(data []byte) (int, error) {
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
		if len(record) < 3 {
			continue // Пропускаем некорректные записи
		}

		// Обработка полей CSV
		// [Name, Lastname, Biography, Birthdate]
		name := record[0]
		lastname := record[1]
		biography := record[2]
		birthdate := time.Now() // По умолчанию текущая дата

		if len(record) > 3 && record[3] != "" {
			parsedDate, err := time.Parse("2006-01-02", record[3])
			if err == nil {
				birthdate = parsedDate
			}
		}

		// Создаем нового автора
		authorInput := &CreateAuthorInput{
			Name:      name,
			Lastname:  lastname,
			Biography: biography,
			Birthdate: birthdate,
		}

		// Проверяем существование автора перед созданием
		existingAuthors, _ := a.repository.GetAuthors()
		exists := false

		for _, existing := range existingAuthors {
			if existing.Name == name && existing.Lastname == lastname {
				exists = true
				break
			}
		}

		if exists {
			continue
		}

		// Создаем автора
		err = a.CreateAuthor(authorInput)
		if err != nil {
			return importedCount, fmt.Errorf("ошибка при создании автора: %w", err)
		}

		importedCount++
	}

	return importedCount, nil
}
