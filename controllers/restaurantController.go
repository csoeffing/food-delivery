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

var restaurant_image = ""

func CreateRestaurant(c *gin.Context) {
	var restaurant models.Restaurant

	err := c.ShouldBindJSON(&restaurant)
	if err != nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, err)
		return
	}

	var dbRestaurant models.Restaurant
	database.DB.Where("restaurant_name = ?", restaurant.Restaurant_name).First(&dbRestaurant)
	if dbRestaurant.ID != 0 {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Business name already exists"))
		return
	}

	restaurant.Registration_status = "PENDING"
	createdMenu := database.DB.Create(&restaurant)
	err = createdMenu.Error

	if err != nil {
		helper.SendErrorPayload(c, http.StatusInternalServerError, err)
		return
	}

	helper.SendDataPayload(c, createdMenu.Value, true)
}

func GetRestaurants(c *gin.Context) {
	var restaurants []models.Restaurant
	var restaurantsHolder []map[string]interface{}

	database.DB.Find(&restaurants)

	for i, _ := range restaurants {

		var profile models.Profile

		database.DB.Model(&restaurants[i]).Related(&profile)

		/*restaurant interface*/
		restaurantData := map[string]interface{}{
			"id":               restaurants[i].ID,
			"restaurant_image": restaurants[i].Restaurant_image,
			"restaurant_name":  restaurants[i].Restaurant_name,
			"phone_number":     restaurants[i].Phone_number,
			"address":          restaurants[i].Address,
			"location":         restaurants[i].LocationID,
			"owner":            profile,
		}

		restaurantsHolder = append(restaurantsHolder, restaurantData)

	}

	helper.SendDataPayload(c, restaurantsHolder, false)
}

func GetRestaurant(c *gin.Context) {
	restaurantIdStr := c.Param("id")
	id, _ := strconv.Atoi(restaurantIdStr)

	var restaurant models.Restaurant
	var profile models.Profile
	var location models.Location

	database.DB.First(&restaurant, id)
	database.DB.Model(&restaurant).Related(&profile)
	database.DB.Model(&restaurant).Related(&location)

	if restaurant.ID == 0 {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Restaurant not found"))
		return
	}

	restaurantData := map[string]interface{}{
		"id":               restaurant.ID,
		"CreatedAt":        restaurant.CreatedAt,
		"UpdatedAt":        restaurant.UpdatedAt,
		"restaurant_image": restaurant.Restaurant_image,
		"restaurant_name":  restaurant.Restaurant_name,
		"phone_number":     restaurant.Phone_number,
		"address":          restaurant.Address,
		"location":         location,
		"owner":            restaurant.ProfileID,
	}

	helper.SendDataPayload(c, restaurantData, false)
}

func UpdateRestaurant(c *gin.Context) {
	restaurantIdStr := c.Param("id")
	id, _ := strconv.Atoi(restaurantIdStr)

	var restaurant models.Restaurant
	var dbRestaurant models.Restaurant
	var profile models.Profile
	var location models.Location

	database.DB.First(&dbRestaurant, id)
	database.DB.Model(&dbRestaurant).Related(&profile)
	database.DB.Model(&dbRestaurant).Related(&location)

	if dbRestaurant.ID == 0 {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Restaurant not found"))
		return
	}

	err := c.ShouldBindJSON(&restaurant)
	if err != nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, err)
		return
	}

	if restaurant.Restaurant_name != "" {
		dbRestaurant.Restaurant_name = restaurant.Restaurant_name
	}

	if restaurant.Phone_number != "" {
		dbRestaurant.Phone_number = restaurant.Phone_number
	}

	if restaurant.Address != "" {
		dbRestaurant.Address = restaurant.Address
	}

	if strconv.Itoa(restaurant.LocationID) != "" {
		dbRestaurant.LocationID = restaurant.LocationID
	}

	file, err := c.FormFile("restaurant_image")

	if err != nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, err)
		return
	}

	if file != nil {
		avatarUrl, err := helper.SingleImageUpload(c, "restaurant_image", config.EnvCloudRestaurantFolder(), "restaurant")
		if err != nil {
			restaurant_image = dbRestaurant.Restaurant_image
		}
		restaurant_image = avatarUrl
	}

	database.DB.Model(&dbRestaurant).Updates(models.Restaurant{
		Restaurant_name:  restaurant.Restaurant_name,
		Address:          restaurant.Address,
		LocationID:       restaurant.LocationID,
		Phone_number:     restaurant.Phone_number,
		Restaurant_image: restaurant_image,
	})

	restaurantData := map[string]interface{}{
		"id":               dbRestaurant.ID,
		"restaurant_image": dbRestaurant.Restaurant_image,
		"restaurant_name":  dbRestaurant.Restaurant_name,
		"phone_number":     dbRestaurant.Phone_number,
		"address":          dbRestaurant.Address,
		"owner":            dbRestaurant.ProfileID,
		"location":         location,
	}

	restaurant_image = ""

	helper.SendDataPayload(c, restaurantData, false)
}
