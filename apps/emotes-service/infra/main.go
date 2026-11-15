package main

import (
	"emotes-service/infra/components"
	"emotes-service/infra/git"
	"emotes-service/infra/stack"

	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		gitRemoteUrl := git.Cli("remote", "get-url", "origin")

		awsConfig := config.New(ctx, "aws")
		awsRegion := awsConfig.Require("region")

		provider, err := aws.NewProvider(ctx, "aws-provider", &aws.ProviderArgs{
			Region: pulumi.String(awsRegion),
			DefaultTags: &aws.ProviderDefaultTagsArgs{
				Tags: pulumi.StringMap{
					"Project":    pulumi.String(ctx.Project()),
					"Stack":      pulumi.String(ctx.Stack()),
					"ManagedBy":  pulumi.String("pulumi"),
					"Repository": pulumi.String(gitRemoteUrl),
					"Commit":     pulumi.String(git.Cli("rev-parse", "HEAD")),
				},
			},
		})
		if err != nil {
			return err
		}

		providerResource := pulumi.Provider(provider)

		appConfig := stack.GetApplicationConfig()

		var integrationComponent *components.IntegrationComponent
		if stack.IsEphemeral(ctx.Stack()) {
			integrationComponent, err = components.NewIntegrationComponent(ctx, providerResource)
			if err != nil {
				return err
			}

			appConfig.Twitch.OauthEndpoint = pulumi.Sprintf("%soauth2/token", integrationComponent.MockTwitchApiUrl)
			appConfig.Twitch.GlobalEmotesEndpoint = pulumi.Sprintf("%shelix/chat/emotes/global", integrationComponent.MockTwitchApiUrl)
			appConfig.Twitch.UsersEndpoint = pulumi.Sprintf("%shelix/users", integrationComponent.MockTwitchApiUrl)
			appConfig.Twitch.ChannelEmotesEndpoint = pulumi.Sprintf("%shelix/chat/emotes", integrationComponent.MockTwitchApiUrl)
		}

		statefulComponent, err := components.NewStatefulComponent(ctx, providerResource)
		if err != nil {
			return err
		}

		statelessComponent, err := components.NewStatelessComponent(ctx, providerResource, appConfig, components.StatefulResource{
			TwitchEmotesEventsStoreTable: statefulComponent.TwitchEmotesEventStoreTable,
			TwitchEmotesProjectionsTable: statefulComponent.TwitchEmotesProjectionsTable,
			TwitchChannelsTable:          statefulComponent.TwitchChannelsTable,
		})
		if err != nil {
			return err
		}

		if integrationComponent != nil {
			ctx.Export("apiInvokeUrl", statelessComponent.ApiInvokeUrl)
			ctx.Export("apiKeyId", statelessComponent.ApiKeyId)
			ctx.Export("syncGlobalEmotesLambdaName", statelessComponent.SyncGlobalEmotesFunction.Name)
			ctx.Export("channelSyncDispatcherLambdaName", statelessComponent.ChannelSyncDispatcherFunction.Name)
			ctx.Export("channelSyncQueueUrl", statelessComponent.ChannelSyncQueueUrl)
			ctx.Export("twitchEmotesEventStoreTable", statefulComponent.TwitchEmotesEventStoreTable.Name)
			ctx.Export("twitchEmotesProjectionsTable", statefulComponent.TwitchEmotesProjectionsTable.Name)
			ctx.Export("twitchChannelsTable", statefulComponent.TwitchChannelsTable.Name)
			ctx.Export("mockTwitchResponsesTableName", integrationComponent.MockTwitchResponsesTable.Name)
		}

		return nil
	})
}
