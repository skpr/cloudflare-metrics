package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	awstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/hasura/go-graphql-client"
	"golang.org/x/oauth2"

	"github.com/skpr/cloudflare-metrics/pkgs/collector/cache_stats"
	"github.com/skpr/cloudflare-metrics/pkgs/pusher"
	"github.com/skpr/cloudflare-metrics/pkgs/types"
	"github.com/skpr/cloudflare-metrics/pkgs/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("failed to load awsconfig:", err)
	}
	fmt.Println("Config:", config)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	httpClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.CloudFlareAPIToken},
	))

	var metricCollectors []types.MetricCollector
	graphQLClient := graphql.NewClient(config.CloudFlareEndpointURL, httpClient)
	cacheStatsCollector := cache_stats.NewCacheStatsCollector(graphQLClient, config)

	metricCollectors = append(metricCollectors, cacheStatsCollector)

	cfg, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("failed to setup aws client: %d", err))
	}

	cloudwatchClient := cloudwatch.NewFromConfig(cfg)
	metricsPusher := pusher.NewPusher(cloudwatchClient)

	ticker := time.NewTicker(time.Second * time.Duration(config.FrequencySeconds))
	for {
		select {
		case <-ctx.Done():
			stop()
			fmt.Println("Stopped")
			return
		case <-ticker.C:
			end := time.Now()
			start := end.Add(-time.Second * time.Duration(config.PeriodSeconds))
			fmt.Println("Fetching cache stats from", start.Format(time.RFC3339), "to", end.Format(time.RFC3339))
			var data []awstypes.MetricDatum
			for _, dataCollector := range metricCollectors {
				d, err := dataCollector.CollectMetrics(ctx, start, end)
				if err != nil {
					log.Println("failed to get data", err)
				}
				data = append(data, d...)
			}
			err := metricsPusher.Push(ctx, "Skpr/CloudFlare", data)
			if err != nil {
				log.Println("failed to push metrics", err)
			}
		}
	}
}
