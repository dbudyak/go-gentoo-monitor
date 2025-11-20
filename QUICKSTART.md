# Quick Start Guide

## One-Command Deployment

```bash
./deploy.sh USER@NAS_IP
```

Example:
```bash
./deploy.sh admin@192.168.1.100
```

This script will:
1. Create directory on NAS
2. Transfer all files
3. Build Docker image
4. Start the service

## Manual Deployment (Step-by-Step)

### On Your Mac

```bash
# 1. Transfer files to NAS
cd /Users/dima/dev/go-gentoo-monitor
rsync -avz --exclude='.git' ./ USER@NAS_IP:~/gentoo-monitor/
```

### On Your NAS

```bash
# 2. SSH into NAS
ssh USER@NAS_IP

# 3. Navigate to project
cd ~/gentoo-monitor

# 4. Deploy with Docker
docker-compose up -d

# 5. Check status
docker-compose ps
docker-compose logs -f
```

### Access Dashboard

Open browser: `http://NAS_IP:8080`

## Common Commands

### On Your NAS

```bash
cd ~/gentoo-monitor

# View logs
docker-compose logs -f

# Stop service
docker-compose down

# Restart service
docker-compose restart

# Rebuild after changes
docker-compose build
docker-compose up -d

# Check status
docker-compose ps

# View resource usage
docker stats gentoo-monitor
```

### From Your Mac

```bash
# Quick update and redeploy
./deploy.sh USER@NAS_IP

# Check service remotely
ssh USER@NAS_IP "cd ~/gentoo-monitor && docker-compose ps"

# View logs remotely
ssh USER@NAS_IP "cd ~/gentoo-monitor && docker-compose logs --tail=50"
```

## Test API

```bash
# From NAS
curl http://localhost:8080/api/metrics | jq

# From Mac
curl http://NAS_IP:8080/api/metrics | jq
```

## Troubleshooting

### Service won't start

```bash
# Check Docker is running
sudo systemctl status docker

# Check logs for errors
docker-compose logs

# Rebuild clean
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

### Can't access from Mac

```bash
# Check firewall on NAS
sudo iptables -L -n | grep 8080

# Test from NAS first
curl http://localhost:8080
```

### Metrics not showing

```bash
# Check container has correct permissions
docker exec gentoo-monitor systemctl --version

# Verify volume mounts
docker inspect gentoo-monitor | grep Mounts -A 20
```

## Next Steps

- See [DEPLOY.md](DEPLOY.md) for detailed deployment options
- Configure auto-start on boot
- Set up reverse proxy with nginx
- Add SSL/TLS for HTTPS access
