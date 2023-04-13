package main

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/lfroomin/restaurant-serverless/internal/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type stubError struct {
	restaurant string
	location   string
}

func Test_RestaurantCreate(t *testing.T) {
	t.Parallel()
	restName := "Rest 1"
	restaurantExp, _ := json.Marshal(model.Restaurant{
		Name: restName,
		Address: &model.Address{
			Location:     &model.Location{},
			TimezoneName: new(string),
		},
	})
	restaurantNoAddressExp, _ := json.Marshal(model.Restaurant{
		Name: restName,
	})

	testCases := []struct {
		name         string
		restaurant   model.Restaurant
		emptyReqBody bool
		responseCode int
		responseBody string
		stubError    stubError
	}{
		{"happy path",
			model.Restaurant{
				Name:    restName,
				Address: &model.Address{},
			},
			false,

			http.StatusCreated,
			string(restaurantExp),
			stubError{},
		},
		{"no address",
			model.Restaurant{
				Name: restName,
			},
			false,
			http.StatusCreated,
			string(restaurantNoAddressExp),
			stubError{},
		},
		{"storage error",
			model.Restaurant{},
			false,
			http.StatusInternalServerError,
			`{"Message":"an error occurred"}`,
			stubError{restaurant: "an error occurred"},
		},
		{"location error",
			model.Restaurant{
				Name:    restName,
				Address: &model.Address{},
			},
			false,
			http.StatusInternalServerError,
			`{"Message":"an error occurred"}`,
			stubError{location: "an error occurred"},
		},
		{"empty request body",
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

			request := events.APIGatewayProxyRequest{}
			if !tc.emptyReqBody {
				body, _ := json.Marshal(tc.restaurant)
				request = events.APIGatewayProxyRequest{Body: string(body)}
			}

			resp := restaurantCreate(h, request)

			assert.Equal(t, tc.responseCode, resp.StatusCode)

			if tc.responseCode != http.StatusCreated {
				assert.Equal(t, tc.responseBody, resp.Body)
			} else {
				// Convert to type Restaurant so comparison can be done
				// without the "Id" field
				expRestaurant := model.Restaurant{}
				_ = json.Unmarshal([]byte(tc.responseBody), &expRestaurant)
				body := model.Restaurant{}
				_ = json.Unmarshal([]byte(resp.Body), &body)

				diff := cmp.Diff(
					expRestaurant,
					body,
					cmpopts.IgnoreFields(model.Restaurant{}, "Id"),
				)
				assert.Empty(t, diff)
			}
		})
	}
}

type restaurantStorerStub struct {
	error string
}

func (s restaurantStorerStub) Save(_ model.Restaurant) error {
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
