package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/micnncim/spanner-monitor/pkg/metrics"
)

var (
	projectID  = flag.String("project-id", "", "The project id for GCP project.")
	instanceID = flag.String("instance-id", "", "The instance id for Cloud Spanner")
)

func main() {
	flag.Parse()
	if *projectID == "" {
		fmt.Fprintln(os.Stderr, "project-id missing")
		os.Exit(1)
	}

	ctx := context.Background()

	client, err := metrics.NewClient(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create metrics client: %v\n", err)
		os.Exit(1)
	}

	if err := client.ReadMetrics(ctx, *projectID, *instanceID); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read metrics: %v\n", err)
		os.Exit(1)
	}
}
