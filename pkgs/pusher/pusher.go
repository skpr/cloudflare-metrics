package pusher

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	awstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"

	"github.com/skpr/cloudflare-metrics/pkgs/types"
)

// Pusher the metrics pusher.
type Pusher struct {
	cloudwatchClient types.CloudwatchInterface
}

// NewPusher creates a new metrics pusher.
func NewPusher(cloudwatchClient types.CloudwatchInterface) *Pusher {
	return &Pusher{
		cloudwatchClient: cloudwatchClient,
	}
}

// Push the metrics.
func (p *Pusher) Push(ctx context.Context, namespace string, metricData []awstypes.MetricDatum) error {
	_, err := p.cloudwatchClient.PutMetricData(ctx, &cloudwatch.PutMetricDataInput{
		MetricData: metricData,
		Namespace:  aws.String(namespace),
	})
	if err != nil {
		return err
	}
	return nil
}
