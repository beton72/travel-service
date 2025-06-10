package models

import (
	"time"
)

type User struct {
	ID             uint      `gorm:"primaryKey"`
	FirstName      string    `gorm:"size:50;not null"`
	LastName       string    `gorm:"size:50;not null"`
	Patronymic     string    `gorm:"size:50"`
	UserEmail      string    `gorm:"size:100;not null;unique"`
	PasswordHash   string    `gorm:"type:text;not null"`
	UserPhone      string    `gorm:"size:15"`
	BirthDate      time.Time `gorm:"not null"`
	Role           string    `gorm:"size:20;not null"`
	Citizenship    string    `gorm:"size:50"`
	HasChildren    bool      `gorm:"default:false"`
	ChildrenInfo   []string  `gorm:"type:jsonb;serializer:json"`
	PassportNumber string    `gorm:"size:15"`
	PhotoURLs      []string  `gorm:"type:jsonb;serializer:json"`

	Bookings    []Booking    `gorm:"foreignKey:UserID"`
	Reviews     []Review     `gorm:"foreignKey:UserID"`
	Agency      *Agency      `gorm:"foreignKey:UserID" json:"agency,omitempty"`
	AdminHotels []AdminHotel `gorm:"foreignKey:UserID"`
}

type Hotel struct {
	ID        uint     `gorm:"primaryKey"`
	Name      string   `gorm:"size:100;not null"`
	Address   string   `gorm:"size:200;not null"`
	INN       string   `gorm:"size:12;not null"`
	Phone     string   `gorm:"size:15;not null"`
	Region    string   `gorm:"size:100;not null"`
	PhotoURLs []string `gorm:"type:jsonb;serializer:json"`
	Amenities []string `gorm:"type:jsonb;serializer:json"`

	Rooms   []Room       `gorm:"foreignKey:HotelID"`
	Admins  []AdminHotel `gorm:"foreignKey:HotelID"`
	Reviews []Review     `gorm:"foreignKey:HotelID"`

	ViewCount   uint    `gorm:"default:0"`
	TotalRating float64 `gorm:"default:0"`
	ReviewCount uint    `gorm:"default:0"`
	Revenue     float64 `gorm:"default:0"`
}

type Room struct {
	ID          uint     `gorm:"primaryKey"`
	HotelID     uint     `gorm:"not null"`
	Type        string   `gorm:"size:50;not null"`
	Description string   `gorm:"type:text"`
	Price       float64  `gorm:"type:numeric"`
	Capacity    int      `gorm:"not null"`
	PhotoURLs   []string `gorm:"type:jsonb;serializer:json"`
	Amenities   []string `gorm:"type:jsonb;serializer:json"`

	Hotel       *Hotel    `gorm:"foreignKey:HotelID"`
	Bookings    []Booking `gorm:"foreignKey:RoomID"`
	StatusToday string    `json:"status_today" gorm:"-"`
}

type Booking struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"not null"`
	RoomID     uint      `gorm:"not null"`
	StartDate  time.Time `gorm:"not null"`
	EndDate    time.Time `gorm:"not null"`
	GuestCount int       `gorm:"not null"`
	Status     string    `gorm:"size:20;default:'new'"`
	Comment    string    `gorm:"type:text"`
	Paid       bool      `gorm:"default:false"`

	User    *User   `gorm:"foreignKey:UserID"`
	Room    *Room   `gorm:"foreignKey:RoomID"`
	Payment Payment `gorm:"foreignKey:BookingID"`
}

type Payment struct {
	ID            uint      `gorm:"primaryKey"`
	BookingID     uint      `gorm:"not null"`
	Amount        float64   `gorm:"type:numeric"`
	Status        string    `gorm:"size:20;not null"`
	PaymentMethod string    `gorm:"size:50"`
	TransactionID string    `gorm:"size:100"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`

	Booking *Booking `gorm:"foreignKey:BookingID"`
}

type Review struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	HotelID   uint      `gorm:"not null"`
	Rating    int       `gorm:"not null"`
	Text      string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	PhotoURLs []string  `gorm:"type:jsonb;serializer:json"`
	Amenities []string  `gorm:"type:jsonb;serializer:json"`

	User  *User `gorm:"foreignKey:UserID"`
	Hotel Hotel `gorm:"foreignKey:HotelID"`
}

type Agency struct {
	ID            uint     `gorm:"primaryKey"`
	Name          string   `gorm:"size:100;not null"`
	LicenseNumber string   `gorm:"size:20;not null"`
	RegionScope   string   `gorm:"size:100;not null"`
	ContactPerson string   `gorm:"size:100;not null"`
	Email         string   `gorm:"size:100;not null"`
	PhotoURLs     []string `gorm:"type:jsonb;serializer:json"`
	UserID        uint     `gorm:"not null"`

	User *User `gorm:"foreignKey:UserID"`
}

type AdminHotel struct {
	ID        uint     `gorm:"primaryKey"`
	UserID    uint     `gorm:"not null"`
	HotelID   uint     `gorm:"not null"`
	PhotoURLs []string `gorm:"type:jsonb;serializer:json"`
	Amenities []string `gorm:"type:jsonb;serializer:json"`

	User  *User  `gorm:"foreignKey:UserID"`
	Hotel *Hotel `gorm:"foreignKey:HotelID"`
}

type RoomWithHotel struct {
	// models.Room
	HotelName string `json:"hotel_name"`
	Address   string `json:"address"`
	Region    string `json:"region"`
}

type HotelAdmin struct {
	ID      uint
	UserID  uint
	HotelID uint
}

func (HotelAdmin) TableName() string {
	return "admin_hotels"
}
