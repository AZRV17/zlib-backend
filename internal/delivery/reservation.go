package delivery

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/AZRV17/zlib-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initReservationRoutes(r *gin.Engine) {
	reservations := r.Group("/reservations")
	{
		reservations.Use(h.AuthMiddleware).GET("/my", h.getUserReservations)
		reservations.Use(h.AuthMiddleware).GET("/:id", h.getReservationByID)
		reservations.Use(h.AuthMiddleware, h.LibrarianMiddleware).GET("/", h.getAllReservations)
		reservations.Use(h.AuthMiddleware, h.LibrarianMiddleware).PATCH("/:id", h.updateReservationStatus)
		reservations.Use(h.AuthMiddleware, h.LibrarianMiddleware).PUT("/:id", h.updateReservationStatus)
		reservations.Use(h.AuthMiddleware, h.LibrarianMiddleware).GET("/export", h.exportReservationsToCSV)
	}
}

func (h *Handler) getUserReservations(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	reservations, err := h.service.ReservationServ.GetReservationsByUserID(userID)
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

	err = h.service.ReservationServ.UpdateReservationStatus(uint(reservationID), input.Status) //nolint:gosec
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createLogInput := &service.CreateLogInput{
		UserID:  userID,
		Action:  "Изменение статуса бронирования",
		Date:    time.Now(),
		Details: "Изменение статуса бронирования: " + input.Status,
	}

	err = h.service.LogServ.CreateLog(createLogInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reservation updated"})
}

func (h *Handler) exportReservationsToCSV(c *gin.Context) {
	reservationData, err := h.service.ReservationServ.ExportReservationsToCSV()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := "books.csv"

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")
	c.Header("Content-Length", fmt.Sprint(len(reservationData)))

	c.Data(http.StatusOK, "text/csv", reservationData)

	userID, err := getUserIDFromContext(c)
	if err != nil {
		log.Printf("Error getting user ID for logging: %v", err)
		return
	}

	createLogInput := &service.CreateLogInput{
		UserID:  userID,
		Action:  "Экспорт бронирований в CSV",
		Date:    time.Now(),
		Details: "Экспорт бронирований в CSV файл",
	}

	err = h.service.LogServ.CreateLog(createLogInput)
	if err != nil {
		log.Printf("Error creating log: %v", err)
	}
}
