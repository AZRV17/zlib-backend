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

	tx := a.DB.Begin()

	if err := tx.First(&author, id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &author, nil
}

func (a AuthorRepository) GetAuthors() ([]*domain.Author, error) {
	var authors []*domain.Author

	tx := a.DB.Begin()

	if err := tx.Preload("Books").Find(&authors).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return authors, nil
}

func (a AuthorRepository) CreateAuthor(author *domain.Author) error {
	tx := a.DB.Begin()

	if err := tx.Create(author).Error; err != nil {
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

func (a AuthorRepository) UpdateAuthor(author *domain.Author) error {
	tx := a.DB.Begin()

	if err := tx.Model(&domain.Author{}).Where("id = ?", author.ID).Save(author).Error; err != nil {
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

func (a AuthorRepository) DeleteAuthor(id uint) error {
	tx := a.DB.Begin()

	if err := tx.Delete(&domain.Author{}, id).Error; err != nil {
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