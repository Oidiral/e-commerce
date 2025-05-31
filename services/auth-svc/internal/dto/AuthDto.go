package dto

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignUpResponse struct {
	Id        string `json:"id"`
	Email     string `json:"email" validate:"required,email"`
	Status    int16  `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
