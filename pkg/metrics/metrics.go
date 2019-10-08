package metrics

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/api/iterator"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// Client is a client for Stackdriver Monitoring.
type Client struct {
	metricClient *monitoring.MetricClient
	projectID    string
}

// NewClient returns a new Client.
func NewClient(ctx context.Context, projectID string) (*Client, error) {
	metricClient, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Client{
		metricClient: metricClient,
		projectID:    projectID,
	}, nil
}

// ReadMetrics reads time series metrics.
// https://cloud.google.com/monitoring/custom-metrics/reading-metrics?hl=ja#monitoring_read_timeseries_fields-go
func (c *Client) ReadMetrics(ctx context.Context) error {
	now := time.Now()
	startTime := now.UTC().Add(-time.Minute * 20)
	endTime := now.UTC()
	req := &monitoringpb.ListTimeSeriesRequest{
		Name: fmt.Sprintf("projects/%s", c.projectID),
		// TODO: Fix metrics type and enable to specify with argument.
		Filter: `metric.type="spanner.googleapis.com/instance/cpu/utilization"`,
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

		log.Printf("%#v\n", resp.GetMetric().Labels["database"])
		log.Printf("\tCPU Unitilization: %.4f\n", resp.GetPoints()[0].GetValue().GetDoubleValue())
	}

	return nil
}
