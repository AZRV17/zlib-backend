package service

import (
	"fmt"
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
	"strings"
)

type ReservationService struct {
	repository repository.ReservationRepo
	bookRepo   repository.BookRepo
}

func NewReservationService(repo repository.ReservationRepo, bookRepo repository.BookRepo) *ReservationService {
	return &ReservationService{repository: repo, bookRepo: bookRepo}
}

func (r ReservationService) GetReservationByID(id uint) (*domain.Reservation, error) {
	return r.repository.GetReservationByID(id)
}

func (r ReservationService) GetReservations() ([]*domain.Reservation, error) {
	return r.repository.GetReservations()
}

func (r ReservationService) CreateReservation(reservationInput *CreateReservationInput) error {
	reservation := &domain.Reservation{
		UserID:       reservationInput.UserID,
		BookID:       reservationInput.BookID,
		Status:       reservationInput.Status,
		DateOfReturn: reservationInput.DateOfReturn,
	}

	return r.repository.CreateReservation(reservation)
}

func (r ReservationService) UpdateReservation(reservationInput *UpdateReservationInput) error {
	reservation, err := r.repository.GetReservationByID(reservationInput.ID)
	if err != nil {
		return err
	}

	reservation.UserID = reservationInput.UserID
	reservation.BookID = reservationInput.BookID
	reservation.Status = reservationInput.Status
	reservation.DateOfReturn = reservationInput.DateOfReturn

	return r.repository.UpdateReservation(reservation)
}

func (r ReservationService) DeleteReservation(id uint) error {
	return r.repository.DeleteReservation(id)
}

func (r ReservationService) GetReservationsByUserID(id uint) ([]*domain.Reservation, error) {
	reservations, err := r.repository.GetUserReservations(id)
	if err != nil {
		return nil, err
	}

	for _, reservation := range reservations {
		if reservation.Book.Picture == "" || strings.HasPrefix(reservation.Book.Picture, "http") {
			continue
		}
		reservation.Book.Picture = "http://localhost:8080/" + reservation.Book.Picture
	}

	return reservations, nil
}

func (r ReservationService) UpdateReservationStatus(id uint, status string) error {
	reservation, err := r.repository.GetReservationByID(id)
	if err != nil {
		return fmt.Errorf("не удалось найти бронирование: %w", err)
	}

	reservation.Status = status
	if err := r.repository.UpdateReservation(reservation); err != nil {
		return fmt.Errorf("ошибка обновления статуса бронирования: %w", err)
	}

	if status == "returned" {
		if err := r.makeUniqueCodeAvailable(reservation.UniqueCodeID); err != nil {
			return err
		}
	}

	return nil
}

func (r ReservationService) makeUniqueCodeAvailable(codeID uint) error {
	code, err := r.bookRepo.GetUniqueCodeByID(codeID)
	if err != nil {
		return fmt.Errorf("ошибка получения уникального кода: %w", err)
	}

	code.IsAvailable = true

	if err := r.bookRepo.UpdateUniqueCode(code); err != nil {
		return fmt.Errorf("ошибка обновления уникального кода: %w", err)
	}

	return nil
}

func (r ReservationService) ExportReservationsToCSV() ([]byte, error) {
	return r.repository.ExportReservationsToCSV()
}
