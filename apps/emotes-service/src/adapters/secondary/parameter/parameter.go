package parameter

import (
	"context"
	"log"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

var ssmClient *ssm.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	ssmClient = ssm.NewFromConfig(cfg)
}

func GetSecret(ctx context.Context, nameOrArn string) (string, error) {
	output, err := ssmClient.GetParameter(ctx,
		&ssm.GetParameterInput{
			Name:           aws.String(nameOrArn),
			WithDecryption: aws.Bool(true),
		})
	if err != nil {
		return "", err
	}

	slog.Info("Got secret value", "nameOrArn", nameOrArn)
	return *output.Parameter.Value, nil
}
