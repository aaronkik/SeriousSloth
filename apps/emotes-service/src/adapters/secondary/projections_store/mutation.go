package projections_store

import (
	"context"
	"crypto/rand"
	"emotes-service/src/adapters/secondary/event_store"
	"emotes-service/src/environment"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var client *dynamodb.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client = dynamodb.NewFromConfig(cfg)
}

type MetadataItem struct {
	/* The Aggregate ID */
	PK string `dynamodbav:"PK"`
	/* Value is METADATA */
	SK              string `dynamodbav:"SK"`
	CurrentSequence int    `dynamodbav:"currentSequence"`
	CreatedAt       string `dynamodbav:"__createdAt"`
	UpdatedAt       string `dynamodbav:"__updatedAt"`
	UpdatedBy       string `dynamodbav:"__updatedBy"`
}

func createMetadataItem(emoteEvent event_store.EmoteServiceEvent) types.TransactWriteItem {
	tableName := environment.GetOrFatal("EVENTS_PROJECTION_TABLE_NAME")
	now := time.Now().UTC().Format(time.RFC3339Nano)

	metadataItem := types.TransactWriteItem{
		Update: &types.Update{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: emoteEvent.AggregateId},
				"SK": &types.AttributeValueMemberS{Value: "METADATA"},
			},
			UpdateExpression:    aws.String("SET #currentSequence = :sequence, #createdAt = if_not_exists(#createdAt, :now), #updatedAt = :now, #updatedBy = :updatedBy"),
			ConditionExpression: aws.String("attribute_not_exists(#currentSequence) OR #currentSequence < :sequence"),
			ExpressionAttributeNames: map[string]string{
				"#currentSequence": "currentSequence",
				"#createdAt":       "__createdAt",
				"#updatedAt":       "__updatedAt",
				"#updatedBy":       "__updatedBy",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":sequence":  &types.AttributeValueMemberN{Value: strconv.Itoa(emoteEvent.Sequence)},
				":now":       &types.AttributeValueMemberS{Value: now},
				":updatedBy": &types.AttributeValueMemberS{Value: emoteEvent.Id},
			},
		},
	}
	return metadataItem
}

type ProjectionItem struct {
	/* The Aggregate ID */
	PK string `dynamodbav:"PK"`
	/* EMOTE#<EMOTE_ID> */
	SK string `dynamodbav:"SK"`
	/* Status can be ACTIVE OR REMOVED */
	Status    string                              `dynamodbav:"status"`
	Id        string                              `dynamodbav:"id"`
	EmoteId   string                              `dynamodbav:"emoteId"`
	RemovedAt *string                             `dynamodbav:"removedAt"`
	Emote     *event_store.EmoteServiceEventEmote `dynamodbav:"emote"`
	CreatedAt string                              `dynamodbav:"__createdAt"`
	UpdatedAt string                              `dynamodbav:"__updatedAt"`
	UpdatedBy string                              `dynamodbav:"__updatedBy"`
}

func createProjectionItem(ctx context.Context, emoteEvent event_store.EmoteServiceEvent) (types.TransactWriteItem, error) {
	tableName := environment.GetOrFatal("EVENTS_PROJECTION_TABLE_NAME")

	switch emoteEvent.EventName {
	case "EmoteAdded":
		{
			now := time.Now().UTC().Format(time.RFC3339Nano)
			projectionItem := ProjectionItem{
				PK:        emoteEvent.AggregateId,
				SK:        fmt.Sprintf("EMOTE#%s", emoteEvent.EmoteId),
				Status:    "ACTIVE",
				Id:        generateId(),
				EmoteId:   emoteEvent.EmoteId,
				RemovedAt: nil,
				Emote:     emoteEvent.Emote,
				CreatedAt: now,
				UpdatedAt: now,
				UpdatedBy: emoteEvent.Id,
			}

			marshalledItem, err := attributevalue.MarshalMap(projectionItem)
			if err != nil {
				slog.ErrorContext(ctx, "Error marshalling item", "error", err, "item", projectionItem)
				log.Fatalf("Error marshalling item")
			}

			transactItem := types.TransactWriteItem{
				Put: &types.Put{
					TableName: aws.String(tableName),
					Item:      marshalledItem,
				},
			}

			return transactItem, nil
		}
	case "EmoteRemoved":
		{
			now := time.Now().UTC().Format(time.RFC3339Nano)
			transactItem := types.TransactWriteItem{
				Update: &types.Update{
					TableName: aws.String(tableName),
					Key: map[string]types.AttributeValue{
						"PK": &types.AttributeValueMemberS{Value: emoteEvent.AggregateId},
						"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("EMOTE#%s", emoteEvent.EmoteId)},
					},
					UpdateExpression: aws.String("SET #status = :status, #removedAt = :removedAt, #updatedAt = :updatedAt, #updatedBy = :updatedBy"),
					ExpressionAttributeNames: map[string]string{
						"#status":    "status",
						"#removedAt": "removedAt",
						"#updatedAt": "__updatedAt",
						"#updatedBy": "__updatedBy",
					},
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":status":    &types.AttributeValueMemberS{Value: "REMOVED"},
						":removedAt": &types.AttributeValueMemberS{Value: now},
						":updatedAt": &types.AttributeValueMemberS{Value: now},
						":updatedBy": &types.AttributeValueMemberS{Value: emoteEvent.Id},
					},
				},
			}

			return transactItem, nil
		}
	}

	slog.ErrorContext(ctx, "Unhandled event", "eventName", emoteEvent.EventName, "event", emoteEvent)
	return types.TransactWriteItem{}, fmt.Errorf("unhandled event: %s", emoteEvent.EventName)
}

func generateId() string {
	bytes := make([]byte, 12)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("es_prj_%s", hex.EncodeToString(bytes))
}

func Persist(ctx context.Context, emoteEvent event_store.EmoteServiceEvent) error {
	projectionItem, err := createProjectionItem(ctx, emoteEvent)
	if err != nil {
		return err
	}

	metadataItem := createMetadataItem(emoteEvent)

	_, err = client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{projectionItem, metadataItem},
	})
	if err != nil {
		if canceled, ok := errors.AsType[*types.TransactionCanceledException](err); ok {
			for _, reason := range canceled.CancellationReasons {
				if reason.Code != nil && *reason.Code == "ConditionalCheckFailed" {
					slog.WarnContext(ctx, "Skipping event",
						"eventId", emoteEvent.Id,
						"sequence", emoteEvent.Sequence,
						"aggregateId", emoteEvent.AggregateId,
						"eventName", emoteEvent.EventName,
						"error", err,
					)
					return nil
				}
			}
		}
		return err
	}

	return nil
}
