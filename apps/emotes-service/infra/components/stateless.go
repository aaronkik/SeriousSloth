package components

import (
	"emotes-service/infra/components/shared"
	"emotes-service/infra/stack"
	"time"

	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/scheduler"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/sqs"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ssm"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type StatelessComponent struct {
	pulumi.ResourceState
	SyncGlobalEmotesFunction *lambda.Function
}

type StatefulResource struct {
	TwitchEmotesEventsStoreTable *dynamodb.Table
	TwitchEmotesProjectionsTable *dynamodb.Table
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
			"TWITCH_EMOTES_EVENT_STORE_TABLE": pulumi.StringInput(statefulResource.TwitchEmotesEventsStoreTable.Name),
			"TWITCH_GLOBAL_EMOTES_ENDPOINT":   applicationConfig.Twitch.GlobalEmotesEndpoint,
			"TWITCH_OAUTH_ENDPOINT":           applicationConfig.Twitch.OauthEndpoint,
			"TWITCH_CLIENT_ID_PARAM_ARN":      pulumi.StringInput(twitchClientIdParam.Arn),
			"TWITCH_CLIENT_SECRET_PARAM_ARN":  pulumi.StringInput(twitchClientSecretParam.Arn),
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
				Effect: pulumi.String("Allow"),
				Actions: pulumi.StringArray{
					pulumi.String("dynamodb:ConditionCheckItem"),
					pulumi.String("dynamodb:PutItem"),
					pulumi.String("dynamodb:Query"),
				},
				Resources: pulumi.StringArray{
					statefulResource.TwitchEmotesEventsStoreTable.Arn,
				},
			},
			&iam.GetPolicyDocumentStatementArgs{
				Effect: pulumi.String("Allow"),
				Actions: pulumi.StringArray{
					pulumi.String("sqs:SendMessage"),
				},
				Resources: pulumi.StringArray{
					statefulResource.TwitchEmotesEventsStoreTable.Arn,
				},
			},
		},
	},
		pulumi.Parent(component),
		providerResource,
		pulumi.DependsOn([]pulumi.Resource{
			twitchClientIdParam,
			twitchClientSecretParam,
		}))
	if err != nil {
		return nil, err
	}

	component.SyncGlobalEmotesFunction = syncGlobalEmotesLambda.Function

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

	schedulerState := "DISABLED"
	if stack.IsEphemeral(ctx.Stack()) {
		schedulerState = "DISABLED"
	}

	_, err = scheduler.NewSchedule(ctx, "sync-global-emotes-scheduler", &scheduler.ScheduleArgs{
		ActionAfterCompletion: pulumi.String("NONE"),
		FlexibleTimeWindow: &scheduler.ScheduleFlexibleTimeWindowArgs{
			Mode: pulumi.String("OFF"),
		},
		ScheduleExpression:         pulumi.String("cron(0 * * * ? *)"),
		ScheduleExpressionTimezone: pulumi.String("UTC"),
		State:                      pulumi.String(schedulerState),
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

	emotesReadModelProducerDlq, err := sqs.NewQueue(ctx, "emotes-read-model-produce-dlq", &sqs.QueueArgs{
		MaxMessageSize:          pulumi.Int(1048576),
		MessageRetentionSeconds: pulumi.Int((time.Hour * 24 * 14) / time.Second),
	},
		pulumi.Parent(component),
		providerResource,
	)
	if err != nil {
		return nil, err
	}

	emotesReadModelProducer, err := components.NewLambda(ctx, "emotes-read-model-producer", &components.LambdaArgs{
		Code: pulumi.NewAssetArchive(map[string]interface{}{
			"bootstrap": pulumi.NewFileAsset("../dist/emotes-read-model-producer/bootstrap"),
		}),
		Environment: map[string]pulumi.StringInput{
			"EVENTS_PROJECTION_TABLE_NAME": pulumi.StringInput(statefulResource.TwitchEmotesProjectionsTable.Name),
		},
		PolicyStatements: iam.GetPolicyDocumentStatementArray{
			&iam.GetPolicyDocumentStatementArgs{
				Effect: pulumi.String("Allow"),
				Actions: pulumi.StringArray{
					pulumi.String("dynamodb:GetRecords"),
					pulumi.String("dynamodb:GetShardIterator"),
					pulumi.String("dynamodb:DescribeStream"),
					pulumi.String("dynamodb:ListStreams"),
				},
				Resources: pulumi.StringArray{
					statefulResource.TwitchEmotesEventsStoreTable.StreamArn,
				},
			},
			&iam.GetPolicyDocumentStatementArgs{
				Effect: pulumi.String("Allow"),
				Actions: pulumi.StringArray{
					pulumi.String("sqs:SendMessage"),
				},
				Resources: pulumi.StringArray{
					emotesReadModelProducerDlq.Arn,
				},
			},
			&iam.GetPolicyDocumentStatementArgs{
				Effect: pulumi.String("Allow"),
				Actions: pulumi.StringArray{
					pulumi.String("dynamodb:*"),
				},
				Resources: pulumi.StringArray{
					statefulResource.TwitchEmotesProjectionsTable.Arn,
				},
			},
		},
	},
		pulumi.Parent(component),
		providerResource,
	)
	if err != nil {
		return nil, err
	}

	_, err = lambda.NewEventSourceMapping(ctx, "emotes-read-model-producer-mapping", &lambda.EventSourceMappingArgs{
		EventSourceArn:       pulumi.StringInput(statefulResource.TwitchEmotesEventsStoreTable.StreamArn),
		FunctionName:         pulumi.StringInput(emotesReadModelProducer.Function.Name),
		StartingPosition:     pulumi.String("TRIM_HORIZON"),
		BatchSize:            pulumi.Int(1),
		MaximumRetryAttempts: pulumi.Int(3),
		DestinationConfig: &lambda.EventSourceMappingDestinationConfigArgs{
			OnFailure: &lambda.EventSourceMappingDestinationConfigOnFailureArgs{
				DestinationArn: pulumi.StringInput(emotesReadModelProducerDlq.Arn),
			},
		},
		FilterCriteria: &lambda.EventSourceMappingFilterCriteriaArgs{
			Filters: &lambda.EventSourceMappingFilterCriteriaFilterArray{
				&lambda.EventSourceMappingFilterCriteriaFilterArgs{
					Pattern: pulumi.String("{ \"eventName\" : [\"INSERT\"] }"),
				},
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
