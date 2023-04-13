package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lfroomin/restaurant-serverless/internal/awsConfig"
	"github.com/lfroomin/restaurant-serverless/internal/dynamo"
	"github.com/lfroomin/restaurant-serverless/internal/geocode"
	"github.com/lfroomin/restaurant-serverless/internal/model"
	"os"
)

type restaurantStorer interface {
	Update(restaurant model.Restaurant) error
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
	lambda.Start(newHandler().restaurantUpdate)
}

// newHandler is used to create service clients, read environments variables,
// read configuration from disk etc.
func newHandler() handler {
	cfg, err := awsConfig.New()
	if err != nil {
		panic(err)
	}

	restaurantsTable := os.Getenv("RestaurantsTable")
	placeIndex := os.Getenv("LocationPlaceIndex")

	fmt.Printf("Env Vars: RestaurantsTable: %s  LocationPlaceIndex: %s\n", restaurantsTable, placeIndex)

	return handler{
		restaurant: dynamo.New(cfg, restaurantsTable),
		location:   geocode.New(cfg, placeIndex),
	}
}
