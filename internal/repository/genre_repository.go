package repository

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/gorm"
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
