package handlers

type loginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type isAdminRequest struct {
	ID int64 `json:"id" validate:"required"`
}

type registerRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=30"`
	LastName    string `json:"last_name" validate:"required,min=2,max=30"`
	Email       string `json:"email" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required"`
}
