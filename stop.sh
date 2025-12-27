#!/bin/bash

# Cloud Disk - Stop All Development Services

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo -e "${YELLOW}🛑 Stopping Cloud Dist services...${NC}\n"

# Stop Stripe CLI
if [ -f stripe_webhook.pid ]; then
    PID=$(cat stripe_webhook.pid)
    if ps -p $PID > /dev/null 2>&1; then
        kill $PID 2>/dev/null
        echo -e "${GREEN}✅ Stopped Stripe CLI (PID: $PID)${NC}"
    fi
    rm stripe_webhook.pid
fi
pkill -f "stripe listen" 2>/dev/null && echo -e "${GREEN}✅ Stopped remaining Stripe CLI processes${NC}"

# Stop Backend
if [ -f backend.pid ]; then
    PID=$(cat backend.pid)
    if ps -p $PID > /dev/null 2>&1; then
        kill $PID 2>/dev/null
        echo -e "${GREEN}✅ Stopped backend (PID: $PID)${NC}"
    fi
    rm backend.pid
fi
pkill -f "cloud-dist" 2>/dev/null && echo -e "${GREEN}✅ Stopped remaining backend processes${NC}"

# Stop Frontend
if [ -f frontend.pid ]; then
    PID=$(cat frontend.pid)
    if ps -p $PID > /dev/null 2>&1; then
        kill $PID 2>/dev/null
        echo -e "${GREEN}✅ Stopped frontend (PID: $PID)${NC}"
    fi
    rm frontend.pid
fi
pkill -f "vite" 2>/dev/null && echo -e "${GREEN}✅ Stopped remaining frontend processes${NC}"

echo -e "\n${GREEN}✅ All services stopped${NC}"

