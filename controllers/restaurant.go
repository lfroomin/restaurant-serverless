package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
	"github.com/lfroomin/restaurant-serverless/internal/dynamo"
	"github.com/lfroomin/restaurant-serverless/internal/geocode"
	"github.com/lfroomin/restaurant-serverless/internal/httpHelper"
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

	response := &events.APIGatewayProxyResponse{
		Headers: httpHelper.CORSHeaders,
	}
	defer print.Json("Response", response)

	restaurant := model.Restaurant{}
	if len(request.Body) > 0 {
		if err := json.Unmarshal([]byte(request.Body), &restaurant); err != nil {
			response.StatusCode = http.StatusInternalServerError
			response.Body = httpHelper.ResponseBodyMsg(fmt.Sprintf("error unmarshalling request body: %s", err.Error()))
			return response, nil
		}
	} else {
		response.StatusCode = http.StatusBadRequest
		response.Body = httpHelper.ResponseBodyMsg("error request body is empty")
		return response, nil
	}

	id := uuid.NewString()
	restaurant.Id = &id
	log.Printf("create restaurantName: %s  restaurantId: %s\n", restaurant.Name, *restaurant.Id)

	// Get the geocode of the restaurant address
	if restaurant.Address != nil {
		location, timezoneName, err := r.Location.Geocode(*restaurant.Address)
		if err != nil {
			response.StatusCode = http.StatusInternalServerError
			response.Body = httpHelper.ResponseBodyMsg(err.Error())
			return response, nil
		}

		restaurant.Address.Location = &location
		restaurant.Address.TimezoneName = &timezoneName
	}

	if err := r.Restaurant.Save(restaurant); err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Body = httpHelper.ResponseBodyMsg(err.Error())
		return response, nil
	}

	data, err := json.Marshal(restaurant)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Body = httpHelper.ResponseBodyMsg(fmt.Sprintf("error marshalling data: %s", err.Error()))
		return response, nil
	}

	response.StatusCode = http.StatusCreated
	response.Body = string(data)
	return response, nil
}

func (r Restaurant) Read(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	print.Json("Request", request)

	response := &events.APIGatewayProxyResponse{
		Headers: httpHelper.CORSHeaders,
	}
	defer print.Json("Response", response)

	restaurantId := request.PathParameters["restaurantId"]

	// Validate input
	if restaurantId == "" {
		response.StatusCode = http.StatusBadRequest
		response.Body = httpHelper.ResponseBodyMsg("restaurantId is empty")
		return response, nil
	}

	log.Printf("read restaurantId: %s\n", restaurantId)

	restaurant, exists, err := r.Restaurant.Get(restaurantId)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Body = httpHelper.ResponseBodyMsg(err.Error())
		return response, nil
	}

	if !exists {
		response.StatusCode = http.StatusNotFound
		return response, nil
	}

	data, err := json.Marshal(restaurant)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Body = httpHelper.ResponseBodyMsg(fmt.Sprintf("error marshalling data: %s", err.Error()))
		return response, nil
	}

	response.StatusCode = http.StatusOK
	response.Body = string(data)
	return response, nil
}

func (r Restaurant) Update(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	print.Json("Request", request)

	response := &events.APIGatewayProxyResponse{
		Headers: httpHelper.CORSHeaders,
	}
	defer print.Json("Response", response)

	restaurantId := request.PathParameters["restaurantId"]

	restaurant := model.Restaurant{}
	if len(request.Body) > 0 {
		if err := json.Unmarshal([]byte(request.Body), &restaurant); err != nil {
			response.StatusCode = http.StatusInternalServerError
			response.Body = httpHelper.ResponseBodyMsg(fmt.Sprintf("error unmarshalling request body: %s", err.Error()))
			return response, nil

		}
	} else {
		response.StatusCode = http.StatusBadRequest
		response.Body = httpHelper.ResponseBodyMsg("error request body is empty")
		return response, nil
	}

	// Validate input
	if restaurant.Id == nil || restaurantId != *restaurant.Id {
		response.StatusCode = http.StatusBadRequest
		response.Body = httpHelper.ResponseBodyMsg("restaurantId in URL path parameters and restaurant in body do not match")
		return response, nil
	}

	log.Printf("update restaurantName: %s  restaurantId: %s\n", restaurant.Name, *restaurant.Id)

	// Get the geocode of the restaurant address
	if restaurant.Address != nil {
		location, timezoneName, err := r.Location.Geocode(*restaurant.Address)
		if err != nil {
			response.StatusCode = http.StatusInternalServerError
			response.Body = httpHelper.ResponseBodyMsg(err.Error())
			return response, nil
		}

		restaurant.Address.Location = &location
		restaurant.Address.TimezoneName = &timezoneName
	}

	if err := r.Restaurant.Update(restaurant); err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Body = httpHelper.ResponseBodyMsg(err.Error())
		return response, nil
	}

	data, err := json.Marshal(restaurant)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Body = httpHelper.ResponseBodyMsg(fmt.Sprintf("error marshalling data: %s", err.Error()))
		return response, nil
	}

	response.StatusCode = http.StatusOK
	response.Body = string(data)
	return response, nil
}

func (r Restaurant) Delete(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	print.Json("Request", request)

	response := &events.APIGatewayProxyResponse{
		Headers: httpHelper.CORSHeaders,
	}
	defer print.Json("Response", response)

	restaurantId := request.PathParameters["restaurantId"]

	// Validate input
	if restaurantId == "" {
		response.StatusCode = http.StatusBadRequest
		response.Body = httpHelper.ResponseBodyMsg("restaurantId is empty")
		return response, nil
	}

	log.Printf("delete restaurantId: %s\n", restaurantId)

	err := r.Restaurant.Delete(restaurantId)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Body = httpHelper.ResponseBodyMsg(err.Error())
		return response, nil
	}

	response.StatusCode = http.StatusOK
	return response, nil
}
