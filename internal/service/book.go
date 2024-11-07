package service

import (
	"errors"
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
	"gorm.io/gorm"
	"time"
)

type BookService struct {
	bookRepo        repository.BookRepo
	reservationRepo repository.ReservationRepo
	db              *gorm.DB
}

func NewBookService(
	bookRepo repository.BookRepo,
	reservationRepo repository.ReservationRepo,
	db *gorm.DB,
) *BookService {
	return &BookService{bookRepo: bookRepo, reservationRepo: reservationRepo, db: db}
}

func (b BookService) GetBookByID(id uint) (*domain.Book, error) {
	return b.bookRepo.GetBookByID(id)
}

func (b BookService) GetBooks() ([]*domain.Book, error) {
	return b.bookRepo.GetBooks()
}

func (b BookService) CreateBook(bookInput *CreateBookInput) error {
	book := &domain.Book{
		Title:             bookInput.Title,
		Description:       bookInput.Description,
		AuthorID:          bookInput.AuthorID,
		GenreID:           bookInput.GenreID,
		PublisherID:       bookInput.PublisherID,
		ISBN:              bookInput.ISBN,
		YearOfPublication: bookInput.YearOfPublication,
		Picture:           bookInput.Picture,
		Rating:            bookInput.Rating,
	}

	return b.bookRepo.CreateBook(book)
}

func (b BookService) UpdateBook(bookInput *UpdateBookInput) error {
	book, err := b.bookRepo.GetBookByID(bookInput.ID)
	if err != nil {
		return err
	}

	book.Title = bookInput.Title
	book.AuthorID = bookInput.AuthorID
	book.GenreID = bookInput.GenreID
	book.PublisherID = bookInput.PublisherID
	book.ISBN = bookInput.ISBN
	book.YearOfPublication = bookInput.YearOfPublication
	book.Picture = bookInput.Picture
	book.Rating = bookInput.Rating

	return b.bookRepo.UpdateBook(book)
}

func (b BookService) DeleteBook(id uint) error {
	return b.bookRepo.DeleteBook(id)
}

func (b BookService) GetBookByTitle(title string) (*domain.Book, error) {
	return b.bookRepo.GetBookByTitle(title)
}

func (b BookService) GetBookByUniqueCode(code uint) (*domain.Book, error) {
	return b.bookRepo.GetBookByUniqueCode(code)
}

func (b BookService) GetGroupedBooksByTitle() ([]*domain.Book, error) {
	return b.bookRepo.GetGroupedBooksByTitle()
}

func (b BookService) CreateUniqueCode(uniqueCode *domain.UniqueCode) error {
	return b.bookRepo.CreateUniqueCode(uniqueCode)
}

func (b BookService) DeleteUniqueCode(id uint) error {
	return b.bookRepo.DeleteUniqueCode(id)
}

func (b BookService) UpdateUniqueCode(uniqueCode *domain.UniqueCode) error {
	return b.bookRepo.UpdateUniqueCode(uniqueCode)
}

func (b BookService) GetUniqueCodes() ([]*domain.UniqueCode, error) {
	return b.bookRepo.GetUniqueCodes()
}

func (b BookService) GetBookUniqueCodes(id uint) ([]*domain.UniqueCode, error) {
	return b.bookRepo.GetBookUniqueCodes(id)
}

func (b BookService) ReserveBook(bookID, userID uint) (*domain.UniqueCode, error) {
	tx := b.db.Begin()

	book, err := b.bookRepo.GetBookByIDWithTransactions(bookID, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	codes, err := b.bookRepo.GetBookUniqueCodesWithTransactions(book.ID, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var code *domain.UniqueCode

	for _, c := range codes {
		if c.IsAvailable {
			code = c
		}
	}

	if code == nil {
		tx.Rollback()
		return nil, errors.New("no available books")
	}

	err = b.reservationRepo.CreateReservationWithTransactions(
		&domain.Reservation{
			BookID:       book.ID,
			UserID:       userID,
			UniqueCodeID: code.ID,
			Status:       "reserved",
			DateOfIssue:  time.Now(),
			DateOfReturn: time.Now().Add(7 * 24 * time.Hour),
		}, tx,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	code.IsAvailable = false

	err = b.bookRepo.UpdateUniqueCodeWithTransactions(code, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return code, nil
}

func (b *BookService) GetUniqueCodeByID(id uint) (*domain.UniqueCode, error) {
	return b.bookRepo.GetUniqueCodeByID(id)
}
