package main

import (
	"emotes-service/util"
	"strings"
	"time"

	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-pulumiservice/sdk/go/pulumiservice"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		gitRemoteUrl := util.GitCli("remote", "get-url", "origin")
		gitRepoSlug, err := util.GetRepositorySlug(gitRemoteUrl)
		if err != nil {
			return err
		}

		deploymentSettings, err := pulumiservice.NewDeploymentSettings(ctx, "deploymentSettingsResource", &pulumiservice.DeploymentSettingsArgs{
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
					Branch:  pulumi.String(util.GitCli("branch", "--show")),
					RepoDir: pulumi.String(strings.TrimSuffix(util.GitCli("rev-parse", "--show-prefix"), "/")),
				},
			},
		})
		if err != nil {
			return err
		}

		now := time.Now()
		stackDestruction := now.Add(time.Hour * 24 * 3)

		if util.IsEphemeral(ctx.Stack()) {
			_, err = pulumiservice.NewTtlSchedule(ctx, "ttlSchedule", &pulumiservice.TtlScheduleArgs{
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

		provider, err := aws.NewProvider(ctx, "awsProvider", &aws.ProviderArgs{
			Region: pulumi.String(awsRegion),
			DefaultTags: &aws.ProviderDefaultTagsArgs{
				Tags: pulumi.StringMap{
					"Project":    pulumi.String(ctx.Project()),
					"Stack":      pulumi.String(ctx.Stack()),
					"ManagedBy":  pulumi.String("pulumi"),
					"Repository": pulumi.String(gitRemoteUrl),
					"Commit":     pulumi.String(util.GitCli("rev-parse", "HEAD")),
				},
			},
		})
		if err != nil {
			return err
		}

		providerResource := pulumi.Provider(provider)

		_, err = NewStatefulComponent(ctx, "stateful", providerResource)
		if err != nil {
			return err
		}

		return nil
	})
}
