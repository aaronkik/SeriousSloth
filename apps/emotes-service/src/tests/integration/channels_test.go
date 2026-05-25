//go:build integration

package integration

import (
	"bytes"
	"context"
	"emotes-service/src/tests/helpers"
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

	clearTable(t, testCtx, ddb, channelsTable)
	t.Cleanup(func() { clearTable(t, testCtx, ddb, channelsTable) })

	channelsUrl := invokeUrl + "/channels"
	testStart := time.Now()

	t.Run("post creates two channels", func(t *testing.T) {
		require := require.New(t)

		inputs := []apiChannel{
			{TwitchId: "100", DisplayName: "alice", ImageUrl: "https://cdn.twitch/100.png"},
			{TwitchId: "200", DisplayName: "bob", ImageUrl: "https://cdn.twitch/200.png"},
		}

		for _, input := range inputs {
			body, status, err := postChannel(testCtx, channelsUrl, key, input)
			require.NoError(err)
			require.Equalf(http.StatusCreated, status, "body: %s", string(body))

			var created apiChannel
			require.NoError(json.Unmarshal(body, &created))
			require.Regexp(`^chnl_[0-9a-f]{24}$`, created.Id)
			require.Equal(input.TwitchId, created.TwitchId)
			require.Equal(input.DisplayName, created.DisplayName)
			require.Equal(input.ImageUrl, created.ImageUrl)

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
		require.Equal("alice", parsed[0].DisplayName)
		require.Regexp(`^chnl_[0-9a-f]{24}$`, parsed[0].Id)
		require.Equal("200", parsed[1].TwitchId)
		require.Equal("bob", parsed[1].DisplayName)
		require.Regexp(`^chnl_[0-9a-f]{24}$`, parsed[1].Id)
		require.NotEqual(parsed[0].Id, parsed[1].Id)
	})

	t.Run("post duplicate returns 409", func(t *testing.T) {
		require := require.New(t)

		_, status, err := postChannel(testCtx, channelsUrl, key, apiChannel{
			TwitchId: "100", DisplayName: "alice", ImageUrl: "https://cdn.twitch/100.png",
		})
		require.NoError(err)
		require.Equal(http.StatusConflict, status)
	})

	t.Run("post with missing fields returns 400", func(t *testing.T) {
		require := require.New(t)

		_, status, err := postChannel(testCtx, channelsUrl, key, apiChannel{
			TwitchId: "300", DisplayName: "", ImageUrl: "https://cdn.twitch/300.png",
		})
		require.NoError(err)
		require.Equal(http.StatusBadRequest, status)
	})
}

func postChannel(ctx context.Context, url, apiKey string, channel apiChannel) ([]byte, int, error) {
	payload, err := json.Marshal(map[string]string{
		"twitchId":    channel.TwitchId,
		"displayName": channel.DisplayName,
		"imageUrl":    channel.ImageUrl,
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
