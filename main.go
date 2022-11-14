package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/hasura/go-graphql-client"
	"golang.org/x/oauth2"

	"github.com/skpr/cloudflare-metrics/pkgs/cloudflare/cache_stats"
	"github.com/skpr/cloudflare-metrics/pkgs/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("failed to load config:", err)
	}
	fmt.Println("Config:", config)

	ctx, done := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		// Block until a signal is received.
		s := <-c
		fmt.Println("Got signal:", s)
		done()
	}()

	httpClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.CloudFlareAPIToken},
	))

	graphQLClient := graphql.NewClient(config.CloudFlareEndpointURL, httpClient)
	cacheStatsFetcher := cache_stats.NewFetcher(graphQLClient)

	ticker := time.NewTicker(time.Second * time.Duration(config.FrequencySeconds))
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Cancelled")
			return
		case <-ticker.C:
			end := time.Now()
			start := end.Add(-time.Second * time.Duration(config.PeriodSeconds))
			fmt.Println("Fetching cache stats from", start.Format(time.RFC3339), "to", end.Format(time.RFC3339))
			cacheStats, err := cacheStatsFetcher.FetchCacheStats(ctx, config.CloudFlareZoneID, start, end)
			if err != nil {
				log.Println("failed to get cache stats", err)
			}
			fmt.Println("Cache stats:", cacheStats)
		}
	}
}
