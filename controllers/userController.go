package controller

import (
	"crunchgarage/restaurant-food-delivery/config"
	"crunchgarage/restaurant-food-delivery/database"
	helper "crunchgarage/restaurant-food-delivery/helpers"
	"crunchgarage/restaurant-food-delivery/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var err error

var profile_image = ""

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword, password string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "Password does not match"
		check = false
	}

	return check, msg
}

/*sign up*/
func SignUp(c *gin.Context) {
	var user models.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, err)
		return
	}

	// check if email exists
	var dbUser models.User
	database.DB.Where("email = ?", user.Email).First(&dbUser)
	if dbUser.ID != 0 {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Email address already exists"))
		return
	}

	// check if phone phone number exists
	database.DB.Where("phone = ?", user.Phone).First(&dbUser)
	if dbUser.ID != 0 {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Phone number  already exists"))
		return
	}

	pro_type := ""

	if user.User_type == "CLIENT" {
		pro_type = ""
	} else {
		pro_type = "CHEF"
	}

	username_email_strip := strings.Split(user.Email, "@")[0]

	user_ := models.User{
		First_name: user.First_name,
		Last_name:  user.Last_name,
		User_type:  user.User_type,
		Email:      user.Email,
		Phone:      user.Phone,
		User_name:  username_email_strip,
		Password:   HashPassword(user.Password),
		Profile: []models.Profile{{
			First_name: user.First_name,
			Last_name:  user.Last_name,
			User_type:  user.User_type,
			User_name:  username_email_strip,
			Pro_type:   pro_type,
		}},
	}

	createdUser := database.DB.Create(&user_)
	err = createdUser.Error

	if err != nil {
		helper.SendErrorPayload(c, http.StatusInternalServerError, err)
		return
	}

	_, err = helper.RegisterEmailAccount(user)

	if err != nil {
		fmt.Println(err)
	}

	// if successful return user profile
	helper.SendDataPayload(c, user_.Profile[0], true)
}

// login
func Login(c *gin.Context) {

	var user models.User
	var dbUser models.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, err)
		return
	}

	// find a user with username and see if that user even exists
	database.DB.Where("user_name = ?", user.User_name).First(&dbUser)

	if dbUser.ID == 0 {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("User does not exist"))
		return
	}

	//check if the password is correct
	passwordIsCorrect, msg := VerifyPassword(user.Password, dbUser.Password)
	if !passwordIsCorrect {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("%s", msg))
		return
	}

	accessTokenExpiration := time.Duration(60) * time.Minute
	refreshTokenExpiration := time.Duration(30*24) * time.Hour

	//user controler
	accessToken, accessTokenExpiresAt, err := helper.GenerateToken(dbUser, accessTokenExpiration)
	if err != nil {
		helper.SendErrorPayload(c, http.StatusInternalServerError, err)
	}

	refreshToken, _, err := helper.GenerateToken(dbUser, refreshTokenExpiration)
	if err != nil {
		helper.SendErrorPayload(c, http.StatusInternalServerError, err)
	}

	signings := &models.Signings{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiration: time.Unix(accessTokenExpiresAt, 0).Format(time.RFC3339),
	}

	helper.SendDataPayload(c, signings, false)
}

// get user info from profile table
func GetUser(c *gin.Context) {
	userIdStr := c.Param("id")
	user_id, _ := strconv.Atoi(userIdStr)

	var profile models.Profile
	var restaurant []models.Restaurant

	database.DB.Where("user_id = ?", user_id).First(&profile)
	if profile.ID == 0 {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("User does not exist"))
		return
	}
	database.DB.Model(&profile).Related(&restaurant)

	profile.Restaurant = restaurant

	helper.SendDataPayload(c, profile, false)
}

func UpdateUser(c *gin.Context) {
	userIdStr := c.Param("id")
	id, _ := strconv.Atoi(userIdStr)

	var user models.User
	var profile models.Profile
	var restaurant []models.Restaurant
	var dbUser models.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, err)
		return
	}

	database.DB.First(&dbUser, id)
	database.DB.Where("user_id = ?", id).First(&profile)
	database.DB.Model(&profile).Related(&restaurant)

	if dbUser.ID == 0 {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("User does not exist"))
		return
	}

	//check for media

	file, err := c.FormFile("profile_image")

	if file != nil {
		avatarUrl, err := helper.SingleImageUpload(c, "profile_image", config.EnvCloudMenuFolder())
		if err != nil {
			profile_image = dbUser.Profile_image
		}
		profile_image = avatarUrl
	}

	dbUser.First_name = user.First_name
	dbUser.Last_name = user.Last_name
	dbUser.User_name = user.User_name
	dbUser.Profile_image = profile_image

	// update user
	updatedUser := database.DB.Save(&dbUser)
	err = updatedUser.Error

	// update profile
	database.DB.Model(&profile).Updates(models.Profile{
		First_name:    dbUser.First_name,
		Last_name:     dbUser.Last_name,
		User_name:     dbUser.User_name,
		Profile_image: profile_image,
	})

	if err != nil {
		helper.SendErrorPayload(c, http.StatusInternalServerError, err)
		return
	}

	profile.Restaurant = restaurant
	dbUser.Profile = []models.Profile{profile}

	profile_image = ""

	helper.SendDataPayload(c, dbUser.Profile, false)
}
