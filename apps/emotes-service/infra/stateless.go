package main

import (
	"emotes-service/infra/components/shared"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type StatelessComponent struct {
	pulumi.ResourceState
}

func NewStatelessComponent(ctx *pulumi.Context, name string, providerResource pulumi.ResourceOption) (*StatelessComponent, error) {
	component := &StatelessComponent{}

	err := ctx.RegisterComponentResourceV2("emotes-service:index:StatelessComponent", name, nil, component)
	if err != nil {
		return nil, err
	}

	_, err = components.NewLambda(ctx, "sync-global-emotes", &components.LambdaArgs{
		Code: pulumi.NewAssetArchive(map[string]interface{}{
			"bootstrap": pulumi.NewFileAsset("../dist/sync-global-emotes/bootstrap"),
		}),
	})
	if err != nil {
		return nil, err
	}

	return component, nil
}
