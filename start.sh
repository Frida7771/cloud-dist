#!/bin/bash

# Cloud Disk - One-Click Development Startup Script
# Starts: Backend (Go), Frontend (Vite), Stripe CLI Webhook Listener

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Cleanup function
cleanup() {
    echo -e "\n${YELLOW}🛑 Shutting down services...${NC}"
    
    # Kill Stripe CLI
    if [ -f stripe_webhook.pid ]; then
        kill $(cat stripe_webhook.pid) 2>/dev/null || true
        rm stripe_webhook.pid
    fi
    pkill -f "stripe listen" 2>/dev/null || true
    
    # Kill backend
    if [ -f backend.pid ]; then
        kill $(cat backend.pid) 2>/dev/null || true
        rm backend.pid
    fi
    pkill -f "cloud-dist" 2>/dev/null || true
    
    # Kill frontend
    if [ -f frontend.pid ]; then
        kill $(cat frontend.pid) 2>/dev/null || true
        rm frontend.pid
    fi
    pkill -f "vite" 2>/dev/null || true
    
    echo -e "${GREEN}✅ All services stopped${NC}"
    exit 0
}

# Trap Ctrl+C
trap cleanup INT TERM

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   Cloud Dist Development Startup      ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}\n"

# Check if services are already running
check_service() {
    local service_name=$1
    local pid_file=$2
    local process_pattern=$3
    
    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        if ps -p $pid > /dev/null 2>&1; then
            echo -e "${YELLOW}⚠️  $service_name is already running (PID: $pid)${NC}"
            return 1
        else
            rm "$pid_file"
        fi
    fi
    
    if pgrep -f "$process_pattern" > /dev/null; then
        echo -e "${YELLOW}⚠️  $service_name is already running${NC}"
        return 1
    fi
    
    return 0
}

# Start Stripe CLI
echo -e "${BLUE}[1/3]${NC} Starting Stripe CLI webhook listener..."
if check_service "Stripe CLI" "stripe_webhook.pid" "stripe listen"; then
    stripe listen --forward-to localhost:8888/api/storage/purchase/webhook > stripe_webhook.log 2>&1 &
    STRIPE_PID=$!
    echo $STRIPE_PID > stripe_webhook.pid
    sleep 2
    if ps -p $STRIPE_PID > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Stripe CLI started (PID: $STRIPE_PID)${NC}"
    else
        echo -e "${RED}❌ Failed to start Stripe CLI${NC}"
        echo -e "${YELLOW}   Check if Stripe CLI is installed: stripe --version${NC}"
        exit 1
    fi
fi

# Start Backend
echo -e "\n${BLUE}[2/3]${NC} Starting backend server..."
if check_service "Backend" "backend.pid" "cloud-dist"; then
    go run ./cmd/cloud-dist/main.go -config configs/config.yaml > backend.log 2>&1 &
    BACKEND_PID=$!
    echo $BACKEND_PID > backend.pid
    sleep 3
    if ps -p $BACKEND_PID > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Backend started (PID: $BACKEND_PID)${NC}"
        echo -e "   ${YELLOW}→ http://localhost:8888${NC}"
    else
        echo -e "${RED}❌ Failed to start backend${NC}"
        echo -e "${YELLOW}   Check backend.log for errors${NC}"
        cleanup
        exit 1
    fi
fi

# Start Frontend
echo -e "\n${BLUE}[3/3]${NC} Starting frontend dev server..."
if check_service "Frontend" "frontend.pid" "vite"; then
    cd frontend
    npm run dev > ../frontend.log 2>&1 &
    FRONTEND_PID=$!
    echo $FRONTEND_PID > ../frontend.pid
    cd ..
    sleep 3
    if ps -p $FRONTEND_PID > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Frontend started (PID: $FRONTEND_PID)${NC}"
        echo -e "   ${YELLOW}→ http://localhost:3000${NC}"
    else
        echo -e "${RED}❌ Failed to start frontend${NC}"
        echo -e "${YELLOW}   Check frontend.log for errors${NC}"
        cleanup
        exit 1
    fi
fi

# Summary
echo -e "\n${GREEN}╔════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║   ✅ All Services Started!            ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════╝${NC}\n"

echo -e "${BLUE}📍 Services:${NC}"
echo -e "   ${GREEN}Backend:${NC}  http://localhost:8888"
echo -e "   ${GREEN}Frontend:${NC} http://localhost:3000"
echo -e "   ${GREEN}Stripe:${NC}   Webhook forwarding active\n"

echo -e "${BLUE}📋 Logs:${NC}"
echo -e "   ${YELLOW}Backend:${NC}  tail -f backend.log"
echo -e "   ${YELLOW}Frontend:${NC} tail -f frontend.log"
echo -e "   ${YELLOW}Stripe:${NC}   tail -f stripe_webhook.log\n"

echo -e "${YELLOW}💡 Press Ctrl+C to stop all services${NC}\n"

# Wait for all background processes
wait

