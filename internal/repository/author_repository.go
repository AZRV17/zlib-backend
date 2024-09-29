package repository

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/gorm"
)

type AuthorRepository struct {
	DB *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) *AuthorRepository {
	return &AuthorRepository{DB: db}
}

func (a AuthorRepository) GetAuthorByID(id uint) (*domain.Author, error) {
	var author domain.Author

	if err := a.DB.First(&author, id).Error; err != nil {
		return nil, err
	}

	return &author, nil
}

func (a AuthorRepository) GetAuthors() ([]*domain.Author, error) {
	var authors []*domain.Author

	if err := a.DB.Find(&authors).Error; err != nil {
		a.DB.Rollback()
		return nil, err
	}

	return authors, nil
}

func (a AuthorRepository) CreateAuthor(author *domain.Author) error {
	if err := a.DB.Create(author).Error; err != nil {
		a.DB.Rollback()
		return err
	}

	return nil
}

func (a AuthorRepository) UpdateAuthor(author *domain.Author) error {
	if err := a.DB.Save(author).Error; err != nil {
		a.DB.Rollback()
		return err
	}

	return nil
}

func (a AuthorRepository) DeleteAuthor(id uint) error {
	if err := a.DB.Delete(&domain.Author{}, id).Error; err != nil {
		a.DB.Rollback()
		return err
	}

	return nil
}

func (a AuthorRepository) GetAuthorBooks(id uint) ([]*domain.Book, error) {
	var books []*domain.Book

	if err := a.DB.Model(&domain.Author{ID: id}).Association("Books").Find(&books); err != nil {
		a.DB.Rollback()
		return nil, err
	}

	return books, nil
}

func (a AuthorRepository) CreateAuthorBook(authorBook *domain.AuthorBook) error {
	if err := a.DB.Model(&domain.AuthorBook{}).Create(authorBook).Error; err != nil {
		a.DB.Rollback()
		return err
	}

	return nil
}

func (a AuthorRepository) DeleteAuthorBook(id uint) error {
	if err := a.DB.Delete(&domain.AuthorBook{}, id).Error; err != nil {
		a.DB.Rollback()
		return err
	}

	return nil
}

func (a AuthorRepository) UpdateAuthorBook(authorBook *domain.AuthorBook) error {
	if err := a.DB.Save(authorBook).Error; err != nil {
		a.DB.Rollback()
		return err
	}

	return nil
}
