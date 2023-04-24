package dynamo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/lfroomin/restaurant-serverless/internal/model"
	"log"
	"time"
)

type dynamoRestaurantStorer interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
}

const key = "RestaurantId"

type RestaurantStorage struct {
	Client dynamoRestaurantStorer
	Table  string
}

type restaurantItem struct {
	RestaurantId string
	Restaurant   model.Restaurant
	Updated      int64
}

func New(cfg aws.Config, table string) RestaurantStorage {
	return RestaurantStorage{
		Client: dynamodb.NewFromConfig(cfg),
		Table:  table,
	}
}

func (rs RestaurantStorage) Save(restaurant model.Restaurant) error {
	log.Printf("RestaurantStorage.Save restaurantId: %s\n", *restaurant.Id)

	r := restaurantItem{
		RestaurantId: *restaurant.Id,
		Restaurant:   restaurant,
		Updated:      time.Now().UnixMilli(),
	}

	av, err := attributevalue.MarshalMap(r)
	if err != nil {
		return fmt.Errorf("error marshalling value: %w", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(rs.Table),
	}

	_, err = rs.Client.PutItem(context.Background(), input)
	if err != nil {
		return fmt.Errorf("error saving restaurant %q in dynamo: %w", *restaurant.Id, err)
	}
	return nil
}

func (rs RestaurantStorage) Get(restaurantId string) (model.Restaurant, bool, error) {
	log.Printf("RestaurantStorage.Get restaurantId: %s\n", restaurantId)

	input := dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			key: &types.AttributeValueMemberS{Value: restaurantId},
		},
		TableName: aws.String(rs.Table),
	}

	item := &restaurantItem{}
	data, err := rs.Client.GetItem(context.Background(), &input)
	if err != nil {
		return model.Restaurant{}, false, fmt.Errorf("error getting restaurant %q in dynamo: %w", restaurantId, err)
	}

	if data.Item != nil {
		if err = attributevalue.UnmarshalMap(data.Item, &item); err != nil {
			return model.Restaurant{}, false, fmt.Errorf("error unmarshalling value: %w", err)
		}
		return item.Restaurant, true, nil
	}

	return model.Restaurant{}, false, nil
}

func (rs RestaurantStorage) Update(restaurant model.Restaurant) error {
	log.Printf("RestaurantStorage.Update restaurantId: %s\n", *restaurant.Id)

	cond := expression.Equal(expression.Name(key), expression.Value(*restaurant.Id))

	update := expression.Set(
		expression.Name("Restaurant"),
		expression.Value(restaurant),
	).Set(
		expression.Name("Updated"),
		expression.Value(time.Now().UnixMilli()),
	)

	expr, err := expression.NewBuilder().WithUpdate(update).WithCondition(cond).Build()
	if err != nil {
		return err
	}

	input := dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			key: &types.AttributeValueMemberS{Value: *restaurant.Id},
		},
		TableName:                 aws.String(rs.Table),
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       expr.Condition(),
	}

	_, err = rs.Client.UpdateItem(context.Background(), &input)
	if err != nil {
		return fmt.Errorf("error updating restaurant %q in dynamo: %w", *restaurant.Id, err)
	}
	return nil
}

func (rs RestaurantStorage) Delete(restaurantId string) error {
	log.Printf("RestaurantStorage.Delete restaurantId: %s\n", restaurantId)

	input := dynamodb.DeleteItemInput{
		TableName: aws.String(rs.Table),
		Key: map[string]types.AttributeValue{
			key: &types.AttributeValueMemberS{Value: restaurantId},
		},
	}

	_, err := rs.Client.DeleteItem(context.Background(), &input)
	if err != nil {
		return fmt.Errorf("error deleting restaurant %q from dynamo: %w", restaurantId, err)
	}

	return nil
}
