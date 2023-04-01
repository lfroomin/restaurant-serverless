package main

import (
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/lfroomin/restaurant-serverless/internal/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_RestaurantRead(t *testing.T) {
	testCases := []struct {
		name         string
		restaurantId string
		exists       bool
		responseCode int
		responseBody string
		stubError    string
	}{
		{"happy path",
			"restId",
			true,
			http.StatusOK,
			`{"name":""}`,
			"",
		},
		{"empty restaurantId",
			"",
			true,
			http.StatusInternalServerError,
			`{"Message":"restaurantId is empty"}`,
			"",
		},
		{"storage error",
			"restId",
			true,
			http.StatusInternalServerError,
			`{"Message":"an error occurred"}`,
			"an error occurred",
		},
		{"restaurant does not exist",
			"restId",
			false,
			http.StatusNotFound,
			"",
			"",
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h := handler{restaurant: restaurantStorerStub{exists: tc.exists, error: tc.stubError}}

			resp := restaurantRead(h, events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"restaurantId": tc.restaurantId},
			})

			assert.Equal(t, tc.responseCode, resp.StatusCode)
			assert.Equal(t, tc.responseBody, resp.Body)
		})
	}
}

type restaurantStorerStub struct {
	exists bool
	error  string
}

func (s restaurantStorerStub) Get(_ string) (model.Restaurant, bool, error) {
	if s.error != "" {
		return model.Restaurant{}, false, errors.New(s.error)
	}
	if !s.exists {
		return model.Restaurant{}, false, nil
	}
	return model.Restaurant{}, true, nil
}
