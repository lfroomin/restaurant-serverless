package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
	"github.com/lfroomin/restaurant-serverless/internal/dynamo"
	"github.com/lfroomin/restaurant-serverless/internal/geocode"
	"github.com/lfroomin/restaurant-serverless/internal/httpResponse"
	"github.com/lfroomin/restaurant-serverless/internal/model"
	"github.com/lfroomin/restaurant-serverless/internal/print"
	"log"
	"net/http"
)

type RestaurantStorer interface {
	Save(restaurant model.Restaurant) error
	Get(restaurantId string) (model.Restaurant, bool, error)
	Update(restaurant model.Restaurant) error
	Delete(restaurantId string) error
}

type Geocoder interface {
	Geocode(address model.Address) (model.Location, string, error)
}

type Restaurant struct {
	Restaurant RestaurantStorer
	Location   Geocoder
}

func (r Restaurant) New(cfg aws.Config, restaurantsTable, placeIndex string) Restaurant {
	return Restaurant{
		Restaurant: dynamo.New(cfg, restaurantsTable),
		Location:   geocode.New(cfg, placeIndex),
	}
}

func (r Restaurant) Create(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	print.Json("Request", request)

	restaurant := model.Restaurant{}
	if len(request.Body) > 0 {
		if err := json.Unmarshal([]byte(request.Body), &restaurant); err != nil {
			return httpResponse.NewServerError(fmt.Sprintf("error unmarshalling request body: %s", err.Error())), nil
		}
	} else {
		return httpResponse.NewBadRequest("error request body is empty"), nil
	}

	id := uuid.NewString()
	restaurant.Id = &id
	log.Printf("create restaurantName: %s  restaurantId: %s\n", restaurant.Name, *restaurant.Id)

	// Get the geocode of the restaurant address
	if restaurant.Address != nil {
		location, timezoneName, err := r.Location.Geocode(*restaurant.Address)
		if err != nil {
			return httpResponse.NewServerError(err.Error()), nil
		}

		restaurant.Address.Location = &location
		restaurant.Address.TimezoneName = &timezoneName
	}

	if err := r.Restaurant.Save(restaurant); err != nil {
		return httpResponse.NewServerError(err.Error()), nil
	}

	return httpResponse.New(http.StatusCreated, restaurant), nil
}

func (r Restaurant) Read(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	print.Json("Request", request)

	restaurantId := request.PathParameters["restaurantId"]

	// Validate input
	if restaurantId == "" {
		return httpResponse.NewBadRequest("restaurantId is empty"), nil
	}

	log.Printf("read restaurantId: %s\n", restaurantId)

	restaurant, exists, err := r.Restaurant.Get(restaurantId)
	if err != nil {
		return httpResponse.NewServerError(err.Error()), nil
	}

	if !exists {
		return httpResponse.New(http.StatusNotFound, nil), nil
	}

	return httpResponse.New(http.StatusOK, restaurant), nil
}

func (r Restaurant) Update(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	print.Json("Request", request)

	restaurantId := request.PathParameters["restaurantId"]

	restaurant := model.Restaurant{}
	if len(request.Body) > 0 {
		if err := json.Unmarshal([]byte(request.Body), &restaurant); err != nil {
			return httpResponse.NewServerError(fmt.Sprintf("error unmarshalling request body: %s", err.Error())), nil
		}
	} else {
		return httpResponse.NewBadRequest("error request body is empty"), nil
	}

	// Validate input
	if restaurant.Id == nil || restaurantId != *restaurant.Id {
		return httpResponse.NewBadRequest("restaurantId in URL path parameters and restaurant in body do not match"), nil
	}

	log.Printf("update restaurantName: %s  restaurantId: %s\n", restaurant.Name, *restaurant.Id)

	// Get the geocode of the restaurant address
	if restaurant.Address != nil {
		location, timezoneName, err := r.Location.Geocode(*restaurant.Address)
		if err != nil {
			return httpResponse.NewServerError(err.Error()), nil
		}

		restaurant.Address.Location = &location
		restaurant.Address.TimezoneName = &timezoneName
	}

	if err := r.Restaurant.Update(restaurant); err != nil {
		return httpResponse.NewServerError(err.Error()), nil
	}

	return httpResponse.New(http.StatusOK, restaurant), nil
}

func (r Restaurant) Delete(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	print.Json("Request", request)

	restaurantId := request.PathParameters["restaurantId"]

	// Validate input
	if restaurantId == "" {
		return httpResponse.NewBadRequest("restaurantId is empty"), nil
	}

	log.Printf("delete restaurantId: %s\n", restaurantId)

	err := r.Restaurant.Delete(restaurantId)
	if err != nil {
		return httpResponse.NewServerError(err.Error()), nil
	}

	return httpResponse.New(http.StatusOK, nil), nil
}
