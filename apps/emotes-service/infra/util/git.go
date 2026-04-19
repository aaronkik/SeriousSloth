package util

import (
	"fmt"
	"net/url"
	"os/exec"
	"strings"
)

func GitCli(args ...string) string {
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return "Unknown"
	}
	if s := strings.TrimSpace(string(out)); s != "" {
		return s
	}
	return "Unknown"
}

// GetRepositorySlug extracts "owner/repo" from SSH or HTTPS git remote URLs.
func GetRepositorySlug(remoteUrl string) (string, error) {
	remoteUrl = strings.TrimSuffix(remoteUrl, ".git")

	// SSH format
	var isSshUrl = strings.Contains(remoteUrl, "git@")
	if isSshUrl {
		slug := remoteUrl[strings.Index(remoteUrl, ":")+1:]
		if slug == "" || !strings.Contains(slug, "/") {
			return "", fmt.Errorf("invalid URL %q: missing owner/repo path", remoteUrl)
		}
		return slug, nil
	}

	// HTTPS format
	parsedUrl, err := url.ParseRequestURI(remoteUrl)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL %q: %w", remoteUrl, err)
	}

	slug := strings.TrimPrefix(parsedUrl.Path, "/")
	if slug == "" || !strings.Contains(slug, "/") {
		return "", fmt.Errorf("invalid URL %q: missing owner/repo path", remoteUrl)
	}
	return slug, nil
}
