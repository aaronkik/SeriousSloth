package main

import (
	"context"
	"emotes-service/src/adapters/secondary/dynamodb"
	"emotes-service/src/adapters/secondary/twitch"
	"emotes-service/src/environment"
	"log/slog"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func handler(ctx context.Context, event events.CloudWatchEvent) error {
	slog.InfoContext(ctx, "canonical log", "event", event)

	accessToken, err := twitch.GetAccessToken()
	if err != nil {
		return err
	}

	globalEmotes, err := twitch.GetGlobalEmotes(accessToken)
	if err != nil {
		return err
	}

	err = dynamodb.PutItem(dynamodb.PutItemInput{
		TableName:  environment.GetOrFatal("TWITCH_EMOTES_SNAPSHOT_TABLE"),
		PK:         "GLOBAL",
		SK:         new(time.Now().UTC().Format(time.RFC3339)),
		Attributes: map[string]interface{}{"globalEmotes": globalEmotes},
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	slog.SetDefault(lambdacontext.NewLogger())
	lambda.Start(handler)
}
