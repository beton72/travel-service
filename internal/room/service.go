package room

import (
	"encoding/json"
	"errors"
	"travel-service/models"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Service interface {
	CreateRoom(input CreateRoomInput, hotelID, userID uint) error
	UpdateRoom(roomID uint, userID uint, input UpdateRoomInput) error
	GetRoomByID(id uint) (*models.Room, error)
}

type service struct {
	db *gorm.DB
}

func NewService(database *gorm.DB) Service {
	return &service{db: database}
}

func (s *service) CreateRoom(input CreateRoomInput, hotelID uint, userID uint) error {
	// Проверка, принадлежит ли пользователь к администраторам указанного отеля
	var adminHotel models.AdminHotel
	err := s.db.
		Where("hotel_id = ? AND user_id = ?", hotelID, userID).
		First(&adminHotel).Error
	if err != nil {
		// Если не найдено соответствие — пользователь не имеет права добавлять номера
		return errors.New("not authorized to add rooms to this hotel")
	}
	// Сериализация массива ссылок на фото в формат JSON
	photoJSON, err := json.Marshal(input.PhotoURLs)
	if err != nil {
		return err
	}
	// Сериализация массива удобств в JSON
	amenitiesJSON, err := json.Marshal(input.Amenities)
	if err != nil {
		return err
	}
	// Создание структуры номера с привязкой к отелю
	room := models.Room{
		HotelID:     hotelID,
		Type:        input.Type,
		Description: input.Description,
		Price:       input.Price,
		Capacity:    input.Capacity,
		Amenities:   datatypes.JSON(amenitiesJSON),
		PhotoURLs:   datatypes.JSON(photoJSON),
	}
	// Сохранение новой записи о номере в базу данных
	return s.db.Create(&room).Error
}

func (s *service) UpdateRoom(roomID uint, userID uint, input UpdateRoomInput) error {
	// Получение записи о номере по ID
	var room models.Room
	if err := s.db.First(&room, roomID).Error; err != nil {
		return err
	}

	// Проверка, является ли пользователь администратором отеля, которому принадлежит номер
	var admin models.AdminHotel
	if err := s.db.Where("user_id = ? AND hotel_id = ?", userID, room.HotelID).First(&admin).Error; err != nil {
		return errors.New("access denied")
	}

	// Обновление данных номера только при наличии новых значений
	if input.Type != nil {
		room.Type = *input.Type // Тип номера
	}
	if input.Description != nil {
		room.Description = *input.Description // Описание номера
	}
	if input.Price != nil {
		room.Price = *input.Price // Цена за ночь
	}
	if input.Capacity != nil {
		room.Capacity = *input.Capacity // Вместимость
	}
	if input.Amenities != nil {
		// Преобразование удобств в JSON
		amenitiesJSON, err := json.Marshal(input.Amenities)
		if err != nil {
			return err
		}
		room.Amenities = datatypes.JSON(amenitiesJSON)
	}

	if input.PhotoURLs != nil {
		// Преобразование ссылок на фото в JSON
		photoJSON, err := json.Marshal(input.PhotoURLs)
		if err != nil {
			return err
		}
		room.PhotoURLs = datatypes.JSON(photoJSON)
	}
	// Сохранение изменений в базе данных
	return s.db.Save(&room).Error
}

func (s *service) GetRoomByID(id uint) (*models.Room, error) {
	var room models.Room
	if err := s.db.Preload("Hotel").First(&room, id).Error; err != nil {
		return nil, err
	}
	return &room, nil
}
