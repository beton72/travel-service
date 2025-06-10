package review

import (
	"errors"
	"travel-service/internal/hotel"
	"travel-service/models"

	"gorm.io/gorm"
)

type Service interface {
	CreateReview(userID uint, input CreateReviewInput) error
	GetHotelReviews(hotelID uint, page int, limit int) ([]models.Review, error)
	GetReviewStats(hotelID uint) (ReviewStats, error)
	GetRandomReview(hotelID uint) (models.Review, error)
}

type service struct {
	db           *gorm.DB
	hotelService hotel.Service
}

func NewService(db *gorm.DB, hotelService hotel.Service) Service {
	return &service{db: db, hotelService: hotelService}
}

func (s *service) CreateReview(userID uint, input CreateReviewInput) error {
	// Валидация рейтинга
	if input.Rating < 1 || input.Rating > 5 {
		return errors.New("рейтинг должен быть от 1 до 5")
	}

	// Проверка бронирований
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

	// Создание отзыва
	review := models.Review{
		UserID:  userID,
		HotelID: input.HotelID,
		Rating:  input.Rating,
		Text:    input.Text,
	}

	if err := s.db.Create(&review).Error; err != nil {
		return err
	}

	// Обновляем рейтинг отеля
	return s.hotelService.UpdateRating(input.HotelID, float64(input.Rating))
}

func (s *service) GetHotelReviews(hotelID uint, page int, limit int) ([]models.Review, error) {
	offset := (page - 1) * limit
	var reviews []models.Review

	err := s.db.
		Where("hotel_id = ?", hotelID).
		Offset(offset).
		Limit(limit).
		Find(&reviews).Error

	return reviews, err
}

func (s *service) GetReviewStats(hotelID uint) (ReviewStats, error) {
	var stats ReviewStats

	// Средний рейтинг
	err := s.db.Model(&models.Review{}).
		Select("AVG(rating) as average_rating, COUNT(*) as total_reviews").
		Where("hotel_id = ?", hotelID).
		Scan(&stats).Error

	if err != nil {
		return ReviewStats{}, err
	}

	// Распределение по оценкам
	var distribution []struct {
		Rating int
		Count  int
	}

	s.db.Model(&models.Review{}).
		Select("rating, COUNT(*) as count").
		Where("hotel_id = ?", hotelID).
		Group("rating").
		Scan(&distribution)

	stats.RatingDistribution = make(map[int]int)
	for _, d := range distribution {
		stats.RatingDistribution[d.Rating] = d.Count
	}

	return stats, nil
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
