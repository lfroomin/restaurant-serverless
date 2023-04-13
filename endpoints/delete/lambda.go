package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/lfroomin/restaurant-serverless/internal/dynamo"
	"os"
)

type restaurantStorer interface {
	Delete(restaurantId string) error
}

type handler struct {
	restaurant restaurantStorer
}

// main is called only once, when the Lambda is initialised (started for the first time).
func main() {
	fmt.Println("Begin main")
	lambda.Start(newHandler().restaurantDelete)
}

// newHandler is used to create service clients, read environments variables,
// read configuration from disk etc.
func newHandler() handler {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(fmt.Errorf("failed loading config: %w", err))
	}

	restaurantsTable := os.Getenv("RestaurantsTable")

	fmt.Printf("Env Vars: RestaurantsTable: %s\n", restaurantsTable)

	return handler{
		restaurant: dynamo.RestaurantStorage{
			Client: dynamodb.NewFromConfig(cfg),
			Table:  restaurantsTable,
		},
	}
}
