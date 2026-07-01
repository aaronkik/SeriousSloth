package main

import (
	"context"
	"emotes-service/src/adapters/secondary/event_store"
	"emotes-service/src/adapters/secondary/twitch"
	"emotes-service/src/environment"
	dispatchchannelsyncs "emotes-service/src/use-cases/dispatch-channel-syncs"
	syncglobalemotes "emotes-service/src/use-cases/sync-global-emotes"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/newrelic/go-agent/v3/integrations/nrlambda"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func handler(ctx context.Context, event events.SQSEvent) (events.SQSEventResponse, error) {
	slog.InfoContext(ctx, "canonical log", "event", event)
	var failures []events.SQSBatchItemFailure

	for _, record := range event.Records {
		if err := handleRecord(ctx, record); err != nil {
			slog.ErrorContext(ctx, "channel-sync record failed", "messageId", record.MessageId, "error", err)
			failures = append(failures, events.SQSBatchItemFailure{ItemIdentifier: record.MessageId})
		}
	}

	return events.SQSEventResponse{BatchItemFailures: failures}, nil
}

func handleRecord(ctx context.Context, record events.SQSMessage) error {
	var body dispatchchannelsyncs.SyncMessage
	if err := json.Unmarshal([]byte(record.Body), &body); err != nil {
		return fmt.Errorf("invalid sync message body: %w", err)
	}

	twitchId := strings.TrimSpace(body.TwitchId)
	if twitchId == "" {
		return fmt.Errorf("sync message missing twitchId")
	}

	if txn := newrelic.FromContext(ctx); txn != nil {
		txn.AddAttribute("channel.twitchId", twitchId)
	}

	aggregateId := event_store.ChannelAggregateId(twitchId)
	slog.InfoContext(ctx, "syncing channel emotes", "twitchId", twitchId, "aggregateId", aggregateId)

	fetch := func(ctx context.Context, accessToken string) ([]twitch.GlobalEmote, error) {
		return twitch.GetChannelEmotes(ctx, accessToken, twitchId)
	}

	return syncglobalemotes.ExecuteForAggregate(ctx, aggregateId, fetch)
}

func main() {
	logger := lambdacontext.NewLogger().With(
		slog.Group("tags",
			"project", environment.GetOrFatal("PROJECT"),
			"stack", environment.GetOrFatal("STACK"),
		),
	)
	slog.SetDefault(logger)

	app, err := newrelic.NewApplication(
		nrlambda.ConfigOption(),
		newrelic.ConfigFromEnvironment(),
	)
	if err != nil {
		slog.Error("error creating app (invalid config)", "err", err)
	}

	nrlambda.Start(handler, app)
}
