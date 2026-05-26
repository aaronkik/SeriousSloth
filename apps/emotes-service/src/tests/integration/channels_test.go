//go:build integration

package integration

import (
	"bytes"
	"context"
	"emotes-service/src/tests/helpers"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/require"
)

type apiChannel struct {
	Id          string `json:"id"`
	TwitchId    string `json:"twitchId"`
	DisplayName string `json:"displayName"`
	ImageUrl    string `json:"imageUrl"`
	AddedAt     string `json:"addedAt"`
	UpdatedAt   string `json:"updatedAt"`
}

func Test_Channels_Api(t *testing.T) {
	req := require.New(t)
	testCtx := context.Background()

	channelsTable := helpers.GetPulumiExport(t, "twitchChannelsTable")
	mockTable := helpers.GetPulumiExport(t, "mockTwitchResponsesTableName")
	invokeUrl := helpers.GetPulumiExport(t, "apiInvokeUrl")
	keyId := helpers.GetPulumiExport(t, "apiKeyId")

	cfg := loadAwsConfig(t)
	ddb := dynamodb.NewFromConfig(cfg)
	apigwClient := apigateway.NewFromConfig(cfg)

	keyOutput, err := apigwClient.GetApiKey(testCtx, &apigateway.GetApiKeyInput{
		ApiKey:       aws.String(keyId),
		IncludeValue: aws.Bool(true),
	})
	req.NoError(err)
	req.NotNil(keyOutput.Value)
	key := *keyOutput.Value

	seedOauthMockResponse(t, testCtx, ddb, mockTable)
	seedChannelsMockResponses(t, testCtx, ddb, mockTable)

	clearTable(t, testCtx, ddb, channelsTable)

	channelsUrl := invokeUrl + "/channels"
	channelUrl := invokeUrl + "/channel"
	testStart := time.Now()

	t.Run("post creates two channels populated from twitch", func(t *testing.T) {
		require := require.New(t)

		cases := []struct {
			twitchId            string
			expectedDisplayName string
			expectedImageUrl    string
		}{
			{"100", "Alice", "https://cdn.twitch/100.png"},
			{"200", "Bob", "https://cdn.twitch/200.png"},
		}

		for _, c := range cases {
			body, status, err := postChannel(testCtx, channelUrl, key, c.twitchId)
			require.NoError(err)
			require.Equalf(http.StatusCreated, status, "body: %s", string(body))

			var created apiChannel
			require.NoError(json.Unmarshal(body, &created))
			require.Regexp(`^chnl_[0-9a-f]{24}$`, created.Id)
			require.Equal(c.twitchId, created.TwitchId)
			require.Equal(c.expectedDisplayName, created.DisplayName)
			require.Equal(c.expectedImageUrl, created.ImageUrl)

			addedAt, err := time.Parse(time.RFC3339Nano, created.AddedAt)
			require.NoError(err)
			require.True(testStart.Before(addedAt))

			updatedAt, err := time.Parse(time.RFC3339Nano, created.UpdatedAt)
			require.NoError(err)
			require.Equal(created.AddedAt, created.UpdatedAt)
			require.True(testStart.Before(updatedAt))
		}
	})

	t.Run("get returns both channels ordered by twitchId", func(t *testing.T) {
		require := require.New(t)

		body, status, err := doGet(testCtx, channelsUrl, key)
		require.NoError(err)
		require.Equalf(http.StatusOK, status, "body: %s", string(body))

		var parsed []apiChannel
		require.NoError(json.Unmarshal(body, &parsed))
		require.Len(parsed, 2)

		ids := []string{parsed[0].TwitchId, parsed[1].TwitchId}
		require.True(sort.StringsAreSorted(ids), "expected SK-ordered results, got %v", ids)
		require.Equal("100", parsed[0].TwitchId)
		require.Equal("Alice", parsed[0].DisplayName)
		require.Equal("https://cdn.twitch/100.png", parsed[0].ImageUrl)
		require.Regexp(`^chnl_[0-9a-f]{24}$`, parsed[0].Id)
		require.Equal("200", parsed[1].TwitchId)
		require.Equal("Bob", parsed[1].DisplayName)
		require.Equal("https://cdn.twitch/200.png", parsed[1].ImageUrl)
		require.Regexp(`^chnl_[0-9a-f]{24}$`, parsed[1].Id)
		require.NotEqual(parsed[0].Id, parsed[1].Id)
	})

	t.Run("post duplicate returns 409", func(t *testing.T) {
		require := require.New(t)

		_, status, err := postChannel(testCtx, channelUrl, key, "100")
		require.NoError(err)
		require.Equal(http.StatusConflict, status)
	})

	t.Run("post with empty twitchId returns 400", func(t *testing.T) {
		require := require.New(t)

		_, status, err := postChannel(testCtx, channelUrl, key, "")
		require.NoError(err)
		require.Equal(http.StatusBadRequest, status)
	})

	t.Run("post with unknown twitchId returns 404", func(t *testing.T) {
		require := require.New(t)

		_, status, err := postChannel(testCtx, channelUrl, key, "999")
		require.NoError(err)
		require.Equal(http.StatusNotFound, status)
	})
}

// seedOauthMockResponse ensures the mock has an OAuth token response, since
// add-channel needs to fetch a token before calling /helix/users. PutItem is
// idempotent so it's safe to run alongside global_emotes_test.go which also
// seeds this row.
func seedOauthMockResponse(t *testing.T, ctx context.Context, client *dynamodb.Client, tableName string) {
	t.Helper()

	oauthFixture, err := os.ReadFile("../fixtures/oauth2-token-response.json")
	if err != nil {
		t.Fatalf("failed to read oauth fixture: %v", err)
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"PK":           &types.AttributeValueMemberS{Value: "/oauth2/token"},
			"responseBody": &types.AttributeValueMemberS{Value: string(oauthFixture)},
		},
	})
	if err != nil {
		t.Fatalf("failed to seed oauth mock: %v", err)
	}
}

// seedChannelsMockResponses inserts /helix/users?id=<id> rows into the mock
// table for each test channel id.
func seedChannelsMockResponses(t *testing.T, ctx context.Context, client *dynamodb.Client, tableName string) {
	t.Helper()

	items := []struct {
		pk          string
		fixturePath string
	}{
		{"/helix/users?id=100", "../fixtures/twitch-user-100.json"},
		{"/helix/users?id=200", "../fixtures/twitch-user-200.json"},
		{"/helix/users?id=999", "../fixtures/twitch-user-empty.json"},
	}

	for _, item := range items {
		fixture, err := os.ReadFile(item.fixturePath)
		if err != nil {
			t.Fatalf("failed to read %s: %v", item.fixturePath, err)
		}
		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]types.AttributeValue{
				"PK":           &types.AttributeValueMemberS{Value: item.pk},
				"responseBody": &types.AttributeValueMemberS{Value: string(fixture)},
			},
		})
		if err != nil {
			t.Fatalf("failed to seed channels mock %s: %v", item.pk, err)
		}
	}
}

func postChannel(ctx context.Context, url, apiKey, twitchId string) ([]byte, int, error) {
	payload, err := json.Marshal(map[string]string{
		"twitchId": twitchId,
	})
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

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
