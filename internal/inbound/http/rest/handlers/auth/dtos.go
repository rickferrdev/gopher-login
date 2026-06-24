package auth

import "time"

type RequestLoginDTO struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=8"`
}

type RequestRegisterDTO struct {
	Username string `json:"username" validate:"min=4,max=16"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=8"`
}

type ResponseLoginDTO struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

type ResponseRegisterDTO struct {
	ID       string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}
