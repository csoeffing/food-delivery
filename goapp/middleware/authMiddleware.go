package middleware

import (
	helper "crunchgarage/restaurant-food-delivery/helpers"
	"crunchgarage/restaurant-food-delivery/models"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtkey = []byte(os.Getenv("JWT_KEY"))

func GetRequestingUserId(c *gin.Context) (int, error) {
	tokenString := c.GetHeader("Authorization")

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (any, error) {
			return jwtkey, nil
		},
	)

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			helper.SendErrorPayload(c, http.StatusUnauthorized, fmt.Errorf("Invalid Token"))
			return -1, err
		}

		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Bad Request"))
		return -1, err
	}

	if !token.Valid {
		helper.SendErrorPayload(c, http.StatusUnauthorized, fmt.Errorf("Invalid Token"))
		return -1, err
	}

	//fmt.Printf("API User: %d\n", claims.User_id)

	return claims.UserId, nil
}

func ApiTokenAuthorization(c *gin.Context) {
	_, err := GetRequestingUserId(c)

	if err != nil {
		c.Abort()
	}

	c.Next()
}
