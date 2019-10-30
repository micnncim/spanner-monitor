package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/api/iterator"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// Client is a client for Stackdriver Monitoring.
type Client struct {
	metricClient *monitoring.MetricClient
}

// NewClient returns a new Client.
func NewClient(ctx context.Context) (*Client, error) {
	metricClient, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Client{
		metricClient: metricClient,
	}, nil
}

// ReadMetrics reads time series metrics.
// https://cloud.google.com/monitoring/custom-metrics/reading-metrics?hl=ja#monitoring_read_timeseries_fields-go
func (c *Client) ReadMetrics(ctx context.Context, projectID, instanceID string) error {
	now := time.Now()
	startTime := now.UTC().Add(-time.Minute * 20)
	endTime := now.UTC()
	filter := `
		metric.type = "spanner.googleapis.com/instance/cpu/utilization_by_priority" AND
		metric.label.priority = "high" AND
		resource.label.instance_id = "%s"
`

	req := &monitoringpb.ListTimeSeriesRequest{
		Name: fmt.Sprintf("projects/%s", projectID),
		// TODO: Fix metrics type and enable to specify with argument.
		Filter: fmt.Sprintf(filter, instanceID),
		Interval: &monitoringpb.TimeInterval{
			StartTime: &timestamp.Timestamp{
				Seconds: startTime.Unix(),
			},
			EndTime: &timestamp.Timestamp{
				Seconds: endTime.Unix(),
			},
		},
		View: monitoringpb.ListTimeSeriesRequest_FULL,
	}

	it := c.metricClient.ListTimeSeries(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		fmt.Printf("%s (priority=%s)\n", resp.GetMetric().Labels["database"], resp.GetMetric().Labels["priority"])
		// fmt.Printf("\tCPU Utilization: %.4f\n", resp.GetPoints()[0].GetValue().GetDoubleValue())
		fmt.Printf("\tCPU Utilization: %d\n", len(resp.GetPoints()))
		for _, p := range resp.GetPoints() {
			fmt.Printf("%.4f ", p.GetValue().GetDoubleValue())
		}
		fmt.Println()
	}

	return nil
}
