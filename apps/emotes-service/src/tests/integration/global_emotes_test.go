//go:build integration

package integration

import (
	"context"
	"emotes-service/src/tests/helpers"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	awslambda "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()
var ddbClient *dynamodb.Client
var lambdaClient *awslambda.Client

var syncLambdaName string
var emotesEventStoreTableName string
var emotesProjectionsTable string
var mockTwitchResponsesTableName string

var apiInvokeUrl string
var apiKey string
var httpClient = &http.Client{Timeout: 10 * time.Second}

var testStartTime = time.Now()
var emoteAddedEventId string

type apiEmote struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Format    []string `json:"format"`
	Scale     []string `json:"scale"`
	ThemeMode []string `json:"theme_mode"`
	Images    struct {
		URL1X string `json:"url_1x"`
		URL2X string `json:"url_2x"`
		URL4X string `json:"url_4x"`
	} `json:"images"`
}

type apiActiveEmote struct {
	Emote   apiEmote `json:"emote"`
	AddedAt string   `json:"addedAt"`
}

type apiRemovedEmote struct {
	Emote     apiEmote `json:"emote"`
	RemovedAt string   `json:"removedAt"`
}

func Test_Emote_Added_Event_Is_Added_To_Event_Store_When_Emote_Exists(t *testing.T) {
	require := require.New(t)

	syncLambdaName = helpers.GetPulumiExport(t, "syncGlobalEmotesLambdaName")
	emotesEventStoreTableName = helpers.GetPulumiExport(t, "twitchEmotesEventStoreTable")
	mockTwitchResponsesTableName = helpers.GetPulumiExport(t, "mockTwitchResponsesTableName")
	emotesProjectionsTable = helpers.GetPulumiExport(t, "twitchEmotesProjectionsTable")

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		t.Fatalf("failed to load AWS config: %v", err)
	}

	ddbClient = dynamodb.NewFromConfig(cfg)
	lambdaClient = awslambda.NewFromConfig(cfg)

	seedMockResponses(t, ctx, ddbClient, mockTwitchResponsesTableName)
	clearTable(t, ctx, ddbClient, emotesEventStoreTableName)
	clearTable(t, ctx, ddbClient, emotesProjectionsTable)

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

	emoteAddedEvent := queryOutput.Items[0]
	emoteAddedEventId = emoteAddedEvent["id"].(*types.AttributeValueMemberS).Value

	require.Equalf(emoteAddedEvent["SK"].(*types.AttributeValueMemberS).Value, "SEQUENCE#0000001", "Event is out of sequence")
	require.Regexp(`^es_.+$`, emoteAddedEventId)

	require.Equal("GLOBAL", emoteAddedEvent["aggregateId"].(*types.AttributeValueMemberS).Value)
	require.Equal("1", emoteAddedEvent["emoteId"].(*types.AttributeValueMemberS).Value)
	require.Equal("EmoteAdded", emoteAddedEvent["eventName"].(*types.AttributeValueMemberS).Value)
	require.Equal("EVENT", emoteAddedEvent["kind"].(*types.AttributeValueMemberS).Value)
	require.Equal("1", emoteAddedEvent["sequence"].(*types.AttributeValueMemberN).Value)

	createdAt, err := time.Parse(time.RFC3339Nano, emoteAddedEvent["__createdAt"].(*types.AttributeValueMemberS).Value)
	require.NoError(err)
	require.True(testStartTime.Before(createdAt))

	assertDDBEmote(t, require, emoteAddedEvent["emote"].(*types.AttributeValueMemberM))
}

