package apigw

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func AnnotateRequest(ctx context.Context, req events.APIGatewayProxyRequest) {
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return
	}

	identity := req.RequestContext.Identity

	attrs := map[string]string{
		"request.http.method":        req.HTTPMethod,
		"request.http.path":          req.RequestContext.ResourcePath,
		"request.http.resource":      req.Resource,
		"request.id":                 req.RequestContext.RequestID,
		"request.identity.apiKeyId":  identity.APIKeyID,
		"request.identity.sourceIp":  identity.SourceIP,
		"request.identity.userAgent": identity.UserAgent,
	}

	for k, v := range attrs {
		if v == "" {
			continue
		}
		txn.AddAttribute(k, v)
	}

	for k, v := range req.PathParameters {
		if v == "" {
			continue
		}
		txn.AddAttribute("request.path."+k, v)
	}
}
