package domain

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	jwt.StandardClaims
	UserID    int64  `json:"user_id"`
	UserEmail string `json:"email"`
}
