package jwt_helpers

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
