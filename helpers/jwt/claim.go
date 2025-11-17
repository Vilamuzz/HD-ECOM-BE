package jwt_helpers

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     uint8  `json:"role"`
	jwt.RegisteredClaims
}
