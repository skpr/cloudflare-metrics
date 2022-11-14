package sync

import (
	"context"
	"fmt"
	"log"
	"time"

	awstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"

	"github.com/skpr/cloudflare-metrics/pkgs/types"
	"github.com/skpr/cloudflare-metrics/pkgs/util"
)

// MetricsSyncer defines the metrics syncer.
type MetricsSyncer struct {
	config            util.Config
	metricsCollectors []types.MetricsCollector
	metricsPusher     types.MetricsPusher
}

// NewMetricsSyncer creates a new metrics syncer.
func NewMetricsSyncer(config util.Config, metricsCollectors []types.MetricsCollector, metricsPusher types.MetricsPusher) *MetricsSyncer {
	return &MetricsSyncer{
		config:            config,
		metricsCollectors: metricsCollectors,
		metricsPusher:     metricsPusher,
	}
}

// Start calls each collector to collect cloudflare metrics then pushes them to aws cloudwatch metrics.
func (s *MetricsSyncer) Start(ctx context.Context, cancelFunc context.CancelFunc) {
	ticker := time.NewTicker(time.Second * time.Duration(s.config.PeriodSeconds))
	for {
		select {
		case <-ctx.Done():
			cancelFunc()
			fmt.Println("Stopped")
			return
		case <-ticker.C:
			end := time.Now()
			start := end.Add(-time.Second * time.Duration(s.config.PeriodSeconds))
			fmt.Println("Fetching cache stats from", start.Format(time.RFC3339), "to", end.Format(time.RFC3339))
			var data []awstypes.MetricDatum
			for _, collector := range s.metricsCollectors {
				d, err := collector.CollectMetrics(ctx, start, end)
				if err != nil {
					log.Println("failed to get data", err)
				}
				data = append(data, d...)
			}
			err := s.metricsPusher.Push(ctx, s.config.MetricsNamespace, data)
			if err != nil {
				log.Println("failed to push metrics", err)
			}
		}
	}
}
