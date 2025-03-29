package service

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
	"strings"
)

type ReservationService struct {
	repository repository.ReservationRepo
}

func NewReservationService(repo repository.ReservationRepo) *ReservationService {
	return &ReservationService{repository: repo}
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
