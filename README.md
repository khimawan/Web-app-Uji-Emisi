# Aplikasi Uji Emisi

Aplikasi web untuk mengelola dan mencatat hasil uji emisi kendaraan. Aplikasi ini memungkinkan pengguna untuk mencatat data kendaraan, melakukan uji emisi, mengimpor data massal dari CSV/Excel, dan melihat hasil uji dalam bentuk tabel serta mengekspor data ke format CSV, Excel, dan PDF.

## Daftar Isi

- [Ringkasan](#ringkasan)
- [Fitur](#fitur)
- [Teknologi yang Digunakan](#teknologi-yang-digunakan)
- [Persyaratan Sistem](#persyaratan-sistem)
- [Instalasi](#instalasi)
- [Deployment dengan Docker](#deployment-dengan-docker)
- [Catatan Penting untuk Deployment](#catatan-penting-untuk-deployment)
- [Menjalankan Aplikasi](#menjalankan-aplikasi)
- [Struktur Aplikasi](#struktur-aplikasi)
- [Akun Default](#akun-default)
- [Panduan Penggunaan](#panduan-penggunaan)
- [Pengujian](#pengujian)
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
- Pemilihan jenis bahan bakar (Bensin M, Bensin N&O, Solar JBB, Solar GVW)
- Import massal dari **CSV dan Excel** memakai template siap-isi
- Template Excel dilengkapi **dropdown** kategori/jenis dan daftar plat nomor terdaftar
- Validasi per baris: baris salah ditolak, baris lain tetap tersimpan

### Uji Emisi
- Interface input hasil uji untuk kendaraan terdaftar
- Parameter berbeda untuk kendaraan bensin dan solar
- Input parameter uji (CO, CO2, HC, O2, Lambda untuk bensin)
- Input Opasitas untuk kendaraan solar
- Pop-up notifikasi hasil emisi
- **Import hasil uji massal dari CSV/Excel**; status Lulus/Tidak Lulus dan catatan
  tetap dihitung sistem sehingga identik dengan input manual

### Analisis dan Laporan
- Dashboard dengan statistik hasil uji
- Pie chart perbandingan data
- Diagram batang jumlah kendaraan
- Export ke format CSV, Excel, dan PDF
- Filter dan pencarian data multi-parameter
- Bulk delete hasil uji

### Manajemen User
- Sistem autentikasi JWT dengan level akses berbeda (Admin, Supervisor, Anggota)
- Username admin **disembunyikan** di halaman manajemen user
- Kredensial akun admin hanya dapat diubah oleh admin itu sendiri
- Akun admin tidak dapat dihapus

### Konfigurasi Sistem
- Pengaturan standar batas emisi per kategori
- Konfigurasi pop-up notifikasi
- Upload **Working Instruction** (PDF) di halaman Home, lengkap dengan preview
  dan tombol unduh — khusus Admin & Supervisor

###;Halaman Beranda
- Deskripsi aplikasi, galeri gambar, dan working instruction PDF
- Semua pengguna dapat melihat & mengunduh; hanya Admin/Supervisor yang mengelola

## Teknologi yang Digunakan

### Backend
- **Go 1.25+** - Bahasa pemrograman (`go.mod` mensyaratkan `go 1.25.0`)
- **Gin** - HTTP web framework
- **GORM** - ORM untuk Go
- **JWT** - JSON Web Token untuk autentikasi
- **Excelize** - Baca/tulis file `.xlsx` untuk template & import data

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
- Go 1.25 atau lebih baru (**wajib**, bukan 1.21 — lihat [Catatan Penting](#catatan-penting-untuk-deployment) butir 1)
- Node.js 18+ dan npm/yarn
- PostgreSQL 14+ (atau Docker, untuk deploy via compose)
- Git

### Hardware
- Minimum: 2GB RAM, 2GHz CPU, 1GB disk space
- Rekomendasi: 4GB+ RAM, 2.5GHz+ CPU, 2GB+ disk space

## Instalasi

### 1. Clone Repository

```bash
git clone https://github.com/khimawan/Web-app-Uji-Emisi.git
cd Web-app-Uji-Emisi
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
npm ci        # gunakan `npm install` hanya bila package.json berubah
```

## Deployment dengan Docker

Deploy produksi memakai `docker compose`. Semua perintah dijalankan dari root repo.

### Prasyarat

- Docker Engine + Docker Compose v2
- Port bebas: `5173` (frontend), `8080` (backend), `5432` (postgres)
- **Ruang disk ≥ 4 GB** — build image Go saja butuh ±1 GB saat proses build

### Langkah

```bash
# 1. Dependency frontend WAJIB di-install di host, bukan di dalam container
cd frontend && npm ci && cd ..

# 2. Folder uploads WAJIB ada sebelum backend start
mkdir -p backend/uploads

# 3. Build dan jalankan
docker compose up -d --build
```

| Service | Alamat |
|---|---|
| Frontend | http://localhost:5173 |
| Backend API | http://localhost:8080 |
| PostgreSQL | localhost:5432 |

### Kenapa `npm ci` harus di host, bukan di container

Compose mem-bind-mount `./frontend:/app` ke dalam container, sehingga **isi host
menimpa isi image**. Kalau `node_modules` tidak ada di host, Vite tidak punya
perintah `vite` untuk dijalankan dan container frontend gagal start.

Efek sampingnya: setiap `npm install` di host langsung terbaca container tanpa
restart, dan `RUN npm install` di `frontend/Dockerfile` selalu sia-sia karena
langsung ditimpa mount. Karena itu langkah tersebut sudah dihapus dari Dockerfile.

### Kenapa `backend/uploads` harus ada sebelum start

`backend/routes/routes.go` memanggil `r.Static("/uploads", "./uploads")`.
Gin **panic saat start** kalau folder itu tidak ada — backend langsung crash,
bukan sekadar membalas 404.

`Dockerfile` sudah menjalankan `mkdir -p`, tetapi volume
`./backend/uploads:/app/uploads` menggantinya dengan folder host yang kosong.

---

## Catatan Penting untuk Deployment

Bagian ini berasal dari kegagalan nyata saat deploy. Melewatkannya bisa membuat
build gagal atau aplikasi crash.

### 1. Versi Go pada image harus cocok dengan `go.mod`

`backend/go.mod` menyatakan `go 1.25.0`, jadi image build **wajib**
`golang:1.25-alpine` atau lebih baru.

| Docker image | Hasil |
|---|---|
| `golang:1.21-alpine` | Build gagal, atau diam-diam mengunduh toolchain 1.25 — boros ruang dan bisa gagal karena disk penuh |
| `golang:1.25-alpine` | Benar |

### 2. Image backend multi-stage: isinya binary, bukan toolchain

Image runtime hanya berisi `alpine` + `./main` (±53 MB). **Tidak ada Go di
dalamnya.** Akibatnya:

- `command: go run main.go` di compose → crash, `go: executable file not found`
- Volume `./backend:/app` → menimpa `/app/main` dengan source, binary hilang

Konsekuensi: **tidak ada hot reload untuk Go.** Setiap kali ada perubahan
`.go`, wajib rebuild:

```bash
docker compose build backend && docker compose up -d backend
```

Frontend tetap punya HMR normal.

### 3. Proxy Vite wajib memakai nama service, bukan `localhost`

`frontend/vite.config.ts` memakai `http://backend:8080` sebagai default.

Di dalam container frontend, `localhost` berarti **container frontend itu
sendiri**, bukan container backend. Kalau dikembalikan menjadi `localhost`,
semua request `/api` dan `/uploads` gagal.

| Gejala | Penyebab |
|---|---|
| `/api` membalas 500, log Vite `ECONNREFUSED ::1:8080` | Proxy masih menunjuk `localhost` |
| `/uploads` mengembalikan HTML | Path `/uploads` tidak diproxy |

Untuk menjalankan Vite **di luar** Docker, set target secara eksplisit:

```bash
VITE_API_TARGET=http://localhost:8080 npm run dev
```

### 4. `allowedHosts` wajib diisi begitu diakses lewat domain atau proxy

Vite ≥ 5.4 menolak request dengan header `Host` yang tidak terdaftar dan
membalas **403**. Di belakang reverse proxy atau Cloudflare tunnel, header
`Host` berisi domain publik, bukan `localhost`. Setiap domain baru harus
ditambahkan ke `server.allowedHosts` di `frontend/vite.config.ts`.

### 5. Working instruction PDF tidak boleh pakai `<iframe>`

Browser mobile tidak merender PDF di dalam `<iframe>` — Chrome Android memicu
unduhan, iOS Safari sering menampilkan halaman kosong. Karena itu PDF
dirender dengan **PDF.js** ke `<canvas>` lewat `src/components/PdfViewer.tsx`.
Jangan dikembalikan ke `<iframe>`.

### 6. Build di mesin kecil

Build Go butuh ±500 MB ruang sementara dan ±1 GB RAM. Untuk mesin kecil,
`backend/Dockerfile` sudah mengarahkan `GOCACHE` dan `$WORK` ke tmpfs serta
menggunakan `go build -p=2`. Kalau masih muncul `no space left on device`:

```bash
docker builder prune -af
```

BuildKit menyimpan salinan layer di build cache **dan** di image store, jadi
kebutuhan disk sementara bisa dua kali lipat dari ukuran image.

### 7. Inisialisasi database hanya jalan sekali

`migrations/001_initial_schema.sql` di-mount ke
`/docker-entrypoint-initdb.d/`, yang **hanya dieksekusi saat volume Postgres
masih kosong**. Setelah volume ada, mengubah file SQL tidak mengubah database.

```bash
docker compose down -v   # HAPUS SELURUH DATA — gunakan dengan sangat hati-hati
```

### 8. Tabel kembar dari AutoMigrate dan SQL init

`controllers.SeedInitialData()` menjalankan GORM AutoMigrate yang membuat nama
tabel bentuk jamak, sementara `migrations/001_initial_schema.sql` membuat
bentuk tunggal. Hasilnya ada pasangan seperti `home_content` dan
`home_contents`. Pastikan query, model, dan migrasi konsisten sebelum menambah
tabel baru.

### 9. Cloudflare tunnel: file config harus ada di host

`cloudflared/config.yml` di-bind-mount read-only. Kalau file ini terhapus dari
host, container yang sedang jalan masih aman karena masih memegang inode lama,
tetapi **restart berikutnya tidak akan start** — mount tidak menemukan
sumbernya.

`docker-compose.yml` juga me-mount `/root/.cloudflared/tunnel.token`. Pastikan
file itu ada sebelum `docker compose up`.

### 10. Setelah `git clone`, selalu jalankan `npm ci`

`node_modules` tidak ikut ter-commit. `frontend/package.json` sekarang memakai
`pdfjs-dist` untuk render PDF.

### 11. Kredensial default masih dipakai

`POSTGRES_PASSWORD` dan `JWT_SECRET` di `docker-compose.yml` masih memakai nilai
bawaan repository. **Ganti keduanya** sebelum dipakai di lingkungan
production yang nyata.

### 12. Pembatas beban sudah aktif secara bawaan

Backend sudah memasang dua lapis perlindungan, keduanya bisa diubah lewat
environment variable di `docker-compose.yml` **tanpa rebuild image**:

| Variabel | Bawaan | Fungsi |
|---|---|---|
| `RATE_LIMIT_PER_MINUTE` | 120 | Request/menit per IP untuk seluruh endpoint |
| `RATE_LIMIT_BURST` | 40 | Token yang boleh ditumpuk sekaligus |
| `LOGIN_RATE_LIMIT_PER_MINUTE` | 10 | Request/menit per IP untuk `/api/auth/login` saja |
| `LOGIN_RATE_LIMIT_BURST` | 5 | Burst untuk login |
| `MAX_CONCURRENT_REQUESTS` | 64 | Request yang boleh diproses bersamaan; sisanya dapat **503** |

Melewati batas menghasilkan **429** (rate limit) atau **503** (beban penuh),
keduanya membawa header `Retry-After`.

State limiter disimpan di memori proses. Kalau someday backend dijalankan lebih
dari satu instance, batas ini berlaku **per instance** dan harus dipindahkan ke
Redis atau reverse proxy.

### 13. `TRUSTED_PROXIES` jangan asal diperluas

Rate limit memakai `c.ClientIP()`, yang di-backup `X-Forwarded-For` hanya dari
proxy yang dipercaya. Nilai bawaan hanya mengizinkan Docker bridge
(`172.16.0.0/12`) dan localhost.

Kalau nilai ini diubah jadi terlalu luas, klien bisa mengirim
`X-Forwarded-For` palsu dan memakai IP baru setiap request sehingga rate limit
tidak pernah menyentuh batasnya.

### 14. Timeout server dan graceful shutdown

`SERVER_READ_TIMEOUT` (15 dtk) dan `SERVER_IDLE_TIMEOUT` (60 dtk) menutup celah
slow-loris — satu koneksi yang dibuat menggantung akan menahan goroutine dan
sebagian kecil RAM selamanya.

`SERVER_WRITE_TIMEOUT` sengaja longgar (120 dtk) karena export dan import CSV
bisa memproses ribuan baris. Kalau nilai ini dikecilkan, export besar akan
terpotong di tengah.

Saat menerima `SIGTERM`, backend menyelesaikan request yang sedang berjalan
lalu keluar rapi, sehingga `docker compose restart` tidak memutus pengguna.

### 15. Pool koneksi database wajib dibatasi

`DB_MAX_OPEN_CONNS` (20) dan `DB_MAX_IDLE_CONNS` (5) membatasi koneksi ke
PostgreSQL. Tanpa batas, setiap request yang bersamaan membuka koneksi baru
sampai kuota `max_connections` habis dan database mulai menolak koneksi.

### 16. Tuning PostgreSQL dan rotasi log

`docker-compose.yml` menyetel `shared_buffers=64MB`, `work_mem=2MB`,
`max_connections=50` — Anglia untuk mesin kecil.

Setiap container punya rotasi log `10 MB x 3`. Tanpa ini, log Gin yang mencatat
setiap request bisa memenuhi disk eMMC ketika server sedang dibebani.

### 17. Kalau server tetap berat, urutannya begini

1. Turunkan `MAX_CONCURRENT_REQUESTS` (misal 64 → 24). Ini tombol paling
   langsung dan tidak mengubah perilaku lain.
2. Naikkan `LOGIN_RATE_LIMIT_*` hanya kalau user legit banyak yang gagal login.
3. Jangan menaikkan `DB_MAX_OPEN_CONNS` melebihi kapasitas RAM — tiap
   koneksi PostgreSQL punya biaya memori.
4. Naikkan batas baru hanya setelah mengukur, bukan menebak.

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

> Saat development **di luar Docker**, `vite.config.ts` memproxy dua path ke backend (`localhost:8080`): `/api` untuk REST API dan `/uploads` untuk file statis (gambar home, working instruction PDF). Tanpa proxy `/uploads`, Vite akan mengembalikan `index.html` sehingga file tidak tampil. Path file dari API selalu berupa path absolut (`/uploads/home/...`).
>
> **Saat memakai Docker, target proxy justru harus `http://backend:8080`** — bukan `localhost`. Lihat [Catatan Penting](#catatan-penting-untuk-deployment) butir 3.

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
Web-app-Uji-Emisi/
├── backend/                    # Backend Go
│   ├── config/                 # Konfigurasi aplikasi & koneksi database
│   ├── controllers/            # Handler HTTP
│   │   ├── kendaraan.go        # CRUD kendaraan
│   │   ├── kendaraan_csv.go    # Import CSV vehicles + generator template CSV
│   │   ├── kendaraan_excel.go  # Generator template Excel kendaraan
│   │   ├── hasil_uji.go        # CRUD hasil uji, statistik, export CSV
│   │   ├── hasil_uji_import.go # Import hasil uji CSV/XLSX + template
│   │   ├── home.go             # Konten beranda, gambar, working instruction
│   │   └── user.go             # Manajemen user
│   ├── middleware/             # Middleware (auth, role, cors)
│   ├── models/                 # Model database
│   ├── routes/                 # Route definitions
│   ├── services/               # Business logic
│   ├── utils/                  # Utility functions
│   ├── uploads/                # File terunggah (diabaikan oleh git)
│   │                            # WAJIB ada sebelum container backend start
│   ├── main.go                 # Entry point
│   ├── go.mod                  # Go module (wajib Go 1.25.0)
│   ├── Dockerfile              # Multi-stage: build binary, runtime alpine
│   └── .env                    # Environment variables
├── frontend/                   # Frontend React/TypeScript
│   ├── src/
│   │   ├── components/
│   │   │   ├── Layout.tsx            # Header, navigasi desktop, sidebar mobile
│   │   │   ├── PdfViewer.tsx         # Working instruction PDF via PDF.js (canvas)
│   │   │   └── TemplateDownloadMenu.tsx  # Dropdown Template CSV/Excel
│   │   ├── pages/              # Page components
│   │   ├── services/           # API services (api.ts)
│   │   ├── types/              # TypeScript types
│   │   ├── hooks/              # Custom hooks
│   │   ├── context/            # React context (auth)
│   │   ├── utils/              # Utility functions (fileUrl)
│   │   └── vite-env.d.ts       # Tipe untuk import ?url (vite/client)
│   ├── public/                 # Static files
│   ├── index.html              # HTML template
│   ├── package.json            # Dependencies
│   ├── package-lock.json       # Lockfile (jangan dihapus)
│   ├── tsconfig.json           # TypeScript config
│   ├── vite.config.ts          # Proxy /api dan /uploads, allowedHosts
│   └── tailwind.config.js     # Tailwind config
├── cloudflared/                # Ingress config untuk Cloudflare tunnel
│   └── config.yml              # hostname -> service (read-only mount)
├── docker-compose.yml          # Orkesrasi seluruh service
├── docs/                       # Dokumentasi
│   ├── PRD.md                  # Product Requirements Document
│   ├── ERD.md                  # Entity Relationship Diagram
│   └── API.md                  # API Documentation
└── README.md                   # Dokumentasi ini
```

## Akun Default

Akun berikut dibuat otomatis saat backend pertama kali dijalankan (seeding).
Informasi ini **tidak ditampilkan** di halaman login.

| Role | Nama | Username | Password |
|---|---|---|---|
| Admin | Superuser Admin | `adminmas` | `adminmas` |
| Supervisor | Kang Supervisor | `supervisoraja` | `supervisoraja` |
| Anggota | Kroco01 | `kroco01` | `kroco01` |

> Ganti password akun-akun ini sebelum dipakai di lingkungan produksi.

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

**Upload CSV/Excel:** klik tombol **Template** lebih dulu, lalu pilih **Template CSV**
atau **Template Excel**.
Keduanya memuat panduan format dan contoh baris untuk seluruh nilai kategori
(`Dinas`, `Umum`) dan jenis (`Bensin M`, `Bensin N&O`, `Solar JBB`,
`Solar GVW`) yang tersedia. Versi Excel dilengkapi sheet `Daftar Nilai Valid`
(daftar nilai + plat nomor yang sudah terdaftar) serta **dropdown** kategori dan
jenis pada sheet `Data Kendaraan` supaya tidak ada salah ketik.
Ganti/hapus baris contoh, isi data, lalu unggah lewat **Upload File**.
Baris yang formatnya salah ditolak per baris dengan nomor baris dan
alasannya; baris lain tetap tersimpan. Format `.xls` (Excel 97-2003) belum
didukung — simpan ulang sebagai `.xlsx` atau `.csv`.
Panduan lengkap ada di `docs/API.md` bagian 2.6.

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

**Import hasil uji (CSV/Excel):** pada panel "Import Data Hasil Uji" klik
**Template CSV** atau **Template Excel** (lewat dropdown tombol **Template**) untuk
mendapatkan file contoh yang
sudah terisi panduan, daftar plat nomor kendaraan terdaftar, dan baris contoh
siap diedit. Isi kolom `plat_nomor` + parameter uji (`co`, `co2`, `hc`, `o2`,
`lambda`, `opasitas`, `tanggal_uji`), lalu unggah lewat tombol **Upload File**.
Hasil **lulus/tidak lulus dan catatan tetap dihitung sistem** dari Parameter
aktif, sehingga hasilnya identik dengan input manual. Baris yang salah
ditolak per baris dengan nomor baris dan alasannya.
Format `.xls` (Excel 97-2003) belum didukung — simpan ulang sebagai `.xlsx`/`.csv`.
Panduan lengkap ada di `docs/API.md` bagian 3.6–3.8.

### Input Nilai Parameter (Page 4)
1. Login sebagai Admin atau Supervisor
2. Klik menu "Parameter"
3. Edit parameter emisi sesuai kebutuhan
4. Simpan perubahan

### Manajemen User (Page 5)
1. Login sebagai Admin atau Supervisor
2. Klik menu "Users"
3. Tambah, edit, atau hapus akun user
4. Catatan:
   - Username admin ditampilkan sebagai `••••••••` dan hanya dapat diubah oleh admin itu sendiri
   - Akun admin tidak dapat dihapus

## Pengujian

```bash
# Unit test backend (parser CSV/XLSX, normalisasi angka, pemetaan kolom, template)
cd backend && go test ./...

# Type check + build production frontend
cd frontend && npx tsc --noEmit && npm run build
```

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
- `POST /api/kendaraan/upload` - Upload CSV/XLSX data kendaraan
- `GET /api/kendaraan/template/csv` - Template input kendaraan (CSV)
- `GET /api/kendaraan/template/excel` - Template input kendaraan (.xlsx)

### Hasil Uji
- `GET /api/hasil-uji` - List hasil uji
- `POST /api/hasil-uji` - Tambah hasil uji
- `PUT /api/hasil-uji/:id` - Update hasil uji
- `DELETE /api/hasil-uji/:id` - Hapus hasil uji
- `GET /api/hasil-uji/statistics` - Statistik hasil uji
- `GET /api/hasil-uji/export/csv` - Export hasil uji ke CSV
- `GET /api/hasil-uji/template/csv` - Template import CSV
- `GET /api/hasil-uji/template/excel` - Template import Excel (.xlsx)
- `POST /api/hasil-uji/upload` - Import hasil uji dari CSV/XLSX

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
- `POST /api/home/working-instruction` - Upload working instruction (PDF, maks 10 MB — admin/supervisor)

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

1. **Halaman putih kosong (blank white page)**
   - `node_modules` rusak atau cache optimizer Vite basi. Perbaiki dengan:
     ```bash
     cd frontend
     rm -rf node_modules node_modules/.vite
     npm ci
     ```
   - Cara cepat mendeteksi masalah module: buka log Vite di terminal tempat `npm run dev` berjalan.

2. **PDF working instruction / gambar tidak tampil**
   - `vite.config.ts` harus memproxy `/uploads`. **Di Docker targetnya `http://backend:8080`**,
     di luar Docker `http://localhost:8080`. Tanpa proxy tersebut Vite mengembalikan
     `index.html` sehingga file tidak dapat dirender.
   - Pastikan `file_path` dari API berupa path absolut (`/uploads/home/...`).
   - Kalau muncul "Halaman X gagal ditampilkan", PDF-nya memang gagal dirender —
     buka lewat tombol **Buka PDF** untuk melihat langsung.
   - Pastikan folder `backend/uploads` ada di host (lihat butir 9 di bawah).

3. **Backend crash saat start dengan panic Static**
   - `r.Static("/uploads", "./uploads")` di `routes.go` akan panic kalau folder
     `backend/uploads` tidak ada. Buat dulu: `mkdir -p backend/uploads`.

4. **Upload CSV/Excel ditolak**
   - Unduh ulang template lewat tombol **Template** lalu pilih **Template CSV** atau
     **Template Excel**, dan pakai strukturnya.
   - Format `.xls` (Excel 97-2003) tidak didukung — simpan ulang sebagai `.xlsx` atau `.csv`.
   - Baca panel merah "N baris ditolak": nomor baris mengacu ke nomor baris fisik file.

5. **Container frontend restart terus-menerus**
   - `node_modules` belum ada di host. Compose mem-bind-mount `./frontend:/app`,
     jadi dependency harus di-install di host: `cd frontend && npm ci`.

6. **Semua request `/api` membalas 500 (ECONNREFUSED)**
   - Proxy Vite masih menunjuk `localhost`. Di dalam container, `localhost` =
     container frontend itu sendiri. Ganti ke `http://backend:8080`.

7. **Halaman membalas 403 padahal diakses lewat domain**
   - Vite ≥ 5.4 memblokir header `Host` yang tidak terdaftar. Tambahkan domainnya
     ke `server.allowedHosts` di `vite.config.ts`.

8. **Database Connection Error**
   - Pastikan PostgreSQL berjalan
   - Periksa konfigurasi di file `.env`
   - Pastikan database sudah dibuat (schema & data awal dibuat otomatis oleh GORM)

9. **Build Docker gagal `no space left on device`**
   - Bersihkan cache build: `docker builder prune -af`
   - Ruang disk yang tersedia perlu lebih besar dari ukuran image,
     karena BuildKit menyimpan dua salinan layer.

10. **Port Sudah Digunakan**
    - Ubah port di file `.env` atau `vite.config.ts`
    - Atau hentikan process yang menggunakan port

11. **Build Error**
    - Jalankan `go mod tidy` di backend
    - Jalankan `npm ci` di frontend
    - Pastikan versi Go dan Node.js sesuai

12. **Login Gagal**
    - Akun awal dibuat otomatis oleh seeding saat backend pertama kali berjalan
    - Periksa username dan password

13. **Cloudflare tunnel tidak start setelah restart**
    - `cloudflared/config.yml` hilang dari host. File itu di-bind-mount read-only,
      jadi container restart tidak akan start kalau sumbernya tidak ada.

## Lisensi

MIT License

## Kontak

Untuk pertanyaan, dukungan, atau kontribusi, silakan hubungi:
- GitHub Issues: https://github.com/khimawan/Web-app-Uji-Emisi/issues
