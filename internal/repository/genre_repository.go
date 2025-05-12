package repository

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/gorm"
	"strconv"
)

type GenreRepository struct {
	DB *gorm.DB
}

func NewGenreRepository(db *gorm.DB) *GenreRepository {
	return &GenreRepository{DB: db}
}

func (g GenreRepository) GetGenreByID(id uint) (*domain.Genre, error) {
	var genre domain.Genre

	tx := g.DB.Begin()

	if err := tx.First(&genre, id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &genre, nil
}

func (g GenreRepository) GetGenres() ([]*domain.Genre, error) {
	var genres []*domain.Genre

	tx := g.DB.Begin()

	if err := tx.Find(&genres).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return genres, nil
}

func (g GenreRepository) CreateGenre(genre *domain.Genre) error {
	tx := g.DB.Begin()

	if err := tx.Create(genre).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (g GenreRepository) UpdateGenre(genre *domain.Genre) error {
	tx := g.DB.Begin()

	if err := tx.Where("id = ?", genre.ID).Save(genre).Error; err != nil {
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

func (g GenreRepository) DeleteGenre(id uint) error {
	tx := g.DB.Begin()

	if err := tx.Delete(&domain.Genre{}, id).Error; err != nil {
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

func (g GenreRepository) ExportGenresToCSV() ([]byte, error) {
	genres, err := g.GetGenres()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	headers := []string{"ID", "Название", "Описание"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	for _, genre := range genres {
		row := []string{
			strconv.FormatUint(uint64(genre.ID), 10),
			genre.Name,
			genre.Description,
		}

		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("error flushing writer: %w", err)
	}

	return buf.Bytes(), nil
}
