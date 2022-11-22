package registry

import (
	"context"

	"github.com/hasura/go-graphql-client"
	"golang.org/x/oauth2"

	"github.com/skpr/cloudflare-metrics/pkgs/collector/cachestats"
	"github.com/skpr/cloudflare-metrics/pkgs/collector/statuscodes"
	"github.com/skpr/cloudflare-metrics/pkgs/collector/totals"
	"github.com/skpr/cloudflare-metrics/pkgs/types"
	"github.com/skpr/cloudflare-metrics/pkgs/util"
)

// CollectorRegistry defines the metrics collector registry.
type CollectorRegistry struct {
	config     util.Config
	collectors []types.MetricsCollector
}

// NewCollectorRegistry creates a new metrics collector registry.
func NewCollectorRegistry(config util.Config) *CollectorRegistry {
	return &CollectorRegistry{
		config: config,
	}
}

// GetCollectors gets the metrics collectors.
func (c *CollectorRegistry) GetCollectors(ctx context.Context) []types.MetricsCollector {
	httpClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.config.CloudFlareAPIToken},
	))

	graphQLClient := graphql.NewClient(c.config.CloudFlareEndpointURL, httpClient)

	cacheStatsCollector := cachestats.NewCollector(c.config, graphQLClient)
	statusCodesCollector := statuscodes.NewCollector(c.config, graphQLClient)
	totalsCollector := totals.NewCollector(c.config, graphQLClient)

	c.collectors = append(c.collectors, cacheStatsCollector, statusCodesCollector, totalsCollector)

	return c.collectors
}
