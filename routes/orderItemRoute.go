package routes

import (
	controller "crunchgarage/restaurant-food-delivery/controllers"

	"github.com/gin-gonic/gin"
)

func OrderItemRouter(router *gin.Engine) {
	router.GET("/api/orderItems", controller.GetOrderItems)
	router.GET("/api/orderItems/restaurant/{id}", controller.GetRestaurantOrderItems)
}
