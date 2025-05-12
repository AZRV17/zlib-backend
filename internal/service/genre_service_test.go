package service

import (
	"errors"
	"testing"

	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGenreRepo - мок для GenreRepo
type MockGenreRepo struct {
	mock.Mock
}

func (m *MockGenreRepo) GetGenreByID(id uint) (*domain.Genre, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Genre), args.Error(1)
}

func (m *MockGenreRepo) GetGenres() ([]*domain.Genre, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Genre), args.Error(1)
}

func (m *MockGenreRepo) CreateGenre(genre *domain.Genre) error {
	args := m.Called(genre)
	return args.Error(0)
}

func (m *MockGenreRepo) UpdateGenre(genre *domain.Genre) error {
	args := m.Called(genre)
	return args.Error(0)
}

func (m *MockGenreRepo) DeleteGenre(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockGenreRepo) ExportGenresToCSV() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func TestGenreService_GetGenreByID(t *testing.T) {
	mockRepo := new(MockGenreRepo)
	genreService := NewGenreService(mockRepo)

	mockGenre := &domain.Genre{
		ID:          1,
		Name:        "Fiction",
		Description: "Fictional literature",
	}

	t.Run("successful get genre", func(t *testing.T) {
		mockRepo.On("GetGenreByID", uint(1)).Return(mockGenre, nil)

		genre, err := genreService.GetGenreByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, genre)
		assert.Equal(t, "Fiction", genre.Name)
		assert.Equal(t, "Fictional literature", genre.Description)
	})

	t.Run("genre not found", func(t *testing.T) {
		mockRepo.On("GetGenreByID", uint(999)).Return(nil, errors.New("genre not found"))

		genre, err := genreService.GetGenreByID(999)
		assert.Error(t, err)
		assert.Nil(t, genre)
		assert.Equal(t, "genre not found", err.Error())
	})
}

func TestGenreService_CreateGenre(t *testing.T) {
	t.Run("successful create genre", func(t *testing.T) {
		mockRepo := new(MockGenreRepo)
		genreService := NewGenreService(mockRepo)

		mockRepo.On("CreateGenre", mock.Anything).Return(nil)

		mockInput := &CreateGenreInput{
			Name:        "Science Fiction",
			Description: "Genre of speculative fiction",
		}

		err := genreService.CreateGenre(mockInput)
		assert.NoError(t, err)
	})

	t.Run("create genre error", func(t *testing.T) {
		mockRepo := new(MockGenreRepo)
		genreService := NewGenreService(mockRepo)

		mockRepo.On("CreateGenre", mock.Anything).Return(errors.New("create genre failed"))

		mockInput := &CreateGenreInput{
			Name:        "Fantasy", // Используем другой жанр для этого теста
			Description: "Genre involving magic and supernatural elements",
		}

		err := genreService.CreateGenre(mockInput)
		assert.Error(t, err)
		assert.Equal(t, "create genre failed", err.Error())
	})
}
