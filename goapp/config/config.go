package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func getEnvVar(name string) string {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv(name)
}

func EnvDBDialect() string {
	return getEnvVar("DIALECT")
}

func EnvDBHost() string {
	return getEnvVar("HOST")
}

func EnvDBPort() string {
	return getEnvVar("DBPORT")
}

func EnvDBUser() string {
	return getEnvVar("USER")
}

func EnvDBName() string {
	return getEnvVar("NAME")
}

func EnvDBPassword() string {
	return getEnvVar("PASSWD")
}

func EnvJwtKey() string {
	return getEnvVar("JWT_KEY")
}

func EnvCloudName() string {
	return getEnvVar("CLOUDINARY_CLOUD_NAME")
}

func EnvCloudAPIKey() string {
	return getEnvVar("CLOUDINARY_API_KEY")
}

func EnvCloudAPISecret() string {
	return getEnvVar("CLOUDINARY_API_SECRET")
}

func EnvCloudUserFolder() string {
	return getEnvVar("CLOUDINARY_USERS_FOLDER")
}

func EnvCloudRestaurantFolder() string {
	return getEnvVar("CLOUDINARY_RESTAURANT_FOLDER")
}

func EnvCloudMenuFolder() string {
	return getEnvVar("CLOUDINARY_MENU_FOLDER")
}

func EnvCloudFoodFolder() string {
	return getEnvVar("CLOUDINARY_FOOD_FOLDER")
}

func SmtpEmailHost() string {
	return getEnvVar("SMTP_EMAIL_HOST")
}

func SmtpEmailUsername() string {
	return getEnvVar("SMTP_EMAIL_USER")
}

func SmtpEmailPassword() string {
	return getEnvVar("SMTP_EMAIL_PASS")
}

func SmtpEmailTestAccount() string {
	return getEnvVar("SMTP_EMAIL_TEST_ACCOUNT")
}
