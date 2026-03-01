# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Security**: Added `.env.example` mapping to `docker-compose.yaml` to secure Grafana default passwords (`GF_SECURITY_ADMIN_PASSWORD`).
- **Alertmanager**: Added `alertmanager_data` volume to persist notification state and silences across container restarts.
- **Alertmanager**: Configured `--config.expand-env=true` to interpolate `SLACK_WEBHOOK_URL` natively inside `alertmanager.yml`.
- **Docker Compose**: Added native Prometheus, Alertmanager, and Grafana `healthcheck` directives to avoid start-up race conditions.
- **Demo App**: Added Prometheus `Gauge` (`active_requests`) and `Histogram` (`request_duration_seconds`) metric types to `server.go` to provide richer demonstration capabilities.
- **Logging**: Integrated Loki 3.x and Promtail for centralized container log aggregation, along with a pre-provisioned "Container Logs" Grafana dashboard.
- **Versioning**: Added `CHANGELOG.md` for proper tracking of repository updates.

### Changed

- **Alertmanager**: Optimized routing in `alertmanager.yml` to group alerts by name, cluster, and service, with proper timeouts for critical versus warning severities.
- **Docker Compose**: Updated all service images to their latest stable version tags, abandoning the deprecated `latest` tag strategy.
- **Repository**: Expanded `.gitignore` to comprehensively ignore compiled binaries, `.env` secure files, and all persistent data volumes.
