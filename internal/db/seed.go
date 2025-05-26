package db

import (
	"errors"
	"log"
	"time"

	"travel-service/models"

	"gorm.io/gorm"
)

func SeedTestData(db *gorm.DB) error {
	var admin models.User
	if err := db.Where("user_email = ?", "admin@example.com").First(&admin).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		// hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		admin = models.User{
			FirstName: "Админ",
			LastName:  "Тестовый",
			UserEmail: "admin@example.com",
			// PasswordHash: string(hashedPassword),
			PasswordHash: "123456",
			Role:         "admin",
		}

		if err := db.Create(&admin).Error; err != nil {
			return err
		}
	}

	var hotel models.Hotel
	if err := db.Where("name = ?", "Тестовый отель").First(&hotel).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		hotel = models.Hotel{
			Name:        "Тестовый отель",
			Address:     "г. Москва, ул. Примерная, д.1",
			Region:      "Москва",
			PhotoURLs:   []string{"https://placehold.co/600x400"},
			TotalRating: 0,
			ReviewCount: 0,
			ViewCount:   0,
		}
		if err := db.Create(&hotel).Error; err != nil {
			return err
		}
	}

	var adminHotel models.AdminHotel
	if err := db.Where("user_id = ? AND hotel_id = ?", admin.ID, hotel.ID).First(&adminHotel).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		adminHotel = models.AdminHotel{
			UserID:  admin.ID,
			HotelID: hotel.ID,
		}
		if err := db.Create(&adminHotel).Error; err != nil {
			return err
		}
	}

	for i := 1; i <= 3; i++ {
		room := models.Room{
			HotelID:     hotel.ID,
			Type:        "Стандарт",
			Description: "Номер с удобствами",
			Price:       float64(3000 + i*500),
			Capacity:    2,
		}
		if err := db.Where("hotel_id = ? AND price = ?", hotel.ID, room.Price).FirstOrCreate(&room).Error; err != nil {
			return err
		}

		for j := 0; j < 3; j++ {
			startDate := time.Date(2025, time.January+time.Month(j), 5, 0, 0, 0, 0, time.UTC)
			endDate := startDate.AddDate(0, 0, 2)

			booking := models.Booking{
				UserID:     admin.ID,
				RoomID:     room.ID,
				StartDate:  startDate,
				EndDate:    endDate,
				Status:     "paid",
				Paid:       true,
				GuestCount: 2,
			}
			if err := db.Create(&booking).Error; err != nil {
				return err
			}

			payment := models.Payment{
				BookingID: booking.ID,
				Amount:    room.Price * 2,
				Status:    "success",
			}
			if err := db.Create(&payment).Error; err != nil {
				return err
			}

			review := models.Review{
				UserID:    admin.ID,
				HotelID:   hotel.ID,
				Rating:    4 + j%2,
				Text:      "Хороший отель!",
				PhotoURLs: []string{},
			}
			if err := db.Create(&review).Error; err != nil {
				return err
			}
		}
	}

	log.Println("Тест данные")
	return nil
}
