package httpResponse

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/lfroomin/restaurant-serverless/internal/print"
	"net/http"
)

var CORSHeaders = map[string]string{
	"Access-Control-Allow-Origin":      "*",
	"Access-Control-Allow-Credentials": "true",
}

func New(statusCode int, data any) *events.APIGatewayProxyResponse {
	response := &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    CORSHeaders,
	}

	defer print.Json("Response", response)

	if data == nil {
		return response
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		fmt.Printf("error marshalling data for API Gateway response: %s", err.Error())
		return response
	}

	response.Body = string(bytes)
	return response
}

func NewNoEncode(statusCode int, data string) *events.APIGatewayProxyResponse {
	response := &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    CORSHeaders,
		Body:       data,
	}

	defer print.Json("Response", response)

	return response
}

func NewMessage(statusCode int, msg string) *events.APIGatewayProxyResponse {
	data := map[string]string{"Message": msg}
	return New(statusCode, data)
}

func NewBadRequest(msg string) *events.APIGatewayProxyResponse {
	return NewMessage(http.StatusBadRequest, msg)
}

func NewServerError(msg string) *events.APIGatewayProxyResponse {
	return NewMessage(http.StatusInternalServerError, msg)
}
