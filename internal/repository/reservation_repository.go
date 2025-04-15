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

	tx := r.DB.Begin()

	if err := tx.First(&reservation, id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &reservation, nil
}

func (r ReservationRepository) GetReservations() ([]*domain.Reservation, error) {
	var reservations []*domain.Reservation

	tx := r.DB.Begin()

	if err := tx.Preload("Book").Preload("UniqueCode").Preload("User").Find(&reservations).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return reservations, nil
}

func (r ReservationRepository) CreateReservation(reservation *domain.Reservation) error {
	tx := r.DB.Begin()

	if err := tx.Create(reservation).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r ReservationRepository) UpdateReservation(reservation *domain.Reservation) error {
	tx := r.DB.Begin()

	if err := tx.Where("id = ?", reservation.ID).Save(reservation).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r ReservationRepository) DeleteReservation(id uint) error {
	tx := r.DB.Begin()

	if err := tx.Delete(&domain.Reservation{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r ReservationRepository) CreateReservationWithTransactions(reservation *domain.Reservation, tx *gorm.DB) error {
	if err := tx.Model(&domain.Reservation{}).Create(reservation).Error; err != nil {
		return err
	}

	return nil
}

func (r ReservationRepository) GetUserReservations(id uint) ([]*domain.Reservation, error) {
	var reservations []*domain.Reservation

	tx := r.DB.Begin()

	if err := tx.Where(
		"user_id = ?",
		id,
	).Preload("Book").Preload("Book.Author").Preload("Book.Genre").Preload("Book.Publisher").
		Find(&reservations).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return reservations, nil
}
