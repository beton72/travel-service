package main

import (
	"net/http"
	"travel-service/internal/analytics"
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

	analyticsService := analytics.NewService(db.DB)
	excelService := analytics.NewExcelService(analyticsService)
	analyticsHandler := analytics.NewHandler(analyticsService, excelService)

	db.SeedTestData(db.DB)

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
	protected.GET("/analytics/hotels/:id/export", analyticsHandler.ExportExcel)

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
	r.Run(":8080")
}
