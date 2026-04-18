package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type StatefulComponent struct {
	pulumi.ResourceState
	EmotesBucket *s3.Bucket
}

func NewStatefulComponent(ctx *pulumi.Context, name string, providerResource pulumi.ResourceOption) (*StatefulComponent, error) {
	component := &StatefulComponent{}

	err := ctx.RegisterComponentResourceV2("emotes-service:index:StatefulComponent", name, nil, component)
	if err != nil {
		return nil, err
	}

	bucket, err := s3.NewBucket(ctx, "my-bucket", nil, pulumi.Parent(component), providerResource)
	if err != nil {
		return nil, err
	}
	component.EmotesBucket = bucket

	return component, nil
}
