package hotel

type CreateHotelInput struct {
	Name      string   `json:"name" binding:"required"`
	Address   string   `json:"address" binding:"required"`
	INN       string   `json:"inn" binding:"required"`
	Phone     string   `json:"phone" binding:"required"`
	Region    string   `json:"region" binding:"required"`
	PhotoURLs []string `json:"photo_urls"`
}

type AddAdminInput struct {
	UserID  uint `json:"user_id" binding:"required"`
	HotelID uint `json:"hotel_id" binding:"required"`
}