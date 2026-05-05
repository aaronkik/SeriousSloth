package main

import (
	"context"
	"emotes-service/src/adapters/secondary/event_store"
	"emotes-service/src/adapters/secondary/projections_store"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func handler(ctx context.Context, event events.DynamoDBEvent) error {
	slog.InfoContext(ctx, "canonical log", "event", event)

	for _, record := range event.Records {
		var emoteEvent event_store.EmoteServiceEvent
		if err := attributevalue.UnmarshalMap(toAttributeValueMap(record.Change.NewImage), &emoteEvent); err != nil {
			slog.ErrorContext(ctx, "unmarshal failed", "error", err)
			return err
		}

		err := projections_store.Persist(ctx, emoteEvent)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	slog.SetDefault(lambdacontext.NewLogger())
	lambda.Start(handler)
}

func toAttributeValue(attributeValue events.DynamoDBAttributeValue) ddbtypes.AttributeValue {
	switch attributeValue.DataType() {
	case events.DataTypeString:
		return &ddbtypes.AttributeValueMemberS{Value: attributeValue.String()}
	case events.DataTypeNumber:
		return &ddbtypes.AttributeValueMemberN{Value: attributeValue.Number()}
	case events.DataTypeBoolean:
		return &ddbtypes.AttributeValueMemberBOOL{Value: attributeValue.Boolean()}
	case events.DataTypeBinary:
		return &ddbtypes.AttributeValueMemberB{Value: attributeValue.Binary()}
	case events.DataTypeNull:
		return &ddbtypes.AttributeValueMemberNULL{Value: true}
	case events.DataTypeList:
		list := attributeValue.List()
		out := make([]ddbtypes.AttributeValue, len(list))
		for i, value := range list {
			out[i] = toAttributeValue(value)
		}
		return &ddbtypes.AttributeValueMemberL{Value: out}
	case events.DataTypeMap:
		return &ddbtypes.AttributeValueMemberM{Value: toAttributeValueMap(attributeValue.Map())}
	case events.DataTypeStringSet:
		return &ddbtypes.AttributeValueMemberSS{Value: attributeValue.StringSet()}
	case events.DataTypeNumberSet:
		return &ddbtypes.AttributeValueMemberNS{Value: attributeValue.NumberSet()}
	case events.DataTypeBinarySet:
		return &ddbtypes.AttributeValueMemberBS{Value: attributeValue.BinarySet()}
	}
	return nil
}

func toAttributeValueMap(m map[string]events.DynamoDBAttributeValue) map[string]ddbtypes.AttributeValue {
	out := make(map[string]ddbtypes.AttributeValue, len(m))
	for k, v := range m {
		out[k] = toAttributeValue(v)
	}
	return out
}
