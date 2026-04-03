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

var menu_image = ""

func CreateMenu(c *gin.Context) {
	var menu models.Menu

	// get formdata
	menu_name := c.PostForm("name")
	file, _ := c.FormFile("menu_image")

	if menu_name == "" {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Menu name is required"))
		return
	}

	if file != nil {
		avatarUrl, err := helper.SingleImageUpload(c, "menu_image", config.EnvCloudMenuFolder(), "menu")
		if err != nil {
			menu_image = ""
		}
		menu_image = avatarUrl
	}

	menu.Name = menu_name
	menu.Menu_image = menu_image

	createdMenu := database.DB.Create(&menu)
	err = createdMenu.Error

	if err != nil {
		helper.SendErrorPayload(c, http.StatusInternalServerError, err)
		return
	}

	helper.SendDataPayload(c, createdMenu.Value, true)
}

func GetMenus(c *gin.Context) {
	var menus []models.Menu

	menuList := database.DB.Find(&menus)
	err = menuList.Error

	if err != nil {
		helper.SendErrorPayload(c, http.StatusInternalServerError, err)
		return
	}

	helper.SendDataPayload(c, menus, false)
}

func GetMenu(c *gin.Context) {
	menuIdStr := c.Param("id")
	id, _ := strconv.Atoi(menuIdStr)

	var menu models.Menu

	database.DB.First(&menu, id)

	if menu.ID == 0 {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("No menu"))
		return
	}

	helper.SendDataPayload(c, menu, false)
}

func UpdateMenu(c *gin.Context) {
	menuIdStr := c.Param("id")
	id, _ := strconv.Atoi(menuIdStr)

	var menu models.Menu
	var dbMenu models.Menu
	database.DB.First(&dbMenu, id)

	if dbMenu.ID == 0 {
		helper.SendErrorPayload(c, http.StatusInternalServerError, fmt.Errorf("Menu not found"))
		return
	}

	err := c.ShouldBindJSON(&menu)
	if err != nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, err)
		return
	}

	if menu.Name != "" {
		dbMenu.Name = menu.Name
	}

	file, err := c.FormFile("menu_image")

	if err != nil {
		helper.SendErrorPayload(c, http.StatusInternalServerError, err)
		return
	}

	if file != nil {
		avatarUrl, err := helper.SingleImageUpload(c, "menu_image", config.EnvCloudMenuFolder(), "menu")
		if err != nil {
			menu_image = dbMenu.Menu_image
		}
		menu_image = avatarUrl
	}

	// update menu
	updatedMenu := database.DB.Model(&dbMenu).Updates(models.Menu{
		Name:       dbMenu.Name,
		Menu_image: menu_image,
	})

	err = updatedMenu.Error

	if err != nil {
		helper.SendErrorPayload(c, http.StatusInternalServerError, err)
		return
	}

	helper.SendDataPayload(c, updatedMenu.Value, false)
}
