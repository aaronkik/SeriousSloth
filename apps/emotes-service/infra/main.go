package main

import (
	"emotes-service/infra/components"
	"emotes-service/infra/git"
	"emotes-service/infra/stack"
	"strings"
	"time"

	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-pulumiservice/sdk/go/pulumiservice"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		gitRemoteUrl := git.Cli("remote", "get-url", "origin")
		gitRepoSlug, err := git.GetRepositorySlug(gitRemoteUrl)
		if err != nil {
			return err
		}

		deploymentSettings, err := pulumiservice.NewDeploymentSettings(ctx, "deployment-settings-resource", &pulumiservice.DeploymentSettingsArgs{
			Organization: pulumi.String(ctx.Organization()),
			Project:      pulumi.String(ctx.Project()),
			Stack:        pulumi.String(ctx.Stack()),
			CacheOptions: &pulumiservice.DeploymentSettingsCacheOptionsArgs{
				Enable: pulumi.Bool(true),
			},
			Vcs: &pulumiservice.DeploymentSettingsVcsArgs{
				DeployCommits:       pulumi.Bool(true),
				PreviewPullRequests: pulumi.Bool(true),
				Provider:            pulumi.String("github"),
				Repository:          pulumi.String(gitRepoSlug),
			},
			SourceContext: &pulumiservice.DeploymentSettingsSourceContextArgs{
				Git: &pulumiservice.DeploymentSettingsGitSourceArgs{
					Branch:  pulumi.String(git.Cli("branch", "--show")),
					RepoDir: pulumi.String(strings.TrimSuffix(git.Cli("rev-parse", "--show-prefix"), "/")),
				},
			},
		})
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		daysUntilSunday := (7 - int(now.Weekday())) % 7
		stackDestruction := time.Date(now.Year(), now.Month(), now.Day()+daysUntilSunday, 23, 59, 59, 0, time.UTC)

		if stack.IsEphemeral(ctx.Stack()) {
			_, err = pulumiservice.NewTtlSchedule(ctx, "ttl-schedule", &pulumiservice.TtlScheduleArgs{
				Organization: pulumi.String(ctx.Organization()),
				Project:      pulumi.String(ctx.Project()),
				Stack:        pulumi.String(ctx.Stack()),
				Timestamp:    pulumi.String(stackDestruction.UTC().Format(time.RFC3339)),
			}, pulumi.DependsOn([]pulumi.Resource{deploymentSettings}))
			if err != nil {
				return err
			}
		}

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

		statefulComponent, err := components.NewStatefulComponent(ctx, providerResource)
		if err != nil {
			return err
		}

		_, err = components.NewStatelessComponent(ctx, providerResource, appConfig, components.StatefulResource{TwitchEmotesSnapshotsTable: statefulComponent.TwitchEmotesSnapshotsTable})
		if err != nil {
			return err
		}

		return nil
	})
}
