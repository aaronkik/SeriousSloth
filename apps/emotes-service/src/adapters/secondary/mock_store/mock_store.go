package mock_store

import (
	"context"
	"emotes-service/src/environment"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var client *dynamodb.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		slog.Error("Error loading AWS config", "error", err)
		panic(err)
	}
	client = dynamodb.NewFromConfig(cfg)
}

func Lookup(ctx context.Context, path string) (string, error) {
	tableName := environment.GetOrFatal("MOCK_RESPONSES_TABLE")

	result, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: path},
		},
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error reading mock response", "path", path, "error", err)
		return "", err
	}
	if result.Item == nil {
		return "", nil
	}
	body, ok := result.Item["responseBody"].(*types.AttributeValueMemberS)
	if !ok {
		return "", nil
	}
	return body.Value, nil
}
