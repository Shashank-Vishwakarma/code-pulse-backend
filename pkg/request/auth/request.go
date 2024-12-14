package request

type RegisterRequest struct {
	Name            string `json:"name" validate:"required,min=5,max=50"`
	Username        string `json:"username" validate:"required,min=3,max=20"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,max=20"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,max=20"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}
