package service

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var tokens []Token

const TokenName = "sw_token"
const Issuer = "Go-GraphQL-Group"
const SecretKey = "StarWars"

type Token struct {
	SW_TOKEN string `json:"sw_token"`
}

type jwtCustomClaims struct {
	jwt.StandardClaims

	Admin bool `json:"admin"`
}

func CreateToken(secretKey []byte, issuer string, isAdmin bool) (token Token, err error) {
	claims := &jwtCustomClaims{
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Add(time.Hour * 1).Unix()),
			Issuer:    issuer,
		},
		isAdmin,
	}

	tokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
	token = Token{
		tokenStr,
	}
	return
}

func ParseToken(tokenStr string, secretKey []byte) (claims jwt.Claims, err error) {
	var token *jwt.Token
	token, err = jwt.Parse(tokenStr, func(*jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	claims = token.Claims
	return
}
