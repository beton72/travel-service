package auth

type RegisterInput struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
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
