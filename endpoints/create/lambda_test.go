package main

import (
	"github.com/lfroomin/restaurant-serverless/internal/dynamo"
	"github.com/lfroomin/restaurant-serverless/internal/geocode"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_NewHandler(t *testing.T) {
	testCases := []struct {
		name             string
		restaurantsTable string
		placeIndex       string
	}{
		{"happy path",
			"RestaurantsTable",
			"LocationPlaceIndex",
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_ = os.Setenv("RestaurantsTable", tc.restaurantsTable)
			_ = os.Setenv("LocationPlaceIndex", tc.placeIndex)

			testHandler := newHandler()

			assert.IsType(t, dynamo.RestaurantStorage{}, testHandler.restaurant)
			assert.Equal(t, tc.restaurantsTable, testHandler.restaurant.(dynamo.RestaurantStorage).Table)
			assert.IsType(t, geocode.LocationService{}, testHandler.location)
			assert.Equal(t, tc.placeIndex, testHandler.location.(geocode.LocationService).PlaceIndex)
		})
	}
}
