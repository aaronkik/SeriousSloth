package main

import (
	"context"
	"emotes-service/src/adapters/secondary/event_store"
	"emotes-service/src/dynamodbstream"
	readmodelproducer "emotes-service/src/use-cases/emotes-read-model-producer"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

func handler(ctx context.Context, event events.DynamoDBEvent) error {
	slog.InfoContext(ctx, "canonical log", "event", event)

	domainEvents := make([]event_store.EmoteServiceEvent, 0, len(event.Records))
	for _, record := range event.Records {
		var emoteEvent event_store.EmoteServiceEvent
		if err := attributevalue.UnmarshalMap(dynamodbstream.ToAttributeValueMap(record.Change.NewImage), &emoteEvent); err != nil {
			slog.ErrorContext(ctx, "unmarshal failed", "error", err)
			return err
		}
		domainEvents = append(domainEvents, emoteEvent)
	}

	return readmodelproducer.Execute(ctx, domainEvents)
}

func main() {
	slog.SetDefault(lambdacontext.NewLogger())
	lambda.Start(handler)
}
