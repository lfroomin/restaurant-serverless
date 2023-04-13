package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/location"
	"github.com/lfroomin/restaurant-serverless/internal/dynamo"
	"github.com/lfroomin/restaurant-serverless/internal/geocode"
	"github.com/lfroomin/restaurant-serverless/internal/model"
	"os"
)

type restaurantStorer interface {
	Save(restaurant model.Restaurant) error
}

type geocoder interface {
	Geocode(address model.Address) (model.Location, string, error)
}

type handler struct {
	restaurant restaurantStorer
	location   geocoder
}

// main is called only once, when the Lambda is initialised (started for the first time).
func main() {
	fmt.Println("Begin main")
	lambda.Start(newHandler().restaurantCreate)
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
	placeIndex := os.Getenv("LocationPlaceIndex")

	fmt.Printf("Env Vars: RestaurantsTable: %s  LocationPlaceIndex: %s\n", restaurantsTable, placeIndex)

	return handler{
		restaurant: dynamo.RestaurantStorage{
			Client: dynamodb.NewFromConfig(cfg),
			Table:  restaurantsTable,
		},
		location: geocode.LocationService{
			Client:     location.NewFromConfig(cfg),
			PlaceIndex: placeIndex,
		},
	}
}
