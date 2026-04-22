package dynamodb

import (
	"context"
	"log"
	"log/slog"

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

type PutItemInput struct {
	TableName  string
	PK         string
	SK         *string
	Attributes any
}

func PutItem(input PutItemInput) error {
	item, err := attributevalue.MarshalMap(input.Attributes)
	if err != nil {
		slog.Error("Error marshalling item", "error", err)
		return err
	}

	item["PK"] = &types.AttributeValueMemberS{Value: input.PK}
	if input.SK != nil {
		item["SK"] = &types.AttributeValueMemberS{Value: *input.SK}
	}

	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &input.TableName,
		Item:      item,
	})
	if err != nil {
		slog.Error("Error putting item", "input", input, "error", err)
		return err
	}

	return nil
}
