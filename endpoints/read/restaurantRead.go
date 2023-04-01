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

func restaurantRead(h handler, request events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	restaurantId := getRequest(request)

	restaurant, exists, err := processRequest(h, restaurantId)

	return getResponse(restaurant, exists, err)
}

func getRequest(request events.APIGatewayProxyRequest) string {
	print.Json("Request", request)

	restaurantId := request.PathParameters["restaurantId"]

	return restaurantId
}

func getResponse(restaurant model.Restaurant, exists bool, err error) *events.APIGatewayProxyResponse {
	response := &events.APIGatewayProxyResponse{
		Headers: httpHelper.CORSHeaders,
	}
	defer print.Json("Response", response)

	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Body = httpHelper.ResponseBodyMsg(err.Error())
		return response
	}

	if !exists {
		response.StatusCode = http.StatusNotFound
		return response
	}

	data, err := json.Marshal(restaurant)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Body = httpHelper.ResponseBodyMsg(fmt.Sprintf("error marshalling data: %s", err.Error()))
		return response
	}

	response.StatusCode = http.StatusOK
	response.Body = string(data)

	return response
}

func processRequest(h handler, restaurantId string) (model.Restaurant, bool, error) {
	fmt.Printf("read restaurantId: %s\n", restaurantId)

	// Validate input
	if restaurantId == "" {
		return model.Restaurant{}, false, errors.New("restaurantId is empty")
	}

	restaurant, exists, err := h.restaurant.Get(restaurantId)
	if err != nil {
		return model.Restaurant{}, false, err
	}

	return restaurant, exists, nil
}
