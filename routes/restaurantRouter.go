package routes

import (
	controller "crunchgarage/restaurant-food-delivery/controllers"
	"crunchgarage/restaurant-food-delivery/middleware"

	"github.com/gin-gonic/gin"
)

func RestaurantRouter(router *gin.Engine) {
	// http://localhost:8134/api/restaurants
	router.GET("/api/restaurants", controller.GetRestaurants)

	// http://localhost:8134/api/restaurant/1
	router.GET("/api/restaurant/:id", controller.GetRestaurant)

	router.PATCH("/api/restaurant/:id", middleware.ApiTokenAuthorization, controller.UpdateRestaurant)

	router.POST("/api/restaurants/create", middleware.ApiTokenAuthorization, controller.CreateRestaurant)
}
