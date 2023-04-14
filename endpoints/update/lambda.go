package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lfroomin/restaurant-serverless/controllers"
	"github.com/lfroomin/restaurant-serverless/internal/awsConfig"
	"os"
)

// main is called only once, when the Lambda is initialised (started for the first time).
func main() {
	fmt.Println("Begin main")
	lambda.Start(newHandler().Update)
}

// newHandler is used to create service clients, read environments variables,
// read configuration from disk etc.
func newHandler() controllers.RestaurantController {
	cfg, err := awsConfig.New()
	if err != nil {
		panic(err)
	}

	restaurantsTable := os.Getenv("RestaurantsTable")
	placeIndex := os.Getenv("LocationPlaceIndex")

	fmt.Printf("Env Vars: RestaurantsTable: %s  LocationPlaceIndex: %s\n", restaurantsTable, placeIndex)

	return controllers.RestaurantController{}.New(cfg, restaurantsTable, placeIndex)
}
