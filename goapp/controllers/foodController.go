package controller

import (
	"crunchgarage/restaurant-food-delivery/config"
	"crunchgarage/restaurant-food-delivery/database"
	helper "crunchgarage/restaurant-food-delivery/helpers"
	"crunchgarage/restaurant-food-delivery/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var food_image = ""

func CreateFood(c *gin.Context) {
	var food models.Food

	// get formdata
	food_name := c.PostForm("name")
	food_price := c.PostForm("price")
	food_menu_id := c.PostForm("menu_id")
	food_restarant_id := c.PostForm("restarant_id")
	food_description := c.PostForm("description")

	file, err := c.FormFile("food_image")

	if food_name == "" {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Food name is required"))
		return
	}

	if food_price == "" {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Food price is required"))
		return
	}

	if food_menu_id == "" {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Menu id is required"))
		return
	}

	if food_restarant_id == "" {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Restaurant id is required"))
		return
	}

	if food_description == "" {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Food Description id is required"))
		return
	}

	if file != nil {
		avatarUrl, err := helper.SingleImageUpload(c, "food_image", config.EnvCloudFoodFolder(), "food")
		if err != nil {
			avatarUrl = ""
		}
		food_image = avatarUrl
	}

	food.Name = food_name
	food.Price, _ = strconv.ParseFloat(food_price, 64)
	food.Food_image = food_image
	food.MenuID, _ = strconv.Atoi(food_menu_id)
	food.RestaurantID, _ = strconv.Atoi(food_restarant_id)

	createdFood := database.DB.Create(&food)
	err = createdFood.Error

	if err != nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, err)
		return
	}

	food_image = ""

	helper.SendDataPayload(c, createdFood.Value, true)
}

func GetFoods(c *gin.Context) {
	var foods []models.Food

	database.DB.Find(&foods)

	helper.SendDataPayload(c, foods, false)
}

func GetFood(c *gin.Context) {
	foodIdStr := c.Param("id")
	id, _ := strconv.Atoi(foodIdStr)

	var food models.Food

	database.DB.First(&food, id)

	if food.ID == 0 {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("No food found"))
		return
	}

	helper.SendDataPayload(c, food, false)
}

func UpdateFood(c *gin.Context) {
	foodIdStr := c.Param("id")
	id, _ := strconv.Atoi(foodIdStr)

	var food models.Food
	var dbFood models.Food

	database.DB.First(&dbFood, id)

	if dbFood.ID == 0 {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Food not found"))
		return
	}

	err := c.ShouldBindJSON(&food)
	if err != nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, err)
		return
	}

	if food.Name != "" {
		dbFood.Name = food.Name
	}

	if strconv.FormatFloat(food.Price, 'E', -1, 32) != "" {
		dbFood.Price = food.Price
	}

	if strconv.Itoa(food.MenuID) != "" {
		dbFood.MenuID = food.MenuID
	}

	if strconv.Itoa(food.RestaurantID) != "" {
		dbFood.RestaurantID = food.RestaurantID
	}

	if food.Description != "" {
		dbFood.Description = food.Description
	}

	file, err := c.FormFile("food_image")

	if err != nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, err)
		return
	}

	if file != nil {
		avatarUrl, err := helper.SingleImageUpload(c, "food_image", config.EnvCloudFoodFolder(), "food")
		if err != nil {
			avatarUrl = dbFood.Food_image
		}
		food_image = avatarUrl
	}

	// update menu
	updatedFood := database.DB.Model(&dbFood).Updates(models.Food{
		Name:         dbFood.Name,
		Price:        dbFood.Price,
		Food_image:   food_image,
		MenuID:       dbFood.MenuID,
		RestaurantID: dbFood.RestaurantID,
		Description:  dbFood.Description,
	})

	err = updatedFood.Error

	if err != nil {
		helper.SendErrorPayload(c, http.StatusInternalServerError, err)
		return
	}

	food_image = ""

	helper.SendDataPayload(c, updatedFood.Value, true)
}
