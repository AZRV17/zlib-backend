package repository

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/gorm"
)

type BookRepository struct {
	DB *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{DB: db}
}

func (b BookRepository) GetBookByID(id uint) (*domain.Book, error) {
	var book domain.Book

	tx := b.DB.Begin()

	if err := tx.Model(&domain.Book{}).Preload("Author").Preload("Genre").Preload("Publisher").First(
		&book,
		id,
	).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (b BookRepository) GetBooks() ([]*domain.Book, error) {
	var books []*domain.Book

	if err := b.DB.Model(&domain.Book{}).Joins("Author").Joins("Genre").Joins("Publisher").Find(&books).Error; err != nil {
		return nil, err
	}

	return books, nil
}

func (b BookRepository) CreateBook(book *domain.Book) error {
	if err := b.DB.Model(&domain.Book{}).Create(book).Error; err != nil {
		return err
	}

	return nil
}

func (b BookRepository) UpdateBook(book *domain.Book) error {
	if err := b.DB.Model(&domain.Book{}).Save(book).Error; err != nil {
		return err
	}

	return nil
}

func (b BookRepository) DeleteBook(id uint) error {
	if err := b.DB.Model(&domain.Book{}).Delete(&domain.Book{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (b BookRepository) GetBookByTitle(title string) (*domain.Book, error) {
	var book domain.Book

	if err := b.DB.Model(&domain.Book{}).Where("title = ?", title).First(&book).Error; err != nil {
		return nil, err
	}

	return &book, nil
}

func (b BookRepository) GetBookByUniqueCode(code uint) (*domain.Book, error) {
	var book domain.Book

	if err := b.DB.Model(&domain.Book{}).Where("unique_code = ?", code).First(&book).Error; err != nil {
		return nil, err
	}

	return &book, nil
}

func (b BookRepository) GetGroupedBooksByTitle() ([]*domain.Book, error) {
	var books []*domain.Book

	if err := b.DB.Select(
		`
				array_agg(id) as id,
				title,
				author_id,
				genre_id,
				publisher_id,
				isbn,
				year_of_publication,
				picture,
				avg(rating) as rating,
				array_agg(unique_code) as unique_code,
				array_egg(created_at) as created_at,
				array_agg(updated_at) as updated_at,
				`,
	).Group("title").Error; err != nil {
		return nil, err
	}

	return books, nil
}

func (b BookRepository) GetBookUniqueCodes(id uint) ([]*domain.UniqueCode, error) {
	var codes []*domain.UniqueCode

	tx := b.DB.Model(&domain.UniqueCode{}).Begin()

	if err := tx.Preload("Book").Preload("Book.Author").Where("book_id = ?", id).Find(&codes).Error; err != nil {
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return codes, nil
}

func (b BookRepository) CreateUniqueCode(uniqueCode *domain.UniqueCode) error {
	tx := b.DB.Begin()

	if err := tx.Create(uniqueCode).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (b BookRepository) DeleteUniqueCode(id uint) error {
	if err := b.DB.Model(&domain.UniqueCode{}).Delete(&domain.UniqueCode{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (b BookRepository) UpdateUniqueCode(uniqueCode *domain.UniqueCode) error {
	tx := b.DB.Begin()

	if err := tx.Model(&domain.UniqueCode{}).Save(uniqueCode).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}
