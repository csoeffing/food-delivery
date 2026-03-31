package routes

import (
	controller "crunchgarage/restaurant-food-delivery/controllers"
	"crunchgarage/restaurant-food-delivery/middleware"

	"github.com/gin-gonic/gin"
)

func InvoiceRouter(router *gin.Engine) {
	router.POST("/api/invoice/create", middleware.ApiTokenAuthorization, controller.CreateInvoice)
}
