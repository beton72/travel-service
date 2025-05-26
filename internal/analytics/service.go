package analytics

import (
	"fmt"
	"time"
	"travel-service/models"

	"gorm.io/gorm"
)

type Service interface {
	GetDashboardStats() (*DashboardStats, error)
	GetHotelStats(req HotelReportRequest) ([]HotelStats, error)
	GetUserStats(req UserReportRequest) ([]UserStats, error)
	GetBookingStats(req BookingReportRequest) ([]BookingStats, error)
	GetRevenueStats(req RevenueReportRequest) ([]RevenueStats, error)
	GetReviewStats(hotelID *uint, dateRange *DateRangeFilter) ([]ReviewStats, error)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}

func (s *service) GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	// Общие счетчики
	if err := s.db.Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&models.Hotel{}).Count(&stats.TotalHotels).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&models.Booking{}).Count(&stats.TotalBookings).Error; err != nil {
		return nil, err
	}

	// Общая выручка
	var totalRevenue float64
	if err := s.db.Model(&models.Payment{}).Where("status = ?", "completed").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue).Error; err != nil {
		return nil, err
	}
	stats.TotalRevenue = totalRevenue

	// Средний рейтинг
	var avgRating float64
	if err := s.db.Model(&models.Review{}).
		Select("COALESCE(AVG(rating), 0)").Scan(&avgRating).Error; err != nil {
		return nil, err
	}
	stats.AverageRating = avgRating

	// Рост за месяц (упрощенно)
	lastMonth := time.Now().AddDate(0, -1, 0)
	var lastMonthBookings int64
	var currentMonthBookings int64

	s.db.Model(&models.Booking{}).Where("created_at >= ? AND created_at < ?",
		lastMonth, time.Now().AddDate(0, -1, 0).AddDate(0, 1, 0)).Count(&lastMonthBookings)
	s.db.Model(&models.Booking{}).Where("created_at >= ?",
		time.Now().AddDate(0, -1, 0).AddDate(0, 1, 0)).Count(&currentMonthBookings)

	if lastMonthBookings > 0 {
		stats.MonthlyGrowth = float64(currentMonthBookings-lastMonthBookings) / float64(lastMonthBookings) * 100
	}

	// Топ отели
	topHotels, err := s.GetHotelStats(HotelReportRequest{})
	if err == nil && len(topHotels) > 5 {
		stats.TopHotels = topHotels[:5]
	} else if err == nil {
		stats.TopHotels = topHotels
	}

	// Недавние бронирования
	recentBookings, err := s.GetBookingStats(BookingReportRequest{})
	if err == nil && len(recentBookings) > 10 {
		stats.RecentBookings = recentBookings[:10]
	} else if err == nil {
		stats.RecentBookings = recentBookings
	}

	return stats, nil
}

func (s *service) GetHotelStats(req HotelReportRequest) ([]HotelStats, error) {
	var stats []HotelStats

	query := s.db.Table("hotels h").
		Select(`
			h.id as hotel_id,
			h.name as hotel_name,
			COUNT(DISTINCT r.id) as total_rooms,
			COUNT(DISTINCT b.id) as total_bookings,
			COALESCE(SUM(p.amount), 0) as revenue,
			COALESCE(AVG(rv.rating), 0) as average_rating,
			COUNT(DISTINCT rv.id) as reviews_count
		`).
		Joins("LEFT JOIN rooms r ON h.id = r.hotel_id").
		Joins("LEFT JOIN bookings b ON r.id = b.room_id").
		Joins("LEFT JOIN payments p ON b.id = p.booking_id AND p.status = 'completed'").
		Joins("LEFT JOIN reviews rv ON h.id = rv.hotel_id").
		Group("h.id, h.name")

	if req.HotelID != 0 {
		query = query.Where("h.id = ?", req.HotelID)
	}

	if req.DateRange != nil {
		query = query.Where("b.created_at BETWEEN ? AND ?",
			req.DateRange.StartDate, req.DateRange.EndDate)
	}

	if err := query.Scan(&stats).Error; err != nil {
		return nil, err
	}

	// Вычисляем occupancy rate отдельно
	for i := range stats {
		stats[i].OccupancyRate = s.calculateOccupancyRate(stats[i].HotelID, req.DateRange)
	}

	return stats, nil
}

func (s *service) GetUserStats(req UserReportRequest) ([]UserStats, error) {
	var stats []UserStats

	query := s.db.Table("users u").
		Select(`
			u.id as user_id,
			u.first_name,
			u.last_name,
			u.user_email as email,
			u.role,
			COUNT(b.id) as bookings_count,
			COALESCE(SUM(p.amount), 0) as total_spent,
			MAX(b.created_at) as last_booking,
			u.created_at as registered_at
		`).
		Joins("LEFT JOIN bookings b ON u.id = b.user_id").
		Joins("LEFT JOIN payments p ON b.id = p.booking_id AND p.status = 'completed'").
		Group("u.id")

	if req.Role != "" && req.Role != "all" {
		query = query.Where("u.role = ?", req.Role)
	}

	if req.DateRange != nil {
		query = query.Where("u.created_at BETWEEN ? AND ?",
			req.DateRange.StartDate, req.DateRange.EndDate)
	}

	return stats, query.Scan(&stats).Error
}

