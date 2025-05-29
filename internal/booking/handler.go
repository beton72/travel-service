package booking

import (
	"net/http"
	"strconv"
	"travel-service/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateBooking(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var booking models.Booking
	roomIDParam := c.Param("id")
	roomIDUint, err := strconv.ParseUint(roomIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID номера"})
		return
	}
	roomID := uint(roomIDUint)

	var input CreateBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные", "details": err.Error()})
		return
	}

	if err := h.service.CreateBooking(roomID, userID, input); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Бронирование успешно создано",
		"booking_id": booking.ID,
	})

}

func (h *Handler) GetUserBookings(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	bookings, err := h.service.GetUserBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить бронирования"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func (h *Handler) CancelBooking(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	bookingIDParam := c.Param("id")
	bookingIDUint, err := strconv.ParseUint(bookingIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID брони"})
		return
	}
	bookingID := uint(bookingIDUint)

	if err := h.service.CancelBooking(bookingID, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Бронирование отменено"})
}
