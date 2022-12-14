# Cloudflare Metrics

This application queries CloudFlare Analytics API to get key metrics, and pushes them to AWS CloudWatch Metrics.

Current metric collectors are:

* Cache Statistics

## Configuration

The application requires the following environment variables be set:

* `CLOUDFLARE_ENDPOINT_URL` The CloudFlare graphql endpoint e.g. https://api.cloudflare.com/client/v4/graphql
* `CLOUDFLARE_API_TOKEN` Your CloudFlare API token with `ReadAnalytics` permissions
* `CLOUDFLARE_ZONE_TAG` The Zone Tag to query.
* `CLOUDFLARE_HOSTNAME` The hostname to filter by.
* `PERIOD_SECONDS` The number of seconds between metric collection (minimum 60 seconds).
* `METRICS_NAMESPACE` The AWS CloudWatch Metric namespace to use.

## Development

Copy `local.env.dist` to `local.env` and set variables.

Use the following make commands:

```
make build
make lint
make vet
make test
```

## Releasing

The application is released as a docker image when a new tag is created.

You can test the release process locally using:
```
goreleaser build --snapshot --rm-dist
```
