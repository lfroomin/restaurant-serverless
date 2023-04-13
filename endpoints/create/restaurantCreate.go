package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/lfroomin/restaurant-serverless/internal/httpHelper"
	"github.com/lfroomin/restaurant-serverless/internal/model"
	"github.com/lfroomin/restaurant-serverless/internal/print"
	"net/http"
)

func (h handler) restaurantCreate(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
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
	fmt.Printf("create restaurantName: %s  restaurantId: %s\n", restaurant.Name, *restaurant.Id)

	// Get the geocode of the restaurant address
	if restaurant.Address != nil {
		location, timezoneName, err := h.location.Geocode(*restaurant.Address)
		if err != nil {
			response.StatusCode = http.StatusInternalServerError
			response.Body = httpHelper.ResponseBodyMsg(err.Error())
			return response, nil
		}

		restaurant.Address.Location = &location
		restaurant.Address.TimezoneName = &timezoneName
	}

	if err := h.restaurant.Save(restaurant); err != nil {
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
