package service

import (
	"errors"
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
	"time"
)

type BookService struct {
	bookRepo        repository.BookRepo
	reservationRepo repository.ReservationRepo
}

func NewBookService(bookRepo repository.BookRepo, reservationRepo repository.ReservationRepo) *BookService {
	return &BookService{bookRepo: bookRepo, reservationRepo: reservationRepo}
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

func (b BookService) GetBookUniqueCodes(id uint) ([]*domain.UniqueCode, error) {
	return b.bookRepo.GetBookUniqueCodes(id)
}

func (b BookService) ReserveBook(bookID, userID uint) (*domain.UniqueCode, error) {
	book, err := b.bookRepo.GetBookByID(bookID)
	if err != nil {
		return nil, err
	}

	codes, err := b.bookRepo.GetBookUniqueCodes(book.ID)
	if err != nil {
		return nil, err
	}

	var code *domain.UniqueCode

	for _, c := range codes {
		if c.IsAvailable {
			code = c
		}
	}

	if code == nil {
		return nil, errors.New("no available books")
	}

	err = b.reservationRepo.CreateReservation(
		&domain.Reservation{
			BookID:       book.ID,
			UserID:       userID,
			UniqueCodeID: code.ID,
			Status:       "reserved",
			DateOfIssue:  time.Now(),
			DateOfReturn: time.Now().Add(7 * 24 * time.Hour),
		},
	)
	if err != nil {
		return nil, err
	}

	code.IsAvailable = false

	err = b.bookRepo.UpdateUniqueCode(code)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//func (b BookService) GetAggregatedBooks() ([]domain.AggregatedBook, error) {
//	return b.repository.GetAggregatedBooks()
//}
