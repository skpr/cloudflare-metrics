package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"

	"github.com/skpr/cloudflare-metrics/pkgs/collector/registry"
	"github.com/skpr/cloudflare-metrics/pkgs/pusher"
	"github.com/skpr/cloudflare-metrics/pkgs/sync"
	"github.com/skpr/cloudflare-metrics/pkgs/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	fmt.Println("Starting metrics sync every", config.FrequencySeconds, "seconds")

	// Handle interrupt signal gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	collectorRegistry := registry.NewCollectorRegistry(config)
	metricCollectors := collectorRegistry.GetCollectors(ctx)

	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to setup aws client: %d", err))
	}

	cloudwatchClient := cloudwatch.NewFromConfig(cfg)
	metricsPusher := pusher.NewPusher(cloudwatchClient)

	syncer := sync.NewMetricsSyncer(config, metricCollectors, metricsPusher)
	syncer.Start(ctx, stop)
}
