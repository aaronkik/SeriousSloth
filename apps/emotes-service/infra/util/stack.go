package util

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
