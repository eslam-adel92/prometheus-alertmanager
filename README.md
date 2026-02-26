# Prometheus Monitoring Stack — Comprehensive Host, Network, Disk & Container Monitoring

A production-ready Docker Compose monitoring stack built with **Prometheus**, **Grafana**, **Alertmanager**, **Loki**, **Node Exporter**, and **cAdvisor**. Designed to be **template-ready** — easily pluggable into any Docker-based environment.

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Docker Host                             │
│                                                             │
│  ┌──────────┐  ┌──────────────┐  ┌──────────┐               │
│  │  Server  │  │ Node Exporter│  │ cAdvisor │               │
│  │(app:8090)│  │   (:9100)    │  │ (:8080)  │               │
│  └─────┬────┘  └──────┬───────┘  └────┬─────┘               │
│        │              │               │                     │
│        └──────────────┬───────────────┘                     │
│                       ▼                                     │
│               ┌────────────────┐                            │
│               │  Prometheus    │                            │
│               │   (:9090)      │                            │
│               └───────┬────────┘                            │
│                       │                                     │
│              ┌────────┴─────────┐                           │
│              ▼                  ▼                           │
│     ┌──────────────┐   ┌──────────────┐                     │
│     │ Alertmanager │   │   Grafana    │                     │
│     │   (:9093)    │   │   (:3000)    │                     │
│     └──────────────┘   └──────────────┘                     │
└─────────────────────────────────────────────────────────────┘
```

---

## 📦 Services

| Service           | Image                              | Port   | Purpose                                                   |
| ----------------- | ---------------------------------- | ------ | --------------------------------------------------------- |
| **server**        | Custom Go app                      | `8090` | Sample app exposing `/ping` counter + `/metrics` endpoint |
| **prometheus**    | `prom/prometheus:v2.45.3`          | `9090` | Metrics collection, storage, and alert evaluation         |
| **node-exporter** | `prom/node-exporter:v1.7.0`        | `9100` | Host-level metrics (CPU, memory, disk, network)           |
| **cadvisor**      | `gcr.io/cadvisor/cadvisor:v0.49.1` | `8080` | Per-container metrics (CPU, memory, network, I/O)         |
| **alertmanager**  | `prom/alertmanager:v0.27.0`        | `9093` | Alert routing, grouping, and notifications                |
| **grafana**       | `grafana/grafana:10.4.0`           | `3000` | Visualization dashboards (6 pre-provisioned)              |
| **loki**          | `grafana/loki:2.9.4`               | `3100` | Log aggregation, storage, and querying via LogQL          |
| **promtail**      | `grafana/promtail:2.9.4`           | —      | Log collector — tails container logs and pushes to Loki   |

---

## 🚀 Quick Start

```bash
# Start all services
docker compose up -d

# Check all services are running
docker compose ps

# View logs
docker compose logs -f
```

### Access the UIs

| UI                                   | URL                                                      |
| ------------------------------------ | -------------------------------------------------------- |
| **Grafana** (dashboards)             | [http://localhost:3000](http://localhost:3000)           |
| **Prometheus** (query & targets)     | [http://localhost:9090](http://localhost:9090)           |
| **Alertmanager** (alerts & silences) | [http://localhost:9093](http://localhost:9093)           |
| **Sample App** (fire test alerts)    | [http://localhost:8090/ping](http://localhost:8090/ping) |
| **cAdvisor** (container stats)       | [http://localhost:8080](http://localhost:8080)           |

---

## 📊 Grafana Dashboards

Five pre-provisioned dashboards are included:

### 1. Host Resources Overview

- **CPU** — Overall usage %, usage by mode (user/system/iowait), stat gauges
- **Memory** — Used vs available, total RAM, usage percentage
- **Load** — 1m/5m/15m load averages with CPU core count reference line
- **Uptime** — System uptime since last boot

### 2. Network Traffic

- **Bandwidth** — Receive/transmit bytes per second per physical interface
- **Errors & Drops** — Network errors and packet drops (indicators of hardware/congestion issues)
- **Packets** — Packet rates for deeper traffic analysis

### 3. Disk Monitoring

- **Space** — Usage % gauges, free space stats, usage trending over time
- **I/O** — Read/write throughput per block device
- **Utilization** — Disk busy % (saturation indicator)

### 4. Container Monitoring

- **Overview** — Running container count, total CPU/memory/network across all containers
- **Per-Container CPU & Memory** — Individual container resource usage
- **Per-Container Network** — Receive/transmit per container
- **Block I/O** — Disk read/write per container

### 5. Ping Request Count (Original)

- Tracks the demo `ping_request_count` metric from the sample app

### 6. Container Logs (Loki)

- **Log Volume** — Bar chart showing log line rate per container over time
- **Log Stream** — Live log viewer with filtering by server and container
- Supports `$server` and `$container` dropdowns for filtering

> **Template Note:** All dashboards use a `$datasource` template variable, so they can be imported into any Grafana instance with any Prometheus datasource.

---

## 🖥️ Remote Server Monitoring (Agent Setup)

To monitor **external servers**, deploy the lightweight agent stack on each remote server. The agent runs only Node Exporter + cAdvisor and exposes metrics for the central Prometheus to scrape.

```bash
# On the remote server:
cd agent/
cp .env.example .env     # Optional: adjust ports
docker compose up -d
```

Then add scrape targets in `prometheus/prometheus.yml` on the central server (see the commented-out examples at the bottom of the file).

👉 **Full setup guide:** [`agent/README.md`](./agent/README.md)

---

## 🔔 Alert Rules

Alert rules are defined in [`prometheus/rules.yml`](./prometheus/rules.yml), organized by category:

| Category      | Alert                           | Severity | Threshold                     |
| ------------- | ------------------------------- | -------- | ----------------------------- |
| **Host**      | `HostHighCpuUsage`              | warning  | CPU > 80% for 5m              |
| **Host**      | `HostCriticalCpuUsage`          | critical | CPU > 95% for 2m              |
| **Host**      | `HostHighMemoryUsage`           | warning  | Memory > 85% for 5m           |
| **Host**      | `HostHighLoad`                  | warning  | Load > CPU cores for 5m       |
| **Disk**      | `HostDiskSpaceLow`              | warning  | Disk > 85% full for 5m        |
| **Disk**      | `HostDiskSpaceCritical`         | critical | Disk > 95% full for 2m        |
| **Disk**      | `HostDiskIOHigh`                | warning  | I/O utilization > 80% for 10m |
| **Disk**      | `HostDiskWillFillIn4Hours`      | critical | Predicted to fill in 4h       |
| **Network**   | `HostNetworkReceiveErrors`      | warning  | > 10 errors/sec for 5m        |
| **Network**   | `HostNetworkTransmitErrors`     | warning  | > 10 errors/sec for 5m        |
| **Network**   | `HostNetworkInterfaceSaturated` | critical | > 80% of 1 Gbps for 5m        |
| **Network**   | `HostNetworkPacketDrops`        | warning  | > 100 drops/sec for 5m        |
| **Container** | `ContainerHighCpuUsage`         | warning  | CPU > 80% for 5m              |
| **Container** | `ContainerHighMemoryUsage`      | warning  | Memory > 85% limit for 5m     |
| **Container** | `ContainerKilled`               | critical | Not seen for 60s              |
| **App**       | `CountGreaterThan5`             | warning  | ping count > 5 (demo)         |

---

## 🔧 Customization Guide (Template Usage)

This stack is designed to be plugged into any Docker environment with minimal changes:

### 1. Add Your Applications

Replace or add to the `server` service in `docker-compose.yaml`. Any service exposing a Prometheus `/metrics` endpoint can be monitored.

### 2. Add Scrape Targets

Add new jobs in `prometheus/prometheus.yml`:

```yaml
# Example: Add a new API service
- job_name: my_api
  static_configs:
    - targets: ["my-api-service:8080"]
      labels:
        component: "application"
