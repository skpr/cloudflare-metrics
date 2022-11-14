package mock

import (
	"context"

	"github.com/hasura/go-graphql-client"

	"github.com/skpr/cloudflare-metrics/pkgs/collector/types"
)

type Client struct {
	types.GraphQLClient
}

func NewMockClient() *Client {
	return &Client{}
}

func (c Client) Query(ctx context.Context, q interface{}, variables map[string]interface{}, options ...graphql.Option) error {
	return nil
}
