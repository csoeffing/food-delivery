package helper

import (
	"crunchgarage/restaurant-food-delivery/config"
	"crunchgarage/restaurant-food-delivery/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtkey = []byte(config.EnvJwtKey())

// Genearte JWT token
func GenerateToken(principal models.User, duration time.Duration) (string, int64, error) {

	claims := &models.Claims{
		UserId:   int(principal.ID),
		UserName: principal.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)

	if err != nil {
		return "", 0, err
	}

	return tokenString, claims.ExpiresAt, nil
}
