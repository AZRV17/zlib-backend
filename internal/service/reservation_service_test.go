package service

import (
	"errors"
	"testing"
	"time"

	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockReservationRepo struct {
	mock.Mock
}

func (m *MockReservationRepo) GetReservationByID(id uint) (*domain.Reservation, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Reservation), args.Error(1)
}

func (m *MockReservationRepo) GetReservations() ([]*domain.Reservation, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Reservation), args.Error(1)
}

func (m *MockReservationRepo) CreateReservation(reservation *domain.Reservation) error {
	args := m.Called(reservation)
	return args.Error(0)
}

func (m *MockReservationRepo) UpdateReservation(reservation *domain.Reservation) error {
	args := m.Called(reservation)
	return args.Error(0)
}

func (m *MockReservationRepo) DeleteReservation(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockReservationRepo) CreateReservationWithTransactions(reservation *domain.Reservation, tx *gorm.DB) error {
	args := m.Called(reservation, tx)
	return args.Error(0)
}

func (m *MockReservationRepo) GetUserReservations(id uint) ([]*domain.Reservation, error) {
	args := m.Called(id)
	return args.Get(0).([]*domain.Reservation), args.Error(1)
}

func (m *MockReservationRepo) ExportReservationsToCSV() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func TestReservationService_GetReservationByID(t *testing.T) {
	mockReservationRepo := new(MockReservationRepo)
	mockBookRepo := new(MockBookRepo) // Используем уже полностью определенный мок

	reservationService := NewReservationService(mockReservationRepo, mockBookRepo)

	issueDate := time.Now().Add(-24 * time.Hour)
	returnDate := time.Now().Add(6 * 24 * time.Hour)

	mockReservation := &domain.Reservation{
		ID:           1,
		UserID:       1,
		BookID:       1,
		Status:       "active",
		DateOfIssue:  issueDate,
		DateOfReturn: returnDate,
	}

	t.Run("successful get reservation", func(t *testing.T) {
		mockReservationRepo.On("GetReservationByID", uint(1)).Return(mockReservation, nil)

		reservation, err := reservationService.GetReservationByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, reservation)
		assert.Equal(t, uint(1), reservation.ID)
		assert.Equal(t, "active", reservation.Status)
	})

	t.Run("reservation not found", func(t *testing.T) {
		mockReservationRepo.On("GetReservationByID", uint(999)).Return(nil, errors.New("reservation not found"))

		reservation, err := reservationService.GetReservationByID(999)
		assert.Error(t, err)
		assert.Nil(t, reservation)
		assert.Equal(t, "reservation not found", err.Error())
	})
}

func TestReservationService_CreateReservation(t *testing.T) {
	issueDate := time.Now()
	returnDate := time.Now().Add(7 * 24 * time.Hour)

	t.Run("successful create reservation", func(t *testing.T) {
		mockReservationRepo := new(MockReservationRepo)
		mockBookRepo := new(MockBookRepo)
		reservationService := NewReservationService(mockReservationRepo, mockBookRepo)

		mockReservationRepo.On("CreateReservation", mock.Anything).Return(nil)

		mockInput := &CreateReservationInput{
			UserID:       1,
			BookID:       1,
			Status:       "active",
			DateOfIssue:  issueDate,
			DateOfReturn: returnDate,
		}

		err := reservationService.CreateReservation(mockInput)
		assert.NoError(t, err)
	})

	t.Run("create reservation error", func(t *testing.T) {
		mockReservationRepo := new(MockReservationRepo)
		mockBookRepo := new(MockBookRepo)
		reservationService := NewReservationService(mockReservationRepo, mockBookRepo)

		mockReservationRepo.On("CreateReservation", mock.Anything).Return(errors.New("create reservation failed"))

		// Используем другие идентификаторы для второго теста
		mockInput := &CreateReservationInput{
			UserID:       2,                            // Другой пользователь
			BookID:       3,                            // Другая книга
			Status:       "pending",                    // Другой статус
			DateOfIssue:  issueDate.Add(1 * time.Hour), // Другое время
			DateOfReturn: returnDate.Add(1 * time.Hour),
		}

		err := reservationService.CreateReservation(mockInput)
		assert.Error(t, err)
		assert.Equal(t, "create reservation failed", err.Error())
	})
}
