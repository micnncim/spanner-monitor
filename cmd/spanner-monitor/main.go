package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/micnncim/spanner-monitor/pkg/metrics"
)

func main() {
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		fmt.Fprintln(os.Stderr, "project-id missing")
		os.Exit(1)
	}

	ctx := context.Background()

	client, err := metrics.NewClient(ctx, projectID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create metrics client: %v\n", err)
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := client.ReadMetrics(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "failed to read metrics: %v\n", err)
			os.Exit(1)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read metrics: %v\n", err)
		os.Exit(1)
	}
}
