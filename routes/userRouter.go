package routes

import (
	controller "crunchgarage/restaurant-food-delivery/controllers"
	"crunchgarage/restaurant-food-delivery/middleware"

	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine) {
	router.POST("/api/user/signup", controller.SignUp)
	router.POST("/api/user/login", controller.Login)

	router.GET("/api/user/:id", middleware.ApiTokenAuthorization, controller.GetUser)
	router.PATCH("/api/user/:id", middleware.ApiTokenAuthorization, controller.UpdateUser)
}
