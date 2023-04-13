package dynamo

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/lfroomin/restaurant-serverless/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Save(t *testing.T) {
	t.Parallel()
	restId := "restId"

	testCases := []struct {
		name       string
		restaurant model.Restaurant
		stubError  string
		errMsg     string
	}{
		{
			name:       "happy path",
			restaurant: model.Restaurant{Id: &restId},
		},
		{
			name:       "error",
			restaurant: model.Restaurant{Id: &restId},
			stubError:  "an error occurred",
			errMsg:     "error saving restaurant \"restId\" in dynamo: an error occurred",
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rs := RestaurantStorage{
				Client: dynamoRestaurantStorerStub{error: tc.stubError},
				Table:  "RestaurantsTable-Test",
			}
			err := rs.Save(tc.restaurant)

			if tc.errMsg != "" {
				if assert.Error(t, err) {
					assert.Equal(t, tc.errMsg, err.Error())
				}
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_Get(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		restId    string
		stubError string
		errMsg    string
	}{
		{
			name:   "happy path",
			restId: "restId",
		},
		{
			name: "unknown restaurantId",
		},
		{
			name:      "error",
			restId:    "restId",
			stubError: "an error occurred",
			errMsg:    "error getting restaurant \"restId\" in dynamo: an error occurred",
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rs := RestaurantStorage{
				Client: dynamoRestaurantStorerStub{restaurantId: tc.restId, error: tc.stubError},
			}
			restaurant, ok, err := rs.Get(tc.restId)

			if tc.errMsg != "" {
				if assert.Error(t, err) {
					assert.Equal(t, tc.errMsg, err.Error())
				}
			} else if tc.restId != "" {
				assert.Nil(t, err)
				assert.Equal(t, model.Restaurant{Id: &tc.restId}, restaurant)
				assert.True(t, ok)
			} else {
				assert.Nil(t, err)
				assert.False(t, ok)
				assert.Equal(t, model.Restaurant{}, restaurant)
			}
		})
	}
}

func Test_Update(t *testing.T) {
	t.Parallel()
	restId := "restId"

	testCases := []struct {
		name       string
		restaurant model.Restaurant
		stubError  string
		errMsg     string
	}{
		{
			name:       "happy path",
			restaurant: model.Restaurant{Id: &restId},
		},
		{
			name:       "error",
			restaurant: model.Restaurant{Id: &restId},
			stubError:  "an error occurred",
			errMsg:     "error updating restaurant \"restId\" in dynamo: an error occurred",
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rs := RestaurantStorage{
				Client: dynamoRestaurantStorerStub{error: tc.stubError},
				Table:  "RestaurantsTable-Test",
			}
			err := rs.Update(tc.restaurant)

			if tc.errMsg != "" {
				if assert.Error(t, err) {
					assert.Equal(t, tc.errMsg, err.Error())
				}
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_Delete(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		restId    string
		stubError string
		errMsg    string
	}{
		{
			name:   "happy path",
			restId: "restId",
		},
		{
			name:      "error",
			restId:    "restId",
			stubError: "an error occurred",
			errMsg:    "error deleting restaurant \"restId\" from dynamo: an error occurred",
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rs := RestaurantStorage{
				Client: dynamoRestaurantStorerStub{error: tc.stubError},
				Table:  "RestaurantsTable-Test",
			}
			err := rs.Delete(tc.restId)

			if tc.errMsg != "" {
				if assert.Error(t, err) {
					assert.Equal(t, tc.errMsg, err.Error())
				}
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

type dynamoRestaurantStorerStub struct {
	restaurantId string
	restaurants  []model.Restaurant
	error        string
}

func (s dynamoRestaurantStorerStub) PutItem(_ context.Context, _ *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	if s.error != "" {
		return nil, errors.New(s.error)
	}
	return nil, nil
}

func (s dynamoRestaurantStorerStub) GetItem(_ context.Context, _ *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	if s.error != "" {
		return nil, errors.New(s.error)
	}
	if s.restaurantId != "" {
		return restaurantItemOutput(s.restaurantId)
	}
	return &dynamodb.GetItemOutput{}, nil
}

func (s dynamoRestaurantStorerStub) UpdateItem(_ context.Context, _ *dynamodb.UpdateItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	if s.error != "" {
		return nil, errors.New(s.error)
	}
	return nil, nil
}

func (s dynamoRestaurantStorerStub) DeleteItem(_ context.Context, _ *dynamodb.DeleteItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	if s.error != "" {
		return nil, errors.New(s.error)
	}
	return nil, nil
}

func restaurantItemOutput(restaurantId string) (*dynamodb.GetItemOutput, error) {
	restaurant := model.Restaurant{
		Id: &restaurantId,
	}
	restaurantItem := restaurantItem{
		RestaurantId: restaurantId,
		Restaurant:   restaurant,
		Updated:      12345,
	}

	av, err := attributevalue.MarshalMap(restaurantItem)
	if err != nil {
		return nil, err
	}
	return &dynamodb.GetItemOutput{Item: av}, nil
}
