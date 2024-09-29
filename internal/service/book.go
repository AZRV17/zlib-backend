package service

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
)

type BookService struct {
	repository repository.BookRepo
}

func NewBookService(repository repository.BookRepo) *BookService {
	return &BookService{repository: repository}
}

func (b BookService) GetBookByID(id uint) (*domain.Book, error) {
	return b.repository.GetBookByID(id)
}

func (b BookService) GetBooks() ([]*domain.Book, error) {
	return b.repository.GetBooks()
}

func (b BookService) CreateBook(bookInput *CreateBookInput) error {
	book := &domain.Book{
		Title:             bookInput.Title,
		AuthorID:          bookInput.AuthorID,
		GenreID:           bookInput.GenreID,
		PublisherID:       bookInput.PublisherID,
		ISBN:              bookInput.ISBN,
		YearOfPublication: bookInput.YearOfPublication,
		Picture:           bookInput.Picture,
		Rating:            bookInput.Rating,
		UniqueCode:        bookInput.UniqueCode,
	}

	return b.repository.CreateBook(book)
}

func (b BookService) UpdateBook(bookInput *UpdateBookInput) error {
	book, err := b.repository.GetBookByID(bookInput.ID)
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
	book.UniqueCode = bookInput.UniqueCode

	return b.repository.UpdateBook(book)
}

func (b BookService) DeleteBook(id uint) error {
	return b.repository.DeleteBook(id)
}

func (b BookService) GetBookByTitle(title string) (*domain.Book, error) {
	return b.repository.GetBookByTitle(title)
}

func (b BookService) GetBookByUniqueCode(code uint) (*domain.Book, error) {
	return b.repository.GetBookByUniqueCode(code)
}

func (b BookService) GetGroupedBooksByTitle() ([]*domain.Book, error) {
	return b.repository.GetGroupedBooksByTitle()
}
