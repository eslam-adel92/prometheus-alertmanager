# Monitoring Agent

Lightweight monitoring agent that runs on remote servers and exposes host + container metrics for the central Prometheus server to scrape.

## Services

| Service           | Port   | Description                                 |
| ----------------- | ------ | ------------------------------------------- |
| **Node Exporter** | `9100` | Host CPU, memory, disk, network metrics     |
| **cAdvisor**      | `8080` | Per-container CPU, memory, disk I/O metrics |

## Quick Start

### 1. Deploy the Agent

```bash
# Copy this agent/ directory to the remote server, then:
cd agent/
cp .env.example .env     # Optional: adjust ports if 9100/8080 are in use
docker compose up -d
```

### 2. Verify Agent is Running

```bash
# Check services are up
docker compose ps

# Test Node Exporter
curl http://localhost:9100/metrics | head

# Test cAdvisor
curl http://localhost:8080/metrics | head
```

### 3. Register with Central Prometheus

On the **central server**, edit `prometheus/prometheus.yml` and add scrape targets for this agent:

```yaml
scrape_configs:
  # ... existing jobs ...

  # Remote Agent: <server-name>
  - job_name: remote_node_exporter_<server-name>
    scrape_interval: 10s
    static_configs:
      - targets: ["<REMOTE_SERVER_IP>:9100"]
        labels:
          server: "<server-name>"
          component: "infrastructure"

  - job_name: remote_cadvisor_<server-name>
    scrape_interval: 10s
    static_configs:
      - targets: ["<REMOTE_SERVER_IP>:8080"]
        labels:
          server: "<server-name>"
          component: "infrastructure"
    metric_relabel_configs:
      - source_labels: [__name__]
        regex: "container_(tasks_state|memory_failures_total)"
        action: drop
```

Then reload Prometheus:

```bash
# On the central server
docker compose restart prometheus
# Or if lifecycle API is enabled:
curl -X POST http://localhost:9090/-/reload
```

### 4. Verify in Grafana

Open the central Grafana at `http://<CENTRAL_SERVER>:3000` and check:

- **Prometheus → Targets**: The new remote targets should appear as `UP`
- **Host Resources** dashboard: Select the remote server from the dropdown (if templated)
- **Container Monitoring** dashboard: Remote containers should appear

## Port Conflicts

If ports `9100` or `8080` are already in use on the remote server, edit `.env`:

```bash
NODE_EXPORTER_PORT=9101
CADVISOR_PORT=8081
```

Then update the corresponding targets in the central `prometheus.yml` to match.

## Firewall Rules

Ensure the central Prometheus server can reach these ports on the remote server:

```bash
# Example: allow Prometheus to scrape (replace with your firewall tool)
sudo firewall-cmd --permanent --add-port=9100/tcp
sudo firewall-cmd --permanent --add-port=8080/tcp
sudo firewall-cmd --reload
```

## Stopping the Agent

```bash
docker compose down
```
