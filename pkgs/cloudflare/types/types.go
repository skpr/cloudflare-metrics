package types

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

type GraphQLClient interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}, options ...graphql.Option) error
}

type CacheStats struct {
	Status         string
	Requests       int32
	ResponseBytes  int32
	SampleInterval float64
}
