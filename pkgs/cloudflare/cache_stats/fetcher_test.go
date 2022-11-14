package cache_stats

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestFetch(t *testing.T) {
	// client := mock.NewMockClient()
	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "54jzo4JmCqn1ziFr6xNYps9OPot8WKJEN33k-7ox"},
	))
	client := graphql.NewClient("https://api.cloudflare.com/client/v4/graphql", httpClient).WithDebug(true)
	fetcher := NewFetcher(client)
	end := time.Now()
	start := end.Add(-time.Hour)

	stats, err := fetcher.FetchCacheStats(context.Background(), "a6f42912cf312b0b199d090dd500aab0", start, end)
	assert.NoError(t, err)
	fmt.Println(stats)
}
