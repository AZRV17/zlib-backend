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

func (a AuthorRepository) GetAuthorByID(id int) (*domain.Author, error) {
	var author domain.Author

	tx := a.DB.Begin()

	if err := tx.First(&author, id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &author, nil
}

func (a AuthorRepository) GetAuthors() ([]*domain.Author, error) {
	var authors []*domain.Author

	tx := a.DB.Begin()

	if err := tx.Find(&authors).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return authors, nil
}

func (a AuthorRepository) CreateAuthor(author *domain.Author) error {
	tx := a.DB.Begin()

	if err := tx.Create(author).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (a AuthorRepository) UpdateAuthor(author *domain.Author) error {
	tx := a.DB.Begin()

	if err := tx.Save(author).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (a AuthorRepository) DeleteAuthor(id int) error {
	tx := a.DB.Begin()

	if err := tx.Delete(&domain.Author{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (a AuthorRepository) GetAuthorBooks(id int) ([]*domain.Book, error) {
	var books []*domain.Book

	tx := a.DB.Begin()

	if err := tx.Model(&domain.Author{ID: id}).Association("Books").Find(&books); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return books, nil
}
