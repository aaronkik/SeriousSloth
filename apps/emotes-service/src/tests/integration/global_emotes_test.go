//go:build integration

package integration

import (
	"context"
	"emotes-service/src/tests/helpers"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	awslambda "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobalEmotes(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	syncLambdaName := helpers.GetPulumiExport(t, "syncGlobalEmotesLambdaName")
	emotesEventStoreTableName := helpers.GetPulumiExport(t, "twitchEmotesEventStoreTable")
	mockTwitchResponsesTableName := helpers.GetPulumiExport(t, "mockTwitchResponsesTableName")

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		t.Fatalf("failed to load AWS config: %v", err)
	}

	ddbClient := dynamodb.NewFromConfig(cfg)
	lambdaClient := awslambda.NewFromConfig(cfg)

	seedMockResponses(t, ctx, ddbClient, mockTwitchResponsesTableName)
	clearTable(t, ctx, ddbClient, emotesEventStoreTableName)

	triggerLambda(t, ctx, lambdaClient, syncLambdaName)

	queryOutput, err := ddbClient.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(emotesEventStoreTableName),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "GLOBAL"},
		},
	})
	if err != nil {
		t.Fatalf("failed to query DynamoDB: %v", err)
	}

	require.Lenf(queryOutput.Items, 1, "expected 1 item, got %d", len(queryOutput.Items))

	emoteServiceEvent := queryOutput.Items[0]

	assert.Equalf(emoteServiceEvent["SK"].(*types.AttributeValueMemberS).Value, "SEQUENCE#0000001", "Event is out of sequence")
	assert.Regexp(`^es_.+$`, emoteServiceEvent["id"].(*types.AttributeValueMemberS).Value)

	assert.Equal("GLOBAL", emoteServiceEvent["aggregateId"].(*types.AttributeValueMemberS).Value)
	assert.Equal("1", emoteServiceEvent["emoteId"].(*types.AttributeValueMemberS).Value)
	assert.Equal("GLOBAL#SEQUENCE#0000001", emoteServiceEvent["eventId"].(*types.AttributeValueMemberS).Value)
	assert.Equal("EmoteAdded", emoteServiceEvent["eventName"].(*types.AttributeValueMemberS).Value)
	assert.Equal("EVENT", emoteServiceEvent["kind"].(*types.AttributeValueMemberS).Value)
	assert.Equal("1", emoteServiceEvent["sequence"].(*types.AttributeValueMemberN).Value)

	emote := emoteServiceEvent["emote"].(*types.AttributeValueMemberM).Value
	assert.Equal("1", emote["id"].(*types.AttributeValueMemberS).Value)
	assert.Equal(":)", emote["name"].(*types.AttributeValueMemberS).Value)

	format := emote["format"].(*types.AttributeValueMemberL).Value
	require.Len(format, 1)
	assert.Equal("static", format[0].(*types.AttributeValueMemberS).Value)

	scale := emote["scale"].(*types.AttributeValueMemberL).Value
	require.Len(scale, 3)
	assert.Equal("1.0", scale[0].(*types.AttributeValueMemberS).Value)
	assert.Equal("2.0", scale[1].(*types.AttributeValueMemberS).Value)
	assert.Equal("3.0", scale[2].(*types.AttributeValueMemberS).Value)

	themeMode := emote["theme_mode"].(*types.AttributeValueMemberL).Value
	require.Len(themeMode, 2)
	assert.Equal("light", themeMode[0].(*types.AttributeValueMemberS).Value)
	assert.Equal("dark", themeMode[1].(*types.AttributeValueMemberS).Value)

	images := emote["images"].(*types.AttributeValueMemberM).Value
	assert.Equal("https://static-cdn.jtvnw.net/emoticons/v2/1/static/light/1.0", images["url_1x"].(*types.AttributeValueMemberS).Value)
	assert.Equal("https://static-cdn.jtvnw.net/emoticons/v2/1/static/light/2.0", images["url_2x"].(*types.AttributeValueMemberS).Value)
	assert.Equal("https://static-cdn.jtvnw.net/emoticons/v2/1/static/light/3.0", images["url_4x"].(*types.AttributeValueMemberS).Value)

	// Remove the emote from the Twitch response
	updateGlobalEmotesResponse(t, ctx, ddbClient, mockTwitchResponsesTableName, "../fixtures/twitch-global-emotes-empty-response.json")
	triggerLambda(t, ctx, lambdaClient, syncLambdaName)

	queryOutput, err = ddbClient.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(emotesEventStoreTableName),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "GLOBAL"},
		},
	})
	if err != nil {
		t.Fatalf("failed to query DynamoDB: %v", err)
	}

	require.Len(queryOutput.Items, 2)
	emoteRemovedEvent := queryOutput.Items[1]

	assert.Equalf(emoteRemovedEvent["SK"].(*types.AttributeValueMemberS).Value, "SEQUENCE#0000002", "Event is out of sequence")

	assert.Equal("GLOBAL", emoteRemovedEvent["aggregateId"].(*types.AttributeValueMemberS).Value)
	assert.Equal("1", emoteRemovedEvent["emoteId"].(*types.AttributeValueMemberS).Value)
	assert.Equal("GLOBAL#SEQUENCE#0000002", emoteRemovedEvent["eventId"].(*types.AttributeValueMemberS).Value)
	assert.Equal("EmoteRemoved", emoteRemovedEvent["eventName"].(*types.AttributeValueMemberS).Value)
	assert.Equal("EVENT", emoteRemovedEvent["kind"].(*types.AttributeValueMemberS).Value)
	assert.Equal("2", emoteRemovedEvent["sequence"].(*types.AttributeValueMemberN).Value)

	emoteRemovedEventEmote := emoteRemovedEvent["emote"].(*types.AttributeValueMemberNULL).Value
	assert.True(emoteRemovedEventEmote)

}

