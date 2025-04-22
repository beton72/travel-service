package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"travel-service/models"
)

var DB *gorm.DB

func InitPostgres() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("POSTGRE ERR ", err)
	}

	log.Println("POSTGRE CONNECT")
	DB = db

	err = db.AutoMigrate(
		&models.User{},
		&models.Hotel{},
		&models.Room{},
		&models.Booking{},
		&models.Review{},
		&models.Agency{},
		&models.AdminHotel{},
		&models.Payment{},
	)
	if err != nil {
		log.Fatal("Ошибка миграции:", err)
	}
}
