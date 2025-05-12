package service

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
	"strings"
)

type FavoriteService struct {
	repository repository.FavoriteRepo
}

func NewFavoriteService(repository repository.FavoriteRepo) *FavoriteService {
	return &FavoriteService{repository: repository}
}

func (f FavoriteService) GetFavoriteByID(id uint) (*domain.Favorite, error) {
	favorite, err := f.repository.GetFavoriteByID(id)
	if err != nil {
		return nil, err
	}

	if favorite.Book.Picture == "" || strings.HasPrefix(favorite.Book.Picture, "http") {
		return favorite, nil
	}

	favorite.Book.Picture = "http://localhost:8080/" + favorite.Book.Picture
	return favorite, err
}

func (f FavoriteService) GetFavorites() ([]*domain.Favorite, error) {
	favorites, err := f.repository.GetFavorites()
	if err != nil {
		return nil, err
	}

	for _, favorite := range favorites {
		if favorite.Book.Picture == "" || strings.HasPrefix(favorite.Book.Picture, "http") {
			continue
		} else {
			favorite.Book.Picture = "http://localhost:8080/" + favorite.Book.Picture
		}
	}

	return favorites, nil
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
	favorites, err := f.repository.GetFavoritesByUserID(userID)
	if err != nil {
		return nil, err
	}

	for _, favorite := range favorites {
		if favorite.Book.Picture == "" || strings.HasPrefix(favorite.Book.Picture, "http") {
			continue
		} else {
			favorite.Book.Picture = "http://localhost:8080/" + favorite.Book.Picture
		}
	}

	return favorites, nil
}

func (f FavoriteService) DeleteFavoriteByUserIDAndBookID(userID uint, bookID uint) (*domain.Favorite, error) {
	return f.repository.DeleteFavoriteByUserIDAndBookID(userID, bookID)
}
