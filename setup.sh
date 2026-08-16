#!/bin/bash

# Setup script for Web App Uji Emisi

set -e

echo "=== Web App Uji Emisi - Setup ==="

# Check if PostgreSQL is running
echo "Checking PostgreSQL..."
if ! command -v psql &> /dev/null; then
    echo "PostgreSQL not found. Installing..."
    sudo apt-get update
    sudo apt-get install -y postgresql postgresql-client
fi

# Start PostgreSQL service
echo "Starting PostgreSQL..."
sudo service postgresql start 2>/dev/null || sudo systemctl start postgresql 2>/dev/null || true

# Create database
echo "Creating database..."
sudo -u postgres psql -c "CREATE DATABASE uji_emisi;" 2>/dev/null || echo "Database might already exist"

# Create user with password
sudo -u postgres psql -c "ALTER USER postgres WITH PASSWORD 'password';" 2>/dev/null || true

# Install Go dependencies
echo "Installing Go dependencies..."
cd backend
go mod tidy
cd ..

# Install Node.js dependencies
echo "Installing Node.js dependencies..."
cd frontend
npm install
cd ..

echo "=== Setup Complete ==="
echo ""
echo "To run the application:"
echo "  Terminal 1 (Backend):  cd backend && go run main.go"
echo "  Terminal 2 (Frontend): cd frontend && npm run dev"
echo ""
echo "Access the app at: http://localhost:5173"
