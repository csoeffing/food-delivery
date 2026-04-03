package main

//jwt "github.com/dgrijalva/jwt-go"
import (
	"crunchgarage/restaurant-food-delivery/database"
	"crunchgarage/restaurant-food-delivery/logging"
	routes "crunchgarage/restaurant-food-delivery/routes"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	database.OpenDB()

	defer database.CloseDB()

	database.AutoMigrate()

	handleRequests()
}

var Port = ":8134"

// Handle API requests
func handleRequests() {
	router := routes.BuildRouter()

	routes.UserRouter(router)
	routes.MenuRouter(router)

	routes.RestaurantRouter(router)
	routes.FoodRouter(router)
	routes.OrderRouter(router)
	routes.OrderItemRouter(router)
	routes.InvoiceRouter(router)
	routes.LocationRouter(router)

	logging.CreateSugared().Debugf("Listening on port %s", Port)

	router.Run(Port)
}
