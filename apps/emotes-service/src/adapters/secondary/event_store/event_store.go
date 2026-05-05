package event_store

import (
	"context"
	"crypto/rand"
	"emotes-service/src/adapters/secondary/twitch"
	"emotes-service/src/environment"
	"encoding/hex"
	"fmt"
	"log"
	"log/slog"
	"sort"
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

const GlobalEmotesAggregateId = "GLOBAL"

type EmoteServiceEventEmoteImages struct {
	URL1X string `dynamodbav:"url_1x"`
	URL2X string `dynamodbav:"url_2x"`
	URL4X string `dynamodbav:"url_4x"`
}

type EmoteServiceEventEmote struct {
	Format    []string                     `dynamodbav:"format"`
	ID        string                       `dynamodbav:"id"`
	Images    EmoteServiceEventEmoteImages `dynamodbav:"images"`
	Name      string                       `dynamodbav:"name"`
	Scale     []string                     `dynamodbav:"scale"`
	ThemeMode []string                     `dynamodbav:"theme_mode"`
}

type EmoteServiceEvent struct {
	PK          string                  `dynamodbav:"PK"`
	SK          string                  `dynamodbav:"SK"`
	AggregateId string                  `dynamodbav:"aggregateId"`
	CreatedAt   string                  `dynamodbav:"__createdAt"`
	Emote       *EmoteServiceEventEmote `dynamodbav:"emote"`
	EmoteId     string                  `dynamodbav:"emoteId"`
	EventName   string                  `dynamodbav:"eventName"`
	Id          string                  `dynamodbav:"id"`
	Kind        string                  `dynamodbav:"kind"`
	Sequence    int                     `dynamodbav:"sequence"`
}

type EmotesAggregate struct {
	/* A set of current emote IDs */
	State          map[string]struct{}
	LatestSequence int
}

func LoadAggregate(ctx context.Context, aggregateId string) (*EmotesAggregate, error) {
	events, err := loadEvents(ctx, aggregateId)

	if err != nil {
		return nil, err
	}

	aggregate := &EmotesAggregate{
		State: make(map[string]struct{}),
	}

	for _, event := range events {
		switch event.EventName {
		case "EmoteAdded":
			aggregate.State[event.EmoteId] = struct{}{}
		case "EmoteRemoved":
			delete(aggregate.State, event.EmoteId)
		}
		aggregate.LatestSequence = event.Sequence
	}

	return aggregate, nil
}

func loadEvents(ctx context.Context, aggregateId string) ([]EmoteServiceEvent, error) {
	paginator := dynamodb.NewQueryPaginator(client, &dynamodb.QueryInput{
		TableName:              aws.String(environment.GetOrFatal("TWITCH_EMOTES_EVENT_STORE_TABLE")),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: aggregateId},
			":sk": &types.AttributeValueMemberS{Value: "SEQUENCE#"},
		},
		ConsistentRead: aws.Bool(true),
	})

	var items []EmoteServiceEvent
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			slog.Error("Error querying", "aggregateId", aggregateId, "error", err)
			return nil, err
		}

		var pageItems []EmoteServiceEvent
		if err := attributevalue.UnmarshalListOfMaps(page.Items, &pageItems); err != nil {
			slog.Error("Error unmarshalling page", "aggregateId", aggregateId, "error", err)
			return nil, err
		}
		items = append(items, pageItems...)
	}

	return items, nil
}

