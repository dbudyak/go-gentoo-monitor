# Deployment Guide

## Deploy from Mac to Gentoo NAS

### Step 1: Prepare the Code

On your Mac, ensure all files are ready:

```bash
cd /Users/dima/dev/go-gentoo-monitor

# Verify project structure
ls -la

# Optional: Test build locally (requires Go)
go build -o monitor ./cmd/monitor
```

### Step 2: Transfer to NAS Server

Replace `USER` and `NAS_IP` with your server credentials:

```bash
# Create directory on NAS
ssh USER@NAS_IP "mkdir -p ~/gentoo-monitor"

# Transfer all files (excluding .git and binaries)
rsync -avz --exclude='.git' --exclude='monitor' \
  /Users/dima/dev/go-gentoo-monitor/ \
  USER@NAS_IP:~/gentoo-monitor/

# Or use scp for simple copy
scp -r /Users/dima/dev/go-gentoo-monitor USER@NAS_IP:~/gentoo-monitor
```

### Step 3: SSH into Your NAS

```bash
ssh USER@NAS_IP
cd ~/gentoo-monitor
```

### Step 4: Build and Deploy with Docker

```bash
# Build the Docker image
docker-compose build

# Start the service
docker-compose up -d

# Check if running
docker-compose ps

# View logs
docker-compose logs -f
```

### Step 5: Verify Deployment

```bash
# Check if service is accessible
curl http://localhost:8080/api/metrics

# Or from your Mac
curl http://NAS_IP:8080/api/metrics
```

Open in browser: `http://NAS_IP:8080`

### Step 6: Manage the Service

```bash
# Stop the service
docker-compose down

# Restart after changes
docker-compose restart

# View resource usage
docker stats gentoo-monitor

# Update after code changes
docker-compose build --no-cache
docker-compose up -d
```

## Troubleshooting

### Permission Issues

If you get permission errors:

```bash
# Ensure Docker socket is accessible
sudo chmod 666 /var/run/docker.sock

# Or add user to docker group
sudo usermod -aG docker $USER
```

### Port Already in Use

Change port in `docker-compose.yml`:

```yaml
ports:
  - "8081:8080"  # Use 8081 instead
```

### Systemd Services Not Showing

Ensure volume mounts are correct:

```bash
# Check if systemd is accessible in container
docker exec gentoo-monitor systemctl list-units --type=service
```

## Auto-start on Boot (Optional)

### Option 1: Docker Compose

```bash
# Enable Docker service
sudo systemctl enable docker

# Service will auto-start with docker-compose restart policy
```

### Option 2: Systemd Service

Create `/etc/systemd/system/gentoo-monitor.service`:

```ini
[Unit]
Description=Gentoo NAS Monitor
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/home/USER/gentoo-monitor
ExecStart=/usr/bin/docker-compose up -d
ExecStop=/usr/bin/docker-compose down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable gentoo-monitor
sudo systemctl start gentoo-monitor
```

## Firewall Configuration

If firewall is enabled:

```bash
# iptables
sudo iptables -A INPUT -p tcp --dport 8080 -j ACCEPT
sudo iptables-save > /etc/iptables/rules.v4

# Or nftables
sudo nft add rule inet filter input tcp dport 8080 accept
```

## Updates and Maintenance

```bash
# Update from Mac
cd /Users/dima/dev/go-gentoo-monitor
rsync -avz --exclude='.git' --exclude='monitor' \
  ./ USER@NAS_IP:~/gentoo-monitor/

# Rebuild on NAS
ssh USER@NAS_IP "cd ~/gentoo-monitor && docker-compose build && docker-compose up -d"
```

## Security Recommendations

1. **Restrict Access**: Use firewall to limit access to trusted IPs
2. **Reverse Proxy**: Consider nginx with authentication
3. **HTTPS**: Add SSL certificate for encrypted access
4. **Monitoring**: Check logs regularly for unusual activity

```bash
# View logs
docker-compose logs --tail=100

# Monitor in real-time
docker-compose logs -f --tail=50
```
