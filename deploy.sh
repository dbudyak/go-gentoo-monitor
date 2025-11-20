#!/bin/bash

# Deployment script for Gentoo NAS Monitor
# Usage: ./deploy.sh USER@NAS_IP

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

if [ $# -eq 0 ]; then
    echo -e "${RED}Error: No server specified${NC}"
    echo "Usage: ./deploy.sh USER@NAS_IP"
    echo "Example: ./deploy.sh admin@192.168.1.100"
    exit 1
fi

SERVER=$1
REMOTE_DIR="~/gentoo-monitor"

echo -e "${GREEN}==> Deploying to ${SERVER}${NC}"

# Step 1: Create directory
echo -e "${YELLOW}[1/4] Creating remote directory...${NC}"
ssh ${SERVER} "mkdir -p ${REMOTE_DIR}"

# Step 2: Transfer files
echo -e "${YELLOW}[2/4] Transferring files...${NC}"
rsync -avz --exclude='.git' --exclude='.idea' --exclude='monitor' --exclude='*.log' \
    --exclude='.DS_Store' \
    ./ ${SERVER}:${REMOTE_DIR}/

# Detect docker compose command
echo -e "${YELLOW}[3/4] Detecting Docker Compose version...${NC}"
if ssh ${SERVER} "command -v docker-compose &> /dev/null"; then
    COMPOSE_CMD="docker-compose"
    echo "Using: docker-compose (legacy)"
elif ssh ${SERVER} "docker compose version &> /dev/null"; then
    COMPOSE_CMD="docker compose"
    echo "Using: docker compose (plugin)"
else
    echo -e "${RED}Error: Docker Compose not found on server${NC}"
    exit 1
fi

# Step 3: Build
echo -e "${YELLOW}[4/5] Building Docker image...${NC}"
ssh ${SERVER} "cd ${REMOTE_DIR} && ${COMPOSE_CMD} build"

# Step 4: Deploy
echo -e "${YELLOW}[5/5] Starting service...${NC}"
ssh ${SERVER} "cd ${REMOTE_DIR} && ${COMPOSE_CMD} up -d"

# Step 5: Verify
echo -e "${GREEN}==> Checking service status...${NC}"
sleep 2
ssh ${SERVER} "cd ${REMOTE_DIR} && ${COMPOSE_CMD} ps"

echo ""
echo -e "${GREEN}==> Deployment complete!${NC}"
echo -e "Access the dashboard at: ${GREEN}http://$(echo ${SERVER} | cut -d'@' -f2):8080${NC}"
echo ""
echo "Useful commands:"
echo "  View logs:    ssh ${SERVER} 'cd ${REMOTE_DIR} && ${COMPOSE_CMD} logs -f'"
echo "  Stop service: ssh ${SERVER} 'cd ${REMOTE_DIR} && ${COMPOSE_CMD} down'"
echo "  Restart:      ssh ${SERVER} 'cd ${REMOTE_DIR} && ${COMPOSE_CMD} restart'"
