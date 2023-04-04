package main

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/lfroomin/restaurant-serverless/internal/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type stubError struct {
	restaurant string
	location   string
}

func Test_RestaurantUpdate(t *testing.T) {
	restId, restName := "Rest1", "Rest 1"

	testCases := []struct {
		name         string
		restaurantId string
		restaurant   model.Restaurant
		emptyReqBody bool
		responseCode int
		responseBody string
		stubError    stubError
	}{
		{"happy path",
			restId,
			model.Restaurant{
				Id:      &restId,
				Name:    restName,
				Address: &model.Address{},
			},
			false,

			http.StatusOK,
			"",
			stubError{},
		},
		{"no address",
			restId,
			model.Restaurant{
				Id:   &restId,
				Name: restName,
			},
			false,
			http.StatusOK,
			"",
			stubError{},
		},
		{"mismatch restaurantId",
			"differentRestId",
			model.Restaurant{
				Id:   &restId,
				Name: restName,
			},
			false,
			http.StatusInternalServerError,
			`{"Message":"restaurantId in URL path parameters and restaurant in body do not match"}`,
			stubError{},
		},
		{"storage error",
			restId,
			model.Restaurant{Id: &restId},
			false,
			http.StatusInternalServerError,
			`{"Message":"an error occurred"}`,
			stubError{restaurant: "an error occurred"},
		},
		{"location error",
			restId,
			model.Restaurant{
				Id:      &restId,
				Name:    restName,
				Address: &model.Address{},
			},
			false,
			http.StatusInternalServerError,
			`{"Message":"an error occurred"}`,
			stubError{location: "an error occurred"},
		},
		{"empty request body",
			"",
			model.Restaurant{},
			true,
			http.StatusInternalServerError,
			`{"Message":"error request body is empty"}`,
			stubError{},
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h := handler{
				restaurant: restaurantStorerStub{error: tc.stubError.restaurant},
				location:   locationServiceStub{error: tc.stubError.location},
			}

			request := events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"restaurantId": tc.restaurantId},
			}
			if !tc.emptyReqBody {
				body, _ := json.Marshal(tc.restaurant)
				request.Body = string(body)
			}

			resp := restaurantUpdate(h, request)

			assert.Equal(t, tc.responseCode, resp.StatusCode)
			assert.Equal(t, tc.responseBody, resp.Body)
		})
	}
}

type restaurantStorerStub struct {
	error string
}

func (s restaurantStorerStub) Update(_ model.Restaurant) error {
	if s.error != "" {
		return errors.New(s.error)
	}
	return nil
}

type locationServiceStub struct {
	error string
}

func (s locationServiceStub) Geocode(_ model.Address) (model.Location, string, error) {
	if s.error != "" {
		return model.Location{}, "", errors.New(s.error)
	}
	return model.Location{}, "", nil
}
