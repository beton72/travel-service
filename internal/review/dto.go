package review

type CreateReviewInput struct {
	HotelID uint   `json:"hotel_id" binding:"required"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Text    string `json:"text" binding:"required"`
}
