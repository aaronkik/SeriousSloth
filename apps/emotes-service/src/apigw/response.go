package apigw

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
)

func JSONResponse(ctx context.Context, status int, body any) (events.APIGatewayProxyResponse, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		slog.ErrorContext(ctx, "marshall error", "error", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: `{"message":"internal error"}`}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(payload),
	}, nil
}
