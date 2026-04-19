package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type StatefulComponent struct {
	pulumi.ResourceState
}

func NewStatefulComponent(ctx *pulumi.Context, name string, providerResource pulumi.ResourceOption) (*StatefulComponent, error) {
	component := &StatefulComponent{}

	err := ctx.RegisterComponentResourceV2("emotes-service:index:StatefulComponent", name, nil, component)
	if err != nil {
		return nil, err
	}

	return component, nil
}
