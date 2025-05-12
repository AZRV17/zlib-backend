package repository

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/gorm"
	"strconv"
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

func (r ReservationRepository) ExportReservationsToCSV() ([]byte, error) {
	reservations, err := r.GetReservations()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	headers := []string{"ID", "Пользователь", "Книга", "Код", "Статус", "Дата бронирования", "Дата возврата"}

	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	for _, reservation := range reservations {
		row := []string{
			strconv.FormatUint(uint64(reservation.ID), 10),
			reservation.User.Login,
			reservation.Book.Title,
			strconv.Itoa(reservation.UniqueCode.Code),
			reservation.Status,
			reservation.DateOfIssue.Format("2006-01-02"),
			reservation.DateOfReturn.Format("2006-01-02"),
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
