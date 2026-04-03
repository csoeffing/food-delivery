package controller

import (
	"crunchgarage/restaurant-food-delivery/config"
	"crunchgarage/restaurant-food-delivery/database"
	helper "crunchgarage/restaurant-food-delivery/helpers"
	"crunchgarage/restaurant-food-delivery/middleware"
	"crunchgarage/restaurant-food-delivery/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bytedance/gopkg/util/logger"
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
	var payloadUser models.User

	err := c.ShouldBind(&payloadUser)
	if err != nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, err)
		return
	}

	// check if email exists
	var dbUser models.User
	database.DB.Where("email = ? OR phone = ? OR user_name = ?", payloadUser.Email, payloadUser.Phone, payloadUser.UserName).First(&dbUser)
	if dbUser.ID != 0 {
		if dbUser.Email == payloadUser.Email {
			helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Email address already exists"))
			return
		} else if dbUser.Phone == payloadUser.Phone {
			helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Phone number already exists"))
			return
		} else if dbUser.UserName == payloadUser.UserName {
			helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("Username already exists"))
			return
		}
	}

	pro_type := ""

	if payloadUser.UserType == "CLIENT" {
		pro_type = ""
	} else {
		pro_type = "CHEF"
	}

	//username_email_strip := strings.Split(user.Email, "@")[0]

	user_ := models.User{
		FirstName: payloadUser.FirstName,
		LastName:  payloadUser.LastName,
		UserType:  payloadUser.UserType,
		Email:     payloadUser.Email,
		Phone:     payloadUser.Phone,
		UserName:  payloadUser.UserName,
		Password:  HashPassword(payloadUser.Password),
		Profiles: []models.Profile{{
			FirstName: payloadUser.FirstName,
			LastName:  payloadUser.LastName,
			UserType:  payloadUser.UserType,
			UserName:  payloadUser.UserName,
			ProType:   pro_type,
		}},
	}

	createdUser := database.DB.Create(&user_)
	err = createdUser.Error

	if err != nil {
		helper.SendErrorPayload(c, http.StatusInternalServerError, err)
		return
	}

	_, err = helper.RegisterEmailAccount(payloadUser)

	if err != nil {
		fmt.Println(err)
	}

	imageFile, err := c.FormFile("profileImage")

	if err == nil {
		logger.Infof("Uploading %s...\n", imageFile.Filename)
		avatarUrl, err := helper.SingleImageUpload(c, "profileImage", config.EnvCloudMenuFolder(), "user_"+payloadUser.UserName)
		if err != nil {
			helper.SendErrorPayload(c, http.StatusInternalServerError, err)
			return
		} else {
			logger.Infof("URL: %s\n", avatarUrl)
		}
	}

	// if successful return user profile
	//helper.SendDataPayload(c, user_.Profiles[0], true)
	helper.SendDataPayload(c, createdUser.Value, true)
}

// login
func Login(c *gin.Context) {
	var payloadUser models.User
	var dbUser models.User

	err := c.ShouldBindJSON(&payloadUser)
	if err != nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, err)
		return
	}

	// find a user with username and see if that user even exists
	database.DB.Where("user_name = ?", payloadUser.UserName).First(&dbUser)

	if dbUser.ID == 0 {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("User does not exist"))
		return
	}

	//check if the password is correct
	passwordIsCorrect, msg := VerifyPassword(payloadUser.Password, dbUser.Password)
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

	requestingUserId, err := middleware.GetRequestingUserId(c)

	if err != nil {
		helper.SendErrorPayload(c, http.StatusBadRequest, err)
		return
	}

	fmt.Printf("User %d is requesting user %d\n", requestingUserId, user_id)

	// TODO: determine if requesting user has appropriate access for user user_id

	var user models.User

	database.DB.Where("id = ?", user_id).First(&user)
	if user.ID == 0 {
		helper.SendErrorPayload(c, http.StatusBadRequest, fmt.Errorf("User does not exist"))
		return
	}

	var profiles []models.Profile
	database.DB.Model(&user).Related(&profiles)
	user.Profiles = profiles

	//var restaurants []models.Restaurant
	//database.DB.Model(&user).Related(&restaurants)
	//user.Restaurant = restaurants

	helper.SendDataPayload(c, user, false)
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
		avatarUrl, err := helper.SingleImageUpload(c, "profile_image", config.EnvCloudMenuFolder(), "user_"+dbUser.UserName)
		if err != nil {
			profile_image = dbUser.ProfileImage
		}
		profile_image = avatarUrl
	}

	dbUser.FirstName = user.FirstName
	dbUser.LastName = user.LastName
	dbUser.UserName = user.UserName
	dbUser.ProfileImage = profile_image

	// update user
	updatedUser := database.DB.Save(&dbUser)
	err = updatedUser.Error

	// update profile
	database.DB.Model(&profile).Updates(models.Profile{
		FirstName:    dbUser.FirstName,
		LastName:     dbUser.LastName,
		UserName:     dbUser.UserName,
		ProfileImage: profile_image,
	})

	if err != nil {
		helper.SendErrorPayload(c, http.StatusInternalServerError, err)
		return
	}

	profile.Restaurant = restaurant
	dbUser.Profiles = []models.Profile{profile}

	profile_image = ""

	helper.SendDataPayload(c, dbUser.Profiles, false)
}
