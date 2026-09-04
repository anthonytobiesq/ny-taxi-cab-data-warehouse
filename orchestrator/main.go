package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type JobConfig struct {
	Year   int
	Month  int
	DBPath string
}

func runIngestJob(ctx context.Context, config JobConfig) error {
	log.Printf("[ORCHESTRATOR] Starting ingestion job for %04d-%02d...", config.Year, config.Month)

	// Path inside container matches WORKDIR /app + COPY data_warehouse/src /app/src
	scriptPath := "/app/src/pipelines/ingest_bronze.py"

	cmd := exec.CommandContext(
		ctx,
		"uv", "run", "python", scriptPath,
		"--year", fmt.Sprintf("%d", config.Year),
		"--month", fmt.Sprintf("%d", config.Month),
		"--db-path", config.DBPath,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	startTime := time.Now()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("job failed for %04d-%02d: %w", config.Year, config.Month, err)
	}

	log.Printf("[ORCHESTRATOR] Successfully ingested %04d-%02d in %s", config.Year, config.Month, time.Since(startTime))
	return nil
}

func main() {
	year := flag.Int("year", 2024, "Year to ingest")
	month := flag.Int("month", 1, "Month to ingest")
	dbPath := flag.String("db-path", "/app/data/yellow_taxi_trip.duckdb", "Path to DuckDB instance")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	config := JobConfig{
		Year:   *year,
		Month:  *month,
		DBPath: *dbPath,
	}

	if err := runIngestJob(ctx, config); err != nil {
		log.Fatalf("[FATAL ERROR] %v", err)
	}
}
