package statuscodes

import (
	"context"
	"fmt"
	"strconv"
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
	fmt.Println("Fetching status code metrics...")
	var q struct {
		Viewer struct {
			Zones []struct {
				Metrics []struct {
					Count      int32
					Dimensions struct {
						EdgeResponseStatus int `graphql:"status: edgeResponseStatus"`
					}
				} `graphql:"metrics: httpRequestsAdaptiveGroups(limit: 15, filter: $filter, orderBy: [edgeResponseStatus_ASC])"`
			} `graphql:"zones(filter: {zoneTag: $zoneTag})"`
		}
	}

	v := variables.NewBuilder().
		WithZoneTag(c.config.CloudFlareZoneTag).
		WithStart(start).
		WithEnd(end).
		WithHostnames(c.config.CloudFlareHostNames).
		Build()
	
	err := c.client.Query(ctx, &q, v, graphql.OperationName("TopPaths"))
	if err != nil {
		return []awstypes.MetricDatum{}, err
	}

	var data []awstypes.MetricDatum
	for _, zone := range q.Viewer.Zones {
		for _, metric := range zone.Metrics {
			d := []awstypes.MetricDatum{
				{
					MetricName: aws.String("statusCodeRequests"),
					Dimensions: []awstypes.Dimension{
						{
							Name:  aws.String("statusCode"),
							Value: aws.String(strconv.Itoa(metric.Dimensions.EdgeResponseStatus)),
						},
					},
					Timestamp: aws.Time(end),
					Value:     aws.Float64(float64(metric.Count)),
					Unit:      awstypes.StandardUnitCount,
				},
			}
			data = append(data, d...)
		}
	}
	fmt.Println("Generated", len(data), "status code metrics")

	return data, nil
}
