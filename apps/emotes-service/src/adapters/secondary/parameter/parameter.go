package parameter

import (
	"context"
	"emotes-service/src/environment"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// GetSecret gets a secret from SSM
// It calls SSM via local host instead of using the SSM client in conjunction with the SSM Lambda Layer.
// The Lambda Layer provides caching.
// See: https://aws.amazon.com/blogs/compute/using-the-aws-parameter-and-secrets-lambda-extension-to-cache-parameters-and-secrets/
func GetSecret(ctx context.Context, nameOrArn string) (string, error) {
	u := url.URL{
		Scheme: "http",
		Host:   "localhost:2773",
		Path:   "/systemsmanager/parameters/get",
	}

	q := u.Query()
	q.Set("name", nameOrArn)
	q.Set("withDecryption", "true")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		slog.ErrorContext(ctx, "Error creating ssm request", "error", err)
		return "", err
	}

	awsSessionToken := environment.GetOrFatal("AWS_SESSION_TOKEN")
	req.Header.Set("X-Aws-Parameters-Secrets-Token", awsSessionToken)

	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: newrelic.NewRoundTripper(nil),
	}

	resp, err := client.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "Error ssm parameter", "error", err)
		return "", err
	}

	defer resp.Body.Close()

	output := &ssm.GetParameterOutput{}
	err = json.NewDecoder(resp.Body).Decode(output)

	if err != nil {
		slog.ErrorContext(ctx, "Error ssm parameter response", "error", err)
		return "", err
	}

	return *output.Parameter.Value, nil
}
