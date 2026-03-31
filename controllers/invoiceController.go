package controller

import (
	"crunchgarage/restaurant-food-delivery/database"
	helper "crunchgarage/restaurant-food-delivery/helpers"
	"crunchgarage/restaurant-food-delivery/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateInvoice(c *gin.Context) {
	var invoice models.Invoice

	err := c.ShouldBindJSON(&invoice)
	if err != nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, err)
		return
	}

	if invoice.OrderID == 0 {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Order id is required"))
		return
	}

	if invoice.UserID == 0 {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Customer id is required"))
		return
	}

	invoice_ := models.Invoice{
		UserID:         invoice.UserID,
		Amount:         invoice.Amount,
		OrderID:        invoice.OrderID,
		Payment_date:   time.Now().Format(time.RFC3339),
		Payment_method: "CARD",
		Payment_status: "PAID",
	}

	createdInvoice := database.DB.Create(&invoice_)
	err = createdInvoice.Error

	if err != nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, err)
		return
	}

	helper.SendDataPayload(c, createdInvoice.Value, true)
}
