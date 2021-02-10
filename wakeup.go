package main

import (
	"context"
	"fmt"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/operation"
	ycsdk "github.com/yandex-cloud/go-sdk"
)

func startComputeInstance(ctx context.Context, sdk *ycsdk.SDK, id string) (*operation.Operation, error) {
	// Run Compute Instance with ID
	return sdk.Compute().Instance().Start(ctx, &compute.StartInstanceRequest{
		InstanceId: id,
	})
}

type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

var (
	folderID = "yourFolderID"
)

func startComputeInstances(ctx context.Context) (*Response, error) {
	// SDK auth via service account
	sdk, err := ycsdk.Build(ctx, ycsdk.Config{
		// Call InstanceServiceAccount get IAM-token and form auth in SDK
		Credentials: ycsdk.InstanceServiceAccount(),
	})
	if err != nil {
		return nil, err
	}
	// Get Compute Instances from folderID
	listInstancesResponse, err := sdk.Compute().Instance().List(ctx, &compute.ListInstancesRequest{
		FolderId: folderID,
	})
	if err != nil {
		return nil, err
	}
	instances := listInstancesResponse.GetInstances()
	count := 0
	// Filter Compute Instance in Running state
	for _, i := range instances {
		if i.Status != compute.Instance_RUNNING {
			_, err := startComputeInstance(ctx, sdk, i.GetId())
			if err != nil {
				return nil, err
			}
			count++
		}
	}
	return &Response{
		StatusCode: 200,
		Body:       fmt.Sprintf("Started %d instances", count),
	}, nil
}
