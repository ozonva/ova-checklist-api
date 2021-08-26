package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"ova-checklist-api/pkg/service"
)

type Client interface {
	service.ChecklistStorageClient

	Close() error
}

type client struct {
	connection *grpc.ClientConn
	impl       service.ChecklistStorageClient
}

func (c *client) CreateChecklist(ctx context.Context, in *service.CreateChecklistRequest, opts ...grpc.CallOption) (*service.CreateChecklistResponse, error) {
	return c.impl.CreateChecklist(ctx, in, opts...)
}

func (c *client) DescribeChecklist(ctx context.Context, in *service.DescribeChecklistRequest, opts ...grpc.CallOption) (*service.DescribeChecklistResponse, error) {
	return c.impl.DescribeChecklist(ctx, in, opts...)
}

func (c *client) ListChecklists(ctx context.Context, in *service.ListChecklistsRequest, opts ...grpc.CallOption) (*service.ListChecklistsResponse, error) {
	return c.impl.ListChecklists(ctx, in, opts...)
}

func (c *client) RemoveChecklist(ctx context.Context, in *service.RemoveChecklistRequest, opts ...grpc.CallOption) (*service.RemoveChecklistResponse, error) {
	return c.impl.RemoveChecklist(ctx, in, opts...)
}

func (c *client) Close() error {
	if c.connection != nil {
		return c.connection.Close()
	}
	return nil
}

func NewClient(host string, port uint16) (Client, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	connection, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	return &client{
		connection: connection,
		impl:       service.NewChecklistStorageClient(connection),
	}, nil
}
