//go:build integration

package integration

import (
	"context"
	"emotes-service/src/tests/helpers"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	awslambda "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Channel_Sync_Populates_Event_Store(t *testing.T) {
	require := require.New(t)
	testCtx := context.Background()

	channelsTable := helpers.GetPulumiExport(t, "twitchChannelsTable")
	eventStoreTable := helpers.GetPulumiExport(t, "twitchEmotesEventStoreTable")
	mockTable := helpers.GetPulumiExport(t, "mockTwitchResponsesTableName")
	dispatcherName := helpers.GetPulumiExport(t, "channelSyncDispatcherLambdaName")
	invokeUrl := helpers.GetPulumiExport(t, "apiInvokeUrl")
	keyId := helpers.GetPulumiExport(t, "apiKeyId")

	cfg := loadAwsConfig(t)
	ddb := dynamodb.NewFromConfig(cfg)
	lambdaCli := awslambda.NewFromConfig(cfg)

	// Reuse the shared apigwClient pattern from channels_test.go to fetch the API key value.
	apiKey := fetchApiKey(t, testCtx, cfg, keyId)

	seedMockResponse(t, testCtx, ddb, mockTable, "/oauth2/token", "../fixtures/oauth2-token-response.json")
	seedMockResponse(t, testCtx, ddb, mockTable, "/helix/users?id=100", "../fixtures/twitch-user-100.json")
	seedMockResponse(t, testCtx, ddb, mockTable, "/helix/users?id=200", "../fixtures/twitch-user-200.json")
	seedMockResponse(t, testCtx, ddb, mockTable, "/helix/chat/emotes?broadcaster_id=100", "../fixtures/twitch-channel-emotes-100.json")
	seedMockResponse(t, testCtx, ddb, mockTable, "/helix/chat/emotes?broadcaster_id=200", "../fixtures/twitch-channel-emotes-200.json")

	// Seed both channels so the dispatcher fan-out is deterministic regardless of
	// leftover state from prior `channels_test` runs (which also writes 100/200).
	seedChannelDirect(t, testCtx, ddb, channelsTable, "100", "Alice", "https://cdn.twitch/100.png")
	seedChannelDirect(t, testCtx, ddb, channelsTable, "200", "Bob", "https://cdn.twitch/200.png")

	aggregateId := "CHANNEL#100"
	aggregateId200 := "CHANNEL#200"
	clearAggregateEvents(t, testCtx, ddb, eventStoreTable, aggregateId)
	clearAggregateEvents(t, testCtx, ddb, eventStoreTable, aggregateId200)

	triggerLambda(t, testCtx, lambdaCli, dispatcherName)

	var eventStoreItems []map[string]types.AttributeValue
	require.EventuallyWithT(func(c *assert.CollectT) {
		out, err := ddb.Query(testCtx, &dynamodb.QueryInput{
			TableName:              aws.String(eventStoreTable),
			KeyConditionExpression: aws.String("PK = :pk"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: aggregateId},
			},
		})
		if !assert.NoError(c, err) {
			return
		}
		assert.Len(c, out.Items, 1)
		eventStoreItems = out.Items
	}, 60*time.Second, 2*time.Second)

	require.Len(eventStoreItems, 1)
	event := eventStoreItems[0]
	require.Equal("SEQUENCE#0000001", event["SK"].(*types.AttributeValueMemberS).Value)
	require.Equal(aggregateId, event["aggregateId"].(*types.AttributeValueMemberS).Value)
	require.Equal("EmoteAdded", event["eventName"].(*types.AttributeValueMemberS).Value)
	require.Equal("ch100_e1", event["emoteId"].(*types.AttributeValueMemberS).Value)

	// Read API takes the raw Twitch id (or "global") — the "CHANNEL#" prefix is
	// internal to the event store and never crosses the API boundary.
	var parsed []apiActiveEmote
	require.EventuallyWithT(func(c *assert.CollectT) {
		body, status, err := doGet(testCtx, invokeUrl+"/emotes/100", apiKey)
		if !assert.NoError(c, err) {
			return
		}
		if !assert.Equal(c, http.StatusOK, status, "body: %s", string(body)) {
			return
		}
		parsed = nil
		if !assert.NoError(c, json.Unmarshal(body, &parsed)) {
			return
		}
		assert.Len(c, parsed, 1)
	}, 60*time.Second, 2*time.Second)

	require.Len(parsed, 1)
	require.Equal("ch100_e1", parsed[0].Emote.Id)
	require.Equal("aliceWave", parsed[0].Emote.Name)

	// Fan-out check: the dispatcher also synced channel 200 via the same tick.
	require.EventuallyWithT(func(c *assert.CollectT) {
		out, err := ddb.Query(testCtx, &dynamodb.QueryInput{
			TableName:              aws.String(eventStoreTable),
			KeyConditionExpression: aws.String("PK = :pk"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: aggregateId200},
			},
		})
		if !assert.NoError(c, err) {
			return
		}
		if !assert.Len(c, out.Items, 1) {
			return
		}
		assert.Equal(c, "ch200_e1", out.Items[0]["emoteId"].(*types.AttributeValueMemberS).Value)
	}, 60*time.Second, 2*time.Second)

	// Re-dispatch and assert idempotency: still exactly one event for this aggregate.
	triggerLambda(t, testCtx, lambdaCli, dispatcherName)
	time.Sleep(10 * time.Second)
	out, err := ddb.Query(testCtx, &dynamodb.QueryInput{
		TableName:              aws.String(eventStoreTable),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: aggregateId},
		},
	})
	require.NoError(err)
	require.Len(out.Items, 1, "expected no duplicate events on second dispatch")
}

