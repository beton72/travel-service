package hotel

import (
	"errors"
	"fmt"
	"time"
	"travel-service/models"

	"gorm.io/gorm"
)

type Service interface {
	CreateHotel(input CreateHotelInput, userID uint) error
	GetHotels() ([]models.Hotel, error)
	AddAdminToHotel(input AddAdminInput) error
	GetHotelByID(id string) (models.Hotel, error) // Новый метод
	IncrementViews(hotelID string) error
	UpdateRating(hotelID uint, newRating float64) error
	FilterHotelsByPriceRange(input FilterHotelsInput) ([]models.Hotel, error)
	SearchAvailableRooms(input SearchRoomsInput) ([]models.RoomWithHotel, error)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}

func (s *service) CreateHotel(input CreateHotelInput, userID uint) error {
	// Создаём отель
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

	// Проверка, существует ли уже такая запись админа отеля
	var existing models.AdminHotel
	err := s.db.Where("user_id = ? AND hotel_id = ?", userID, hotel.ID).First(&existing).Error
	if err == nil {
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// Добавляем админа, если он ещё не добавлен
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

func (s *service) GetHotelByID(id string) (models.Hotel, error) {
	var hotel models.Hotel
	err := s.db.Preload("Rooms").First(&hotel, id).Error
	return hotel, err
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

func (s *service) IncrementViews(hotelID string) error {
	return s.db.Model(&models.Hotel{}).
		Where("id = ?", hotelID).
		Update("view_count", gorm.Expr("view_count + 1")).
		Error
}

func (s *service) UpdateRating(hotelID uint, newRating float64) error {
	return s.db.Model(&models.Hotel{}).
		Where("id = ?", hotelID).
		Updates(map[string]interface{}{
			"total_rating": gorm.Expr("total_rating + ?", newRating),
			"review_count": gorm.Expr("review_count + 1"),
		}).Error
}

func (s *service) SearchAvailableRooms(input SearchRoomsInput) ([]models.RoomWithHotel, error) {
	// Парсим даты
	checkIn, err := time.Parse("2006-01-02", input.CheckIn)
	if err != nil {
		return nil, errors.New("неверный формат даты заезда")
	}

	checkOut, err := time.Parse("2006-01-02", input.CheckOut)
	if err != nil {
		return nil, errors.New("неверный формат даты выезда")
	}

	// Базовый запрос
	query := s.db.Model(&models.Room{}).
		Select("rooms.*, hotels.name as hotel_name, hotels.address, hotels.region").
		Joins("JOIN hotels ON hotels.id = rooms.hotel_id").
		Where("hotels.region = ?", input.Region).
		Where("NOT EXISTS (SELECT 1 FROM bookings WHERE bookings.room_id = rooms.id AND ? < bookings.end_date AND ? > bookings.start_date)",
			checkOut, checkIn)

	// Фильтр по цене (если указан)
	if input.MinPrice != nil {
		query = query.Where("rooms.price >= ?", *input.MinPrice)
	}
	if input.MaxPrice != nil {
		query = query.Where("rooms.price <= ?", *input.MaxPrice)
	}

	var results []models.RoomWithHotel
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (s *service) FilterHotelsByPriceRange(input FilterHotelsInput) ([]models.Hotel, error) {
	// Валидация ценового диапазона
	if input.MaxPrice <= input.MinPrice {
		return nil, errors.New("максимальная цена должна быть больше минимальной")
	}

	query := s.db.Model(&models.Hotel{}).
		Select("hotels.*, AVG(rooms.price) as avg_price").
		Joins("JOIN rooms ON rooms.hotel_id = hotels.id").
		Group("hotels.id").
		Having("AVG(rooms.price) BETWEEN ? AND ?", input.MinPrice, input.MaxPrice)

	// Фильтр по датам (если указаны)
	if input.CheckIn != "" && input.CheckOut != "" {
		checkIn, err := time.Parse("2006-01-02", input.CheckIn)
		if err != nil {
			return nil, errors.New("неверный формат даты заезда")
		}

		checkOut, err := time.Parse("2006-01-02", input.CheckOut)
		if err != nil {
			return nil, errors.New("неверный формат даты выезда")
		}

		query = query.Where(`
            NOT EXISTS (
                SELECT 1 FROM bookings 
                JOIN rooms ON bookings.room_id = rooms.id
                WHERE rooms.hotel_id = hotels.id
                AND bookings.start_date < ? 
                AND bookings.end_date > ?
            )`, checkOut, checkIn)
	}

	var hotels []models.Hotel
	if err := query.
		Preload("Rooms").
		Order("avg_price ASC"). // Сортировка по цене
		Find(&hotels).Error; err != nil {
		return nil, err
	}

	return hotels, nil
}
