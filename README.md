NYC Taxi Data Warehouse & Orchestrator (WIP)
---

```markdown
# NYC Taxi Data Warehouse & Orchestrator

A local-first, multi-language data engineering pipeline that ingests, cleans, and analyzes NYC TLC Trip Record Data. The project combines a Go-based orchestration CLI, a modular Python data ingestion engine (`dlt`), and an embedded DuckDB database within a unified Docker environment.

---

## Architecture Overview

```text
├── main.go (Go Orchestrator)
│     └── Triggers CLI Subprocesses
├── data_warehouse/ (Python Package)
│     ├── dlt Ingestion Pipeline (Parquet -> DuckDB)
│     └── DuckDB Analytics Layer
└── docker-compose.yml (Container Stack & Jupyter Service)

```

* **Orchestrator:** Written in Go (`golang:1.26-alpine`), handling execution parameters, date backfills, and container process lifecycles.
* **Data Processing:** Python 3.13 package (`data_warehouse`) using `dlt`, `duckdb`, `pandas`, and `pyarrow`.
* **Storage:** Local DuckDB database persisted to `./data/yellow_taxi_trip.duckdb`.
* **Environment Management:** Package dependencies managed via `uv`.

---

## Getting Started

### Prerequisites

* Docker Engine & Docker Compose
* `uv` (for local Python development)
* Go 1.26+ (if running the orchestrator outside Docker)

### Running via Docker Compose

1. **Build the container images:**
```bash
docker compose build

```


2. **Run an ingestion job for January 2024:**
```bash
docker compose run --rm orchestrator

```


3. **Spin up Jupyter Lab for data exploration:**
```bash
docker compose up notebook

```


Navigate to `http://localhost:8888` to query `yellow_taxi_trip.duckdb`.

---

## Repository Structure

```text
.
├── Dockerfile             # Multi-stage build (Go + Python/uv)
├── docker-compose.yml     # Container services setup
├── main.go                # Go orchestrator CLI entrypoint
├── data/                  # Mount point for local DuckDB files
└── data_warehouse/        # Dedicated Python package subfolder

```

```

---

### 2. Package `README.md`
*Place at `/data_warehouse/README.md`*

```markdown
# Data Warehouse Python Package (`data_warehouse`)

A standalone Python package containing the data extraction, transformation, and loading (ETL) logic for the NYC Taxi pipeline.

---

## Package Details

* **Build System:** `hatchling` (configured with `src/` layout)
* **Dependency Management:** `uv` (`pyproject.toml` and `uv.lock`)
* **Primary Libraries:** `dlt`, `duckdb`, `pandas`, `pyarrow`

---

## Local Development & Setup

### 1. Install Package Dependencies
From the `data_warehouse/` directory:

```bash
uv sync

```

### 2. Direct Execution

Run individual pipeline steps directly using `uv`:

```bash
# Ingest bronze raw data directly to DuckDB
uv run python -m src.pipelines.ingest_bronze --year 2024 --month 1 --db-path ../data/yellow_taxi_trip.duckdb

```

---

## Package Structure

```text
data_warehouse/
├── pyproject.toml
├── uv.lock
├── README.md
└── src/
    ├── __init__.py
    └── pipelines/
        ├── ingest_bronze.py
        └── transform_silver.py

```

```

<ElicitationsGroup message="Would you like to refine the documentation further?">
  <Elicitation label="Add Docker environment variables section" query="Add a configuration section to the root README explaining environment variables and volume mounts."/>
  <Elicitation label="Document Jupyter query notebooks" query="Create a starter Jupyter notebook guide showing how to connect to DuckDB and query the bronze/silver tables."/>
</ElicitationsGroup>

```