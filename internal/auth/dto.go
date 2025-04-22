package auth

import "time"

type RegisterInput struct {
	FirstName      string    `json:"first_name" binding:"required"`
	LastName       string    `json:"last_name" binding:"required"`
	Patronymic     string    `json:"patronymic"` // 👈 вот это добавь
	Email          string    `json:"email" binding:"required,email"`
	Password       string    `json:"password" binding:"required"`
	Phone          string    `json:"phone"`
	BirthDate      time.Time `json:"birth_date"`
	Role           string    `json:"role"`
	Citizenship    string    `json:"citizenship"`
	HasChildren    bool      `json:"has_children"`
	ChildrenInfo   []string  `json:"children_info"`
	PassportNumber string    `json:"passport_number"`
	PhotoURLs      []string  `json:"photo_urls"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserInput struct {
	FirstName      *string   `json:"first_name,omitempty"`
	LastName       *string   `json:"last_name,omitempty"`
	Patronymic     *string   `json:"patronymic,omitempty"`
	UserPhone      *string   `json:"user_phone,omitempty"`
	BirthDate      *string   `json:"birth_date,omitempty"`
	Citizenship    *string   `json:"citizenship,omitempty"`
	HasChildren    *bool     `json:"has_children,omitempty"`
	ChildrenInfo   *[]string `json:"children_info,omitempty"`
	PassportNumber *string   `json:"passport_number,omitempty"`
	PhotoURLs      *[]string `json:"photo_urls,omitempty"`
	Role           *string   `json:"role"`
}
