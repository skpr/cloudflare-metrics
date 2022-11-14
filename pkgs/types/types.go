package types

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	awstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

type MetricCollector interface {
	CollectMetrics(ctx context.Context, start, end time.Time) ([]awstypes.MetricDatum, error)
}

// CloudwatchInterface provides and interface for Cloudwatch.
type CloudwatchInterface interface {
	PutMetricData(ctx context.Context, params *cloudwatch.PutMetricDataInput, optFns ...func(options *cloudwatch.Options)) (*cloudwatch.PutMetricDataOutput, error)
}
