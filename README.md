# certsuite-overview

CertSuite Overview is a CLI tool that consolidates data from Quay (image pulls), DCI (test suite runs), and the CertSuite Collector. It automatically updates a centralized MySQL database, offering real-time insights into the usage of CertSuite tools and services. By providing a single, unified view of key metrics, it simplifies monitoring and reporting for stakeholders.

## Prerequisites

- **Go 1.26+** (the project uses Go 1.26.1 with toolchain go1.26.3)
- **MySQL database** for storing consolidated metrics
- **golangci-lint** (for running the linter)
- **Quay credentials** (`CLIENTID`, `APISECRET`) for fetching image pull data
- **DCI credentials** (`BEARERTOKEN`) for fetching test suite run data

## Build

```bash
# Build all packages
make build

# Run static analysis
make vet

# Run linter
make lint
```

The build produces a `certsuite-overview` binary from `cmd/main.go`.

## Configuration

Configuration is loaded from environment variables using [Viper](https://github.com/spf13/viper). All variables are required.

| Variable | Description |
|---|---|
| `DB_USER` | MySQL database username |
| `DB_PASSWORD` | MySQL database password |
| `DB_URL` | MySQL database host/URL |
| `DB_PORT` | MySQL database port |
| `CLIENTID` | Quay API client ID |
| `APISECRET` | Quay API secret |
| `BEARERTOKEN` | DCI API bearer token |
| `NAMESPACE` | Quay namespace to query |
| `REPOSITORY` | Quay repository to query |

Example:

```bash
export DB_USER=certsuite
export DB_PASSWORD=secret
export DB_URL=localhost
export DB_PORT=3306
export CLIENTID=quay-client-id
export APISECRET=quay-api-secret
export BEARERTOKEN=dci-bearer-token
export NAMESPACE=redhat-best-practices-for-k8s
export REPOSITORY=certsuite
```

## Usage

The CLI uses [Cobra](https://github.com/spf13/cobra) for command management.

```bash
# Fetch data from Quay and DCI and store it in MySQL
./certsuite-overview fetch
```

The `fetch` command pulls image pull statistics from Quay and test suite run data from DCI, then writes the results to the configured MySQL database.

## Grafana Dashboard

The `grafana/` directory contains provisioning files for a Grafana dashboard that visualizes the collected data:

- **`grafana/dashboard/dashboard.json`** -- Dashboard definition with six panels:
  - Quay Pull Events Over Time
  - Quay Pull Events by Month
  - Quay Pull Events by Kind
  - DCI Test Runs Over Time
  - DCI Test Runs by Month
  - DCI Test Cases Ranked by Failures
- **`grafana/datasource/datasource.yaml`** -- MySQL datasource configuration template (replace `DB_URL`, `DB_PORT`, `DB_USER`, and `DB_PASSWORD` placeholders with actual values)

To use, import the dashboard JSON into your Grafana instance and configure the MySQL datasource to point at the same database used by the CLI.

## Key Features

1. **Quay Image Pull Tracking** -- Tracks and consolidates image pull data from Quay to monitor repository usage over time.
2. **DCI Test Suite Run Tracking** -- Partners running CertSuite via DCI (Distributed CI) have their test suite executions tracked.
3. **CertSuite Collector Data Visualization** -- Data collected from the CertSuite Collector is visualized through the Grafana dashboard, including historical tracking.
4. **Automated MySQL Integration** -- All data streams from Quay, DCI, and the Collector are automatically consolidated into a regularly updated MySQL database.

## Data Sources

| Source | Description | Library |
|---|---|---|
| Quay Image Pulls | Tracks image pull counts from the CertSuite Quay repository | [go-quay](https://github.com/sebrandon1/go-quay) |
| DCI Test Suite Runs | Monitors CertSuite test executions performed via DCI | [go-dci](https://github.com/sebrandon1/go-dci) |
| CertSuite Collector | Detailed test metrics and execution statistics | Grafana dashboard |
