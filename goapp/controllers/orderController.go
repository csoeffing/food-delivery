package controller

import (
	"crunchgarage/restaurant-food-delivery/database"
	helper "crunchgarage/restaurant-food-delivery/helpers"
	"crunchgarage/restaurant-food-delivery/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	var order models.Order

	err := c.ShouldBindJSON(&order)
	if err != nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, err)
		return
	}

	if order.ProfileID == 0 {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("UserID field is required"))
		return
	}

	if order.Delivery_address == "" {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Delivery address field is required"))
		return
	}

	if order.OrderItem == nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Order Item(s) field is required"))
		return
	}

	// get sum of order items
	total_price_sum := 0.0
	for _, orderItem := range order.OrderItem {
		total_price_sum += orderItem.Unit_price * float64(orderItem.Quantity)
	}

	// get estimated delivery charge from calaculate delivery charge based on location ai
	delivery_charge_estimate := 102.0

	// get total amount of order items and deleivery charge
	total_amount_sum := total_price_sum + delivery_charge_estimate

	order_ := models.Order{
		ProfileID:        order.ProfileID,
		Delivery_address: order.Delivery_address,
		Order_status:     "PENDING",
		Order_Date:       time.Now().Format(time.RFC3339),
		Total_price:      total_price_sum,
		Delivery_charge:  delivery_charge_estimate,
		Total_amount:     total_amount_sum,
	}

	createdOrder := database.DB.Create(&order_)

	err = createdOrder.Error

	if err != nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, err)
		return
	}

	for _, food_item := range order.OrderItem {
		order_item_ := models.OrderItem{
			Quantity:     food_item.Quantity,
			Unit_price:   food_item.Unit_price,
			OrderID:      int(order_.ID),
			FoodID:       food_item.FoodID,
			RestaurantID: food_item.RestaurantID,
		}

		database.DB.Create(&order_item_)
	}

	var order_item []models.OrderItem

	database.DB.Model(&order_).Related(&order_item)
	order_.OrderItem = order_item

	helper.SendDataPayload(c, order_, true)
}

func GetOrder(c *gin.Context) {
	orderIdStr := c.Param("id")
	id, _ := strconv.Atoi(orderIdStr)

	var order models.Order
	var orderItem []models.OrderItem
	var invoice models.Invoice

	database.DB.First(&order, id)
	database.DB.Model(&order).Related(&orderItem)
	database.DB.Model(&order).Related(&invoice)

	if order.ID == 0 {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Order not found"))
		return
	}

	var orderItemHolder []map[string]interface{}

	for i := range orderItem {

		var food models.Food

		database.DB.Model(&orderItem[i]).Related(&food)

		foodData := map[string]interface{}{
			"id":          food.ID,
			"name":        food.Name,
			"price":       food.Price,
			"image":       food.Food_image,
			"description": food.Description,
			"status":      food.Status,
		}

		orderItemData := map[string]interface{}{
			"id":          orderItem[i].ID,
			"quantity":    orderItem[i].Quantity,
			"unit_price":  orderItem[i].Unit_price,
			"food":        foodData,
			"total_price": orderItem[i].Unit_price * float64(orderItem[i].Quantity),
		}

		orderItemHolder = append(orderItemHolder, orderItemData)
	}

	// invoice interface
	invoiceData := map[string]interface{}{
		"id":             invoice.ID,
		"payment_date":   invoice.Payment_date,
		"payment_status": invoice.Payment_status,
		"payment_method": invoice.Payment_method,
		"payment_amount": invoice.Amount,
	}

	orderData := map[string]interface{}{
		"id":               order.ID,
		"CreatedAt":        order.CreatedAt,
		"UpdatedAt":        order.UpdatedAt,
		"customer_id":      order.ProfileID,
		"order_items":      orderItemHolder,
		"delivery_address": order.Delivery_address,
		"order_status":     order.Order_status,
		"driver_id":        order.DriverID,
		"order_date":       order.Order_Date,
		"total_price":      order.Total_price,
		"delivery_charge":  order.Delivery_charge,
		"total_amount":     order.Total_amount,
		"payment_details":  invoiceData,
	}

	helper.SendDataPayload(c, orderData, false)
}

