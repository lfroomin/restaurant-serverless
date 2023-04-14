package main

import (
	"github.com/lfroomin/restaurant-serverless/internal/dynamo"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_NewHandler(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		restaurantsTable string
	}{
		{
			name:             "happy path",
			restaurantsTable: "RestaurantsTable",
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_ = os.Setenv("RestaurantsTable", tc.restaurantsTable)

			testHandler := newHandler()

			assert.IsType(t, dynamo.RestaurantStorage{}, testHandler.Restaurant)
			assert.Equal(t, tc.restaurantsTable, testHandler.Restaurant.(dynamo.RestaurantStorage).Table)
		})
	}
}
