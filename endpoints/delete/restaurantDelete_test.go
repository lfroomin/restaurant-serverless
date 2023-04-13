package main

import (
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_RestaurantDelete(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		restaurantId string
		responseCode int
		responseBody string
		stubError    string
	}{
		{"happy path",
			"restId",
			http.StatusOK,
			"",
			"",
		},
		{"empty restaurantId",
			"",
			http.StatusInternalServerError,
			`{"Message":"restaurantId is empty"}`,
			"",
		},
		{"storage error",
			"restId",
			http.StatusInternalServerError,
			`{"Message":"an error occurred"}`,
			"an error occurred",
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h := handler{restaurant: restaurantStorerStub{error: tc.stubError}}

			resp := restaurantDelete(h, events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"restaurantId": tc.restaurantId},
			})

			assert.Equal(t, tc.responseCode, resp.StatusCode)
			assert.Equal(t, tc.responseBody, resp.Body)
		})
	}
}

type restaurantStorerStub struct {
	error string
}

func (s restaurantStorerStub) Delete(_ string) error {
	if s.error != "" {
		return errors.New(s.error)
	}
	return nil
}
