package main

import (
	"log"
	"net/http"
	"os"
	"time"
	"travel-service/internal/auth"
	"travel-service/internal/booking"
	"travel-service/internal/db"
	"travel-service/internal/hotel"
	"travel-service/internal/middleware"
	"travel-service/internal/payment"
	"travel-service/internal/review"
	"travel-service/internal/room"
	"travel-service/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.LoadEnv()
	db.InitPostgres()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	roomService := room.NewService(db.DB)
	roomHandler := room.NewHandler(roomService)

	bookingService := booking.NewService(db.DB)
	bookingHandler := booking.NewHandler(bookingService)

	authService := auth.NewService()
	authHandler := auth.NewHandler(authService)

	protected := r.Group("/")
	protected.Use(middleware.AuthRequired())

	hotelService := hotel.NewService(db.DB)
	hotelHandler := hotel.NewHandler(hotelService)

	reviewService := review.NewService(db.DB, hotelService)
	reviewHandler := review.NewHandler(reviewService)

	paymentService := payment.NewService(db.DB)
	paymentHandler := payment.NewHandler(paymentService)

	protected.GET("/me", authHandler.GetMe)
	protected.PATCH("/me", authHandler.UpdateMe)
	protected.GET("/me/bookings", bookingHandler.GetUserBookings)
	protected.POST("/hotels", hotelHandler.CreateHotel)
	protected.POST("/hotel-admins", hotelHandler.AddAdminToHotel)
	protected.POST("/hotels/:id/rooms", roomHandler.CreateRoom)
	protected.PATCH("/rooms/:id", roomHandler.UpdateRoom)
	protected.POST("/rooms/:id/book", bookingHandler.CreateBooking)
	protected.DELETE("/bookings/:id/cancel", bookingHandler.CancelBooking)
	protected.POST("/bookings/:id/pay", paymentHandler.PayBooking)
	protected.GET("/me/payments", paymentHandler.GetUserPayments)
	protected.POST("/reviews", reviewHandler.CreateReview)
	protected.GET("/me/hotels", hotelHandler.GetMyHotels)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/rooms/search", hotelHandler.SearchAvailableRooms)
	r.GET("/hotels", hotelHandler.GetHotels)
	r.GET("/hotels/:id", hotelHandler.GetHotelByID)
	r.GET("/rooms/:id", roomHandler.GetRoom)

	r.GET("/hotels/:id/reviews", reviewHandler.GetHotelReviews)
	r.GET("/hotels/:id/reviews/stats", reviewHandler.GetReviewStats)
	r.POST("/hotels/filter-by-price", hotelHandler.FilterHotelsByPrice)
	seedUsers()
	sqlContent, err := os.ReadFile("final_seed_data_from_models.sql")
	if err != nil {
		log.Fatalf("не удалось прочитать final_seed_data_from_models.sql: %v", err)
	}

	if err := db.DB.Exec(string(sqlContent)).Error; err != nil {
		log.Fatalf("ошибка выполнения SQL сидов: %v", err)
	}
	r.Run(":8080")
}

func seedUsers() {
	authService := auth.NewService()

	users := []auth.RegisterInput{
		{
			FirstName:      "Ирина",
			LastName:       "Козлова",
			Patronymic:     "Васильевна",
			Email:          "user@example.com",
			Password:       "qwerty123",
			Phone:          "89991234567",
			BirthDate:      time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			Role:           "client",
			Citizenship:    "Россия",
			PassportNumber: "1234567890",
		},
		{
			FirstName:      "Алексей",
			LastName:       "Смирнов",
			Patronymic:     "Юрьевич",
			Email:          "admin@example.com",
			Password:       "admin123",
			Phone:          "89999887766",
			BirthDate:      time.Date(1985, 6, 15, 0, 0, 0, 0, time.UTC),
			Role:           "admin",
			Citizenship:    "Россия",
			PassportNumber: "9876543210",
		},
	}

	for _, u := range users {
		_, err := authService.Register(u)
		if err != nil {
			log.Printf("Ошибка при создании пользователя %s: %v\n", u.Email, err)
		} else {
			log.Printf("Пользователь создан: email=%s, пароль=%s, роль=%s\n", u.Email, u.Password, u.Role)
		}
	}
}
