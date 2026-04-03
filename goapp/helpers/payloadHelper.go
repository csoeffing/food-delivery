package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendErrorPayload(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{
		"success": false,
		"message": err.Error(),
	})
}

func SendDataPayload(c *gin.Context, data any, wasCreated bool) {
	if wasCreated {
		c.JSON(http.StatusCreated, data)
	} else {
		c.JSON(http.StatusOK, data)
	}
}
