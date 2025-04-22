package payment

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) PayBooking(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	bookingIDParam := c.Param("id")
	bookingID, err := strconv.ParseUint(bookingIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID брони"})
		return
	}

	err = h.service.ProcessMockPayment(uint(bookingID), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Оплата прошла успешно"})
}

func (h *Handler) GetUserPayments(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	payments, err := h.service.GetUserPayments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить платежи"})
		return
	}

	c.JSON(http.StatusOK, payments)
}
