package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/lfroomin/restaurant-serverless/internal/httpHelper"
	"github.com/lfroomin/restaurant-serverless/internal/print"
	"net/http"
)

func restaurantDelete(h handler, request events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	restaurantId := getRequest(request)

	err := processRequest(h, restaurantId)

	return getResponse(err)
}

func getRequest(request events.APIGatewayProxyRequest) string {
	print.Json("Request", request)

	restaurantId := request.PathParameters["restaurantId"]

	return restaurantId
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

func processRequest(h handler, restaurantId string) error {
	fmt.Printf("delete restaurantId: %s\n", restaurantId)

	// Validate input
	if restaurantId == "" {
		return errors.New("restaurantId is empty")
	}

	err := h.restaurant.Delete(restaurantId)
	if err != nil {
		return err
	}

	return nil
}
