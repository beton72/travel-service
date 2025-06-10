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
	// Получаем бронирование по ID
	var booking models.Booking
	if err := s.db.First(&booking, bookingID).Error; err != nil {
		return errors.New("бронирование не найдено")
	}

	// Проверка, принадлежит ли бронь пользователю
	if booking.UserID != userID {
		return errors.New("нельзя оплатить чужое бронирование")
	}

	// Проверка, не оплачена ли бронь ранее
	if booking.Paid {
		return errors.New("бронирование уже оплачено")
	}
	// Создаём запись об оплате (эмуляция оплаты)
	payment := models.Payment{
		BookingID:     bookingID,                                 // Связь с бронью
		Status:        "success",                                 // Статус оплаты
		PaymentMethod: "mock",                                    // Тип оплаты
		TransactionID: fmt.Sprintf("MOCK-%d", time.Now().Unix()), // Уникальный ID транзакции
	}
	// Сохраняем платёж в базу
	if err := s.db.Create(&payment).Error; err != nil {
		return err
	}

	// Обновляем статус бронирования
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