func AppendEvents(ctx context.Context, events []EmoteServiceEvent) error {
	eventsLength := len(events)
	if eventsLength == 0 {
		slog.InfoContext(ctx, "No events to append")
		return nil
	}

	slog.InfoContext(ctx, "Events to append", "length", eventsLength)

	table := environment.GetOrFatal("TWITCH_EMOTES_EVENT_STORE_TABLE")

	const transactWriteMaxItems = 100
	for start := 0; start < eventsLength; start += transactWriteMaxItems {
		end := min(start+transactWriteMaxItems, eventsLength)
		chunk := events[start:end]

		items := make([]types.TransactWriteItem, 0, len(chunk))
		for _, event := range chunk {
			marshalledEvent, err := attributevalue.MarshalMap(event)
			if err != nil {
				slog.ErrorContext(ctx, "Error marshalling event", "error", err, "event", event)
				return err
			}
			items = append(items, types.TransactWriteItem{
				Put: &types.Put{
					TableName:           aws.String(table),
					Item:                marshalledEvent,
					ConditionExpression: aws.String("attribute_not_exists(SK)"),
				},
			})
		}

		_, err := client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: items,
		})

		if err != nil {
			slog.ErrorContext(ctx, "AppendEvents failed",
				"chunkStart", start,
				"chunkEnd", end,
				"items", items,
				"error", err,
			)
			return err
		}
	}
	return nil
}

func DecideSyncEvents(aggregateId string, emotesAggregate *EmotesAggregate, globalEmotes []twitch.GlobalEmote) []EmoteServiceEvent {
	twitchState := make(map[string]twitch.GlobalEmote, len(globalEmotes))
	for _, emote := range globalEmotes {
		twitchState[emote.ID] = emote
	}

	var removedIds []string
	for id := range emotesAggregate.State {
		if _, exists := twitchState[id]; !exists {
			removedIds = append(removedIds, id)
		}
	}
	sort.Strings(removedIds)

	var addedIds []string
	for id := range twitchState {
		if _, exists := emotesAggregate.State[id]; !exists {
			addedIds = append(addedIds, id)
		}
	}
	sort.Strings(addedIds)

	currentSequence := emotesAggregate.LatestSequence
	events := make([]EmoteServiceEvent, 0, len(removedIds)+len(addedIds))

	for _, id := range removedIds {
		currentSequence++
		events = append(events, createEmoteEvent(createEmoteEventInput{
			AggregateId: aggregateId,
			Sequence:    currentSequence,
			CreatedAt:   time.Now().UTC().Format(time.RFC3339Nano),
			EventName:   "EmoteRemoved",
			EmoteId:     id,
			Emote:       nil,
		}))
	}

	for _, id := range addedIds {
		currentSequence++
		src := twitchState[id]
		emote := &EmoteServiceEventEmote{
			Format:    src.Format,
			ID:        src.ID,
			Name:      src.Name,
			Scale:     src.Scale,
			ThemeMode: src.ThemeMode,
			Images: EmoteServiceEventEmoteImages{
				URL1X: src.Images.URL1X,
				URL2X: src.Images.URL2X,
				URL4X: src.Images.URL4X,
			},
		}

		events = append(events, createEmoteEvent(createEmoteEventInput{
			AggregateId: aggregateId,
			Sequence:    currentSequence,
			CreatedAt:   time.Now().UTC().Format(time.RFC3339Nano),
			EventName:   "EmoteAdded",
			EmoteId:     id,
			Emote:       emote,
		}))
	}

	return events
}

func generateEventSequence(n int) string {
	return fmt.Sprintf("%07d", n)
}

func generateId() string {
	bytes := make([]byte, 12)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("es_%s", hex.EncodeToString(bytes))
}

type createEmoteEventInput struct {
	AggregateId string
	Sequence    int
	CreatedAt   string
	EventName   string
	EmoteId     string
	Emote       *EmoteServiceEventEmote
}

func createEmoteEvent(eventInput createEmoteEventInput) EmoteServiceEvent {
	pk := eventInput.AggregateId
	sk := fmt.Sprintf("SEQUENCE#%s", generateEventSequence(eventInput.Sequence))
	return EmoteServiceEvent{
		PK:          pk,
		SK:          sk,
		AggregateId: eventInput.AggregateId,
		CreatedAt:   eventInput.CreatedAt,
		Emote:       eventInput.Emote,
		EmoteId:     eventInput.EmoteId,
		EventName:   eventInput.EventName,
		Id:          generateId(),
		Kind:        "EVENT",
		Sequence:    eventInput.Sequence,
	}
}
