package projections_store

import (
	"context"
	"emotes-service/src/environment"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func QueryActiveEmotes(ctx context.Context, aggregateId string) ([]ProjectionItem, error) {
	paginator := dynamodb.NewQueryPaginator(client, &dynamodb.QueryInput{
		TableName:              aws.String(environment.GetOrFatal("EVENTS_PROJECTION_TABLE_NAME")),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		FilterExpression:       aws.String("#status = :active"),
		ExpressionAttributeNames: map[string]string{
			"#status": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: aggregateId},
			":sk":     &types.AttributeValueMemberS{Value: "EMOTE#"},
			":active": &types.AttributeValueMemberS{Value: "ACTIVE"},
		},
	})

	var items []ProjectionItem
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "Error querying projections", "aggregateId", aggregateId, "error", err)
			return nil, err
		}

		var pageItems []ProjectionItem
		if err := attributevalue.UnmarshalListOfMaps(page.Items, &pageItems); err != nil {
			slog.ErrorContext(ctx, "Error unmarshalling projections page", "aggregateId", aggregateId, "error", err)
			return nil, err
		}
		items = append(items, pageItems...)
	}

	return items, nil
}
