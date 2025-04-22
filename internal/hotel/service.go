package hotel

import (
	"fmt"
	"travel-service/models"

	"gorm.io/gorm"
)

type Service interface {
	CreateHotel(input CreateHotelInput, userID uint) error
	GetHotels() ([]models.Hotel, error)
	AddAdminToHotel(input AddAdminInput) error
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}

func (s *service) CreateHotel(input CreateHotelInput, userID uint) error {
	// 1. Создаём отель
	hotel := models.Hotel{
		Name:      input.Name,
		Address:   input.Address,
		INN:       input.INN,
		Phone:     input.Phone,
		Region:    input.Region,
		PhotoURLs: input.PhotoURLs,
	}

	// Сохраняем отель
	if err := s.db.Create(&hotel).Error; err != nil {
		return err
	}

	// 2. Проверка: существует ли уже такая запись админа отеля
	var existing models.AdminHotel
	err := s.db.Where("user_id = ? AND hotel_id = ?", userID, hotel.ID).First(&existing).Error
	if err == nil {
		// Уже существует — не вставляем второй раз
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		// Какая-то другая ошибка при проверке
		return err
	}

	// 3. Добавляем админа, если он ещё не добавлен
	adminHotel := models.AdminHotel{
		UserID:    userID,
		HotelID:   hotel.ID,
		PhotoURLs: []string{},
	}

	return s.db.Create(&adminHotel).Error
}

func (s *service) GetHotels() ([]models.Hotel, error) {
	var hotels []models.Hotel
	if err := s.db.
		Preload("Admins").
		Preload("Admins.User").
		Find(&hotels).Error; err != nil {
		return nil, err
	}
	return hotels, nil
}

func (s *service) AddAdminToHotel(input AddAdminInput) error {
	var user models.User
	if err := s.db.First(&user, input.UserID).Error; err != nil {
		return fmt.Errorf("пользователь не найден")
	}

	var hotel models.Hotel
	if err := s.db.First(&hotel, input.HotelID).Error; err != nil {
		return fmt.Errorf("отель не найден")
	}

	admin := models.AdminHotel{
		UserID:  input.UserID,
		HotelID: input.HotelID,
	}

	if err := s.db.Create(&admin).Error; err != nil {
		return err
	}

	return nil
}
