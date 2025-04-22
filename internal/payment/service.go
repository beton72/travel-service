package payment

import (
	"errors"
	"fmt"
	"time"
	"travel-service/models"

	"gorm.io/gorm"
)

type Service interface {
	ProcessMockPayment(bookingID, userID uint) error
	GetUserPayments(userID uint) ([]models.Payment, error)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}

func (s *service) ProcessMockPayment(bookingID, userID uint) error {
	var booking models.Booking
	if err := s.db.First(&booking, bookingID).Error; err != nil {
		return errors.New("бронирование не найдено")
	}
	if booking.UserID != userID {
		return errors.New("нельзя оплатить чужое бронирование")
	}
	if booking.Paid {
		return errors.New("бронирование уже оплачено")
	}

	payment := models.Payment{
		BookingID:     bookingID,
		Amount:        1000, // ты можешь тут рассчитать по дням и цене номера
		Status:        "success",
		PaymentMethod: "mock",
		TransactionID: fmt.Sprintf("MOCK-%d", time.Now().Unix()),
		CreatedAt:     time.Now(),
	}

	if err := s.db.Create(&payment).Error; err != nil {
		return err
	}

	booking.Paid = true
	booking.Status = "paid"
	return s.db.Save(&booking).Error
}

func (s *service) GetUserPayments(userID uint) ([]models.Payment, error) {
	var payments []models.Payment

	err := s.db.
		Joins("JOIN bookings ON bookings.id = payments.booking_id").
		Where("bookings.user_id = ?", userID).
		Preload("Booking").
		Find(&payments).Error

	if err != nil {
		return nil, err
	}

	return payments, nil
}
