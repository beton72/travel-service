package booking

import (
	"errors"
	"time"
	"travel-service/models"

	"gorm.io/gorm"
)

type Service interface {
	CreateBooking(roomID uint, userID uint, input CreateBookingInput) error
	GetUserBookings(userID uint) ([]models.Booking, error)
	CancelBooking(bookingID uint, userID uint) error
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}

func (s *service) CreateBooking(roomID uint, userID uint, input CreateBookingInput) error {
	// Проверка существования номера в базе
	var room models.Room
	if err := s.db.First(&room, roomID).Error; err != nil {
		return errors.New("номер не найден")
	}
	// Парсинг и проверка дат
	start, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return errors.New("неверный формат даты начала")
	}
	end, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		return errors.New("неверный формат даты окончания")
	}
	if !start.Before(end) {
		return errors.New("дата начала должна быть раньше даты окончания")
	}

	// Проверка, есть ли брони, пересекающиеся с заданным периодом
	var count int64
	err = s.db.Model(&models.Booking{}).
		Where("room_id = ? AND start_date < ? AND end_date > ?", roomID, end, start).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("номер занят в выбранные даты")
	}
	// Создание новой записи бронирования
	booking := models.Booking{
		UserID:     userID,           // Идентификатор пользователя
		RoomID:     roomID,           // Идентификатор номера
		StartDate:  start,            // Дата начала
		EndDate:    end,              // Дата окончания
		GuestCount: input.GuestCount, // Количество гостей
		Comment:    input.Comment,    // Комментарий к брони
		Paid:       false,            // По умолчанию — не оплачено
	}
	// Сохранение брони в базе данных
	return s.db.Create(&booking).Error
}

func (s *service) GetUserBookings(userID uint) ([]models.Booking, error) {
	var bookings []models.Booking
	err := s.db.
		Preload("Room").
		Preload("Room.Hotel").
		Where("user_id = ?", userID).
		Order("start_date DESC").
		Find(&bookings).Error
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *service) CancelBooking(bookingID uint, userID uint) error {
	// Получаем запись бронирования по ID
	var booking models.Booking
	if err := s.db.First(&booking, bookingID).Error; err != nil {
		return errors.New("бронирование не найдено")
	}

	// Проверяем, принадлежит ли бронь текущему пользователю
	if booking.UserID != userID {
		return errors.New("вы не можете отменить чужую бронь")
	}

	// Устанавливаем статус "отменено"
	booking.Status = "cancelled"

	// Сохраняем изменения
	return s.db.Save(&booking).Error
}
