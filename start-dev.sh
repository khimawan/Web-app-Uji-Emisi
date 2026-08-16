#!/bin/bash

# Quick start script for Devin environment

set -e

echo "=========================================="
echo "  Web App Uji Emisi - Quick Start"
echo "=========================================="

# Check if setup was run
if [ ! -d "frontend/node_modules" ]; then
    echo "Running setup first..."
    ./devin-setup.sh
fi

# Check if PostgreSQL is running
if ! pg_isready -q 2>/dev/null; then
    echo "Starting PostgreSQL..."
    sudo service postgresql start 2>/dev/null || sudo systemctl start postgresql 2>/dev/null || true
    sleep 2
fi

# Check if database exists
if ! sudo -u postgres psql -d uji_emisi -c "SELECT 1" &>/dev/null; then
    echo "Creating database..."
    sudo -u postgres psql -c "CREATE DATABASE uji_emisi;" 2>/dev/null || true
    sudo -u postgres psql -d uji_emisi -f migrations/001_initial_schema.sql 2>/dev/null || true
fi

# Start the application
echo ""
echo "Starting application..."
./run.sh
