package main

import (
	"net/http"
	"travel-service/internal/auth"
	"travel-service/internal/db"
	"travel-service/internal/hotel"
	"travel-service/internal/middleware"
	"travel-service/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.LoadEnv()
	db.InitPostgres()

	authService := auth.NewService()
	authHandler := auth.NewHandler(authService)

	protected := r.Group("/")
	protected.Use(middleware.AuthRequired())

	hotelService := hotel.NewService(db.DB)
	hotelHandler := hotel.NewHandler(hotelService)

	protected.GET("/me", authHandler.GetMe)
	protected.PATCH("/me", authHandler.UpdateMe)
	protected.POST("/hotels", hotelHandler.CreateHotel)
	protected.POST("/hotel-admins", hotelHandler.AddAdminToHotel)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.GET("/hotels", hotelHandler.GetHotels)
	r.Run(":8080")
}
