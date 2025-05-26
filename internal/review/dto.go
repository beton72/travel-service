package review

import "time"

type CreateReviewInput struct {
	HotelID uint   `json:"hotel_id" binding:"required"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Text    string `json:"text" binding:"required,max=1000"`
}

type ReviewResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Rating    int       `json:"rating"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type ReviewStats struct {
	AverageRating      float64     `json:"average_rating"`
	TotalReviews       int         `json:"total_reviews"`
	RatingDistribution map[int]int `json:"rating_distribution"`
}
