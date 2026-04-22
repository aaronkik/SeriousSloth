package main

import (
	"context"
	"log/slog"

	"emotes-service/src/adapters/secondary/twitch"

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

	_, err = twitch.GetGlobalEmotes(accessToken)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	slog.SetDefault(lambdacontext.NewLogger())
	lambda.Start(handler)
}
