package delivery

import (
	"github.com/AZRV17/zlib-backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) initReservationRoutes(r *gin.Engine) {
	reservations := r.Group("/reservations")
	{
		reservations.Use(h.AuthMiddleware).GET("/cookie", h.getUserReservations)
		reservations.Use(h.AuthMiddleware).GET("/:id", h.getReservationByID)
		reservations.Use(h.AuthMiddleware, h.LibrarianMiddleware).GET("/", h.getAllReservations)
		reservations.Use(h.AuthMiddleware, h.LibrarianMiddleware).PATCH("/:id", h.updateReservationStatus)
	}
}

func (h *Handler) getUserReservations(c *gin.Context) {
	userID, err := c.Cookie("id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if userIDInt == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	reservations, err := h.service.ReservationServ.GetReservationsByUserID(uint(userIDInt)) //nolint:gosec
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reservations)
}

func (h *Handler) getAllReservations(c *gin.Context) {
	reservations, err := h.service.ReservationServ.GetReservations()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reservations)
}

func (h *Handler) getReservationByID(c *gin.Context) {
	reservationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservation, err := h.service.ReservationServ.GetReservationByID(uint(reservationID)) //nolint:gosec
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

type updateReservationStatusInput struct {
	Status string `json:"status"`
}

func (h *Handler) updateReservationStatus(c *gin.Context) {
	var input updateReservationStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservation, err := h.service.ReservationServ.GetReservationByID(uint(reservationID)) //nolint:gosec
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservation.Status = input.Status
	err = h.service.ReservationServ.UpdateReservation(
		&service.UpdateReservationInput{
			ID:           reservation.ID,
			UserID:       reservation.UserID,
			BookID:       reservation.BookID,
			Status:       reservation.Status,
			DateOfReturn: reservation.DateOfReturn,
			DateOfIssue:  reservation.DateOfIssue,
		},
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
