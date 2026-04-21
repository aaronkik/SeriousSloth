package main

import (
	"emotes-service/infra/stack"

	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type StatefulComponent struct {
	pulumi.ResourceState
	TwitchEmotesSnapshotsTable *dynamodb.Table
}

func NewStatefulComponent(ctx *pulumi.Context, name string, providerResource pulumi.ResourceOption) (*StatefulComponent, error) {
	component := &StatefulComponent{}

	err := ctx.RegisterComponentResourceV2("emotes-service:index:StatefulComponent", name, nil, component)
	if err != nil {
		return nil, err
	}

	twitchEmotesSnapshotsTable, err := dynamodb.NewTable(ctx, "twitch-emotes-snapshots", &dynamodb.TableArgs{
		BillingMode: pulumi.String("PAY_PER_REQUEST"),
		HashKey:     pulumi.String("PK"),
		RangeKey:    pulumi.String("SK"),
		Attributes: dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("PK"),
				Type: pulumi.String("S"),
			},
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("SK"),
				Type: pulumi.String("S"),
			},
		},
		DeletionProtectionEnabled: pulumi.BoolPtr(stack.IsProduction(ctx.Stack())),
	}, pulumi.Parent(component), providerResource)
	if err != nil {
		return nil, err
	}

	component.TwitchEmotesSnapshotsTable = twitchEmotesSnapshotsTable

	return component, nil
}
