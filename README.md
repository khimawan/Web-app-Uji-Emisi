# Aplikasi Uji Emisi

Aplikasi web untuk mengelola dan mencatat hasil uji emisi kendaraan. Aplikasi ini memungkinkan pengguna untuk mencatat data kendaraan, melakukan uji emisi, mengimpor data massal dari CSV/Excel, dan melihat hasil uji dalam bentuk tabel serta mengekspor data ke format CSV, Excel, dan PDF.

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
- **Go 1.21+** - Bahasa pemrograman
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

> Saat development, `vite.config.ts` memproxy dua path ke backend (`localhost:8080`): `/api` untuk REST API dan `/uploads` untuk file statis (gambar home, working instruction PDF). Tanpa proxy `/uploads`, Vite akan mengembalikan `index.html` sehingga file tidak tampil. Path file dari API selalu berupa path absolut (`/uploads/home/...`).

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
│   ├── main.go                 # Entry point
│   ├── go.mod                  # Go module
│   └── .env                    # Environment variables
├── frontend/                   # Frontend React/TypeScript
│   ├── src/
│   │   ├── components/         # Reusable components
│   │   ├── pages/              # Page components
│   │   ├── services/           # API services (api.ts)
│   │   ├── types/              # TypeScript types
│   │   ├── hooks/              # Custom hooks
│   │   ├── context/            # React context (auth)
│   │   └── utils/              # Utility functions (fileUrl)
│   ├── public/                 # Static files
│   ├── index.html              # HTML template
│   ├── package.json            # Dependencies
│   ├── package-lock.json       # Lockfile (jangan dihapus)
│   ├── tsconfig.json           # TypeScript config
│   ├── vite.config.ts          # Vite config (proxy /api dan /uploads ke backend)
│   └── tailwind.config.js     # Tailwind config
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

**Upload CSV/Excel:** klik **Template CSV** atau **Template Excel** lebih dulu.
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
**Template CSV** atau **Template Excel** untuk mendapatkan file contoh yang
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
   - `vite.config.ts` harus memproxy `/uploads` ke `http://localhost:8080`. Tanpa proxy
     tersebut Vite mengembalikan `index.html` sehingga file tidak dapat dirender.
   - Pastikan `file_path` dari API berupa path absolut (`/uploads/home/...`).

3. **Upload CSV/Excel ditolak**
   - Unduh ulang template lewat tombol **Template CSV** / **Template Excel** dan pakai
     strukturnya.
   - Format `.xls` (Excel 97-2003) tidak didukung — simpan ulang sebagai `.xlsx` atau `.csv`.
   - Baca panel merah "N baris ditolak": nomor baris mengacu ke nomor baris fisik file.

4. **Database Connection Error**
   - Pastikan PostgreSQL berjalan
   - Periksa konfigurasi di file `.env`
   - Pastikan database sudah dibuat (schema & data awal dibuat otomatis oleh GORM)

5. **Port Sudah Digunakan**
   - Ubah port di file `.env` atau `vite.config.ts`
   - Atau hentikan process yang menggunakan port

6. **Build Error**
   - Jalankan `go mod tidy` di backend
   - Jalankan `npm ci` di frontend
   - Pastikan versi Go dan Node.js sesuai

7. **Login Gagal**
   - Akun awal dibuat otomatis oleh seeding saat backend pertama kali berjalan
   - Periksa username dan password

## Lisensi

MIT License

## Kontak

Untuk pertanyaan, dukungan, atau kontribusi, silakan hubungi:
- GitHub Issues: https://github.com/khimawan/Web-app-Uji-Emisi/issues
