package models

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserId   int    `json:"userId"`
	UserName string `json:"userName"`
	jwt.StandardClaims
}

type Signings struct {
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	AccessTokenExpiration string `json:"access_token_expiration"`
}
