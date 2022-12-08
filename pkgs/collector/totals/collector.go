package totals

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/hasura/go-graphql-client"

	"github.com/skpr/cloudflare-metrics/pkgs/collector/variables"
	"github.com/skpr/cloudflare-metrics/pkgs/types"
	"github.com/skpr/cloudflare-metrics/pkgs/util"
)

// Collector defines a collector.
type Collector struct {
	types.MetricsCollector
	config util.Config
	client types.GraphQLClient
}

// NewCollector creates a new collector.
func NewCollector(config util.Config, client types.GraphQLClient) *Collector {
	return &Collector{
		config: config,
		client: client,
	}
}

// CollectMetrics calls the graphQL endpoint to collect metrics.
func (c *Collector) CollectMetrics(ctx context.Context, start, end time.Time) ([]awstypes.MetricDatum, error) {
	fmt.Println("Fetching request totals metrics...")
	var q struct {
		Viewer struct {
			Zones []struct {
				Metrics []struct {
					Count int32
					Sum   struct {
						EdgeResponseBytes int64
					}
				} `graphql:"metrics: httpRequestsAdaptiveGroups(limit: 1, filter: $filter)"`
			} `graphql:"zones(filter: {zoneTag: $zoneTag})"`
		}
	}

	v := variables.NewBuilder().
		WithZoneTag(c.config.CloudFlareZoneTag).
		WithStart(start).
		WithEnd(end).
		WithHostnames(c.config.CloudFlareHostNames).
		Build()

	err := c.client.Query(ctx, &q, v, graphql.OperationName("Metrics"))
	if err != nil {
		return []awstypes.MetricDatum{}, err
	}

	var data []awstypes.MetricDatum
	for _, zone := range q.Viewer.Zones {
		for _, total := range zone.Metrics {
			d := []awstypes.MetricDatum{
				{
					MetricName: aws.String("totalRequests"),
					Timestamp:  aws.Time(end),
					Value:      aws.Float64(float64(total.Count)),
					Unit:       awstypes.StandardUnitCount,
				},
				{
					MetricName: aws.String("totalResponseBytes"),
					Timestamp:  aws.Time(end),
					Value:      aws.Float64(float64(total.Sum.EdgeResponseBytes)),
					Unit:       awstypes.StandardUnitBytes,
				},
			}
			data = append(data, d...)
		}
	}
	fmt.Println("Generated", len(data), "total metrics")

	return data, nil
}
