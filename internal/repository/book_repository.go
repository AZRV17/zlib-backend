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

	if err := b.DB.First(&book, id).Error; err != nil {
		return nil, err
	}

	return &book, nil
}

func (b BookRepository) GetBooks() ([]*domain.Book, error) {
	var books []*domain.Book

	if err := b.DB.Find(&books).Error; err != nil {
		return nil, err
	}

	return books, nil
}

func (b BookRepository) CreateBook(book *domain.Book) error {
	if err := b.DB.Create(book).Error; err != nil {
		return err
	}

	return nil
}

func (b BookRepository) UpdateBook(book *domain.Book) error {
	if err := b.DB.Save(book).Error; err != nil {
		return err
	}

	return nil
}

func (b BookRepository) DeleteBook(id uint) error {
	if err := b.DB.Delete(&domain.Book{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (b BookRepository) GetBookByTitle(title string) (*domain.Book, error) {
	var book domain.Book

	if err := b.DB.Where("title = ?", title).First(&book).Error; err != nil {
		return nil, err
	}

	return &book, nil
}

func (b BookRepository) GetBookByUniqueCode(code uint) (*domain.Book, error) {
	var book domain.Book

	if err := b.DB.Where("unique_code = ?", code).First(&book).Error; err != nil {
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
