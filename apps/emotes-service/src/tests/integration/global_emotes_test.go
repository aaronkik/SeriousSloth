//go:build integration

package integration

import (
	"context"
	"emotes-service/src/tests/helpers"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	awslambda "github.com/aws/aws-sdk-go-v2/service/lambda"
)

func TestGlobalEmotes(t *testing.T) {
	syncLambdaName := helpers.GetPulumiExport(t, "syncGlobalEmotesLambdaName")
	snapshotsTableName := helpers.GetPulumiExport(t, "twitchEmotesSnapshotsTableName")
	mockTwitchResponsesTableName := helpers.GetPulumiExport(t, "mockTwitchResponsesTableName")

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		t.Fatalf("failed to load AWS config: %v", err)
	}

	ddbClient := dynamodb.NewFromConfig(cfg)
	lambdaClient := awslambda.NewFromConfig(cfg)

	seedMockResponses(t, ctx, ddbClient, mockTwitchResponsesTableName)
	clearSnapshotsTable(t, ctx, ddbClient, snapshotsTableName)

	beforeInvoke := time.Now().UTC()

	invokeOutput, err := lambdaClient.Invoke(ctx, &awslambda.InvokeInput{
		FunctionName: aws.String(syncLambdaName),
		Payload:      []byte(`{}`),
	})
	if err != nil {
		t.Fatalf("failed to invoke Lambda: %v", err)
	}
	if invokeOutput.FunctionError != nil {
		t.Fatalf("Lambda returned error: %s, payload: %s", *invokeOutput.FunctionError, string(invokeOutput.Payload))
	}

	queryOutput, err := ddbClient.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(snapshotsTableName),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "GLOBAL"},
		},
	})
	if err != nil {
		t.Fatalf("failed to query DynamoDB: %v", err)
	}

	if len(queryOutput.Items) == 0 {
		t.Fatal("expected at least one record in DynamoDB after Lambda invocation, got 0")
	}

	item := queryOutput.Items[0]

	skAttr, ok := item["SK"].(*types.AttributeValueMemberS)
	if !ok {
		t.Fatal("record missing 'SK' attribute or not a string")
	}
	sk, err := time.Parse(time.RFC3339, skAttr.Value)
	if err != nil {
		t.Fatalf("SK %q is not RFC3339: %v", skAttr.Value, err)
	}
	if sk.Before(beforeInvoke) {
		t.Fatalf("SK %s is before beforeInvoke %s", sk, beforeInvoke)
	}

	globalEmotesAttr, ok := item["globalEmotes"]
	if !ok {
		t.Fatal("record missing 'globalEmotes' attribute")
	}

	listAttr, ok := globalEmotesAttr.(*types.AttributeValueMemberL)
	if !ok {
		t.Fatalf("globalEmotes is not a list, got %T", globalEmotesAttr)
	}

	if len(listAttr.Value) != 1 {
		t.Fatalf("expected 1 emote in globalEmotes, got %d", len(listAttr.Value))
	}

	t.Logf("verified %d emote(s) written to DynamoDB at SK=%s", len(listAttr.Value), sk)
}

func clearSnapshotsTable(t *testing.T, ctx context.Context, client *dynamodb.Client, tableName string) {
	t.Helper()

	var lastKey map[string]types.AttributeValue
	for {
		scanOutput, err := client.Scan(ctx, &dynamodb.ScanInput{
			TableName: aws.String(tableName),
			//ProjectionExpression: aws.String("PK, SK"),
			ExclusiveStartKey: lastKey,
		})
		if err != nil {
			t.Fatalf("failed to scan snapshots table: %v", err)
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
				t.Fatalf("failed to delete snapshot item: %v", err)
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