func Test_Added_Emote_Is_Added_To_Projection(t *testing.T) {
	require := require.New(t)

	var projectionQueryOutput *dynamodb.QueryOutput

	require.EventuallyWithT(func(c *assert.CollectT) {
		queryOutput, err := ddbClient.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(emotesProjectionsTable),
			KeyConditionExpression: aws.String("PK = :pk"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: "GLOBAL"},
			},
		})
		projectionQueryOutput = queryOutput

		if err != nil {
			t.Fatalf("failed to query DynamoDB: %v", err)
		}

		assert.Len(c, projectionQueryOutput.Items, 1, len(projectionQueryOutput.Items))
	}, 60*time.Second, 2*time.Second)

	emoteProjection := projectionQueryOutput.Items[0]

	require.Equal("EMOTE#1", emoteProjection["SK"].(*types.AttributeValueMemberS).Value)
	require.Regexp(`^es_prj_.+$`, emoteProjection["id"].(*types.AttributeValueMemberS).Value)

	require.Equal("ACTIVE", emoteProjection["status"].(*types.AttributeValueMemberS).Value)
	require.Equal("1", emoteProjection["emoteId"].(*types.AttributeValueMemberS).Value)
	require.Equal(true, emoteProjection["removedAt"].(*types.AttributeValueMemberNULL).Value)
	require.Equal("1", emoteProjection["__lastEventSequence"].(*types.AttributeValueMemberN).Value)

	require.Equal(emoteAddedEventId, emoteProjection["__updatedBy"].(*types.AttributeValueMemberS).Value)

	createdAt, err := time.Parse(time.RFC3339Nano, emoteProjection["__createdAt"].(*types.AttributeValueMemberS).Value)
	require.NoError(err)
	require.True(testStartTime.Before(createdAt))

	updatedAt, err := time.Parse(time.RFC3339Nano, emoteProjection["__updatedAt"].(*types.AttributeValueMemberS).Value)
	require.NoError(err)
	require.True(testStartTime.Before(updatedAt))

	assertDDBEmote(t, require, emoteProjection["emote"].(*types.AttributeValueMemberM))
}

func Test_Api_Returns_Active_Emote_For_Channel(t *testing.T) {
	require := require.New(t)

	apiInvokeUrl = helpers.GetPulumiExport(t, "apiInvokeUrl")
	apiKeyId := helpers.GetPulumiExport(t, "apiKeyId")

	apigwClient := apigateway.NewFromConfig(loadAwsConfig(t))
	keyOutput, err := apigwClient.GetApiKey(ctx, &apigateway.GetApiKeyInput{
		ApiKey:       aws.String(apiKeyId),
		IncludeValue: aws.Bool(true),
	})
	require.NoError(err)
	require.NotNil(keyOutput.Value)
	apiKey = *keyOutput.Value

	var parsed []apiActiveEmote
	require.EventuallyWithT(func(c *assert.CollectT) {
		body, status, err := getEmotes(ctx, apiInvokeUrl, apiKey, "GLOBAL")
		if !assert.NoError(c, err) {
			return
		}
		if !assert.Equal(c, http.StatusOK, status) {
			return
		}

		parsed = nil
		if !assert.NoError(c, json.Unmarshal(body, &parsed)) {
			return
		}
		assert.Len(c, parsed, 1)
	}, 30*time.Second, 2*time.Second)

	require.Len(parsed, 1)
	emote := parsed[0]
	require.Equal("1", emote.Emote.Id)
	require.Equal(":)", emote.Emote.Name)
	require.Equal([]string{"static"}, emote.Emote.Format)
	require.Equal([]string{"1.0", "2.0", "3.0"}, emote.Emote.Scale)
	require.Equal([]string{"light", "dark"}, emote.Emote.ThemeMode)
	require.Equal("https://static-cdn.jtvnw.net/emoticons/v2/1/static/light/1.0", emote.Emote.Images.URL1X)
	require.Equal("https://static-cdn.jtvnw.net/emoticons/v2/1/static/light/2.0", emote.Emote.Images.URL2X)
	require.Equal("https://static-cdn.jtvnw.net/emoticons/v2/1/static/light/3.0", emote.Emote.Images.URL4X)

	addedAt, err := time.Parse(time.RFC3339Nano, emote.AddedAt)
	require.NoError(err)
	require.True(testStartTime.Before(addedAt))
}

func Test_Api_Returns_Forbidden_Without_Api_Key(t *testing.T) {
	require := require.New(t)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiInvokeUrl+"/emotes/GLOBAL", nil)
	require.NoError(err)

	resp, err := httpClient.Do(req)
	require.NoError(err)
	defer resp.Body.Close()

	require.Equal(http.StatusForbidden, resp.StatusCode)
}

var emoteRemovedEventId string

