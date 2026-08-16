# Running on Devin

This guide explains how to run the Web App Uji Emisi on Devin environment.

## Quick Start

### Option 1: Automated Setup

```bash
# Run the setup script
./devin-setup.sh

# Start the application
./start-dev.sh
```

### Option 2: Manual Setup

```bash
# 1. Setup PostgreSQL
sudo service postgresql start
sudo -u postgres psql -c "CREATE DATABASE uji_emisi;"
sudo -u postgres psql -c "ALTER USER postgres WITH PASSWORD 'password';"

# 2. Setup Backend
cd backend
go mod tidy
cd ..

# 3. Setup Frontend
cd frontend
npm install
cd ..

# 4. Run the application
./run.sh
```

## Environment Variables

### Backend (.env)

Create `backend/.env` file:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=uji_emisi
JWT_SECRET=devin-secret-key-2024
PORT=8080
```

### Frontend

No additional configuration needed. The frontend uses Vite proxy to connect to the backend.

## Application URLs

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080/api

## Default Accounts

| Role | Username | Password |
|------|----------|----------|
| Admin | adminmas | adminmas |
| Supervisor | supervisoraja | supervisoraja |
| Anggota | kroco01 | kroco01 |

---

# Tutorial Auto Open di Devin

## Langkah 1: Buka Devin

Buka https://app.devin.ai dan login dengan akun Anda.

## Langkah 2: Buat New Session

Klik **"New Session"** atau **"Start Session"**

## Langkah 3: Clone Repository

Ketik di chat Devin:

```
Clone repository git@github.com:khimawan/Web-app-Uji-Emisi.git
```

Atau:

```
Clone https://github.com/khimawan/Web-app-Uji-Emisi.git
```

## Langkah 4: Setup Environment

Setelah cloned, ketik:

```
Run the setup script: ./devin-setup.sh
```

Atau manual:

```
Setup this project:
1. Start PostgreSQL
2. Create database uji_emisi
3. Run backend: cd backend && go run main.go
4. Run frontend: cd frontend && npm install && npm run dev
```

## Langkah 5: Jalankan Aplikasi

```
Start the application using: ./start-dev.sh
```

Atau:

```
Run both backend and frontend servers
```

## Contoh Prompt Lengkap

```
I want to run this web application locally.

1. Clone the repository from git@github.com:khimawan/Web-app-Uji-Emisi.git
2. Setup PostgreSQL database named "uji_emisi"
3. Start the Go backend on port 8080
4. Start the React frontend on port 5173
5. Verify both services are running

Use the setup scripts provided in the repository.
```

## Shortcut: Direct URL

Anda bisa langsung buka URL ini di Devin:

```
https://devin.ai/new?repo=git@github.com:khimawan/Web-app-Uji-Emisi.git
```

## Tips

| Kebutuhan | Prompt |
|-----------|--------|
| Setup saja | `Run ./devin-setup.sh` |
| Jalankan app | `Run ./start-dev.sh` |
| Hanya backend | `Start Go backend: cd backend && go run main.go` |
| Hanya frontend | `Start React: cd frontend && npm run dev` |
| Reset database | `Reset database: make db-reset` |
| Build production | `Build frontend: cd frontend && npm run build` |

## Credentials Default

```
Admin:      adminmas / adminmas
Supervisor: supervisoraja / supervisoraja
Anggota:    kroco01 / kroco01
```

## URLs

- Frontend: http://localhost:5173
- Backend API: http://localhost:8080/api

---

## Troubleshooting

### PostgreSQL Issues

```bash
# Check PostgreSQL status
sudo service postgresql status

# Restart PostgreSQL
sudo service postgresql restart

# Reset database
sudo -u postgres psql -c "DROP DATABASE IF EXISTS uji_emisi;"
sudo -u postgres psql -c "CREATE DATABASE uji_emisi;"
```

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or change port in backend/.env
PORT=8081
```

### Go Module Issues

```bash
cd backend
go clean -modcache
go mod tidy
```

### Node.js Issues

```bash
cd frontend
rm -rf node_modules
npm install
```

## Project Structure

```
web-app-uji-emisi/
├── backend/           # Go backend
├── frontend/          # React frontend
├── docs/              # Documentation
├── migrations/        # SQL migrations
├── setup.sh           # Basic setup
├── devin-setup.sh     # Devin environment setup
├── run.sh             # Run both services
└── start-dev.sh       # Quick start with checks
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/auth/login | Login |
| GET | /api/kendaraan | List kendaraan |
| POST | /api/kendaraan | Create kendaraan |
| GET | /api/hasil-uji | List hasil uji |
| POST | /api/hasil-uji | Create hasil uji |
| GET | /api/parameters | List parameters |
| GET | /api/users | List users |

For full API documentation, see [docs/API.md](docs/API.md)
