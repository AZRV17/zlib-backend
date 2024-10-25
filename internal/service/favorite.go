package service

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
)

type FavoriteService struct {
	repository repository.FavoriteRepo
}

func NewFavoriteService(repository repository.FavoriteRepo) *FavoriteService {
	return &FavoriteService{repository: repository}
}

func (f FavoriteService) GetFavoriteByID(id uint) (*domain.Favorite, error) {
	return f.repository.GetFavoriteByID(id)
}

func (f FavoriteService) GetFavorites() ([]*domain.Favorite, error) {
	return f.repository.GetFavorites()
}

func (f FavoriteService) CreateFavorite(favoriteInput *CreateFavoriteInput) error {
	favorite := domain.Favorite{
		UserID: favoriteInput.UserID,
		BookID: favoriteInput.BookID,
	}

	return f.repository.CreateFavorite(&favorite)
}

func (f FavoriteService) DeleteFavorite(id uint) error {
	return f.repository.DeleteFavorite(id)
}

func (f FavoriteService) GetFavoriteByUserID(userID uint) ([]*domain.Favorite, error) {
	return f.repository.GetFavoritesByUserID(userID)
}

func (f FavoriteService) DeleteFavoriteByUserIDAndBookID(userID uint, bookID uint) (*domain.Favorite, error) {
	return f.repository.DeleteFavoriteByUserIDAndBookID(userID, bookID)
}
