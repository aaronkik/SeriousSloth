package stack

const (
	production = "prod"
	staging    = "staging"
)

func IsEphemeral(stack string) bool {
	return stack != production && stack != staging
}

func IsProduction(stack string) bool {
	return stack == production
}

func EmoteSyncSchedule(stack string) string {
	if stack == staging {
		return "cron(0 0 * * ? *)"
	}
	return "cron(0 * * * ? *)"
}
