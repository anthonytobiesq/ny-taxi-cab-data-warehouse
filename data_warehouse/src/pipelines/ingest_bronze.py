import argparse
import sys
import requests
import dlt



@dlt.resource(name="yellow_trip_data", write_disposition="append")
def fetch_yellow_taxi_data(year: int, month: int):
    url = f"https://d37ci6vzurychx.cloudfront.net/trip-data/yellow_tripdata_{year:04d}-{month:02d}.parquet"
    print(f"Fetching data from: {url}")

    response = requests.get(url)
    if response.status_code != 200:
        raise FileNotFoundError(f"Failed to fetch data from {url}. Status code: {response.status_code}")

    yield url

def run_pipeline(year: int, month: int, db_path: str):
    pipeline = dlt.pipeline(
        pipeline_name="ny_taxi_bronze",
        destination=dlt.destinations.duckdb(db_path),
        dataset_name="bronze",
    )

    load_info = pipeline.run(fetch_yellow_taxi_data(year, month))
    print(load_info)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Ingest Yellow Taxi Parquet files into DuckDB.")
    parser.add_argument("--year", type=int, required=True, help="Year of the data to ingest (e.g., 2024).")
    parser.add_argument("--month", type=int, required=True, help="Month of the data to ingest (1-12).")
    parser.add_argument("--db-path", type=str, required=False, default="yellow_taxi_trip.duckdb", help="Path to the DuckDB database file.")

    args = parser.parse_args()

    try:
        run_pipeline(args.year, args.month, args.db_path)
    except Exception as e:
        print(f"Error occurred: {e}", file=sys.stderr)
        sys.exit(1)