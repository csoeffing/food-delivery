package routes

import (
	controller "crunchgarage/restaurant-food-delivery/controllers"
	"crunchgarage/restaurant-food-delivery/middleware"

	"github.com/gin-gonic/gin"
)

func MenuRouter(router *gin.Engine) {
	router.GET("/api/menus", controller.GetMenus)
	router.GET("/api/menu/:id", controller.GetMenu)

	router.POST("/api/menus", middleware.ApiTokenAuthorization, controller.CreateMenu)
	router.PATCH("/api/menu/:id", middleware.ApiTokenAuthorization, controller.UpdateMenu)
}
