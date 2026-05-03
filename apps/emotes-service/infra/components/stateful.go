package components

import (
	"emotes-service/infra/stack"

	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type StatefulComponent struct {
	pulumi.ResourceState
	TwitchEmotesEventStoreTable  *dynamodb.Table
	TwitchEmotesProjectionsTable *dynamodb.Table
}

func NewStatefulComponent(ctx *pulumi.Context, providerResource pulumi.ResourceOption) (*StatefulComponent, error) {
	component := &StatefulComponent{}

	err := ctx.RegisterComponentResourceV2("emotes-service:index:StatefulComponent", "stateful", nil, component)
	if err != nil {
		return nil, err
	}

	twitchEmotesEventStoreTable, err := dynamodb.NewTable(ctx, "twitch-emotes-event-store", &dynamodb.TableArgs{
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
		StreamEnabled:             pulumi.BoolPtr(true),
		StreamViewType:            pulumi.String("NEW_IMAGE"),
	},
		pulumi.Parent(component),
		providerResource,
	)
	if err != nil {
		return nil, err
	}

	twitchEmotesProjectionsTable, err := dynamodb.NewTable(ctx, "twitch-emotes-projections", &dynamodb.TableArgs{
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
	},
		pulumi.Parent(component),
		providerResource,
	)
	if err != nil {
		return nil, err
	}

	component.TwitchEmotesEventStoreTable = twitchEmotesEventStoreTable
	component.TwitchEmotesProjectionsTable = twitchEmotesProjectionsTable

	return component, nil
}
