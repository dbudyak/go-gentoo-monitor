# Gentoo NAS Monitor

A lightweight system monitoring tool for Gentoo Linux, built with Go and designed for NAS servers. Provides real-time metrics through a clean web interface.

## Features

- **Memory Usage**: Real-time memory statistics with visual indicators
- **Disk Space**: Multi-drive monitoring with usage percentages
- **System Info**: Kernel version and system uptime
- **Systemd Services**: Active services status tracking
- **Docker Containers**: Running containers with port mappings

## API

### GET /api/metrics

Returns JSON with all system metrics:

```json
{
  "memory": {...},
  "disks": [...],
  "uptime": 123456789000,
  "kernel": "6.1.67-gentoo",
  "services": [...],
  "containers": [...],
  "timestamp": "2024-01-01T00:00:00Z"
}
```