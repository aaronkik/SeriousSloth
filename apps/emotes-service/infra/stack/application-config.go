package stack

type TwitchConfig struct {
	OauthEndpoint        string
	GlobalEmotesEndpoint string
}

type ApplicationConfig struct {
	Twitch TwitchConfig
}

var applicationConfig = ApplicationConfig{
	Twitch: TwitchConfig{
		OauthEndpoint:        "https://id.twitch.tv/oauth2/token",
		GlobalEmotesEndpoint: "https://api.twitch.tv/helix/chat/emotes/global",
	},
}

func GetApplicationConfig() ApplicationConfig {
	return applicationConfig
}