func UpdateOrder(c *gin.Context) {
	orderIdStr := c.Param("id")
	id, _ := strconv.Atoi(orderIdStr)

	var order models.Order
	var dbOrder models.Order
	var orderItem []models.OrderItem

	database.DB.First(&dbOrder, id)
	database.DB.Model(&dbOrder).Related(&orderItem)

	if dbOrder.ID == 0 {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Order not found"))
		return
	}

	err := c.ShouldBindJSON(&order)
	if err != nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, err)
		return
	}

	if order.ProfileID == 0 {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("customer id field is required"))
		return
	}

	if order.Delivery_address == "" {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Delivery address field is required"))
		return
	}

	if order.OrderItem == nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Order Item(s) field is required"))
		return
	}

	// get sum of order items
	total_price_sum := 0.0
	for _, orderItem := range order.OrderItem {
		total_price_sum += orderItem.Unit_price * float64(orderItem.Quantity)
	}

	// get estimated delivery charge from calaculate delivery charge based on location ai
	delivery_charge_estimate := 102.0

	// get total amount of order items and deleivery charge
	total_amount_sum := total_price_sum + delivery_charge_estimate

	dbOrder.ProfileID = order.ProfileID
	dbOrder.Delivery_address = order.Delivery_address
	dbOrder.Order_status = order.Order_status
	dbOrder.Order_Date = time.Now().Format(time.RFC3339)
	dbOrder.Total_price = total_price_sum
	dbOrder.Delivery_charge = delivery_charge_estimate
	dbOrder.Total_amount = total_amount_sum

	updated_order := database.DB.Save(&dbOrder)

	err = updated_order.Error

	if err != nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, err)
		return
	}

	for i, _ := range order.OrderItem {
		UpdateOrderItemFunc(order.OrderItem[i], id)
	}

	var order_item []models.OrderItem

	database.DB.Model(&dbOrder).Related(&order_item)
	dbOrder.OrderItem = order_item

	// var orderItemHolder []map[string]interface{}

	// for i, _ := range order.OrderItem {

	// 	var food models.Food

	// 	database.DB.Model(&order.OrderItem[i]).Related(&food)

	// 	foodData := map[string]interface{}{
	// 		"id":          food.ID,
	// 		"name":        food.Name,
	// 		"price":       food.Price,
	// 		"image":       food.Food_image,
	// 		"description": food.Description,
	// 		"status":      food.Status,
	// 	}

	// 	orderItemData := map[string]interface{}{
	// 		"id":          order.OrderItem[i].ID,
	// 		"quantity":    order.OrderItem[i].Quantity,
	// 		"unit_price":  order.OrderItem[i].Unit_price,
	// 		"food":        foodData,
	// 		"total_price": order.OrderItem[i].Unit_price * float64(order.OrderItem[i].Quantity),
	// 	}

	// 	orderItemHolder = append(orderItemHolder, orderItemData)
	// }

	// orderData := map[string]interface{}{
	// 	"id":               dbOrder.ID,
	// 	"customer_id":      dbOrder.ProfileID,
	// 	"order_items":      orderItemHolder,
	// 	"delivery_address": dbOrder.Delivery_address,
	// 	"order_status":     dbOrder.Order_status,
	// 	"driver_id":        dbOrder.DriverID,
	// 	"order_date":       dbOrder.Order_Date,
	// 	"total_price":      dbOrder.Total_price,
	// 	"delivery_charge":  dbOrder.Delivery_charge,
	// 	"total_amount":     dbOrder.Total_amount,
	// }

	helper.SendDataPayload(c, dbOrder, false)
}
