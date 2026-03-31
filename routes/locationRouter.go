package routes

import (
	controller "crunchgarage/restaurant-food-delivery/controllers"

	"github.com/gin-gonic/gin"
)

func LocationRouter(router *gin.Engine) {
	router.GET("/api/location", controller.GetLocations)
}
