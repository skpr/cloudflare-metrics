package cache_stats

import (
	"context"
	"fmt"
	"time"

	"github.com/hasura/go-graphql-client"

	"github.com/skpr/cloudflare-metrics/pkg/cloudflare/types"
)

type Fetcher struct {
	client types.GraphQLClient
}

func NewFetcher(client types.GraphQLClient) *Fetcher {
	return &Fetcher{
		client: client,
	}
}

func (f *Fetcher) FetchCacheStats(ctx context.Context, zoneTag string, start, end time.Time) ([]types.CacheStats, error) {
	var q struct {
		Viewer struct {
			Zones []struct {
				CacheStatus []struct {
					Avg struct {
						SampleInterval float64
					}
					Count      int32
					Dimensions struct {
						CacheStatus string
					}
					Sum struct {
						EdgeResponseBytes int32
					}
				} `graphql:"httpRequestsAdaptiveGroups(limit: 50, filter: $filter, orderBy: [count_DESC]) "`
			} `graphql:"zones(filter: {zoneTag: $zoneTag})"`
		}
	}

	variables := map[string]interface{}{
		"zoneTag": zoneTag,
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
					"clientRequestHTTPHost": "bond.edu.au",
				},
			},
		},
	}
	query, err := graphql.ConstructQuery(&q, variables, graphql.OperationName("CacheStatusQuery"))
	if err != nil {
		return []types.CacheStats{}, err
	}
	fmt.Println(query)
	err = f.client.Query(ctx, &q, variables, graphql.OperationName("CacheStatusQuery"))
	if err != nil {
		return []types.CacheStats{}, err
	}

	var stats []types.CacheStats
	if len(q.Viewer.Zones) > 0 {
		zone := q.Viewer.Zones[0]
		for _, c := range zone.CacheStatus {
			stat := types.CacheStats{
				Status:         c.Dimensions.CacheStatus,
				Requests:       c.Count,
				ResponseBytes:  c.Sum.EdgeResponseBytes,
				SampleInterval: c.Avg.SampleInterval,
			}
			stats = append(stats, stat)
		}
	}

	return stats, nil
}
