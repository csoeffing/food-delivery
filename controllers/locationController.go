package controller

import (
	"crunchgarage/restaurant-food-delivery/database"
	helper "crunchgarage/restaurant-food-delivery/helpers"
	"crunchgarage/restaurant-food-delivery/models"

	"github.com/gin-gonic/gin"
)

func GetLocations(c *gin.Context) {

	var location []models.Location

	database.DB.Find(&location)

	helper.SendDataPayload(c, location, false)
}
