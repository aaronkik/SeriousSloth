package environment

import (
	"log/slog"
	"os"
)

func GetOrFatal(env string) string {
	envVar := os.Getenv(env)
	if envVar == "" {
		slog.Error("Environment variable not set", "env", env)
		os.Exit(1)
	}
	return envVar
}
