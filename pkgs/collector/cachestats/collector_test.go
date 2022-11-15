package cachestats

import (
	"context"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/hasura/go-graphql-client"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"

	"github.com/skpr/cloudflare-metrics/pkgs/util"
)

func TestCollect(t *testing.T) {
	t.Skip("Do not test using real config")
	// client := mock.NewMockClient()
	config, err := util.LoadConfig("../../..")
	assert.NoError(t, err)
	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.CloudFlareAPIToken},
	))
	client := graphql.NewClient(config.CloudFlareEndpointURL, httpClient).WithDebug(true)

	fetcher := NewCacheStatsCollector(config, client)
	end := time.Now()
	start := end.Add(-time.Minute)

	data, err := fetcher.CollectMetrics(context.Background(), start, end)
	assert.NoError(t, err)
	spew.Dump(data)
}
