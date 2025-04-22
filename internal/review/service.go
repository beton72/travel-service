package review

import (
	"errors"
	"travel-service/models"

	"gorm.io/gorm"
)

type Service interface {
	CreateReview(userID uint, input CreateReviewInput) error
	GetHotelReviews(hotelID uint) ([]models.Review, error)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}

func (s *service) CreateReview(userID uint, input CreateReviewInput) error {
	var count int64
	err := s.db.
		Model(&models.Booking{}).
		Joins("JOIN rooms ON bookings.room_id = rooms.id").
		Where("bookings.user_id = ? AND bookings.paid = true AND rooms.hotel_id = ?", userID, input.HotelID).
		Count(&count).Error

	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("отзыв можно оставить только после оплаты проживания")
	}

	review := models.Review{
		UserID:  userID,
		HotelID: input.HotelID,
		Rating:  input.Rating,
		Text:    input.Text,
	}

	return s.db.Create(&review).Error
}

func (s *service) GetHotelReviews(hotelID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := s.db.Where("hotel_id = ?", hotelID).Find(&reviews).Error
	return reviews, err
}
