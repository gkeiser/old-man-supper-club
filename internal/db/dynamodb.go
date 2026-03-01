package db

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/grant/old-man-supper-club/internal/models"
)

type Repository struct {
	client    *dynamodb.Client
	tableName string
}

func NewRepository(client *dynamodb.Client, tableName string) *Repository {
	return &Repository{
		client:    client,
		tableName: tableName,
	}
}

// GetConfig fetches the global scoring weights.
func (r *Repository) GetConfig(ctx context.Context) (*models.Config, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "CONFIG#GLOBAL"},
			"SK": &types.AttributeValueMemberS{Value: "METADATA"},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, fmt.Errorf("config not found")
	}

	var config models.Config
	err = attributevalue.UnmarshalMap(out.Item, &config)
	return &config, err
}

// ListRestaurants fetches all restaurants from the table.
func (r *Repository) ListRestaurants(ctx context.Context) ([]models.Restaurant, error) {
	out, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(r.tableName),
		FilterExpression: aws.String("begins_with(PK, :pk) AND SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "RESTAURANT#"},
			":sk": &types.AttributeValueMemberS{Value: "METADATA"},
		},
	})
	if err != nil {
		return nil, err
	}

	var restaurants []models.Restaurant
	err = attributevalue.UnmarshalListOfMaps(out.Items, &restaurants)
	return restaurants, err
}

// SaveRestaurant persists a new restaurant to DynamoDB.
func (r *Repository) SaveRestaurant(ctx context.Context, restaurant models.Restaurant) error {
	item, err := attributevalue.MarshalMap(restaurant)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	return err
}

// SaveReview persists a member's review to DynamoDB.
func (r *Repository) SaveReview(ctx context.Context, review models.Review) error {
	item, err := attributevalue.MarshalMap(review)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	return err
}

// GetRestaurantData fetches a restaurant and all its reviews in ONE query.
func (r *Repository) GetRestaurantData(ctx context.Context, id string) (*models.Restaurant, []models.Review, error) {
	pk := fmt.Sprintf("RESTAURANT#%s", id)

	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pk},
		},
	})
	if err != nil {
		return nil, nil, err
	}

	var restaurant *models.Restaurant
	var reviews []models.Review

	for _, item := range out.Items {
		var sk string
		err = attributevalue.Unmarshal(item["SK"], &sk)
		if err != nil {
			continue
		}

		if sk == "METADATA" {
			err = attributevalue.UnmarshalMap(item, &restaurant)
		} else {
			var review models.Review
			err = attributevalue.UnmarshalMap(item, &review)
			if err == nil {
				reviews = append(reviews, review)
			}
		}
	}

	if restaurant == nil {
		return nil, nil, fmt.Errorf("restaurant not found")
	}

	return restaurant, reviews, nil
}
