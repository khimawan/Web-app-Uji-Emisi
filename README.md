# Aplikasi Uji Emisi

Aplikasi web untuk mengelola dan mencatat hasil uji emisi kendaraan. Aplikasi ini memungkinkan pengguna untuk mencatat data kendaraan, melakukan uji emisi, dan melihat hasil uji dalam bentuk tabel serta mengekspor data ke format CSV, Excel, dan PDF.

## Daftar Isi

- [Ringkasan](#ringkasan)
- [Fitur](#fitur)
- [Teknologi yang Digunakan](#teknologi-yang-digunakan)
- [Persyaratan Sistem](#persyaratan-sistem)
- [Instalasi](#instalasi)
- [Menjalankan Aplikasi](#menjalankan-aplikasi)
- [Struktur Aplikasi](#struktur-aplikasi)
- [Akun Default](#akun-default)
- [Panduan Penggunaan](#panduan-penggunaan)
- [API Reference](#api-reference)
- [Dokumentasi](#dokumentasi)
- [Troubleshooting](#troubleshooting)
- [Lisensi](#lisensi)

## Ringkasan

Aplikasi Uji Emisi adalah solusi web-based untuk mengelola dan mencatat proses uji emisi kendaraan bermotor. Aplikasi ini dibangun dengan backend Golang dan frontend React/TypeScript untuk mempermudah petugas uji emisi dalam mencatat, mengevaluasi, dan menghasilkan laporan hasil uji emisi secara efisien. Sistem secara otomatis memvalidasi hasil uji berdasarkan standar yang telah ditetapkan dan mengategorikan kendaraan sebagai "Lulus" atau "Tidak Lulus".

## Fitur

### Input Data Kendaraan
- Formulir input kendaraan dengan validasi real-time
- Klasifikasi kendaraan (Dinas/Umum)
- Pencatatan detail kendaraan (plat nomor, merek, tipe, tahun)
- Pemilihan jenis bahan bakar (Bensin/Solar)
- Input batch kendaraan via CSV

### Uji Emisi
- Interface input hasil uji untuk kendaraan terdaftar
- Parameter berbeda untuk kendaraan bensin dan solar
- Input parameter uji (CO, CO2, HC, O2, Lambda untuk bensin)
- Input Opasitas untuk kendaraan solar
- Pop-up notifikasi hasil emisi

### Analisis dan Laporan
- Dashboard dengan statistik hasil uji
- Pie chart perbandingan data
- Diagram batang jumlah kendaraan
- Export ke format CSV, Excel, dan PDF
- Filter dan pencarian data multi-parameter

### Manajemen User
- Sistem autentikasi dengan level akses berbeda (Admin, Supervisor, Anggota)
- Manajemen akun user

### Konfigurasi Sistem
- Pengaturan standar batas emisi per kategori
- Konfigurasi pop-up notifikasi

## Teknologi yang Digunakan

### Backend
- **Go 1.21+** - Bahasa pemrograman
- **Gin** - HTTP web framework
- **GORM** - ORM untuk Go
- **JWT** - JSON Web Token untuk autentikasi

### Database
- **PostgreSQL** - Database relasional

### Frontend
- **React 18** - Library untuk UI
- **TypeScript** - Type-safe JavaScript
- **Vite** - Build tool
- **Tailwind CSS** - CSS framework
- **Recharts** - Library chart untuk React
- **Axios** - HTTP client

## Persyaratan Sistem

### Software
- Go 1.21 atau lebih baru
- Node.js 18+ dan npm/yarn
- PostgreSQL 14+
- Git

### Hardware
- Minimum: 2GB RAM, 2GHz CPU, 1GB disk space
- Rekomendasi: 4GB+ RAM, 2.5GHz+ CPU, 2GB+ disk space

## Instalasi

### 1. Clone Repository

```bash
git clone https://github.com/username/web-app-uji-emisi.git
cd web-app-uji-emisi
```

### 2. Setup Backend

```bash
cd backend
go mod download
```

### 3. Setup Database

```bash
# Buat database PostgreSQL
psql -U postgres -c "CREATE DATABASE uji_emisi;"
```

### 4. Konfigurasi Environment

Buat file `.env` di folder backend:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=uji_emisi
JWT_SECRET=your-secret-key
PORT=8080
```

### 5. Setup Frontend

```bash
cd ../frontend
npm install
```

## Menjalankan Aplikasi

### Development Mode

1. Jalankan Backend:
```bash
cd backend
go run main.go
```

2. Jalankan Frontend (terminal terpisah):
```bash
cd frontend
npm run dev
```

3. Buka browser dan akses:
```
http://localhost:5173
```

### Production Mode

1. Build Frontend:
```bash
cd frontend
npm run build
```

2. Jalankan Backend (akan serve frontend juga):
```bash
cd backend
go run main.go
```

3. Akses aplikasi:
```
http://localhost:8080
```

## Struktur Aplikasi

```
web-app-uji-emisi/
├── backend/                    # Backend Go
│   ├── config/                 # Konfigurasi aplikasi
│   ├── controllers/            # Handler HTTP
│   ├── middleware/              # Middleware (auth, cors)
│   ├── models/                 # Model database
│   ├── routes/                 # Route definitions
│   ├── services/               # Business logic
│   ├── utils/                  # Utility functions
│   ├── main.go                 # Entry point
│   ├── go.mod                  # Go module
│   └── .env                    # Environment variables
├── frontend/                   # Frontend React/TypeScript
│   ├── src/
│   │   ├── components/         # Reusable components
│   │   ├── pages/              # Page components
│   │   ├── services/           # API services
│   │   ├── types/              # TypeScript types
│   │   ├── hooks/              # Custom hooks
│   │   ├── context/            # React context
│   │   └── utils/              # Utility functions
│   ├── public/                 # Static files
│   ├── index.html              # HTML template
│   ├── package.json            # Dependencies
│   ├── tsconfig.json           # TypeScript config
│   ├── vite.config.ts          # Vite config
│   └── tailwind.config.js      # Tailwind config
├── docs/                       # Dokumentasi
│   ├── PRD.md                  # Product Requirements Document
│   ├── ERD.md                  # Entity Relationship Diagram
│   └── API.md                  # API Documentation
└── README.md                   # Dokumentasi ini
```

## Akun Default

### Admin
- **Nama:** Superuser Admin
- **Username:** adminmas
- **Password:** adminmas

### Supervisor
- **Nama:** Kang Supervisor
- **Username:** supervisoraja
- **Password:** supervisoraja

### Anggota
- **Nama:** Kroco01
- **Username:** kroco01
- **Password:** kroco01

## Panduan Penggunaan

### Login
1. Buka halaman login
2. Masukkan username dan password
3. Klik tombol "Login"

### Input Data Kendaraan (Page 1)
1. Login sebagai Admin, Supervisor, atau Anggota
2. Klik menu "Input Data Kendaraan"
3. Isi form dengan data kendaraan
4. Klik "Simpan" atau upload CSV untuk input batch
5. Data akan muncul di tabel history

### Input Data Emisi (Page 2)
1. Login sebagai Admin, Supervisor, atau Anggota
2. Klik menu "Input Data Emisi"
3. Pilih kendaraan dari daftar
4. Isi parameter emisi sesuai jenis kendaraan
5. Klik "Simpan Hasil"
6. Pop-up notifikasi akan muncul dengan hasil uji

### Data Hasil Uji (Page 3)
1. Login sebagai Admin, Supervisor, atau Anggota
2. Klik menu "Data Hasil Uji"
3. Lihat statistik dan chart di bagian atas
4. Lihat tabel data kendaraan di bagian bawah
5. Gunakan filter dan pencarian jika diperlukan
6. Download data dalam format CSV, Excel, atau PDF

### Input Nilai Parameter (Page 4)
1. Login sebagai Admin atau Supervisor
2. Klik menu "Parameter"
3. Edit parameter emisi sesuai kebutuhan
4. Simpan perubahan

### Manajemen User (Page 5)
1. Login sebagai Admin atau Supervisor
2. Klik menu "Users"
3. Tambah, edit, atau hapus akun user
4. Catatan: Akun admin tidak dapat dihapus

## API Reference

### Authentication
- `POST /api/auth/login` - Login
- `POST /api/auth/logout` - Logout
- `GET /api/auth/me` - Get current user

### Kendaraan
- `GET /api/kendaraan` - List kendaraan
- `POST /api/kendaraan` - Tambah kendaraan
- `GET /api/kendaraan/:id` - Detail kendaraan
- `PUT /api/kendaraan/:id` - Update kendaraan
- `DELETE /api/kendaraan/:id` - Hapus kendaraan
- `POST /api/kendaraan/upload` - Upload CSV

### Hasil Uji
- `GET /api/hasil-uji` - List hasil uji
- `POST /api/hasil-uji` - Tambah hasil uji
- `PUT /api/hasil-uji/:id` - Update hasil uji
- `DELETE /api/hasil-uji/:id` - Hapus hasil uji
- `GET /api/hasil-uji/statistics` - Statistik hasil uji

### Parameter
- `GET /api/parameters` - List parameter
- `POST /api/parameters` - Tambah parameter
- `PUT /api/parameters/:id` - Update parameter
- `DELETE /api/parameters/:id` - Hapus parameter

### Pop-up Notes
- `GET /api/popup-notes` - List pop-up notes
- `POST /api/popup-notes` - Tambah pop-up note
- `PUT /api/popup-notes/:id` - Update pop-up note
- `DELETE /api/popup-notes/:id` - Hapus pop-up note

### Home Page
- `GET /api/home/description` - Get deskripsi
- `PUT /api/home/description` - Update deskripsi
- `POST /api/home/image` - Upload gambar
- `DELETE /api/home/image/:id` - Hapus gambar
- `GET /api/home/working-instruction` - Get working instruction
- `POST /api/home/working-instruction` - Upload working instruction

### Users (Admin)
- `GET /api/users` - List users
- `POST /api/users` - Tambah user
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Hapus user

## Dokumentasi

Dokumentasi lengkap dapat ditemukan di folder `docs/`:

- [PRD (Product Requirements Document)](docs/PRD.md)
- [ERD (Entity Relationship Diagram)](docs/ERD.md)
- [API Documentation](docs/API.md)

## Troubleshooting

1. **Database Connection Error**
   - Pastikan PostgreSQL berjalan
   - Periksa konfigurasi di file `.env`
   - Pastikan database sudah dibuat

2. **Port Sudah Digunakan**
   - Ubah port di file `.env` atau `vite.config.ts`
   - Atau hentikan process yang menggunakan port

3. **Build Error**
   - Jalankan `go mod tidy` di backend
   - Jalankan `npm install` di frontend
   - Pastikan versi Go dan Node.js sesuai

4. **Login Gagal**
   - Pastikan akun sudah dibuat (jalankan seeding)
   - Periksa username dan password

## Lisensi

MIT License

## Kontak

Untuk pertanyaan, dukungan, atau kontribusi, silakan hubungi:
- GitHub Issues: https://github.com/username/web-app-uji-emisi/issues
