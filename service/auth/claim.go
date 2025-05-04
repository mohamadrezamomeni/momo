package auth

import "github.com/golang-jwt/jwt/v5"

type Claim struct {
	UserID    string `json:"user_id"`
	IsAdmin   bool   `json:"is_admin"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Username  string `json:"username"`
}

type EncryptedClaim struct {
	Data string
	jwt.RegisteredClaims
}
