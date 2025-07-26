package dto

type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type CreateUserRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type LoginResponse struct {
	Message string `json:"message"`
}

type VerifyTokenResponse struct {
	Token string `json:"token"`
	User  interface{} `json:"user"`
}