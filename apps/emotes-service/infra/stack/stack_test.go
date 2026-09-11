package stack

import "testing"

func TestEmoteSyncSchedule(t *testing.T) {
	for name, want := range map[string]string{
		"staging": "cron(0 0 * * ? *)",
		"prod":    "cron(0 * * * ? *)",
		"dev":     "cron(0 * * * ? *)",
	} {
		if got := EmoteSyncSchedule(name); got != want {
			t.Errorf("EmoteSyncSchedule(%q) = %q, want %q", name, got, want)
		}
	}
}
