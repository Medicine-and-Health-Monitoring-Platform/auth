package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	PhoneNumber  string    `json:"phone_number"`
	Role         string    `json:"role"` // 'customer', 'provider', 'admin'
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    int64     `json:"deleted_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AccessToken struct {
	Token string `json:"access_token"`
}

type RefreshToken struct {
	Token string `json:"refresh_token"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenDetails struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
	Expiry string `json:"expiry"`
}