func (s *service) GetBookingStats(req BookingReportRequest) ([]BookingStats, error) {
	var stats []BookingStats

	query := s.db.Table("bookings b").
		Select(`
			b.id as booking_id,
			CONCAT(u.first_name, ' ', u.last_name) as user_name,
			u.user_email as user_email,
			h.name as hotel_name,
			r.type as room_type,
			b.start_date,
			b.end_date,
			b.guest_count,
			b.status,
			r.price * EXTRACT(DAY FROM (b.end_date - b.start_date)) as total_price,
			b.paid,
			b.created_at
		`).
		Joins("JOIN users u ON b.user_id = u.id").
		Joins("JOIN rooms r ON b.room_id = r.id").
		Joins("JOIN hotels h ON r.hotel_id = h.id").
		Order("b.created_at DESC")

	if req.HotelID != nil {
		query = query.Where("h.id = ?", *req.HotelID)
	}

	if req.UserID != nil {
		query = query.Where("u.id = ?", *req.UserID)
	}

	if req.Status != "" && req.Status != "all" {
		query = query.Where("b.status = ?", req.Status)
	}

	if req.DateRange != nil {
		query = query.Where("b.created_at BETWEEN ? AND ?",
			req.DateRange.StartDate, req.DateRange.EndDate)
	}

	return stats, query.Scan(&stats).Error
}

func (s *service) GetRevenueStats(req RevenueReportRequest) ([]RevenueStats, error) {
	var stats []RevenueStats

	var dateFormat string
	var groupBy string

	switch req.Period {
	case "daily":
		dateFormat = "DATE(p.created_at)"
		groupBy = "DATE(p.created_at)"
	case "weekly":
		dateFormat = "DATE_TRUNC('week', p.created_at)"
		groupBy = "DATE_TRUNC('week', p.created_at)"
	case "monthly":
		dateFormat = "DATE_TRUNC('month', p.created_at)"
		groupBy = "DATE_TRUNC('month', p.created_at)"
	default:
		dateFormat = "DATE(p.created_at)"
		groupBy = "DATE(p.created_at)"
	}

	query := s.db.Table("payments p").
		Select(fmt.Sprintf(`
			%s as date,
			'%s' as period,
			h.id as hotel_id,
			h.name as hotel_name,
			SUM(p.amount) as revenue,
			COUNT(p.id) as bookings_count,
			AVG(p.amount) as average_price
		`, dateFormat, req.Period)).
		Joins("JOIN bookings b ON p.booking_id = b.id").
		Joins("JOIN rooms r ON b.room_id = r.id").
		Joins("JOIN hotels h ON r.hotel_id = h.id").
		Where("p.status = 'completed'").
		Group(groupBy + ", h.id, h.name").
		Order("date DESC")

	if req.HotelID != nil {
		query = query.Where("h.id = ?", *req.HotelID)
	}

	if req.DateRange != nil {
		query = query.Where("p.created_at BETWEEN ? AND ?",
			req.DateRange.StartDate, req.DateRange.EndDate)
	}

	return stats, query.Scan(&stats).Error
}

func (s *service) GetReviewStats(hotelID *uint, dateRange *DateRangeFilter) ([]ReviewStats, error) {
	var stats []ReviewStats

	query := s.db.Table("reviews r").
		Select(`
			r.id as review_id,
			CONCAT(u.first_name, ' ', u.last_name) as user_name,
			h.name as hotel_name,
			r.rating,
			r.comment,
			r.created_at
		`).
		Joins("JOIN users u ON r.user_id = u.id").
		Joins("JOIN hotels h ON r.hotel_id = h.id").
		Order("r.created_at DESC")

	if hotelID != nil {
		query = query.Where("r.hotel_id = ?", *hotelID)
	}

	if dateRange != nil {
		query = query.Where("r.created_at BETWEEN ? AND ?",
			dateRange.StartDate, dateRange.EndDate)
	}

	return stats, query.Scan(&stats).Error
}

func (s *service) calculateOccupancyRate(hotelID uint, dateRange *DateRangeFilter) float64 {
	// Упрощенный расчет загруженности
	var totalRooms int64
	var bookedRooms int64

	s.db.Model(&models.Room{}).Where("hotel_id = ?", hotelID).Count(&totalRooms)

	query := s.db.Table("bookings b").
		Joins("JOIN rooms r ON b.room_id = r.id").
		Where("r.hotel_id = ? AND b.status IN ('confirmed', 'completed')", hotelID)

	if dateRange != nil {
		query = query.Where("b.start_date <= ? AND b.end_date >= ?",
			dateRange.EndDate, dateRange.StartDate)
	}

	query.Count(&bookedRooms)

	if totalRooms == 0 {
		return 0
	}

	return float64(bookedRooms) / float64(totalRooms) * 100
}
