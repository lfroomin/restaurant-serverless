package controllers

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

func Test_Create(t *testing.T) {
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
		{
			name: "happy path",
			restaurant: model.Restaurant{
				Name:    restName,
				Address: &model.Address{},
			},
			responseCode: http.StatusCreated,
			responseBody: string(restaurantExp),
		},
		{
			name: "no address",
			restaurant: model.Restaurant{
				Name: restName,
			},
			responseCode: http.StatusCreated,
			responseBody: string(restaurantNoAddressExp),
		},
		{
			name:         "storage error",
			restaurant:   model.Restaurant{},
			responseCode: http.StatusInternalServerError,
			responseBody: `{"Message":"an error occurred"}`,
			stubError:    stubError{restaurant: "an error occurred"},
		},
		{
			name: "location error",
			restaurant: model.Restaurant{
				Name:    restName,
				Address: &model.Address{},
			},
			responseCode: http.StatusInternalServerError,
			responseBody: `{"Message":"an error occurred"}`,
			stubError:    stubError{location: "an error occurred"},
		},
		{
			name:         "empty request body",
			emptyReqBody: true,
			responseCode: http.StatusBadRequest,
			responseBody: `{"Message":"error request body is empty"}`,
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rc := RestaurantController{
				Restaurant: restaurantStorerStub{error: tc.stubError.restaurant},
				Location:   locationServiceStub{error: tc.stubError.location},
			}

			request := events.APIGatewayProxyRequest{}
			if !tc.emptyReqBody {
				body, _ := json.Marshal(tc.restaurant)
				request = events.APIGatewayProxyRequest{Body: string(body)}
			}

			resp, _ := rc.Create(request)

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

func Test_Read(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		restaurantId string
		notExist     bool
		responseCode int
		responseBody string
		stubError    string
	}{
		{
			name:         "happy path",
			restaurantId: "restId",
			responseCode: http.StatusOK,
			responseBody: `{"name":""}`,
		},
		{
			name:         "empty restaurantId",
			responseCode: http.StatusBadRequest,
			responseBody: `{"Message":"restaurantId is empty"}`,
		},
		{
			name:         "storage error",
			restaurantId: "restId",
			responseCode: http.StatusInternalServerError,
			responseBody: `{"Message":"an error occurred"}`,
			stubError:    "an error occurred",
		},
		{
			name:         "restaurant does not exist",
			restaurantId: "restId",
			notExist:     true,
			responseCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rc := RestaurantController{Restaurant: restaurantStorerStub{notExist: tc.notExist, error: tc.stubError}}

			resp, _ := rc.Read(events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"restaurantId": tc.restaurantId},
			})

			assert.Equal(t, tc.responseCode, resp.StatusCode)
			assert.Equal(t, tc.responseBody, resp.Body)
		})
	}
}

func Test_Update(t *testing.T) {
	t.Parallel()
	restId, restName := "Rest1", "Rest 1"
	restaurantExp, _ := json.Marshal(model.Restaurant{
		Id:   &restId,
		Name: restName,
		Address: &model.Address{
			Location:     &model.Location{},
			TimezoneName: new(string),
		},
	})
	restaurantNoAddressExp, _ := json.Marshal(model.Restaurant{
		Id:   &restId,
		Name: restName,
	})

	testCases := []struct {
		name         string
		restaurantId string
		restaurant   model.Restaurant
		emptyReqBody bool
		responseCode int
		responseBody string
		stubError    stubError
	}{
		{
			name:         "happy path",
			restaurantId: restId,
			restaurant: model.Restaurant{
				Id:      &restId,
				Name:    restName,
				Address: &model.Address{},
			},
			responseCode: http.StatusOK,
			responseBody: string(restaurantExp),
		},
		{
			name:         "no address",
			restaurantId: restId,
			restaurant: model.Restaurant{
				Id:   &restId,
				Name: restName,
			},
			responseCode: http.StatusOK,
			responseBody: string(restaurantNoAddressExp),
		},
		{
			name:         "restaurantId is nil",
			restaurantId: restId,
			restaurant: model.Restaurant{
				Name: restName,
			},
			responseCode: http.StatusBadRequest,
			responseBody: `{"Message":"restaurantId in URL path parameters and restaurant in body do not match"}`,
		},
		{
			name:         "mismatch restaurantId",
			restaurantId: "differentRestId",
			restaurant: model.Restaurant{
				Id:   &restId,
				Name: restName,
			},
			responseCode: http.StatusBadRequest,
			responseBody: `{"Message":"restaurantId in URL path parameters and restaurant in body do not match"}`,
		},
		{
			name:         "storage error",
			restaurantId: restId,
			restaurant:   model.Restaurant{Id: &restId},
			responseCode: http.StatusInternalServerError,
			responseBody: `{"Message":"an error occurred"}`,
			stubError:    stubError{restaurant: "an error occurred"},
		},
		{
			name:         "location error",
			restaurantId: restId,
			restaurant: model.Restaurant{
				Id:      &restId,
				Name:    restName,
				Address: &model.Address{},
			},
			responseCode: http.StatusInternalServerError,
			responseBody: `{"Message":"an error occurred"}`,
			stubError:    stubError{location: "an error occurred"},
		},
		{
			name:         "empty request body",
			emptyReqBody: true,
			responseCode: http.StatusBadRequest,
			responseBody: `{"Message":"error request body is empty"}`,
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rc := RestaurantController{
				Restaurant: restaurantStorerStub{error: tc.stubError.restaurant},
				Location:   locationServiceStub{error: tc.stubError.location},
			}

			request := events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"restaurantId": tc.restaurantId},
			}
			if !tc.emptyReqBody {
				body, _ := json.Marshal(tc.restaurant)
				request.Body = string(body)
			}

			resp, _ := rc.Update(request)

			assert.Equal(t, tc.responseCode, resp.StatusCode)
			assert.Equal(t, tc.responseBody, resp.Body)
		})
	}
}

func Test_Delete(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		restaurantId string
		responseCode int
		responseBody string
		stubError    string
	}{
		{
			name:         "happy path",
			restaurantId: "restId",
			responseCode: http.StatusOK,
		},
		{
			name:         "empty restaurantId",
			responseCode: http.StatusBadRequest,
			responseBody: `{"Message":"restaurantId is empty"}`,
		},
		{
			name:         "storage error",
			restaurantId: "restId",
			responseCode: http.StatusInternalServerError,
			responseBody: `{"Message":"an error occurred"}`,
			stubError:    "an error occurred",
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rc := RestaurantController{Restaurant: restaurantStorerStub{error: tc.stubError}}

			resp, _ := rc.Delete(events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"restaurantId": tc.restaurantId},
			})

			assert.Equal(t, tc.responseCode, resp.StatusCode)
			assert.Equal(t, tc.responseBody, resp.Body)
		})
	}
}

type restaurantStorerStub struct {
	notExist bool
	error    string
}

func (s restaurantStorerStub) Save(_ model.Restaurant) error {
	if s.error != "" {
		return errors.New(s.error)
	}
	return nil
}

func (s restaurantStorerStub) Get(_ string) (model.Restaurant, bool, error) {
	if s.error != "" {
		return model.Restaurant{}, false, errors.New(s.error)
	}
	if s.notExist {
		return model.Restaurant{}, false, nil
	}
	return model.Restaurant{}, true, nil
}

func (s restaurantStorerStub) Update(_ model.Restaurant) error {
	if s.error != "" {
		return errors.New(s.error)
	}
	return nil
}

func (s restaurantStorerStub) Delete(_ string) error {
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
