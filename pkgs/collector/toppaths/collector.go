package toppaths

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
	fmt.Println("Fetching top paths metrics...")
	var q struct {
		Viewer struct {
			Zones []struct {
				ZoneTag  string
				TopPaths []struct {
					Avg struct {
						SampleInterval float64
					}
					Count      int32
					Dimensions struct {
						ClientRequestPath     string `graphql:"metric: clientRequestPath"`
						ClientRequestHTTPHost string `graphql:"clientRequestHTTPHost"`
					}
					Sum struct {
						EdgeResponseBytes int32
					}
				} `graphql:"httpRequestsAdaptiveGroups(limit: 15, filter: $filter, orderBy: [count_DESC])"`
			} `graphql:"zones(filter: {zoneTag_in: $zoneTags})"`
		}
	}
	variables := map[string]interface{}{
		"zoneTags": c.config.CloudFlareZoneTags,
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
	err := c.client.Query(ctx, &q, variables, graphql.OperationName("TopPaths"))
	if err != nil {
		return []awstypes.MetricDatum{}, err
	}

	var data []awstypes.MetricDatum
	for _, zone := range q.Viewer.Zones {
		for _, topPath := range zone.TopPaths {
			d := []awstypes.MetricDatum{
				{
					MetricName: aws.String("requests"),
					Dimensions: []awstypes.Dimension{
						{
							Name:  aws.String("requestPath"),
							Value: aws.String(topPath.Dimensions.ClientRequestPath),
						},
						{
							Name:  aws.String("host"),
							Value: aws.String(topPath.Dimensions.ClientRequestHTTPHost),
						},
						{
							Name:  aws.String("zone"),
							Value: aws.String(zone.ZoneTag),
						},
					},
					Timestamp: aws.Time(end),
					Value:     aws.Float64(float64(topPath.Count)),
				},
			}
			data = append(data, d...)
		}
	}
	return data, nil
}