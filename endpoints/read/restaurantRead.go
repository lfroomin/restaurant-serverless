package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/lfroomin/restaurant-serverless/internal/httpHelper"
	"github.com/lfroomin/restaurant-serverless/internal/print"
	"net/http"
)

func restaurantRead(h handler, request events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
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
		return response
	}

	fmt.Printf("read restaurantId: %s\n", restaurantId)

	restaurant, exists, err := h.restaurant.Get(restaurantId)
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
