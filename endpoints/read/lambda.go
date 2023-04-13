package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lfroomin/restaurant-serverless/internal/awsConfig"
	"github.com/lfroomin/restaurant-serverless/internal/dynamo"
	"github.com/lfroomin/restaurant-serverless/internal/model"
	"os"
)

type restaurantStorer interface {
	Get(restaurantId string) (model.Restaurant, bool, error)
}

type handler struct {
	restaurant restaurantStorer
}

// main is called only once, when the Lambda is initialised (started for the first time).
func main() {
	fmt.Println("Begin main")
	lambda.Start(newHandler().restaurantRead)
}

// newHandler is used to create service clients, read environments variables,
// read configuration from disk etc.
func newHandler() handler {
	cfg, err := awsConfig.New()
	if err != nil {
		panic(err)
	}

	restaurantsTable := os.Getenv("RestaurantsTable")

	fmt.Printf("Env Vars: RestaurantsTable: %s\n", restaurantsTable)

	return handler{
		restaurant: dynamo.New(cfg, restaurantsTable),
	}
}