func fetchApiKey(t *testing.T, ctx context.Context, cfg aws.Config, keyId string) string {
	t.Helper()
	apigwCli := apigateway.NewFromConfig(cfg)
	out, err := apigwCli.GetApiKey(ctx, &apigateway.GetApiKeyInput{
		ApiKey:       aws.String(keyId),
		IncludeValue: aws.Bool(true),
	})
	if err != nil {
		t.Fatalf("failed to fetch API key: %v", err)
	}
	if out.Value == nil {
		t.Fatalf("API key value missing")
	}
	return *out.Value
}

func seedMockResponse(t *testing.T, ctx context.Context, client *dynamodb.Client, tableName, pk, fixturePath string) {
	t.Helper()
	body, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", fixturePath, err)
	}
	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"PK":           &types.AttributeValueMemberS{Value: pk},
			"responseBody": &types.AttributeValueMemberS{Value: string(body)},
		},
	})
	if err != nil {
		t.Fatalf("failed to seed mock PK=%s: %v", pk, err)
	}
}

func seedChannelDirect(t *testing.T, ctx context.Context, client *dynamodb.Client, tableName, twitchId, displayName, imageUrl string) {
	t.Helper()
	now := time.Now().UTC().Format(time.RFC3339Nano)
	_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"PK":          &types.AttributeValueMemberS{Value: "CHANNEL"},
			"SK":          &types.AttributeValueMemberS{Value: twitchId},
			"id":          &types.AttributeValueMemberS{Value: "chnl_seed" + twitchId},
			"twitchId":    &types.AttributeValueMemberS{Value: twitchId},
			"displayName": &types.AttributeValueMemberS{Value: displayName},
			"imageUrl":    &types.AttributeValueMemberS{Value: imageUrl},
			"addedAt":     &types.AttributeValueMemberS{Value: now},
			"updatedAt":   &types.AttributeValueMemberS{Value: now},
		},
	})
	if err != nil {
		t.Fatalf("failed to seed channel %s: %v", twitchId, err)
	}
}

func clearAggregateEvents(t *testing.T, ctx context.Context, client *dynamodb.Client, tableName, aggregateId string) {
	t.Helper()
	out, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: aggregateId},
		},
	})
	if err != nil {
		t.Fatalf("failed to query aggregate events: %v", err)
	}
	for _, item := range out.Items {
		_, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"PK": item["PK"],
				"SK": item["SK"],
			},
		})
		if err != nil {
			t.Fatalf("failed to delete event row: %v", err)
		}
	}
}
