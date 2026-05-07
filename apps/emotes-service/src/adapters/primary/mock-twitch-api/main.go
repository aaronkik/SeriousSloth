package main

import (
	"context"
	mocktwitchapi "emotes-service/src/use-cases/mock-twitch-api"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func handler(ctx context.Context, event events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	path := event.RequestContext.HTTP.Path
	slog.InfoContext(ctx, "mock-twitch-api request", "path", path, "method", event.RequestContext.HTTP.Method)

	body, err := mocktwitchapi.Execute(ctx, path)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       body,
	}, nil
}

func main() {
	slog.SetDefault(lambdacontext.NewLogger())
	lambda.Start(handler)
}
