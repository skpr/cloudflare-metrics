package cachestats

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/hasura/go-graphql-client"

	"github.com/skpr/cloudflare-metrics/pkgs/types"
	"github.com/skpr/cloudflare-metrics/pkgs/util"
)

// Collector defines a cache stat collector.
type Collector struct {
	types.MetricsCollector
	config util.Config
	client types.GraphQLClient
}

// NewCollector creates a new cache stat collector.
func NewCollector(config util.Config, client types.GraphQLClient) *Collector {
	return &Collector{
		config: config,
		client: client,
	}
}

// CollectMetrics calls the graphQL endpoint to collect cache stats.
func (c *Collector) CollectMetrics(ctx context.Context, start, end time.Time) ([]awstypes.MetricDatum, error) {
	fmt.Println("Fetching cache stat Metrics...")
	var q struct {
		Viewer struct {
			Zones []struct {
				Metrics []struct {
					Count      int32
					Dimensions struct {
						CacheStatus string
					}
				} `graphql:"metrics: httpRequestsAdaptiveGroups(limit: 16, filter: $filter, orderBy: [count_DESC])"`
			} `graphql:"zones(filter: {zoneTag: $zoneTag})"`
		}
	}

	variables := map[string]interface{}{
		"zoneTag": c.config.CloudFlareZoneTag,
		"filter": map[string]interface{}{
			"AND": []map[string]interface{}{
				{
					"datetime_geq": start.Format(time.RFC3339),
					"datetime_leq": end.Format(time.RFC3339),
				},
				{
					"requestSource": "eyeball",
				},
				{
					"clientRequestHTTPHost": c.config.CloudFlareHostName,
				},
			},
		},
	}
	err := c.client.Query(ctx, &q, variables, graphql.OperationName("CacheStatus"))
	if err != nil {
		return []awstypes.MetricDatum{}, err
	}

	var data []awstypes.MetricDatum
	for _, zone := range q.Viewer.Zones {
		for _, cacheStatus := range zone.Metrics {
			d := []awstypes.MetricDatum{
				{
					MetricName: aws.String("cacheStatusRequests"),
					Dimensions: []awstypes.Dimension{
						{
							Name:  aws.String("cacheStatus"),
							Value: aws.String(cacheStatus.Dimensions.CacheStatus),
						},
					},
					Timestamp: aws.Time(end),
					Value:     aws.Float64(float64(cacheStatus.Count)),
				},
			}
			data = append(data, d...)
		}
	}

	return data, nil
}
