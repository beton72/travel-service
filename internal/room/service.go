package room

import (
	"errors"
	"travel-service/models"

	"gorm.io/gorm"
)

type Service interface {
	CreateRoom(input CreateRoomInput, hotelID, userID uint) error
	UpdateRoom(roomID uint, userID uint, input UpdateRoomInput) error
	GetRoomByID(id uint) (*models.Room, error)
	GetRandomReview(hotelID uint) (models.Review, error)
}

type service struct {
	db *gorm.DB
}

func NewService(database *gorm.DB) Service {
	return &service{db: database}
}

func (s *service) CreateRoom(input CreateRoomInput, hotelID uint, userID uint) error {
	var adminHotel models.AdminHotel
	err := s.db.
		Where("hotel_id = ? AND user_id = ?", hotelID, userID).
		First(&adminHotel).Error
	if err != nil {
		return errors.New("not authorized to add rooms to this hotel")
	}

	room := models.Room{
		HotelID:     hotelID,
		Type:        input.Type,
		Description: input.Description,
		Price:       input.Price,
		Capacity:    input.Capacity,
		Amenities:   input.Amenities,
		PhotoURLs:   input.PhotoURLs,
	}
	return s.db.Create(&room).Error
}

func (s *service) UpdateRoom(roomID uint, userID uint, input UpdateRoomInput) error {
	var room models.Room
	if err := s.db.First(&room, roomID).Error; err != nil {
		return err
	}

	var admin models.AdminHotel
	if err := s.db.Where("user_id = ? AND hotel_id = ?", userID, room.HotelID).First(&admin).Error; err != nil {
		return errors.New("access denied")
	}

	if input.Type != nil {
		room.Type = *input.Type
	}
	if input.Description != nil {
		room.Description = *input.Description
	}
	if input.Price != nil {
		room.Price = *input.Price
	}
	if input.Capacity != nil {
		room.Capacity = *input.Capacity
	}
	if input.Amenities != nil {
		room.Amenities = *input.Amenities
	}
	if input.PhotoURLs != nil {
		room.PhotoURLs = *input.PhotoURLs
	}
	return s.db.Save(&room).Error
}

func (s *service) GetRoomByID(id uint) (*models.Room, error) {
	var room models.Room
	if err := s.db.Preload("Hotel").First(&room, id).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (s *service) GetRandomReview(hotelID uint) (models.Review, error) {
	var review models.Review
	err := s.db.Where("hotel_id = ?", hotelID).
		Order("RANDOM()").
		Limit(1).
		Take(&review).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Review{}, nil
	}
	return review, err
}
