package analytics

import (
	"time"

	"gorm.io/gorm"
)

type MonthlyBookings struct {
	Month string `json:"month"`
	Count int    `json:"count"`
}

type MonthlyRevenue struct {
	Month   string  `json:"month"`
	Revenue float64 `json:"revenue"`
}

type MonthlyOccupancy struct {
	Month          string  `json:"month"`
	OccupancyRate  float64 `json:"occupancy_rate"`
	BookedDays     int     `json:"booked_days"`
	TotalAvailable int     `json:"total_available"`
}

type MonthlyReviews struct {
	Month         string  `json:"month"`
	AverageRating float64 `json:"average_rating"`
	ReviewCount   int     `json:"review_count"`
}

type AverageCheck struct {
	Month        string  `json:"month"`
	AvgCheck     float64 `json:"avg_check"`
	TotalAmount  float64 `json:"total_amount"`
	BookingCount int     `json:"booking_count"`
}

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetMonthlyBookings(hotelID uint, startDate, endDate time.Time) ([]MonthlyBookings, error) {
	var results []MonthlyBookings
	query := `
	SELECT TO_CHAR(b.start_date, 'YYYY-MM') as month, COUNT(*) as count
	FROM bookings b
	JOIN rooms r ON b.room_id = r.id
	WHERE r.hotel_id = ? AND b.start_date BETWEEN ? AND ? AND b.status != 'cancelled'
	GROUP BY month ORDER BY month
`
	err := s.db.Raw(query, hotelID, startDate, endDate).Scan(&results).Error
	return results, err
}

func (s *Service) GetMonthlyRevenue(hotelID uint, startDate, endDate time.Time) ([]MonthlyRevenue, error) {
	var results []MonthlyRevenue
	query := `
        SELECT TO_CHAR(b.start_date, 'YYYY-MM') as month, COALESCE(SUM(p.amount), 0) as revenue
        FROM bookings b
        LEFT JOIN payments p ON b.id = p.booking_id AND p.status = 'success'
        WHERE b.hotel_id = ? AND b.start_date BETWEEN ? AND ? AND b.status != 'cancelled'
        GROUP BY month ORDER BY month
    `
	err := s.db.Raw(query, hotelID, startDate, endDate).Scan(&results).Error
	return results, err
}

func (s *Service) GetMonthlyOccupancy(hotelID uint, startDate, endDate time.Time) ([]MonthlyOccupancy, error) {
	var results []MonthlyOccupancy

	query := `
	WITH monthly_stats AS (
		SELECT 
			TO_CHAR(b.start_date, 'YYYY-MM') as month,
			SUM(EXTRACT(DAY FROM (b.end_date - b.start_date))) as booked_days,
			COUNT(DISTINCT r.id) * 30 as total_available
		FROM bookings b
		JOIN rooms r ON b.room_id = r.id
		WHERE b.hotel_id = ?
			AND b.start_date BETWEEN ? AND ?
			AND b.status != 'cancelled'
		GROUP BY TO_CHAR(b.start_date, 'YYYY-MM')
	)
	SELECT 
		month,
		CASE 
			WHEN total_available > 0 
			THEN ROUND((booked_days::decimal / total_available::decimal) * 100, 2)
			ELSE 0 
		END as occupancy_rate,
		booked_days,
		total_available
	FROM monthly_stats
	ORDER BY month;
	`

	err := s.db.Raw(query, hotelID, startDate, endDate).Scan(&results).Error
	return results, err
}

func (s *Service) GetMonthlyReviews(hotelID uint, startDate, endDate time.Time) ([]MonthlyReviews, error) {
	var results []MonthlyReviews

	query := `
		SELECT 
			TO_CHAR(created_at, 'YYYY-MM') as month,
			ROUND(AVG(rating::decimal), 2) as average_rating,
			COUNT(*) as review_count
		FROM reviews 
		WHERE hotel_id = ? 
			AND created_at BETWEEN ? AND ?
		GROUP BY month
		ORDER BY month
	`

	err := s.db.Raw(query, hotelID, startDate, endDate).Scan(&results).Error
	return results, err
}

func (s *Service) GetAverageCheck(hotelID uint, startDate, endDate time.Time) ([]AverageCheck, error) {
	var results []AverageCheck

	query := `
		SELECT 
			TO_CHAR(b.start_date, 'YYYY-MM') as month,
			ROUND(AVG(p.amount), 2) as avg_check,
			SUM(p.amount) as total_amount,
			COUNT(b.id) as booking_count
		FROM bookings b
		JOIN payments p ON b.id = p.booking_id AND p.status = 'success'
		WHERE b.hotel_id = ? 
			AND b.start_date BETWEEN ? AND ? 
			AND b.status != 'cancelled'
		GROUP BY month
		ORDER BY month
	`

	err := s.db.Raw(query, hotelID, startDate, endDate).Scan(&results).Error
	return results, err
}
