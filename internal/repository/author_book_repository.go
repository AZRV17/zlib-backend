package repository

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/gorm"
)

type AuthorBookRepository struct {
	DB *gorm.DB
}

func NewAuthorBookRepository(db *gorm.DB) *AuthorBookRepository {
	return &AuthorBookRepository{
		DB: db,
	}
}

func (a *AuthorBookRepository) GetAuthorBookByID(id uint) (*domain.AuthorBook, error) {
	var authorBook domain.AuthorBook

	if err := a.DB.First(&authorBook, id).Error; err != nil {
		return nil, err
	}

	return &authorBook, nil
}

func (a AuthorBookRepository) GetAuthorBooks() ([]*domain.AuthorBook, error) {
	var authorBooks []*domain.AuthorBook

	if err := a.DB.Find(&authorBooks).Error; err != nil {
		return nil, err
	}

	return authorBooks, nil
}

func (a AuthorBookRepository) CreateAuthorBook(authorBook *domain.AuthorBook) error {
	if err := a.DB.Create(authorBook).Error; err != nil {
		return err
	}

	return nil
}

func (a AuthorBookRepository) UpdateAuthorBook(authorBook *domain.AuthorBook) error {
	if err := a.DB.Save(authorBook).Error; err != nil {
		return err
	}

	return nil
}

func (a AuthorBookRepository) DeleteAuthorBook(id uint) error {
	if err := a.DB.Delete(&domain.AuthorBook{}, id).Error; err != nil {
		return err
	}

	return nil
}
