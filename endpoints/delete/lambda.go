package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lfroomin/restaurant-serverless/controllers"
	"github.com/lfroomin/restaurant-serverless/internal/awsConfig"
	"log"
	"os"
)

// main is called only once, when the Lambda is initialised (started for the first time).
func main() {
	cfg, err := awsConfig.New()
	if err != nil {
		log.Fatal(err)
	}

	restaurantsTable := os.Getenv("RestaurantsTable")

	log.Printf("Env Vars: RestaurantsTable: %s\n", restaurantsTable)

	c := controllers.Restaurant{}.New(cfg, restaurantsTable, "")

	lambda.Start(c.Delete)
}
