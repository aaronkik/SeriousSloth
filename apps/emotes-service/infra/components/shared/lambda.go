package components

import (
	"emotes-service/infra/util"
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type LambdaArgs struct {
	Code        pulumi.ArchiveInput
	MemorySize  pulumi.IntPtrInput
	Timeout     pulumi.IntPtrInput
	Environment pulumi.StringMap
}

type Lambda struct {
	pulumi.ResourceState
	Function *lambda.Function
	Role     *iam.Role
}

func NewLambda(ctx *pulumi.Context, name string, args *LambdaArgs, opts ...pulumi.ResourceOption) (*Lambda, error) {
	if args == nil {
		return nil, fmt.Errorf("args is required")
	}
	if args.Code == nil {
		return nil, fmt.Errorf("args.Code is required")
	}

	component := &Lambda{}
	err := ctx.RegisterComponentResourceV2("components-shared:index:Lambda", name, nil, component, opts...)
	if err != nil {
		return nil, err
	}

	logRetentionDays := 1
	if util.IsProduction(ctx.Stack()) {
		logRetentionDays = 7
	}

	lambdaLogGroup, err := cloudwatch.NewLogGroup(ctx, fmt.Sprintf("%s-log-group", name), &cloudwatch.LogGroupArgs{
		RetentionInDays: pulumi.Int(logRetentionDays),
	}, pulumi.Parent(component))
	if err != nil {
		return nil, err
	}

	assumeRolePolicy, err := iam.GetPolicyDocument(ctx, &iam.GetPolicyDocumentArgs{
		Statements: []iam.GetPolicyDocumentStatement{
			{
				Actions: []string{"sts:AssumeRole"},
				Principals: []iam.GetPolicyDocumentStatementPrincipal{
					{
						Type:        "Service",
						Identifiers: []string{"lambda.amazonaws.com"},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	lambdaRole, err := iam.NewRole(ctx, fmt.Sprintf("%s-role", name), &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(assumeRolePolicy.Json),
	}, pulumi.Parent(component))
	if err != nil {
		return nil, err
	}

	_, err = iam.NewRolePolicy(ctx, fmt.Sprintf("%s-log-policy", name), &iam.RolePolicyArgs{
		Role: lambdaRole.Name,
		Policy: lambdaLogGroup.Arn.ApplyT(func(arn string) (string, error) {
			policy := map[string]interface{}{
				"Version": "2012-10-17",
				"Statement": []map[string]interface{}{
					{
						"Effect": "Allow",
						"Action": []string{
							"logs:CreateLogStream",
							"logs:PutLogEvents",
						},
						"Resource": []string{fmt.Sprintf("%s:*", arn)},
					},
				},
			}
			policyJSON, marshallError := json.Marshal(policy)
			if marshallError != nil {
				return "", marshallError
			}

			return string(policyJSON), nil
		}).(pulumi.StringOutput),
	}, pulumi.Parent(component), pulumi.DependsOn([]pulumi.Resource{lambdaLogGroup}))
	if err != nil {
		return nil, err
	}

	memorySize := args.MemorySize
	if memorySize == nil {
		memorySize = pulumi.Int(512)
	}

	timeout := args.Timeout
	if timeout == nil {
		timeout = pulumi.Int(10)
	}

	logLevel := "INFO"
	if util.IsEphemeral(ctx.Stack()) {
		logLevel = "DEBUG"
	}

	envVars := pulumi.StringMap{
		"AWS_LAMBDA_LOG_LEVEL": pulumi.String(logLevel),
	}
	for k, v := range args.Environment {
		envVars[k] = v
	}

	function, err := lambda.NewFunction(ctx, name, &lambda.FunctionArgs{
		Role:        lambdaRole.Arn,
		Runtime:     pulumi.String("provided.al2023"),
		Code:        args.Code,
		Handler:     pulumi.String("bootstrap"),
		PackageType: pulumi.String("Zip"),
		MemorySize:  memorySize,
		Timeout:     timeout,
		Environment: &lambda.FunctionEnvironmentArgs{
			Variables: envVars,
		},
		Architectures: pulumi.StringArray{
			pulumi.String("arm64"),
		},
		LoggingConfig: &lambda.FunctionLoggingConfigArgs{
			ApplicationLogLevel: pulumi.String(logLevel),
			LogFormat:           pulumi.String("JSON"),
			LogGroup:            lambdaLogGroup.Name,
			SystemLogLevel:      pulumi.String(logLevel),
		},
	}, pulumi.Parent(component), pulumi.DependsOn([]pulumi.Resource{lambdaLogGroup}))
	if err != nil {
		return nil, err
	}

	component.Function = function
	component.Role = lambdaRole

	return component, nil
}
