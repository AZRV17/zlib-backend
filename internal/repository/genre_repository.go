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

	if err := g.DB.First(&genre, id).Error; err != nil {
		return nil, err
	}

	return &genre, nil
}

func (g GenreRepository) GetGenres() ([]*domain.Genre, error) {
	var genres []*domain.Genre

	if err := g.DB.Find(&genres).Error; err != nil {
		return nil, err
	}

	return genres, nil
}

func (g GenreRepository) CreateGenre(genre *domain.Genre) error {
	if err := g.DB.Create(genre).Error; err != nil {
		return err
	}

	return nil
}

func (g GenreRepository) UpdateGenre(genre *domain.Genre) error {
	if err := g.DB.Save(genre).Error; err != nil {
		return err
	}

	return nil
}

func (g GenreRepository) DeleteGenre(id uint) error {
	if err := g.DB.Delete(&domain.Genre{}, id).Error; err != nil {
		return err
	}

	return nil
}
