package routes

import (
	controller "crunchgarage/restaurant-food-delivery/controllers"
	"crunchgarage/restaurant-food-delivery/middleware"

	"github.com/gin-gonic/gin"
)

func FoodRouter(router *gin.Engine) {
	router.GET("/api/foods", controller.GetFoods)
	router.GET("/api/food/:id", controller.GetFood)
	router.PATCH("/api/food/:id", middleware.ApiTokenAuthorization, controller.UpdateFood)
	router.POST("/api/foods", middleware.ApiTokenAuthorization, controller.CreateFood)
}
