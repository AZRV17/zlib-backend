package service

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
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
