package main

import (
	"context"
	"emotes-service/src/adapters/secondary/event_store"
	"emotes-service/src/adapters/secondary/twitch"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func handler(ctx context.Context, event events.CloudWatchEvent) error {
	slog.InfoContext(ctx, "canonical log", "event", event)

	accessToken, err := twitch.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	globalEmotes, err := twitch.GetGlobalEmotes(ctx, accessToken)
	if err != nil {
		return err
	}

	aggregate, err := event_store.LoadAggregate(ctx, event_store.GlobalEmotesAggregateId)
	if err != nil {
		return err
	}

	syncEvents := event_store.DecideSyncEvents(event_store.GlobalEmotesAggregateId, aggregate, globalEmotes)

	err = event_store.AppendEvents(ctx, syncEvents)
	if err == nil {
		return nil
	}

	return err
}

func main() {
	slog.SetDefault(lambdacontext.NewLogger())
	lambda.Start(handler)
}
