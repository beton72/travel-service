package hotel

type CreateHotelInput struct {
	Name          string   `json:"name" binding:"required"`
	Address       string   `json:"address" binding:"required"`
	INN           string   `json:"inn" binding:"required"`
	Phone         string   `json:"phone" binding:"required"`
	Region        string   `json:"region" binding:"required"`
	PhotoURLs     []string `json:"photo_urls"`
	InitialRating float64  `json:"initial_rating"`
}

type AddAdminInput struct {
	UserID  uint `json:"user_id" binding:"required"`
	HotelID uint `json:"hotel_id" binding:"required"`
}

type FilterHotelsInput struct {
	MinPrice float64 `json:"min_price" binding:"min=0"`
	MaxPrice float64 `json:"max_price" binding:"required,gtfield=MinPrice"`
	CheckIn  string  `json:"check_in,omitempty"`  // Опционально
	CheckOut string  `json:"check_out,omitempty"` // Опционально
}

type SearchRoomsInput struct {
	Region   string   `json:"region" binding:"required"`
	CheckIn  string   `json:"check_in" binding:"required"`  // Формат: "2006-01-02"
	CheckOut string   `json:"check_out" binding:"required"` // Формат: "2006-01-02"
	MinPrice *float64 `json:"min_price,omitempty"`
	MaxPrice *float64 `json:"max_price,omitempty"`
}
