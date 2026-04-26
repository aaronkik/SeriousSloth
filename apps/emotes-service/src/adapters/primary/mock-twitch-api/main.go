package main

import (
	"context"
	"emotes-service/src/environment"
	"fmt"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var ddbClient *dynamodb.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		slog.Error("Error loading AWS config", "error", err)
		panic(err)
	}
	ddbClient = dynamodb.NewFromConfig(cfg)
}

func handler(ctx context.Context, event events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	path := event.RequestContext.HTTP.Path
	slog.InfoContext(ctx, "mock-twitch-api request", "path", path, "method", event.RequestContext.HTTP.Method)

	tableName := environment.GetOrFatal("MOCK_RESPONSES_TABLE")

	result, err := ddbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: path},
		},
	})
	if err != nil {
		slog.Error("Error reading mock response", "path", path, "error", err)
		return events.LambdaFunctionURLResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	if result.Item == nil {
		return events.LambdaFunctionURLResponse{StatusCode: 404, Body: fmt.Sprintf("no mock response for path: %s", path)}, nil
	}

	responseBody, ok := result.Item["responseBody"].(*types.AttributeValueMemberS)
	if !ok {
		return events.LambdaFunctionURLResponse{StatusCode: 500, Body: "responseBody attribute not found or not a string"}, nil
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       responseBody.Value,
	}, nil
}

func main() {
	slog.SetDefault(lambdacontext.NewLogger())
	lambda.Start(handler)
}
