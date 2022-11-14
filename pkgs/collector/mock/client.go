package mock

import (
	"context"

	"github.com/hasura/go-graphql-client"

	"github.com/skpr/cloudflare-metrics/pkgs/types"
)

// Client defines the mock client.
type Client struct {
	types.GraphQLClient
}

// NewMockClient creates a new mock graphql client.
func NewMockClient() *Client {
	return &Client{}
}

// Query implements the interface.
func (c Client) Query(ctx context.Context, q interface{}, variables map[string]interface{}, options ...graphql.Option) error {
	return nil
}
