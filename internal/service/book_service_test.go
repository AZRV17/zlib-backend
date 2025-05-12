package service

import (
	"errors"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBookRepo struct {
	mock.Mock
}

func (m *MockBookRepo) GetBookByID(id uint) (*domain.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Book), args.Error(1)
}

func (m *MockBookRepo) GetBooks() ([]*domain.Book, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Book), args.Error(1)
}

func (m *MockBookRepo) CreateBook(book *domain.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookRepo) UpdateBook(book *domain.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookRepo) DeleteBook(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBookRepo) GetBookByTitle(title string) (*domain.Book, error) {
	args := m.Called(title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Book), args.Error(1)
}

func (m *MockBookRepo) GetBookByUniqueCode(code uint) (*domain.Book, error) {
	args := m.Called(code)
	return args.Get(0).(*domain.Book), args.Error(1)
}

func (m *MockBookRepo) GetGroupedBooksByTitle() ([]*domain.Book, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Book), args.Error(1)
}

func (m *MockBookRepo) GetBookUniqueCodes(id uint) ([]*domain.UniqueCode, error) {
	args := m.Called(id)
	return args.Get(0).([]*domain.UniqueCode), args.Error(1)
}

func (m *MockBookRepo) CreateUniqueCode(uniqueCode *domain.UniqueCode) error {
	args := m.Called(uniqueCode)
	return args.Error(0)
}

func (m *MockBookRepo) DeleteUniqueCode(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBookRepo) UpdateUniqueCode(uniqueCode *domain.UniqueCode) error {
	args := m.Called(uniqueCode)
	return args.Error(0)
}

func (m *MockBookRepo) GetBookByIDWithTransactions(id uint, tx *gorm.DB) (*domain.Book, error) {
	args := m.Called(id, tx)
	return args.Get(0).(*domain.Book), args.Error(1)
}

func (m *MockBookRepo) GetBookUniqueCodesWithTransactions(id uint, tx *gorm.DB) ([]*domain.UniqueCode, error) {
	args := m.Called(id, tx)
	return args.Get(0).([]*domain.UniqueCode), args.Error(1)
}

func (m *MockBookRepo) UpdateUniqueCodeWithTransactions(uniqueCode *domain.UniqueCode, tx *gorm.DB) error {
	args := m.Called(uniqueCode, tx)
	return args.Error(0)
}

func (m *MockBookRepo) GetUniqueCodes() ([]*domain.UniqueCode, error) {
	args := m.Called()
	return args.Get(0).([]*domain.UniqueCode), args.Error(1)
}

func (m *MockBookRepo) GetUniqueCodeByID(id uint) (*domain.UniqueCode, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.UniqueCode), args.Error(1)
}

func (m *MockBookRepo) GetBooksWithPagination(limit int, offset int) ([]*domain.Book, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]*domain.Book), args.Error(1)
}

func (m *MockBookRepo) FindBookByTitle(limit int, offset int, title string) ([]*domain.Book, error) {
	args := m.Called(limit, offset, title)
	return args.Get(0).([]*domain.Book), args.Error(1)
}

func (m *MockBookRepo) ExportBooksToCSV() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockBookRepo) GetAudiobookFilesByBookID(bookID uint) ([]*domain.AudiobookFile, error) {
	args := m.Called(bookID)
	return args.Get(0).([]*domain.AudiobookFile), args.Error(1)
}

func (m *MockBookRepo) GetAudiobookFileByID(id uint) (*domain.AudiobookFile, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.AudiobookFile), args.Error(1)
}

func (m *MockBookRepo) CreateAudiobookFile(file *domain.AudiobookFile) error {
	args := m.Called(file)
	return args.Error(0)
}

func (m *MockBookRepo) UpdateAudiobookFile(file *domain.AudiobookFile) error {
	args := m.Called(file)
	return args.Error(0)
}

func (m *MockBookRepo) DeleteAudiobookFile(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBookRepo) FindBooks(limit int, offset int, query string) ([]*domain.Book, error) {
	args := m.Called(limit, offset, query)
	return args.Get(0).([]*domain.Book), args.Error(1)
}

func TestBookService_GetBookByID(t *testing.T) {
	mockBookRepo := new(MockBookRepo)
	mockReservationRepo := new(MockReservationRepo)
	bookService := NewBookService(mockBookRepo, mockReservationRepo, nil)

	mockBook := &domain.Book{
		ID:                1,
		Title:             "War and Peace",
		Description:       "Epic novel by Leo Tolstoy",
		YearOfPublication: time.Date(1869, 1, 1, 0, 0, 0, 0, time.UTC),
		Rating:            4.5,
	}

	t.Run(
		"successful get book", func(t *testing.T) {
			mockBookRepo.On("GetBookByID", uint(1)).Return(mockBook, nil)

			book, err := bookService.GetBookByID(1)
			assert.NoError(t, err)
			assert.NotNil(t, book)
			assert.Equal(t, "War and Peace", book.Title)
			assert.Equal(t, 4.5, float64(book.Rating))
		},
	)

	t.Run(
		"book not found", func(t *testing.T) {
			mockBookRepo.On("GetBookByID", uint(999)).Return(nil, errors.New("book not found"))

			book, err := bookService.GetBookByID(999)
			assert.Error(t, err)
			assert.Nil(t, book)
			assert.Equal(t, "book not found", err.Error())
		},
	)
}

