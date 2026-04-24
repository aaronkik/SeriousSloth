package components

import (
	"emotes-service/infra/components/shared"

	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type IntegrationComponent struct {
	pulumi.ResourceState
	MockTwitchApiUrl         pulumi.StringOutput
	MockTwitchResponsesTable *dynamodb.Table
}

func NewIntegrationComponent(ctx *pulumi.Context, providerResource pulumi.ResourceOption) (*IntegrationComponent, error) {
	component := &IntegrationComponent{}

	err := ctx.RegisterComponentResourceV2("emotes-service:index:IntegrationComponent", "integration", nil, component)
	if err != nil {
		return nil, err
	}

	twitchMockResponsesTable, err := dynamodb.NewTable(ctx, "twitch-mock-responses", &dynamodb.TableArgs{
		BillingMode: pulumi.String("PAY_PER_REQUEST"),
		HashKey:     pulumi.String("PK"),
		Attributes: dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("PK"),
				Type: pulumi.String("S"),
			},
		},
	}, pulumi.Parent(component), providerResource)
	if err != nil {
		return nil, err
	}

	mockTwitchApiLambda, err := components.NewLambda(ctx, "mock-twitch-api", &components.LambdaArgs{
		Code: pulumi.NewAssetArchive(map[string]interface{}{
			"bootstrap": pulumi.NewFileAsset("../dist/mock-twitch-api/bootstrap"),
		}),
		Environment: pulumi.StringMap{
			"MOCK_RESPONSES_TABLE": pulumi.StringInput(twitchMockResponsesTable.Name),
		},
		PolicyStatements: iam.GetPolicyDocumentStatementArray{
			&iam.GetPolicyDocumentStatementArgs{
				Effect:  pulumi.String("Allow"),
				Actions: pulumi.StringArray{pulumi.String("dynamodb:GetItem")},
				Resources: pulumi.StringArray{
					twitchMockResponsesTable.Arn,
				},
			},
		},
	}, pulumi.Parent(component), providerResource)
	if err != nil {
		return nil, err
	}

	mockTwitchApiUrlFunctionUrl, err := lambda.NewFunctionUrl(ctx, "mock-twitch-api-url", &lambda.FunctionUrlArgs{
		FunctionName:      mockTwitchApiLambda.Function.Name,
		AuthorizationType: pulumi.String("NONE"),
	}, pulumi.Parent(component), providerResource)
	if err != nil {
		return nil, err
	}

	_, err = lambda.NewPermission(ctx, "mock-twitch-api-url-permission", &lambda.PermissionArgs{
		Action:              pulumi.String("lambda:InvokeFunctionUrl"),
		Function:            mockTwitchApiLambda.Function.Name,
		Principal:           pulumi.String("*"),
		FunctionUrlAuthType: pulumi.String("NONE"),
	}, pulumi.Parent(component), providerResource)
	if err != nil {
		return nil, err
	}

	component.MockTwitchApiUrl = mockTwitchApiUrlFunctionUrl.FunctionUrl
	component.MockTwitchResponsesTable = twitchMockResponsesTable

	return component, nil
}
