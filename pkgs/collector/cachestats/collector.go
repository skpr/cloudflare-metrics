package cachestats

import (
	"context"
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

// NewCacheStatsCollector creates a new cache stat collector.
func NewCacheStatsCollector(config util.Config, client types.GraphQLClient) *Collector {
	return &Collector{
		config: config,
		client: client,
	}
}

// CollectMetrics calls the graphQL endpoint to collect cache stats.
func (c *Collector) CollectMetrics(ctx context.Context, start, end time.Time) ([]awstypes.MetricDatum, error) {
	var q struct {
		Viewer struct {
			Zones []struct {
				CacheStatus []struct {
					Avg struct {
						SampleInterval float64
					}
					Count      int32
					Dimensions struct {
						CacheStatus           string
						ClientRequestHTTPHost string `graphql:"clientRequestHTTPHost"`
					}
					Sum struct {
						EdgeResponseBytes int32
					}
				} `graphql:"httpRequestsAdaptiveGroups(limit: 50, filter: $filter, orderBy: [count_DESC])"`
			} `graphql:"zones(filter: {zoneTag: $zoneTag})"`
		}
	}

	variables := map[string]interface{}{
		"zoneTag": c.config.CloudFlareZoneID,
		"filter": map[string]interface{}{
			"AND": []map[string]interface{}{
				{
					"datetime_geq": start.Format(time.RFC3339),
					"datetime_leq": end.Format(time.RFC3339),
				},
				{
					"requestSource": "eyeball",
				},
			},
		},
	}
	err := c.client.Query(ctx, &q, variables, graphql.OperationName("CacheStatusQuery"))
	if err != nil {
		return []awstypes.MetricDatum{}, err
	}

	var data []awstypes.MetricDatum
	for _, zone := range q.Viewer.Zones {
		for _, cacheStatus := range zone.CacheStatus {
			d := []awstypes.MetricDatum{
				{
					MetricName: aws.String("requests"),
					Dimensions: []awstypes.Dimension{
						{
							Name:  aws.String("cacheStatus"),
							Value: aws.String(cacheStatus.Dimensions.CacheStatus),
						},
						{
							Name:  aws.String("host"),
							Value: aws.String(cacheStatus.Dimensions.ClientRequestHTTPHost),
						},
						{
							Name:  aws.String("zone"),
							Value: aws.String(c.config.CloudFlareZoneID),
						},
					},
					Timestamp: aws.Time(end),
					Value:     aws.Float64(float64(cacheStatus.Count)),
				},
				{
					MetricName: aws.String("responseBytes"),
					Dimensions: []awstypes.Dimension{
						{
							Name:  aws.String("cacheStatus"),
							Value: aws.String(cacheStatus.Dimensions.CacheStatus),
						},
						{
							Name:  aws.String("host"),
							Value: aws.String(cacheStatus.Dimensions.ClientRequestHTTPHost),
						},
						{
							Name:  aws.String("zone"),
							Value: aws.String(c.config.CloudFlareZoneID),
						},
					},
					Timestamp: aws.Time(end),
					Value:     aws.Float64(float64(cacheStatus.Sum.EdgeResponseBytes)),
				},
			}
			data = append(data, d...)
		}
	}

	return data, nil
}
