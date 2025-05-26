package analytics

import "time"

// DateRangeFilter фильтр по диапазону дат
type DateRangeFilter struct {
	StartDate time.Time `json:"start_date" form:"start_date"`
	EndDate   time.Time `json:"end_date" form:"end_date"`
}

// HotelReportRequest запрос на отчет по отелю
type HotelReportRequest struct {
	HotelID   uint              `json:"hotel_id" form:"hotel_id"`
	DateRange *DateRangeFilter  `json:"date_range,omitempty"`
	Format    string            `json:"format" form:"format"` // "excel", "json"
}

// UserReportRequest запрос на отчет по пользователям
type UserReportRequest struct {
	Role      string           `json:"role" form:"role"` // "client", "admin", "all"
	DateRange *DateRangeFilter `json:"date_range,omitempty"`
	Format    string           `json:"format" form:"format"`
}

// BookingReportRequest запрос на отчет по бронированиям
type BookingReportRequest struct {
	HotelID   *uint            `json:"hotel_id,omitempty" form:"hotel_id"`
	UserID    *uint            `json:"user_id,omitempty" form:"user_id"`
	Status    string           `json:"status" form:"status"` // "new", "confirmed", "cancelled", "all"
	DateRange *DateRangeFilter `json:"date_range,omitempty"`
	Format    string           `json:"format" form:"format"`
}

// RevenueReportRequest запрос на отчет по доходам
type RevenueReportRequest struct {
	HotelID   *uint            `json:"hotel_id,omitempty" form:"hotel_id"`
	Period    string           `json:"period" form:"period"` // "daily", "weekly", "monthly"
	DateRange *DateRangeFilter `json:"date_range,omitempty"`
	Format    string           `json:"format" form:"format"`
}

// HotelStats статистика по отелю
type HotelStats struct {
	HotelID        uint    `json:"hotel_id"`
	HotelName      string  `json:"hotel_name"`
	TotalRooms     int     `json:"total_rooms"`
	TotalBookings  int64   `json:"total_bookings"`
	Revenue        float64 `json:"revenue"`
	AverageRating  float64 `json:"average_rating"`
	OccupancyRate  float64 `json:"occupancy_rate"`
	ReviewsCount   int64   `json:"reviews_count"`
}

// UserStats статистика по пользователям
type UserStats struct {
	UserID       uint    `json:"user_id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Email        string  `json:"email"`
	Role         string  `json:"role"`
	BookingsCount int64  `json:"bookings_count"`
	TotalSpent   float64 `json:"total_spent"`
	LastBooking  *time.Time `json:"last_booking,omitempty"`
	RegistredAt  time.Time `json:"registered_at"`
}

// BookingStats статистика по бронированиям
type BookingStats struct {
	BookingID    uint      `json:"booking_id"`
	UserName     string    `json:"user_name"`
	UserEmail    string    `json:"user_email"`
	HotelName    string    `json:"hotel_name"`
	RoomType     string    `json:"room_type"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	GuestCount   int       `json:"guest_count"`
	Status       string    `json:"status"`
	TotalPrice   float64   `json:"total_price"`
	Paid         bool      `json:"paid"`
	CreatedAt    time.Time `json:"created_at"`
}

// RevenueStats статистика по доходам
type RevenueStats struct {
	Date         time.Time `json:"date"`
	Period       string    `json:"period"`
	HotelID      *uint     `json:"hotel_id,omitempty"`
	HotelName    string    `json:"hotel_name,omitempty"`
	Revenue      float64   `json:"revenue"`
	BookingsCount int64    `json:"bookings_count"`
	AveragePrice float64   `json:"average_price"`
}

// ReviewStats статистика по отзывам
type ReviewStats struct {
	ReviewID    uint      `json:"review_id"`
	UserName    string    `json:"user_name"`
	HotelName   string    `json:"hotel_name"`
	Rating      int       `json:"rating"`
	Comment     string    `json:"comment"`
	CreatedAt   time.Time `json:"created_at"`
}

// DashboardStats общая статистика для дашборда
type DashboardStats struct {
	TotalUsers     int64   `json:"total_users"`
	TotalHotels    int64   `json:"total_hotels"`
	TotalBookings  int64   `json:"total_bookings"`
	TotalRevenue   float64 `json:"total_revenue"`
	AverageRating  float64 `json:"average_rating"`
	MonthlyGrowth  float64 `json:"monthly_growth"`
	TopHotels      []HotelStats `json:"top_hotels"`
	RecentBookings []BookingStats `json:"recent_bookings"`
}