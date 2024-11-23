package models

type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name" validate:"required,min=2,max=30"`
	LastName     string `json:"last_name" validate:"required,min=2,max=30"`
	Email        string `json:"email" validate:"required"`
	PhoneNumber  string `json:"phone_number" validate:"required"`
	HashPassword []byte `json:"hash_password"`
	RefreshToken string `json:"refresh_token"`

	IP string `json:"ip"`
}
