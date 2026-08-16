#!/bin/bash

# Run script for Web App Uji Emisi

set -e

echo "Starting Web App Uji Emisi..."

# Colors
GREEN='\033[0;32m'
NC='\033[0m'

# Function to cleanup on exit
cleanup() {
    echo ""
    echo "Stopping services..."
    kill $BACKEND_PID $FRONTEND_PID 2>/dev/null
    exit 0
}

trap cleanup SIGINT SIGTERM

# Start backend in background
echo -e "${GREEN}Starting backend server...${NC}"
cd backend
go run main.go &
BACKEND_PID=$!
cd ..

# Wait a bit for backend to start
sleep 3

# Start frontend in background
echo -e "${GREEN}Starting frontend server...${NC}"
cd frontend
npm run dev &
FRONTEND_PID=$!
cd ..

echo ""
echo "=========================================="
echo "  Services are running!"
echo "=========================================="
echo ""
echo "  Frontend: http://localhost:5173"
echo "  Backend:  http://localhost:8080"
echo ""
echo "  Press Ctrl+C to stop all services"
echo ""

# Wait for background processes
wait
