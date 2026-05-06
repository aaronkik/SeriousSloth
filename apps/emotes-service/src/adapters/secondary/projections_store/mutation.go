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

type ProjectionItem struct {
	/* The Aggregate ID */
	PK string `dynamodbav:"PK"`
	/* EMOTE#<EMOTE_ID> */
	SK string `dynamodbav:"SK"`
	/* Status can be ACTIVE OR REMOVED */
	Status            string                              `dynamodbav:"status"`
	Id                string                              `dynamodbav:"id"`
	EmoteId           string                              `dynamodbav:"emoteId"`
	RemovedAt         *string                             `dynamodbav:"removedAt"`
	Emote             *event_store.EmoteServiceEventEmote `dynamodbav:"emote"`
	LastEventSequence int                                 `dynamodbav:"__lastEventSequence"`
	CreatedAt         string                              `dynamodbav:"__createdAt"`
	UpdatedAt         string                              `dynamodbav:"__updatedAt"`
	UpdatedBy         string                              `dynamodbav:"__updatedBy"`
}

func buildProjectionUpdate(ctx context.Context, emoteEvent event_store.EmoteServiceEvent) (*dynamodb.UpdateItemInput, error) {
	tableName := environment.GetOrFatal("EVENTS_PROJECTION_TABLE_NAME")
	now := time.Now().UTC().Format(time.RFC3339Nano)
	seq := strconv.Itoa(emoteEvent.Sequence)

	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: emoteEvent.AggregateId},
		"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("EMOTE#%s", emoteEvent.EmoteId)},
	}

	switch emoteEvent.EventName {
	case "EmoteAdded":
		emoteAttr, err := attributevalue.Marshal(emoteEvent.Emote)
		if err != nil {
			slog.ErrorContext(ctx, "marshal emote failed", "error", err, "event", emoteEvent)
			return nil, err
		}

		return &dynamodb.UpdateItemInput{
			TableName: aws.String(tableName),
			Key:       key,
			UpdateExpression: aws.String("SET " +
				"#status = :active, " +
				"#id = if_not_exists(#id, :id), " +
				"#emoteId = :emoteId, " +
				"#emote = :emote, " +
				"#removedAt = :null, " +
				"#createdAt = if_not_exists(#createdAt, :now), " +
				"#updatedAt = :now, " +
				"#updatedBy = :eventId, " +
				"#lastSeq = :seq"),
			ConditionExpression: aws.String("attribute_not_exists(#lastSeq) OR #lastSeq < :seq"),
			ExpressionAttributeNames: map[string]string{
				"#status":    "status",
				"#id":        "id",
				"#emoteId":   "emoteId",
				"#emote":     "emote",
				"#removedAt": "removedAt",
				"#createdAt": "__createdAt",
				"#updatedAt": "__updatedAt",
				"#updatedBy": "__updatedBy",
				"#lastSeq":   "__lastEventSequence",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":active":  &types.AttributeValueMemberS{Value: "ACTIVE"},
				":id":      &types.AttributeValueMemberS{Value: generateId()},
				":emoteId": &types.AttributeValueMemberS{Value: emoteEvent.EmoteId},
				":emote":   emoteAttr,
				":null":    &types.AttributeValueMemberNULL{Value: true},
				":now":     &types.AttributeValueMemberS{Value: now},
				":eventId": &types.AttributeValueMemberS{Value: emoteEvent.Id},
				":seq":     &types.AttributeValueMemberN{Value: seq},
			},
		}, nil

	case "EmoteRemoved":
		return &dynamodb.UpdateItemInput{
			TableName: aws.String(tableName),
			Key:       key,
			UpdateExpression: aws.String("SET " +
				"#status = :removed, " +
				"#removedAt = :now, " +
				"#updatedAt = :now, " +
				"#updatedBy = :eventId, " +
				"#lastSeq = :seq"),
			ConditionExpression: aws.String("attribute_not_exists(#lastSeq) OR #lastSeq < :seq"),
			ExpressionAttributeNames: map[string]string{
				"#status":    "status",
				"#removedAt": "removedAt",
				"#updatedAt": "__updatedAt",
				"#updatedBy": "__updatedBy",
				"#lastSeq":   "__lastEventSequence",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":removed": &types.AttributeValueMemberS{Value: "REMOVED"},
				":now":     &types.AttributeValueMemberS{Value: now},
				":eventId": &types.AttributeValueMemberS{Value: emoteEvent.Id},
				":seq":     &types.AttributeValueMemberN{Value: seq},
			},
		}, nil
	}

	slog.ErrorContext(ctx, "Unhandled event", "eventName", emoteEvent.EventName, "event", emoteEvent)
	return nil, fmt.Errorf("unhandled event: %s", emoteEvent.EventName)
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
	update, err := buildProjectionUpdate(ctx, emoteEvent)
	if err != nil {
		return err
	}

	_, err = client.UpdateItem(ctx, update)
	if err != nil {
		if _, ok := errors.AsType[*types.ConditionalCheckFailedException](err); ok {
			slog.InfoContext(ctx, "idempotent skip",
				"eventId", emoteEvent.Id,
				"sequence", emoteEvent.Sequence,
				"aggregateId", emoteEvent.AggregateId,
				"emoteId", emoteEvent.EmoteId,
				"eventName", emoteEvent.EventName,
			)
			return nil
		}
		slog.ErrorContext(ctx, "projection update failed",
			"eventId", emoteEvent.Id,
			"sequence", emoteEvent.Sequence,
			"aggregateId", emoteEvent.AggregateId,
			"eventName", emoteEvent.EventName,
			"error", err,
		)
		return err
	}

	return nil
}
