package repository

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/gorm"
)

type ReservationRepository struct {
	DB *gorm.DB
}

func NewReservationRepository(db *gorm.DB) *ReservationRepository {
	return &ReservationRepository{DB: db}
}

func (r ReservationRepository) GetReservationByID(id uint) (*domain.Reservation, error) {
	var reservation domain.Reservation

	if err := r.DB.First(&reservation, id).Error; err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (r ReservationRepository) GetReservations() ([]*domain.Reservation, error) {
	var reservations []*domain.Reservation

	if err := r.DB.Find(&reservations).Error; err != nil {
		return nil, err
	}

	return reservations, nil
}

func (r ReservationRepository) CreateReservation(reservation *domain.Reservation) error {
	if err := r.DB.Create(reservation).Error; err != nil {
		return err
	}

	return nil
}

func (r ReservationRepository) UpdateReservation(reservation *domain.Reservation) error {
	if err := r.DB.Save(reservation).Error; err != nil {
		return err
	}

	return nil
}

func (r ReservationRepository) DeleteReservation(id uint) error {
	if err := r.DB.Delete(&domain.Reservation{}, id).Error; err != nil {
		return err
	}

	return nil
}
