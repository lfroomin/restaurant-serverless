package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/lfroomin/restaurant-serverless/internal/httpHelper"
	"github.com/lfroomin/restaurant-serverless/internal/model"
	"github.com/lfroomin/restaurant-serverless/internal/print"
	"net/http"
)

func restaurantUpdate(h handler, request events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	restaurant, err := getRequest(request)
	if err != nil {
		return getResponse(err)
	}

	err = processRequest(h, restaurant)

	return getResponse(err)
}

func getRequest(request events.APIGatewayProxyRequest) (model.Restaurant, error) {
	print.Json("Request", request)

	restaurantId := request.PathParameters["restaurantId"]

	restaurant := model.Restaurant{}
	if len(request.Body) > 0 {
		if err := json.Unmarshal([]byte(request.Body), &restaurant); err != nil {
			return model.Restaurant{}, fmt.Errorf("error unmarshalling request body: %s", err.Error())
		}
	} else {
		return model.Restaurant{}, errors.New("error request body is empty")
	}

	if restaurantId != *restaurant.Id {
		return model.Restaurant{}, errors.New("restaurantId in URL path parameters and restaurant in body do not match")
	}

	return restaurant, nil
}

func getResponse(err error) *events.APIGatewayProxyResponse {
	response := &events.APIGatewayProxyResponse{
		Headers: httpHelper.CORSHeaders,
	}
	defer print.Json("Response", response)

	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Body = httpHelper.ResponseBodyMsg(err.Error())
		return response
	}

	response.StatusCode = http.StatusOK

	return response
}

func processRequest(h handler, restaurant model.Restaurant) error {
	fmt.Printf("update restaurantName: %s  restaurantId: %s\n", restaurant.Name, *restaurant.Id)

	// Get the geocode of the restaurant address
	if restaurant.Address != nil {
		location, timezoneName, err := h.location.Geocode(*restaurant.Address)
		if err != nil {
			return err
		}

		restaurant.Address.Location = &location
		restaurant.Address.TimezoneName = &timezoneName
	}

	if err := h.restaurant.Update(restaurant); err != nil {
		return err
	}

	return nil
}