func Test_Emote_Remove_Event_Is_Added_To_Event_Store_When_Emote_No_Longer_Exists(t *testing.T) {
	require := require.New(t)

	updateGlobalEmotesResponse(t, ctx, ddbClient, mockTwitchResponsesTableName, "../fixtures/twitch-global-emotes-empty-response.json")
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

	require.Len(queryOutput.Items, 2)
	emoteRemovedEvent := queryOutput.Items[1]
	emoteRemovedEventId = emoteRemovedEvent["id"].(*types.AttributeValueMemberS).Value

	require.Equalf(emoteRemovedEvent["SK"].(*types.AttributeValueMemberS).Value, "SEQUENCE#0000002", "Event is out of sequence")

	require.Equal("GLOBAL", emoteRemovedEvent["aggregateId"].(*types.AttributeValueMemberS).Value)
	require.Equal("1", emoteRemovedEvent["emoteId"].(*types.AttributeValueMemberS).Value)
	require.Equal("EmoteRemoved", emoteRemovedEvent["eventName"].(*types.AttributeValueMemberS).Value)
	require.Equal("EVENT", emoteRemovedEvent["kind"].(*types.AttributeValueMemberS).Value)
	require.Equal("2", emoteRemovedEvent["sequence"].(*types.AttributeValueMemberN).Value)

	emoteRemovedEventEmote := emoteRemovedEvent["emote"].(*types.AttributeValueMemberNULL).Value
	require.True(emoteRemovedEventEmote)
}

func Test_Emote_Remove_Event_Is_Reflected_In_Projection(t *testing.T) {
	require := require.New(t)

	var projectionQueryOutput *dynamodb.QueryOutput

	require.EventuallyWithT(func(c *assert.CollectT) {
		queryOutput, err := ddbClient.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(emotesProjectionsTable),
			KeyConditionExpression: aws.String("PK = :pk"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: "GLOBAL"},
			},
		})
		projectionQueryOutput = queryOutput

		if err != nil {
			t.Fatalf("failed to query DynamoDB: %v", err)
		}

		if !assert.Len(c, projectionQueryOutput.Items, 1) {
			return
		}

		emoteProjection := projectionQueryOutput.Items[0]
		assert.Equal(c, "REMOVED", emoteProjection["status"].(*types.AttributeValueMemberS).Value)
		assert.Equal(c, "2", emoteProjection["__lastEventSequence"].(*types.AttributeValueMemberN).Value)
	}, 60*time.Second, 2*time.Second)

	emoteProjection := projectionQueryOutput.Items[0]

	require.Equal("EMOTE#1", emoteProjection["SK"].(*types.AttributeValueMemberS).Value)
	require.Regexp(`^es_prj_.+$`, emoteProjection["id"].(*types.AttributeValueMemberS).Value)

	require.Equal("REMOVED", emoteProjection["status"].(*types.AttributeValueMemberS).Value)
	require.Equal("1", emoteProjection["emoteId"].(*types.AttributeValueMemberS).Value)
	require.Equal("2", emoteProjection["__lastEventSequence"].(*types.AttributeValueMemberN).Value)

	removedAt, err := time.Parse(time.RFC3339Nano, emoteProjection["removedAt"].(*types.AttributeValueMemberS).Value)
	require.NoError(err)
	require.True(testStartTime.Before(removedAt))

	require.Equal(emoteRemovedEventId, emoteProjection["__updatedBy"].(*types.AttributeValueMemberS).Value)

	createdAt, err := time.Parse(time.RFC3339Nano, emoteProjection["__createdAt"].(*types.AttributeValueMemberS).Value)
	require.NoError(err)

	updatedAt, err := time.Parse(time.RFC3339Nano, emoteProjection["__updatedAt"].(*types.AttributeValueMemberS).Value)
	require.NoError(err)
	require.True(createdAt.Before(updatedAt))

	assertDDBEmote(t, require, emoteProjection["emote"].(*types.AttributeValueMemberM))
}

func Test_Api_Returns_No_Emotes_When_Emote_Has_Been_Removed(t *testing.T) {
	require := require.New(t)

	require.EventuallyWithT(func(c *assert.CollectT) {
		body, status, err := getEmotes(ctx, apiInvokeUrl, apiKey, "GLOBAL")
		if !assert.NoError(c, err) {
			return
		}
		if !assert.Equal(c, http.StatusOK, status) {
			return
		}

		var parsed []apiActiveEmote
		if !assert.NoError(c, json.Unmarshal(body, &parsed)) {
			return
		}
		assert.Empty(c, parsed)
	}, 30*time.Second, 2*time.Second)
}

