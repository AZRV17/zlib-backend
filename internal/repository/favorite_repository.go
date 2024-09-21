package repository

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/gorm"
)

type FavoriteRepository struct {
	DB *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{DB: db}
}

func (f FavoriteRepository) GetFavoriteByID(id uint) (*domain.Favorite, error) {
	var favorite domain.Favorite

	if err := f.DB.First(&favorite, id).Error; err != nil {
		return nil, err
	}

	return &favorite, nil
}

func (f FavoriteRepository) GetFavorites() ([]*domain.Favorite, error) {
	var favorites []*domain.Favorite

	if err := f.DB.Find(&favorites).Error; err != nil {
		return nil, err
	}

	return favorites, nil
}

func (f FavoriteRepository) CreateFavorite(favorite *domain.Favorite) error {
	if err := f.DB.Create(favorite).Error; err != nil {
		return err
	}

	return nil
}

func (f FavoriteRepository) DeleteFavorite(id uint) error {
	if err := f.DB.Delete(&domain.Favorite{}, id).Error; err != nil {
		return err
	}

	return nil
}