func TestBookService_CreateBook(t *testing.T) {
	publicationDate := time.Date(1869, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run(
		"successful create book", func(t *testing.T) {
			mockBookRepo := new(MockBookRepo)
			mockReservationRepo := new(MockReservationRepo)
			bookService := NewBookService(mockBookRepo, mockReservationRepo, nil)

			mockBookRepo.On("CreateBook", mock.Anything).Return(nil)

			mockInput := &CreateBookInput{
				Title:             "War and Peace",
				Description:       "Epic novel by Leo Tolstoy",
				AuthorID:          1,
				GenreID:           2,
				PublisherID:       3,
				ISBN:              9780140447934,
				YearOfPublication: publicationDate,
				Rating:            4.5,
			}

			err := bookService.CreateBook(mockInput)
			assert.NoError(t, err)
		},
	)

	t.Run(
		"create book error", func(t *testing.T) {
			mockBookRepo := new(MockBookRepo)
			mockReservationRepo := new(MockReservationRepo)
			bookService := NewBookService(mockBookRepo, mockReservationRepo, nil)

			mockBookRepo.On("CreateBook", mock.Anything).Return(errors.New("create book failed"))

			mockInput := &CreateBookInput{
				Title:             "Anna Karenina",
				Description:       "Tragic novel by Leo Tolstoy",
				AuthorID:          1,
				GenreID:           2,
				PublisherID:       3,
				ISBN:              9780198748847,
				YearOfPublication: time.Date(1877, 1, 1, 0, 0, 0, 0, time.UTC),
				Rating:            4.7,
			}

			err := bookService.CreateBook(mockInput)
			assert.Error(t, err)
			assert.Equal(t, "create book failed", err.Error())
		},
	)
}

func TestBookService_GetBookByTitle(t *testing.T) {
	mockBookRepo := new(MockBookRepo)
	mockReservationRepo := new(MockReservationRepo)
	bookService := NewBookService(mockBookRepo, mockReservationRepo, nil)

	mockBook := &domain.Book{
		ID:                1,
		Title:             "War and Peace",
		Description:       "Epic novel by Leo Tolstoy",
		YearOfPublication: time.Date(1869, 1, 1, 0, 0, 0, 0, time.UTC),
		Rating:            4.5,
	}

	t.Run(
		"successful get book by title", func(t *testing.T) {
			mockBookRepo.On("GetBookByTitle", "War and Peace").Return(mockBook, nil)

			book, err := bookService.GetBookByTitle("War and Peace")
			assert.NoError(t, err)
			assert.NotNil(t, book)
			assert.Equal(t, "War and Peace", book.Title)
		},
	)

	t.Run(
		"book not found by title", func(t *testing.T) {
			mockBookRepo.On("GetBookByTitle", "Unknown Book").
				Return(nil, errors.New("book not found"))

			book, err := bookService.GetBookByTitle("Unknown Book")
			assert.Error(t, err)
			assert.Nil(t, book)
			assert.Equal(t, "book not found", err.Error())
		},
	)
}

func TestBookService_FindBookByTitle(t *testing.T) {
	mockBookRepo := new(MockBookRepo)
	mockReservationRepo := new(MockReservationRepo)
	bookService := NewBookService(mockBookRepo, mockReservationRepo, nil)

	mockBooks := []*domain.Book{
		{
			ID:                1,
			Title:             "War and Peace",
			Description:       "Epic novel by Leo Tolstoy",
			YearOfPublication: time.Date(1869, 1, 1, 0, 0, 0, 0, time.UTC),
			Rating:            4.5,
		},
		{
			ID:                2,
			Title:             "War of the Worlds",
			Description:       "Science fiction novel by H.G. Wells",
			YearOfPublication: time.Date(1898, 1, 1, 0, 0, 0, 0, time.UTC),
			Rating:            4.2,
		},
	}

	t.Run(
		"successful find books by title", func(t *testing.T) {
			mockBookRepo.On("FindBookByTitle", 10, 0, "War").Return(mockBooks, nil)

			books, err := bookService.FindBookByTitle(10, 0, "War")
			assert.NoError(t, err)
			assert.NotNil(t, books)
			assert.Equal(t, 2, len(books))
			assert.Equal(t, "War and Peace", books[0].Title)
			assert.Equal(t, "War of the Worlds", books[1].Title)
		},
	)

	t.Run(
		"no books found", func(t *testing.T) {
			mockBookRepo.On("FindBookByTitle", 10, 0, "Unknown").Return([]*domain.Book{}, nil)

			books, err := bookService.FindBookByTitle(10, 0, "Unknown")
			assert.NoError(t, err)
			assert.Empty(t, books)
		},
	)
}
