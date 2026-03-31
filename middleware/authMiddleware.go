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

func ApiTokenAuthorization(c *gin.Context) {
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
			c.Abort()
			return
		}

		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Bad Request"))
		c.Abort()
		return
	}

	if !token.Valid {
		helper.SendErrorPayload(c, http.StatusUnauthorized, fmt.Errorf("Invalid Token"))
		c.Abort()
		return
	}

	c.Next()
}
