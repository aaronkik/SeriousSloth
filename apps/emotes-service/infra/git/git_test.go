package git

import (
	"strings"
	"testing"
)

func TestGetRepositorySlug(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			"ssh url",
			"git@github.com:SeriousSloth/emotes-service.git",
			"SeriousSloth/emotes-service",
		},
		{
			"https url",
			"https://github.com/SeriousSloth/emotes-service.git",
			"SeriousSloth/emotes-service",
		},
		{
			"ssh url without .git suffix",
			"git@github.com:SeriousSloth/emotes-service",
			"SeriousSloth/emotes-service",
		},
		{
			"https url without .git suffix",
			"https://github.com/SeriousSloth/emotes-service",
			"SeriousSloth/emotes-service",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result, err := GetRepositorySlug(test.url)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != test.expected {
				t.Errorf("GetRepositorySlug(%q) = %q, want %q", test.url, result, test.expected)
			}
		})
	}
}

func TestGetRepositorySlugErrors(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{"empty string", "", "failed to parse URL"},
		{"missing path", "https://github.com", "missing owner/repo path"},
		{"invalid https url", "https://", "invalid URL"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			_, err := GetRepositorySlug(test.url)
			if err == nil {
				t.Fatalf("GetRepositorySlug(%q) expected error, got nil", test.url)
			}
			if !strings.Contains(err.Error(), test.expected) {
				t.Errorf("expected error containing %q, got %q", test.expected, err.Error())
			}
		})
	}
}
