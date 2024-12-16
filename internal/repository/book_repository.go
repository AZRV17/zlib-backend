package repository

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/gorm"
	"strconv"
)

type BookRepository struct {
	DB *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{DB: db}
}

func (b *BookRepository) GetBookByID(id uint) (*domain.Book, error) {
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

func (b *BookRepository) GetBooks() ([]*domain.Book, error) {
	var books []*domain.Book

	tx := b.DB.Begin()

	if err := tx.Model(&domain.Book{}).
		Preload("Author").
		Preload("Genre").
		Preload("Publisher").
		Preload("UniqueCodes").
		Find(&books).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (b *BookRepository) CreateBook(book *domain.Book) error {
	if err := b.DB.Model(&domain.Book{}).Create(book).Error; err != nil {
		return err
	}

	return nil
}

func (b *BookRepository) UpdateBook(book *domain.Book) error {
	if err := b.DB.Model(&domain.Book{}).Where("id = ?", book.ID).Save(book).Error; err != nil {
		return err
	}

	return nil
}

func (b *BookRepository) DeleteBook(id uint) error {
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
	tx := b.DB.Begin()

	if err := tx.Model(&domain.UniqueCode{}).Delete(&domain.UniqueCode{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (b BookRepository) UpdateUniqueCode(uniqueCode *domain.UniqueCode) error {
	tx := b.DB.Begin()

	if err := tx.Model(&domain.UniqueCode{}).Where("id = ?", uniqueCode.ID).Save(uniqueCode).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (b BookRepository) GetBookByIDWithTransactions(id uint, tx *gorm.DB) (*domain.Book, error) {
	var book domain.Book

	if err := tx.Model(&domain.Book{}).Preload("Author").Preload("Genre").Preload("Publisher").First(
		&book,
		id,
	).Error; err != nil {
		return nil, err
	}

	return &book, nil
}

func (b BookRepository) GetBookUniqueCodesWithTransactions(id uint, tx *gorm.DB) ([]*domain.UniqueCode, error) {
	var codes []*domain.UniqueCode

	if err := tx.Model(&domain.UniqueCode{}).Preload("Book").Preload("Book.Author").Where(
		"book_id = ?",
		id,
	).Find(&codes).Error; err != nil {
		return nil, err
	}

	return codes, nil
}

func (b BookRepository) UpdateUniqueCodeWithTransactions(uniqueCode *domain.UniqueCode, tx *gorm.DB) error {
	if err := tx.Model(&domain.UniqueCode{}).Where("id = ?", uniqueCode.ID).Save(uniqueCode).Error; err != nil {
		return err
	}

	return nil
}

func (b BookRepository) GetUniqueCodes() ([]*domain.UniqueCode, error) {
	var codes []*domain.UniqueCode

	tx := b.DB.Begin()

	if err := tx.Model(&domain.UniqueCode{}).Preload("Book").Find(&codes).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return codes, nil
}

func (b BookRepository) GetUniqueCodeByID(id uint) (*domain.UniqueCode, error) {
	var code domain.UniqueCode

	tx := b.DB.Begin()

	if err := tx.Model(&domain.UniqueCode{}).Preload("Book").First(&code, id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &code, nil
}

func (b BookRepository) GetBooksWithPagination(limit int, offset int) ([]*domain.Book, error) {
	var books []*domain.Book

	tx := b.DB.Begin()

	if err := tx.Model(&domain.Book{}).Limit(limit).Offset(offset).Find(&books).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (b BookRepository) FindBookByTitle(limit int, offset int, title string) ([]*domain.Book, error) {
	var books []*domain.Book

	tx := b.DB.Begin()

	if err := tx.Model(&domain.Book{}).Where(
		"lower(title) like lower(?)",
		"%"+title+"%",
	).Limit(limit).Offset(offset).Find(&books).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return books, nil
}

func (b BookRepository) ExportBooksToCSV() ([]byte, error) {
	books, err := b.GetBooks()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	headers := []string{"ID", "Название", "Автор", "Жанр", "Издательство", "Дата выхода"}

	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	for _, book := range books {
		row := []string{
			strconv.FormatUint(uint64(book.ID), 10),
			book.Title,
			book.Author.Name + " " + book.Author.Lastname,
			book.Genre.Name,
			book.Publisher.Name,
			book.YearOfPublication.Format("2006-01-02"),
		}

		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("error flushing writer: %w", err)
	}

	return buf.Bytes(), nil
}
