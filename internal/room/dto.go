package room

import "gorm.io/datatypes"

type CreateRoomInput struct {
	Type        string         `json:"type" binding:"required"`
	Description string         `json:"description"`
	Price       float64        `json:"price" binding:"required"`
	Capacity    int            `json:"capacity" binding:"required"`
	Amenities   datatypes.JSON `gorm:"type:jsonb"`
	PhotoURLs   datatypes.JSON `gorm:"type:jsonb"`
}

type UpdateRoomInput struct {
	Type        *string   `json:"type,omitempty"`
	Description *string   `json:"description,omitempty"`
	Price       *float64  `json:"price,omitempty"`
	Capacity    *int      `json:"capacity,omitempty"`
	Amenities   *[]string `json:"amenities,omitempty"`
	PhotoURLs   *[]string `json:"photo_urls,omitempty"`
}

