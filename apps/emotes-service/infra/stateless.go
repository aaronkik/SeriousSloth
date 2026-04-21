package main

import (
	"emotes-service/infra/components/shared"
	"emotes-service/infra/stack"

	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/scheduler"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type StatelessComponent struct {
	pulumi.ResourceState
}

func NewStatelessComponent(ctx *pulumi.Context, name string, providerResource pulumi.ResourceOption, applicationConfig stack.ApplicationConfig) (*StatelessComponent, error) {
	component := &StatelessComponent{}

	err := ctx.RegisterComponentResourceV2("emotes-service:index:StatelessComponent", name, nil, component)
	if err != nil {
		return nil, err
	}

	syncGlobalEmotesLambda, err := components.NewLambda(ctx, "sync-global-emotes", &components.LambdaArgs{
		Code: pulumi.NewAssetArchive(map[string]interface{}{
			"bootstrap": pulumi.NewFileAsset("../dist/sync-global-emotes/bootstrap"),
		}),
		Environment: map[string]pulumi.StringInput{
			"TWITCH_GLOBAL_EMOTES_ENDPOINT": pulumi.String(applicationConfig.Twitch.GlobalEmotesEndpoint),
			"TWITCH_OAUTH_ENDPOINT":         pulumi.String(applicationConfig.Twitch.OauthEndpoint),
		},
	}, pulumi.Parent(component), providerResource)
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
	}, pulumi.Parent(component), providerResource)
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
	}, pulumi.Parent(component), pulumi.DependsOn([]pulumi.Resource{syncGlobalEmotesLambda}))
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
	}, pulumi.Parent(component), providerResource)
	if err != nil {
		return nil, err
	}

	return component, nil
}
