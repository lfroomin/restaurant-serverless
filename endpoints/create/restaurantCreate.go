package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/lfroomin/restaurant-serverless/internal/httpHelper"
	"github.com/lfroomin/restaurant-serverless/internal/model"
	"github.com/lfroomin/restaurant-serverless/internal/print"
	"net/http"
)

func restaurantCreate(h handler, request events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	restaurant, err := getRequest(request)
	if err != nil {
		return getResponse(model.Restaurant{}, err)
	}

	createdRestaurant, err := processRequest(h, restaurant)

	return getResponse(createdRestaurant, err)
}

func getRequest(request events.APIGatewayProxyRequest) (model.Restaurant, error) {
	print.Json("Request", request)

	restaurant := model.Restaurant{}
	if len(request.Body) > 0 {
		if err := json.Unmarshal([]byte(request.Body), &restaurant); err != nil {
			return model.Restaurant{}, fmt.Errorf("error unmarshalling request body: %s", err.Error())
		}
	} else {
		return model.Restaurant{}, errors.New("error request body is empty")
	}

	return restaurant, nil
}

func getResponse(createdRestaurant model.Restaurant, err error) *events.APIGatewayProxyResponse {
	response := &events.APIGatewayProxyResponse{
		Headers: httpHelper.CORSHeaders,
	}
	defer print.Json("Response", response)

	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Body = httpHelper.ResponseBodyMsg(err.Error())
		return response
	}

	data, err := json.Marshal(createdRestaurant)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Body = httpHelper.ResponseBodyMsg(fmt.Sprintf("error marshalling data: %s", err.Error()))
		return response
	}

	response.StatusCode = http.StatusCreated
	response.Body = string(data)

	return response
}

func processRequest(h handler, restaurant model.Restaurant) (model.Restaurant, error) {
	id := uuid.NewString()
	restaurant.Id = &id
	fmt.Printf("create restaurantName: %s  restaurantId: %s\n", restaurant.Name, *restaurant.Id)

	// Get the geocode of the restaurant address
	if restaurant.Address != nil {
		location, timezoneName, err := h.location.Geocode(*restaurant.Address)
		if err != nil {
			return model.Restaurant{}, err
		}

		restaurant.Address.Location = &location
		restaurant.Address.TimezoneName = &timezoneName
	}

	if err := h.restaurant.Save(restaurant); err != nil {
		return model.Restaurant{}, err
	}

	return restaurant, nil
}
