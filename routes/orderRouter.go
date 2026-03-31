package routes

import (
	controller "crunchgarage/restaurant-food-delivery/controllers"
	"crunchgarage/restaurant-food-delivery/middleware"

	"github.com/gin-gonic/gin"
)

func OrderRouter(router *gin.Engine) {
	router.POST("/api/order/create", middleware.ApiTokenAuthorization, controller.CreateOrder)

	// http://localhost:8080/api/order/1
	router.GET("/api/order/:id", controller.GetOrder)

	router.PATCH("/api/order/{id}", middleware.ApiTokenAuthorization, controller.UpdateOrder)
}
