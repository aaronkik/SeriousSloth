package projections_store

import (
	"context"
	"emotes-service/src/environment"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func QueryActiveEmotes(ctx context.Context, aggregateId string) ([]ProjectionItem, error) {
	return queryEmotesByStatus(ctx, aggregateId, "ACTIVE")
}

func QueryRemovedEmotes(ctx context.Context, aggregateId string) ([]ProjectionItem, error) {
	return queryEmotesByStatus(ctx, aggregateId, "REMOVED")
}

func queryEmotesByStatus(ctx context.Context, aggregateId, status string) ([]ProjectionItem, error) {
	txn := newrelic.FromContext(ctx)
	tableName := environment.GetOrFatal("EVENTS_PROJECTION_TABLE_NAME")

	paginator := dynamodb.NewQueryPaginator(client, &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		FilterExpression:       aws.String("#status = :status"),
		ExpressionAttributeNames: map[string]string{
			"#status": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: aggregateId},
			":sk":     &types.AttributeValueMemberS{Value: "EMOTE#"},
			":status": &types.AttributeValueMemberS{Value: status},
		},
	})

	var items []ProjectionItem
	for paginator.HasMorePages() {
		ddbSeg := newrelic.DatastoreSegment{
			StartTime:          txn.StartSegmentNow(),
			Product:            newrelic.DatastoreDynamoDB,
			Collection:         tableName,
			Operation:          "Query",
			ParameterizedQuery: "PK = :pk AND begins_with(SK, :sk)",
			QueryParameters: map[string]any{
				":pk":     aggregateId,
				":sk":     "EMOTE#",
				":status": status,
			},
		}
		page, err := paginator.NextPage(ctx)
		ddbSeg.End()
		if err != nil {
			slog.ErrorContext(ctx, "Error querying projections", "aggregateId", aggregateId, "status", status, "error", err)
			return nil, err
		}

		var pageItems []ProjectionItem
		if err := attributevalue.UnmarshalListOfMaps(page.Items, &pageItems); err != nil {
			slog.ErrorContext(ctx, "Error unmarshalling projections page", "aggregateId", aggregateId, "status", status, "error", err)
			return nil, err
		}
		items = append(items, pageItems...)
	}

	return items, nil
}
