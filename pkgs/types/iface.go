package types

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	awstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/hasura/go-graphql-client"
)

// MetricsCollector provides an interface for a metrics collector.
type MetricsCollector interface {
	CollectMetrics(ctx context.Context, start, end time.Time) ([]awstypes.MetricDatum, error)
}

// GraphQLClient provides an interface for the graphql client.
type GraphQLClient interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}, options ...graphql.Option) error
}

// MetricsPusher provides an interface for a metrics pusher.
type MetricsPusher interface {
	Push(ctx context.Context, namespace string, metricData []awstypes.MetricDatum) error
}

// CloudwatchInterface provides an interface for Cloudwatch.
type CloudwatchInterface interface {
	PutMetricData(ctx context.Context, params *cloudwatch.PutMetricDataInput, optFns ...func(options *cloudwatch.Options)) (*cloudwatch.PutMetricDataOutput, error)
}
