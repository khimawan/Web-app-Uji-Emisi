.PHONY: setup run build clean db-setup db-reset help

# Default target
help:
	@echo ""
	@echo "Web App Uji Emisi - Available Commands:"
	@echo ""
	@echo "  make setup        - Setup development environment"
	@echo "  make run          - Run both backend and frontend"
	@echo "  make backend      - Run backend only"
	@echo "  make frontend     - Run frontend only"
	@echo "  make build        - Build frontend for production"
	@echo "  make db-setup     - Setup database"
	@echo "  make db-reset     - Reset database"
	@echo "  make clean        - Clean temporary files"
	@echo ""

# Setup environment
setup:
	@echo "Setting up development environment..."
	./devin-setup.sh

# Run all services
run:
	@echo "Starting all services..."
	./run.sh

# Run backend only
backend:
	@echo "Starting backend server..."
	cd backend && go run main.go

# Run frontend only
frontend:
	@echo "Starting frontend server..."
	cd frontend && npm run dev

# Build frontend
build:
	@echo "Building frontend..."
	cd frontend && npm run build

# Setup database
db-setup:
	@echo "Setting up database..."
	sudo -u postgres psql -c "CREATE DATABASE uji_emisi;" 2>/dev/null || echo "Database might exist"
	sudo -u postgres psql -c "ALTER USER postgres WITH PASSWORD 'password';" 2>/dev/null || true
	sudo -u postgres psql -d uji_emisi -f migrations/001_initial_schema.sql 2>/dev/null || echo "Tables might exist"

# Reset database
db-reset:
	@echo "Resetting database..."
	sudo -u postgres psql -c "DROP DATABASE IF EXISTS uji_emisi;"
	sudo -u postgres psql -c "CREATE DATABASE uji_emisi;"
	sudo -u postgres psql -d uji_emisi -f migrations/001_initial_schema.sql

# Clean temporary files
clean:
	@echo "Cleaning temporary files..."
	rm -rf frontend/dist
	rm -rf backend/tmp
	rm -rf backend/*.exe
	rm -rf frontend/.vite
	find . -type d -name "node_modules" -exec rm -rf {} + 2>/dev/null || true
	find . -type d -name "vendor" -exec rm -rf {} + 2>/dev/null || true
	@echo "Clean complete!"
