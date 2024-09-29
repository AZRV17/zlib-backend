package service

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
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
