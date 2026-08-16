#!/bin/bash

# Devin Environment Setup Script for Web App Uji Emisi

set -e

echo "=========================================="
echo "  Web App Uji Emisi - Devin Setup"
echo "=========================================="

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

# Check and install dependencies
echo ""
echo "Step 1: Checking system dependencies..."

# Check Go
if command -v go &> /dev/null; then
    print_status "Go $(go version | cut -d' ' -f3) is installed"
else
    print_warning "Go not found. Installing..."
    wget -q https://go.dev/dl/go1.21.5.linux-amd64.tar.gz -O /tmp/go.tar.gz
    sudo tar -C /usr/local -xzf /tmp/go.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    export PATH=$PATH:/usr/local/go/bin
    print_status "Go installed"
fi

# Check Node.js
if command -v node &> /dev/null; then
    print_status "Node.js $(node -v) is installed"
else
    print_warning "Node.js not found. Installing..."
    curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
    sudo apt-get install -y nodejs
    print_status "Node.js installed"
fi

# Check npm
if command -v npm &> /dev/null; then
    print_status "npm $(npm -v) is installed"
else
    print_error "npm not found"
    exit 1
fi

# Check PostgreSQL
if command -v psql &> /dev/null; then
    print_status "PostgreSQL client is installed"
else
    print_warning "PostgreSQL not found. Installing..."
    sudo apt-get update
    sudo apt-get install -y postgresql postgresql-client
    print_status "PostgreSQL installed"
fi

# Step 2: Setup PostgreSQL
echo ""
echo "Step 2: Setting up PostgreSQL..."

# Start PostgreSQL
sudo service postgresql start 2>/dev/null || sudo systemctl start postgresql 2>/dev/null || true
print_status "PostgreSQL service started"

# Create database and user
echo "Creating database..."
sudo -u postgres psql -c "ALTER USER postgres WITH PASSWORD 'password';" 2>/dev/null || true
sudo -u postgres psql -c "CREATE DATABASE uji_emisi;" 2>/dev/null || print_warning "Database might already exist"
print_status "Database configured"

# Step 3: Setup Backend
echo ""
echo "Step 3: Setting up Go backend..."

cd backend
go mod tidy
print_status "Go dependencies installed"
cd ..

# Step 4: Setup Frontend
echo ""
echo "Step 4: Setting up React frontend..."

cd frontend
npm install
print_status "Node.js dependencies installed"
cd ..

# Step 5: Create uploads directory
echo ""
echo "Step 5: Creating required directories..."

mkdir -p backend/uploads/home
mkdir -p backend/uploads/temp
print_status "Upload directories created"

# Step 6: Make scripts executable
chmod +x setup.sh
chmod +x run.sh
chmod +x start-dev.sh
print_status "Scripts made executable"

echo ""
echo "=========================================="
echo "  Setup Complete!"
echo "=========================================="
echo ""
echo "To start the application:"
echo ""
echo "  Option 1 (Quick start):"
echo "    ./start-dev.sh"
echo ""
echo "  Option 2 (Manual):"
echo "    Terminal 1: cd backend && go run main.go"
echo "    Terminal 2: cd frontend && npm run dev"
echo ""
echo "Access the app at: http://localhost:5173"
echo ""
echo "Default accounts:"
echo "  Admin:      adminmas / adminmas"
echo "  Supervisor: supervisoraja / supervisoraja"
echo "  Anggota:    kroco01 / kroco01"
echo ""
