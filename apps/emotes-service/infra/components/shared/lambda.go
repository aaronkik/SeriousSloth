package components

import (
	"emotes-service/infra/stack"
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type LambdaArgs struct {
	Code             pulumi.ArchiveInput
	MemorySize       pulumi.IntPtrInput
	Timeout          pulumi.IntPtrInput
	Environment      pulumi.StringMap
	PolicyStatements iam.GetPolicyDocumentStatementArray
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
	if stack.IsProduction(ctx.Stack()) {
		logRetentionDays = 7
	}

	lambdaLogGroup, err := cloudwatch.NewLogGroup(ctx, fmt.Sprintf("%s-log-group", name), &cloudwatch.LogGroupArgs{
		RetentionInDays: pulumi.Int(logRetentionDays),
	}, pulumi.Parent(component))
	if err != nil {
		return nil, err
	}

	lambdaRole, err := iam.NewRole(ctx, fmt.Sprintf("%s-role", name), &iam.RoleArgs{
		AssumeRolePolicy: iam.GetPolicyDocumentOutput(ctx, iam.GetPolicyDocumentOutputArgs{
			Statements: iam.GetPolicyDocumentStatementArray{
				&iam.GetPolicyDocumentStatementArgs{
					Actions: pulumi.StringArray{pulumi.String("sts:AssumeRole")},
					Principals: iam.GetPolicyDocumentStatementPrincipalArray{
						&iam.GetPolicyDocumentStatementPrincipalArgs{
							Type:        pulumi.String("Service"),
							Identifiers: pulumi.StringArray{pulumi.String("lambda.amazonaws.com")},
						},
					},
				},
			},
		}).Json(),
	}, pulumi.Parent(component))
	if err != nil {
		return nil, err
	}

	policyStatements := iam.GetPolicyDocumentStatementArray{
		&iam.GetPolicyDocumentStatementArgs{
			Effect: pulumi.String("Allow"),
			Actions: pulumi.StringArray{
				pulumi.String("logs:CreateLogStream"),
				pulumi.String("logs:PutLogEvents"),
			},
			Resources: pulumi.StringArray{
				pulumi.Sprintf("%s:*", lambdaLogGroup.Arn),
			},
		},
	}
	policyStatements = append(policyStatements, args.PolicyStatements...)

	_, err = iam.NewRolePolicy(ctx, fmt.Sprintf("%s-policy", name), &iam.RolePolicyArgs{
		Role: lambdaRole.Name,
		Policy: iam.GetPolicyDocumentOutput(ctx, iam.GetPolicyDocumentOutputArgs{
			Statements: policyStatements,
		}).Json(),
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
	if stack.IsEphemeral(ctx.Stack()) {
		logLevel = "DEBUG"
	}

	envVars := pulumi.StringMap{
		"AWS_LAMBDA_LOG_FORMAT": pulumi.String("JSON"),
		"AWS_LAMBDA_LOG_LEVEL":  pulumi.String(logLevel),
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
