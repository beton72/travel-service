package booking

type CreateBookingInput struct {
	StartDate  string `json:"start_date" binding:"required"` // формат: "2025-05-01"
	EndDate    string `json:"end_date" binding:"required"`
	GuestCount int    `json:"guest_count" binding:"required"`
	Comment    string `json:"comment"`
}