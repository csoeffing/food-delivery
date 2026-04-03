package controller

import (
	"crunchgarage/restaurant-food-delivery/database"
	helper "crunchgarage/restaurant-food-delivery/helpers"
	"crunchgarage/restaurant-food-delivery/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetOrderItems(c *gin.Context) {
	var orderItems []models.OrderItem

	order_items := database.DB.Find(&orderItems)
	err := order_items.Error

	if err != nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, err)
		return
	}

	helper.SendDataPayload(c, orderItems, false)
}

func UpdateOrderItemFunc(item models.OrderItem, order_id int) {

	id := item.ID

	var orderItem models.OrderItem

	database.DB.First(&orderItem, id)

	if orderItem.ID == 0 {

		return
	}

	orderItem.Quantity = item.Quantity
	orderItem.Unit_price = item.Unit_price
	orderItem.FoodID = item.FoodID
	orderItem.RestaurantID = item.RestaurantID
	orderItem.OrderID = order_id

	database.DB.Save(&orderItem)

}

func GetRestaurantOrderItems(c *gin.Context) {
	restaurantIdStr := c.Param("id")
	id, _ := strconv.Atoi(restaurantIdStr)

	var orderItems []models.OrderItem
	var orderItemHolder []map[string]interface{}

	database.DB.Where("restaurant_id = ?", id).Find(&orderItems)

	if len(orderItems) == 0 {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Order not found"))
		return
	}

	for i, _ := range orderItems {

		var food models.Food
		var order models.Order
		var invoice models.Invoice

		database.DB.Model(&orderItems[i]).Related(&food)
		database.DB.Model(&orderItems[i]).Related(&order)
		database.DB.Model(&order).Related(&invoice)

		/*food interface*/
		foodData := map[string]interface{}{
			"id":          food.ID,
			"name":        food.Name,
			"price":       food.Price,
			"image":       food.Food_image,
			"description": food.Description,
			"status":      food.Status,
		}

		/*order interface*/
		orderData := map[string]interface{}{
			"id":               order.ID,
			"order_date":       order.Order_Date,
			"delivery_address": order.Delivery_address,
			"status":           order.Order_status,
			"customer_id":      order.ProfileID,
		}

		/**invoice interface*/
		invoiceData := map[string]interface{}{
			"id":             invoice.ID,
			"payment_date":   invoice.Payment_date,
			"payment_status": invoice.Payment_status,
		}

		/*order item interface*/
		orderItemData := map[string]interface{}{
			"id":              orderItems[i].ID,
			"quantity":        orderItems[i].Quantity,
			"unit_price":      orderItems[i].Unit_price,
			"food_details":    foodData,
			"restaurant_id":   orderItems[i].RestaurantID,
			"order_details":   orderData,
			"total_price":     float64(orderItems[i].Quantity) * orderItems[i].Unit_price,
			"payment_details": invoiceData,
		}

		orderItemHolder = append(orderItemHolder, orderItemData)

	}

	helper.SendDataPayload(c, orderItemHolder, false)
}
