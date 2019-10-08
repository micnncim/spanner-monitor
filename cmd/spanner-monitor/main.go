package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/micnncim/spanner-monitor/pkg/metrics"
)

var (
	projectID = flag.String("project-id", "", "Project ID for GCP project.")
)

func main() {
	flag.Parse()
	if *projectID == "" {
		fmt.Fprintln(os.Stderr, "project-id missing")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	client, err := metrics.NewClient(ctx, *projectID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create metrics client: %v\n", err)
		os.Exit(1)
	}

	if err := client.ReadMetrics(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read metrics: %v\n", err)
		os.Exit(1)
	}
}
