package models

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserName string `json:"user_name"`
	jwt.StandardClaims
}
