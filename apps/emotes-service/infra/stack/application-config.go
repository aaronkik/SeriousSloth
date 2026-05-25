package stack

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

type TwitchConfig struct {
	OauthEndpoint        pulumi.StringInput
	GlobalEmotesEndpoint pulumi.StringInput
	UsersEndpoint        pulumi.StringInput
}

type ApplicationConfig struct {
	Twitch TwitchConfig
}

var applicationConfig = ApplicationConfig{
	Twitch: TwitchConfig{
		OauthEndpoint:        pulumi.String("https://id.twitch.tv/oauth2/token"),
		GlobalEmotesEndpoint: pulumi.String("https://api.twitch.tv/helix/chat/emotes/global"),
		UsersEndpoint:        pulumi.String("https://api.twitch.tv/helix/users"),
	},
}

func GetApplicationConfig() ApplicationConfig {
	return applicationConfig
}
