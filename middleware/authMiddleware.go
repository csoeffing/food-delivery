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
		func(token *jwt.Token) (interface{}, error) {
			return jwtkey, nil
		})
	//	json.NewEncoder(w).Encode(token.Claims)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			//w.WriteHeader(http.StatusUnauthorized)
			//json.NewEncoder(w).Encode("Invalid Token")
			helper.SendErrorPayload(c, http.StatusUnauthorized, fmt.Errorf("Invalid Token"))
			c.Abort()
			return
		}
		//w.WriteHeader(http.StatusBadRequest)
		//json.NewEncoder(w).Encode("Bad Request")
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Bad Request"))
		c.Abort()
		return
	}

	if !token.Valid {
		//w.WriteHeader(http.StatusUnauthorized)
		//json.NewEncoder(w).Encode("Invalid Token")
		helper.SendErrorPayload(c, http.StatusUnauthorized, fmt.Errorf("Invalid Token"))
		c.Abort()
		return
	}

	//endpoint(w, r)
	c.Next()
}
