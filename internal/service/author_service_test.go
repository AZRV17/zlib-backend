package service

import (
	"errors"
	"testing"
	"time"

	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthorRepo - мок для AuthorRepo
type MockAuthorRepo struct {
	mock.Mock
}

func (m *MockAuthorRepo) GetAuthorByID(id uint) (*domain.Author, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Author), args.Error(1)
}

func (m *MockAuthorRepo) GetAuthors() ([]*domain.Author, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Author), args.Error(1)
}

func (m *MockAuthorRepo) CreateAuthor(author *domain.Author) error {
	args := m.Called(author)
	return args.Error(0)
}

func (m *MockAuthorRepo) UpdateAuthor(author *domain.Author) error {
	args := m.Called(author)
	return args.Error(0)
}

func (m *MockAuthorRepo) DeleteAuthor(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAuthorRepo) GetAuthorBooks(id uint) ([]*domain.Book, error) {
	args := m.Called(id)
	return args.Get(0).([]*domain.Book), args.Error(1)
}

func (m *MockAuthorRepo) CreateAuthorBook(authorBook *domain.AuthorBook) error {
	args := m.Called(authorBook)
	return args.Error(0)
}

func (m *MockAuthorRepo) DeleteAuthorBook(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAuthorRepo) UpdateAuthorBook(authorBook *domain.AuthorBook) error {
	args := m.Called(authorBook)
	return args.Error(0)
}

func (m *MockAuthorRepo) ExportAuthorsToCSV() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func TestAuthorService_GetAuthorByID(t *testing.T) {
	mockRepo := new(MockAuthorRepo)
	authorService := NewAuthorService(mockRepo)

	mockAuthor := &domain.Author{
		ID:        1,
		Name:      "Leo",
		Lastname:  "Tolstoy",
		Biography: "Famous Russian writer",
		Birthdate: time.Date(1828, 9, 9, 0, 0, 0, 0, time.UTC),
	}

	t.Run("successful get author", func(t *testing.T) {
		mockRepo.On("GetAuthorByID", uint(1)).Return(mockAuthor, nil)

		author, err := authorService.GetAuthorByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, author)
		assert.Equal(t, "Leo", author.Name)
		assert.Equal(t, "Tolstoy", author.Lastname)
	})

	t.Run("author not found", func(t *testing.T) {
		mockRepo.On("GetAuthorByID", uint(999)).Return(nil, errors.New("author not found"))

		author, err := authorService.GetAuthorByID(999)
		assert.Error(t, err)
		assert.Nil(t, author)
		assert.Equal(t, "author not found", err.Error())
	})
}

func TestAuthorService_CreateAuthor(t *testing.T) {
	birthdate := time.Date(1828, 9, 9, 0, 0, 0, 0, time.UTC)

	t.Run("successful create author", func(t *testing.T) {
		mockRepo := new(MockAuthorRepo)
		authorService := NewAuthorService(mockRepo)

		mockRepo.On("CreateAuthor", mock.Anything).Return(nil)

		mockInput := &CreateAuthorInput{
			Name:      "Leo",
			Lastname:  "Tolstoy",
			Biography: "Famous Russian writer",
			Birthdate: birthdate,
		}

		err := authorService.CreateAuthor(mockInput)
		assert.NoError(t, err)
	})

	t.Run("create author error", func(t *testing.T) {
		mockRepo := new(MockAuthorRepo)
		authorService := NewAuthorService(mockRepo)

		mockRepo.On("CreateAuthor", mock.Anything).Return(errors.New("create author failed"))

		mockInput := &CreateAuthorInput{
			Name:      "Leo",
			Lastname:  "Tolstoy",
			Biography: "Famous Russian writer",
			Birthdate: birthdate,
		}

		err := authorService.CreateAuthor(mockInput)
		assert.Error(t, err)
		assert.Equal(t, "create author failed", err.Error())
	})
}
