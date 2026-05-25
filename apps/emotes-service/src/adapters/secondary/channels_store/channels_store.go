package channels_store

import (
	"context"
	"emotes-service/src/environment"
	"errors"
	"log"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/newrelic/go-agent/v3/newrelic"
)

const ChannelsPartitionKey = "CHANNEL"

var ErrAlreadyExists = errors.New("channel already exists")

var client *dynamodb.Client

func init() {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client = dynamodb.NewFromConfig(cfg)
}

type ChannelItem struct {
	PK          string `dynamodbav:"PK"`
	SK          string `dynamodbav:"SK"`
	Id          string `dynamodbav:"id"`
	TwitchId    string `dynamodbav:"twitchId"`
	DisplayName string `dynamodbav:"displayName"`
	ImageUrl    string `dynamodbav:"imageUrl"`
	AddedAt     string `dynamodbav:"addedAt"`
	UpdatedAt   string `dynamodbav:"updatedAt"`
}

func QueryAll(ctx context.Context) ([]ChannelItem, error) {
	txn := newrelic.FromContext(ctx)
	tableName := environment.GetOrFatal("CHANNELS_TABLE_NAME")

	paginator := dynamodb.NewQueryPaginator(client, &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: ChannelsPartitionKey},
		},
	})

	var items []ChannelItem
	for paginator.HasMorePages() {
		ddbSeg := newrelic.DatastoreSegment{
			StartTime:          txn.StartSegmentNow(),
			Product:            newrelic.DatastoreDynamoDB,
			Collection:         tableName,
			Operation:          "Query",
			ParameterizedQuery: "PK = :pk",
			QueryParameters: map[string]any{
				":pk": ChannelsPartitionKey,
			},
		}
		page, err := paginator.NextPage(ctx)
		ddbSeg.End()
		if err != nil {
			slog.ErrorContext(ctx, "channels_store Query failed", "error", err)
			return nil, err
		}

		var pageItems []ChannelItem
		if err := attributevalue.UnmarshalListOfMaps(page.Items, &pageItems); err != nil {
			slog.ErrorContext(ctx, "channels_store unmarshal failed", "error", err)
			return nil, err
		}
		items = append(items, pageItems...)
	}

	return items, nil
}

func Put(ctx context.Context, item ChannelItem) error {
	txn := newrelic.FromContext(ctx)
	tableName := environment.GetOrFatal("CHANNELS_TABLE_NAME")

	attrs, err := attributevalue.MarshalMap(item)
	if err != nil {
		slog.ErrorContext(ctx, "channels_store marshal failed", "error", err)
		return err
	}

	ddbSeg := newrelic.DatastoreSegment{
		StartTime:  txn.StartSegmentNow(),
		Product:    newrelic.DatastoreDynamoDB,
		Collection: tableName,
		Operation:  "PutItem",
	}
	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(tableName),
		Item:                attrs,
		ConditionExpression: aws.String("attribute_not_exists(SK)"),
	})
	ddbSeg.End()

	if err != nil {
		if _, ok := errors.AsType[*types.ConditionalCheckFailedException](err); ok {
			return ErrAlreadyExists
		}
		slog.ErrorContext(ctx, "channels_store PutItem failed", "error", err)
		return err
	}

	return nil
}