func clearTable(t *testing.T, ctx context.Context, client *dynamodb.Client, tableName string) {
	t.Helper()

	var lastKey map[string]types.AttributeValue
	for {
		scanOutput, err := client.Scan(ctx, &dynamodb.ScanInput{
			TableName:         aws.String(tableName),
			ExclusiveStartKey: lastKey,
		})
		if err != nil {
			t.Fatalf("failed to scan table: %v", err)
		}

		for _, scanned := range scanOutput.Items {
			_, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
				TableName: aws.String(tableName),
				Key: map[string]types.AttributeValue{
					"PK": scanned["PK"],
					"SK": scanned["SK"],
				},
			})
			if err != nil {
				t.Fatalf("failed to delete item: %v", err)
			}
		}

		if scanOutput.LastEvaluatedKey == nil {
			break
		}
		lastKey = scanOutput.LastEvaluatedKey
	}
}

func seedMockResponses(t *testing.T, ctx context.Context, client *dynamodb.Client, tableName string) {
	t.Helper()

	oauthFixture, err := os.ReadFile("../fixtures/oauth2-token-response.json")
	if err != nil {
		t.Fatalf("failed to read oauth fixture: %v", err)
	}

	emotesFixture, err := os.ReadFile("../fixtures/twitch-global-emotes-starting-response.json")
	if err != nil {
		t.Fatalf("failed to read emotes fixture: %v", err)
	}

	items := []struct {
		pk   string
		body string
	}{
		{"/oauth2/token", string(oauthFixture)},
		{"/helix/chat/emotes/global", string(emotesFixture)},
	}

	for _, item := range items {
		_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]types.AttributeValue{
				"PK":           &types.AttributeValueMemberS{Value: item.pk},
				"responseBody": &types.AttributeValueMemberS{Value: item.body},
			},
		})
		if err != nil {
			t.Fatalf("failed to seed mock response PK=%s: %v", item.pk, err)
		}
	}
}

func updateGlobalEmotesResponse(t *testing.T, ctx context.Context, client *dynamodb.Client, tableName string, fixturePath string) {
	t.Helper()

	fixture, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("failed to read emotes fixture: %v", err)
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{
				Value: "/helix/chat/emotes/global",
			},
			"responseBody": &types.AttributeValueMemberS{
				Value: string(fixture),
			},
		},
	})
	if err != nil {
		t.Fatal("failed to put", err)
	}
}

func triggerLambda(t *testing.T, ctx context.Context, client *awslambda.Client, lambdaName string) {
	t.Helper()

	invokeOutput, err := client.Invoke(ctx, &awslambda.InvokeInput{
		FunctionName: aws.String(lambdaName),
		Payload:      []byte("{}"),
	})

	if err != nil {
		t.Fatalf("failed to invoke Lambda: %v", err)
	}
	if invokeOutput.FunctionError != nil {
		t.Fatalf("Lambda returned error: %s, payload: %s", *invokeOutput.FunctionError, string(invokeOutput.Payload))
	}
}
