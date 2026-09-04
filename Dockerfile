# --- Stage 1: Build Go Orchestrator ---
FROM golang:1.26-alpine AS go-builder
WORKDIR /build
COPY go.mod ./
COPY orchestrator/main.go ./
RUN go build -o /build/bin/orchestrator main.go

# --- Stage 2: Runtime Environment with uv & Python ---
FROM python:3.13-slim

# Install system build dependencies if C extensions need compilation
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Install uv binary from official distribution
COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/
WORKDIR /app

# Copy dependency manifests from data_warehouse folder
COPY data_warehouse/pyproject.toml data_warehouse/uv.lock ./
RUN uv sync --frozen --no-dev --no-install-project

# Copy orchestrator binary and ingestion logic
COPY --from=go-builder /build/bin/orchestrator /app/bin/orchestrator
COPY data_warehouse/src /app/src

# Create volume mount target for DuckDB storage
RUN mkdir -p /app/data

ENTRYPOINT ["/app/bin/orchestrator"]
CMD ["--year=2024", "--month=1", "--db-path=/app/data/yellow_taxi_trip.duckdb"]