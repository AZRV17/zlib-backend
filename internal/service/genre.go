package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
	"io"
)

type GenreService struct {
	repository repository.GenreRepo
}

func NewGenreService(repo repository.GenreRepo) *GenreService {
	return &GenreService{repository: repo}
}

func (g GenreService) GetGenreByID(id uint) (*domain.Genre, error) {
	return g.repository.GetGenreByID(id)
}

func (g GenreService) GetGenres() ([]*domain.Genre, error) {
	return g.repository.GetGenres()
}

func (g GenreService) CreateGenre(genreInput *CreateGenreInput) error {
	genre := domain.Genre{
		Name:        genreInput.Name,
		Description: genreInput.Description,
	}

	return g.repository.CreateGenre(&genre)
}

func (g GenreService) UpdateGenre(genreInput *UpdateGenreInput) error {
	genre := domain.Genre{
		ID:          genreInput.ID,
		Name:        genreInput.Name,
		Description: genreInput.Description,
	}

	return g.repository.UpdateGenre(&genre)
}

func (g GenreService) DeleteGenre(id uint) error {
	return g.repository.DeleteGenre(id)

}

func (g GenreService) ExportGenresToCSV() ([]byte, error) {
	return g.repository.ExportGenresToCSV()
}

func (g GenreService) ImportGenresFromCSV(data []byte) (int, error) {
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
		// [Name, Description]
		name := record[0]
		description := ""
		if len(record) > 1 {
			description = record[1]
		}

		// Создаем новый жанр
		genreInput := &CreateGenreInput{
			Name:        name,
			Description: description,
		}

		// Проверяем существование жанра перед созданием
		existingGenres, _ := g.repository.GetGenres()
		exists := false

		for _, existing := range existingGenres {
			if existing.Name == name {
				exists = true
				break
			}
		}

		if exists {
			continue
		}

		// Создаем жанр
		err = g.CreateGenre(genreInput)
		if err != nil {
			return importedCount, fmt.Errorf("ошибка при создании жанра: %w", err)
		}

		importedCount++
	}

	return importedCount, nil
}