func Test_Api_Returns_Removed_Emote_For_Channel(t *testing.T) {
	require := require.New(t)

	var parsed []apiRemovedEmote
	require.EventuallyWithT(func(c *assert.CollectT) {
		body, status, err := getRemovedEmotes(ctx, apiInvokeUrl, apiKey, "GLOBAL")
		if !assert.NoError(c, err) {
			return
		}
		if !assert.Equal(c, http.StatusOK, status) {
			return
		}

		parsed = nil
		if !assert.NoError(c, json.Unmarshal(body, &parsed)) {
			return
		}
		assert.Len(c, parsed, 1)
	}, 30*time.Second, 2*time.Second)

	require.Len(parsed, 1)
	emote := parsed[0]
	require.Equal("1", emote.Emote.Id)
	require.Equal(":)", emote.Emote.Name)
	require.Equal([]string{"static"}, emote.Emote.Format)
	require.Equal([]string{"1.0", "2.0", "3.0"}, emote.Emote.Scale)
	require.Equal([]string{"light", "dark"}, emote.Emote.ThemeMode)
	require.Equal("https://static-cdn.jtvnw.net/emoticons/v2/1/static/light/1.0", emote.Emote.Images.URL1X)
	require.Equal("https://static-cdn.jtvnw.net/emoticons/v2/1/static/light/2.0", emote.Emote.Images.URL2X)
	require.Equal("https://static-cdn.jtvnw.net/emoticons/v2/1/static/light/3.0", emote.Emote.Images.URL4X)

	removedAt, err := time.Parse(time.RFC3339Nano, emote.RemovedAt)
	require.NoError(err)
	require.True(testStartTime.Before(removedAt))
}

func loadAwsConfig(t *testing.T) aws.Config {
	t.Helper()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		t.Fatalf("failed to load AWS config: %v", err)
	}
	return cfg
}

func getEmotes(ctx context.Context, baseUrl, apiKey, channelId string) ([]byte, int, error) {
	return doGet(ctx, baseUrl+"/emotes/"+channelId, apiKey)
}

func getRemovedEmotes(ctx context.Context, baseUrl, apiKey, channelId string) ([]byte, int, error) {
	return doGet(ctx, baseUrl+"/emotes/"+channelId+"/removed", apiKey)
}

func doGet(ctx context.Context, url, apiKey string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("x-api-key", apiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return body, resp.StatusCode, nil
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

	t.Logf("cleared table %s", tableName)
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

	t.Logf("seeded mock responses to table %s", tableName)
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

func assertDDBEmote(t *testing.T, require *require.Assertions, emoteAttr *types.AttributeValueMemberM) {
	t.Helper()

	emote := emoteAttr.Value
	require.Equal("1", emote["id"].(*types.AttributeValueMemberS).Value)
	require.Equal(":)", emote["name"].(*types.AttributeValueMemberS).Value)

	format := emote["format"].(*types.AttributeValueMemberL).Value
	require.Len(format, 1)
	require.Equal("static", format[0].(*types.AttributeValueMemberS).Value)

	scale := emote["scale"].(*types.AttributeValueMemberL).Value
	require.Len(scale, 3)
	require.Equal("1.0", scale[0].(*types.AttributeValueMemberS).Value)
	require.Equal("2.0", scale[1].(*types.AttributeValueMemberS).Value)
	require.Equal("3.0", scale[2].(*types.AttributeValueMemberS).Value)

	themeMode := emote["theme_mode"].(*types.AttributeValueMemberL).Value
	require.Len(themeMode, 2)
	require.Equal("light", themeMode[0].(*types.AttributeValueMemberS).Value)
	require.Equal("dark", themeMode[1].(*types.AttributeValueMemberS).Value)

	images := emote["images"].(*types.AttributeValueMemberM).Value
	require.Equal("https://static-cdn.jtvnw.net/emoticons/v2/1/static/light/1.0", images["url_1x"].(*types.AttributeValueMemberS).Value)
	require.Equal("https://static-cdn.jtvnw.net/emoticons/v2/1/static/light/2.0", images["url_2x"].(*types.AttributeValueMemberS).Value)
	require.Equal("https://static-cdn.jtvnw.net/emoticons/v2/1/static/light/3.0", images["url_4x"].(*types.AttributeValueMemberS).Value)

}