```

### 3. Configure Alert Notifications

Edit `alertmanager/alertmanager.yml` to set up your notification channels:

```yaml
# Example: Slack notifications
receivers:
  - name: slack_notifications
    slack_configs:
      - api_url: "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
        channel: "#alerts"
```

### 4. Adjust Alert Thresholds

All thresholds in `prometheus/rules.yml` can be adjusted to match your SLOs:

- CPU: 80%/95% → adjust for your baseline
- Disk: 85%/95% → adjust based on partition sizes
- Network: 100MB/s → adjust based on NIC speed (1G/10G/25G)

### 5. Environment Labels

Update `external_labels` in `prometheus/prometheus.yml` to identify your deployment:

```yaml
external_labels:
  environment: "production" # development, staging, production
  region: "us-east-1" # your cloud region
  cluster: "web-cluster-01" # your cluster name
```

---

## 📁 File Structure

```
.
├── docker-compose.yaml                    # Central server service definitions
├── Dockerfile                             # Sample Go app build
├── server.go                              # Sample app source
├── go.mod / go.sum                        # Go dependencies
├── agent/                                 # 🖥️ Remote server agent
│   ├── docker-compose.yaml                # Agent services (Node Exporter + cAdvisor + Promtail)
│   ├── .env.example                       # Configurable ports, LOKI_URL, SERVER_NAME
│   ├── README.md                          # Agent deployment guide
│   └── promtail/
│       └── promtail-config.yaml           # Remote Promtail config
├── prometheus/
│   ├── prometheus.yml                     # Scrape configs & targets
│   └── rules.yml                          # Alert rules (host/disk/net/container)
├── alertmanager/
│   └── alertmanager.yml                   # Alert routing & receivers
├── loki/
│   └── loki-config.yaml                   # Loki server config (storage, retention)
├── promtail/
│   └── promtail-config.yaml               # Local Promtail config (Docker log discovery)
└── grafana/
    └── provisioning/
        ├── datasources/
        │   ├── datasource.yaml             # Prometheus datasource config
        │   └── loki.yaml                   # Loki datasource config
        └── dashboards/
            ├── dashboard.yaml              # Dashboard provisioner config
            └── definitions/
                ├── ping_request_count.json  # App metrics dashboard
                ├── host_resources.json      # Host CPU/memory/load dashboard
                ├── network_traffic.json     # Network bandwidth/errors dashboard
                ├── disk_monitoring.json     # Disk space/I/O dashboard
                ├── container_monitoring.json # Container resources dashboard
                └── container_logs.json      # Container logs dashboard (Loki)
```

---

## 🛑 Stopping and Cleanup

```bash
# Stop all services (preserves data)
docker compose down

# Stop all services AND delete all stored data
docker compose down -v
```

---

## 📚 References

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Prometheus Alerting Rules](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/)
- [Alertmanager Configuration](https://prometheus.io/docs/alerting/latest/configuration/)
- [Node Exporter Metrics](https://prometheus.io/docs/guides/node-exporter/)
- [cAdvisor Documentation](https://github.com/google/cadvisor)
- [Grafana Provisioning](https://grafana.com/docs/grafana/latest/administration/provisioning/)
