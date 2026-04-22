package main

import (
	"emotes-service/infra/components/shared"
	"emotes-service/infra/stack"

	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/scheduler"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ssm"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type StatelessComponent struct {
	pulumi.ResourceState
}

type StatefulResource struct {
	twitchEmotesSnapshotsTable *dynamodb.Table
}

func NewStatelessComponent(ctx *pulumi.Context, providerResource pulumi.ResourceOption, applicationConfig stack.ApplicationConfig, statefulResource StatefulResource) (*StatelessComponent, error) {
	component := &StatelessComponent{}

	err := ctx.RegisterComponentResourceV2("emotes-service:index:StatelessComponent", "stateless", nil, component)
	if err != nil {
		return nil, err
	}

	appConfig := config.New(ctx, "app")
	twitchClientId := appConfig.RequireSecret("twitch-client-id")
	twitchClientSecret := appConfig.RequireSecret("twitch-client-secret")

	twitchClientIdParam, err := ssm.NewParameter(ctx, "twitch-client-id", &ssm.ParameterArgs{
		Type:  pulumi.String(ssm.ParameterTypeSecureString),
		Value: pulumi.StringInput(twitchClientId),
	},
		pulumi.Parent(component),
		providerResource,
	)
	if err != nil {
		return nil, err
	}

	twitchClientSecretParam, err := ssm.NewParameter(ctx, "twitch-client-secret", &ssm.ParameterArgs{
		Type:  pulumi.String(ssm.ParameterTypeSecureString),
		Value: pulumi.StringInput(twitchClientSecret),
	},
		pulumi.Parent(component),
		providerResource,
	)
	if err != nil {
		return nil, err
	}

	syncGlobalEmotesLambda, err := components.NewLambda(ctx, "sync-global-emotes", &components.LambdaArgs{
		Code: pulumi.NewAssetArchive(map[string]interface{}{
			"bootstrap": pulumi.NewFileAsset("../dist/sync-global-emotes/bootstrap"),
		}),
		Environment: map[string]pulumi.StringInput{
			"TWITCH_EMOTES_SNAPSHOT_TABLE":   pulumi.StringInput(statefulResource.twitchEmotesSnapshotsTable.Name),
			"TWITCH_GLOBAL_EMOTES_ENDPOINT":  pulumi.String(applicationConfig.Twitch.GlobalEmotesEndpoint),
			"TWITCH_OAUTH_ENDPOINT":          pulumi.String(applicationConfig.Twitch.OauthEndpoint),
			"TWITCH_CLIENT_ID_PARAM_ARN":     pulumi.StringInput(twitchClientIdParam.Arn),
			"TWITCH_CLIENT_SECRET_PARAM_ARN": pulumi.StringInput(twitchClientSecretParam.Arn),
		},
		PolicyStatements: iam.GetPolicyDocumentStatementArray{
			&iam.GetPolicyDocumentStatementArgs{
				Effect:  pulumi.String("Allow"),
				Actions: pulumi.StringArray{pulumi.String("ssm:GetParameter")},
				Resources: pulumi.StringArray{
					twitchClientIdParam.Arn,
					twitchClientSecretParam.Arn,
				},
			},
			&iam.GetPolicyDocumentStatementArgs{
				Effect:  pulumi.String("Allow"),
				Actions: pulumi.StringArray{pulumi.String("dynamodb:PutItem")},
				Resources: pulumi.StringArray{
					statefulResource.twitchEmotesSnapshotsTable.Arn,
				},
			},
		},
	},
		pulumi.Parent(component),
		providerResource,
		pulumi.DependsOn([]pulumi.Resource{
			twitchClientIdParam,
			twitchClientSecretParam,
			statefulResource.twitchEmotesSnapshotsTable,
		}))
	if err != nil {
		return nil, err
	}

	caller, err := aws.GetCallerIdentity(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	accountId := caller.AccountId

	syncGlobalEmotesSchedulerRole, err := iam.NewRole(ctx, "sync-global-emotes-scheduler-role", &iam.RoleArgs{
		AssumeRolePolicy: iam.GetPolicyDocumentOutput(ctx, iam.GetPolicyDocumentOutputArgs{
			Statements: iam.GetPolicyDocumentStatementArray{
				&iam.GetPolicyDocumentStatementArgs{
					Actions: pulumi.StringArray{pulumi.String("sts:AssumeRole")},
					Conditions: &iam.GetPolicyDocumentStatementConditionArray{
						&iam.GetPolicyDocumentStatementConditionArgs{
							Test:     pulumi.String("StringEquals"),
							Variable: pulumi.String("aws:SourceAccount"),
							Values: pulumi.StringArray{
								pulumi.String(accountId),
							},
						},
					},
					Principals: iam.GetPolicyDocumentStatementPrincipalArray{
						&iam.GetPolicyDocumentStatementPrincipalArgs{
							Type: pulumi.String("Service"),
							Identifiers: pulumi.StringArray{
								pulumi.String("scheduler.amazonaws.com"),
							},
						},
					},
				},
			},
		}).Json(),
	},
		pulumi.Parent(component),
		providerResource,
	)
	if err != nil {
		return nil, err
	}

	_, err = iam.NewRolePolicy(ctx, "sync-global-emotes-scheduler-role-policy", &iam.RolePolicyArgs{
		Role: syncGlobalEmotesSchedulerRole.Name,
		Policy: iam.GetPolicyDocumentOutput(ctx, iam.GetPolicyDocumentOutputArgs{
			Statements: iam.GetPolicyDocumentStatementArray{
				&iam.GetPolicyDocumentStatementArgs{
					Effect:  pulumi.String("Allow"),
					Actions: pulumi.StringArray{pulumi.String("lambda:InvokeFunction")},
					Resources: pulumi.StringArray{
						syncGlobalEmotesLambda.Function.Arn,
						pulumi.Sprintf("%s:*", syncGlobalEmotesLambda.Function.Arn),
					},
				},
			},
		}).Json(),
	},
		pulumi.Parent(component),
		providerResource,
		pulumi.DependsOn([]pulumi.Resource{
			syncGlobalEmotesLambda,
		},
		),
	)
	if err != nil {
		return nil, err
	}

	_, err = scheduler.NewSchedule(ctx, "sync-global-emotes-scheduler", &scheduler.ScheduleArgs{
		ActionAfterCompletion: pulumi.String("NONE"),
		FlexibleTimeWindow: &scheduler.ScheduleFlexibleTimeWindowArgs{
			Mode: pulumi.String("OFF"),
		},
		ScheduleExpression:         pulumi.String("cron(0 * * * ? *)"),
		ScheduleExpressionTimezone: pulumi.String("UTC"),
		State:                      pulumi.String("ENABLED"),
		Target: &scheduler.ScheduleTargetArgs{
			Arn:     pulumi.StringInput(syncGlobalEmotesLambda.Function.Arn),
			RoleArn: pulumi.StringInput(syncGlobalEmotesSchedulerRole.Arn),
			RetryPolicy: &scheduler.ScheduleTargetRetryPolicyArgs{
				MaximumEventAgeInSeconds: nil,
				MaximumRetryAttempts:     pulumi.Int(0),
			},
		},
	},
		pulumi.Parent(component),
		providerResource,
	)
	if err != nil {
		return nil, err
	}

	return component, nil
}
