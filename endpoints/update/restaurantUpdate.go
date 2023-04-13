package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/lfroomin/restaurant-serverless/internal/httpHelper"
	"github.com/lfroomin/restaurant-serverless/internal/model"
	"github.com/lfroomin/restaurant-serverless/internal/print"
	"net/http"
)

func restaurantUpdate(h handler, request events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
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
			return response

		}
	} else {
		response.StatusCode = http.StatusBadRequest
		response.Body = httpHelper.ResponseBodyMsg("error request body is empty")
		return response
	}

	// Validate input
	if restaurantId != *restaurant.Id {
		response.StatusCode = http.StatusBadRequest
		response.Body = httpHelper.ResponseBodyMsg("restaurantId in URL path parameters and restaurant in body do not match")
		return response
	}

	fmt.Printf("update restaurantName: %s  restaurantId: %s\n", restaurant.Name, *restaurant.Id)

	// Get the geocode of the restaurant address
	if restaurant.Address != nil {
		location, timezoneName, err := h.location.Geocode(*restaurant.Address)
		if err != nil {
			response.StatusCode = http.StatusInternalServerError
			response.Body = httpHelper.ResponseBodyMsg(err.Error())
			return response
		}

		restaurant.Address.Location = &location
		restaurant.Address.TimezoneName = &timezoneName
	}

	if err := h.restaurant.Update(restaurant); err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Body = httpHelper.ResponseBodyMsg(err.Error())
		return response
	}

	response.StatusCode = http.StatusOK
	return response
}
