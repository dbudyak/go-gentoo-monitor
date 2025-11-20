# Gentoo NAS Monitor

A lightweight system monitoring tool for Gentoo Linux, built with Go and designed for NAS servers. Provides real-time metrics through a clean web interface.

## Features

- **Memory Usage**: Real-time memory statistics with visual indicators
- **Disk Space**: Multi-drive monitoring with usage percentages
- **System Info**: Kernel version and system uptime
- **Systemd Services**: Active services status tracking
- **Docker Containers**: Running containers with port mappings

## Quick Start

### Local Development

```bash
make run
# Or: go run ./cmd/monitor
```

### Deploy to Server

See [DEPLOY.md](DEPLOY.md) for complete deployment instructions.

**Quick deployment:**

```bash
# From Mac: Transfer files
rsync -avz --exclude='.git' ./ USER@NAS_IP:~/gentoo-monitor/

# On NAS: Deploy
ssh USER@NAS_IP "cd ~/gentoo-monitor && docker-compose up -d"
```

Access the dashboard at `http://NAS_IP:8080`

## Configuration

Environment variables:
- `PORT`: HTTP server port (default: 8080)

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

## Requirements

- Docker and Docker Compose
- Host access to systemd and Docker socket
- Read access to system filesystems

## Architecture

```
cmd/monitor/          # Application entry point
internal/
  ├── collector/      # System metrics collectors
  ├── metrics/        # Data type definitions
  └── server/         # HTTP server and frontend
```

## License

MIT
