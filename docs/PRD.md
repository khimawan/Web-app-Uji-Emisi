# Product Requirements Document (PRD)
## Aplikasi Web Uji Emisi

**Versi:** 1.0  
**Tanggal:** Agustus 2026  
**Status:** Draft

---

## 1. Ringkasan Proyek

Aplikasi Web Uji Emisi adalah sistem informasi berbasis web yang dirancang untuk mengelola dan mencatat hasil uji emisi kendaraan bermotor. Sistem ini memungkinkan petugas untuk mencatat data kendaraan, melakukan pengujian emisi, mengevaluasi hasil berdasarkan standar yang berlaku, dan menghasilkan laporan analitis.

## 2. Tujuan

### 2.1 Tujuan Bisnis
- Mempermudah pencatatan data kendaraan dan hasil uji emisi
- Menyediakan sistem validasi otomatis berdasarkan standar emisi
- Menghasilkan laporan dan analisis data secara real-time
- Mendukung pengelolaan data yang efisien dan akurat

### 2.2 Tujuan Teknis
- Membangun aplikasi web modern dengan performa tinggi
- Mengimplementasikan arsitektur microservice-ready
- Menyediakan API yang RESTful untuk integrasi sistem lain
- Memastikan keamanan data dan autentikasi yang robust

## 3. Pengguna Target

| Role | Deskripsi | Akses |
|------|-----------|-------|
| Admin | Administrator sistem | Full access |
| Supervisor | Pengawas | Manage anggota, edit parameter |
| Anggota | Petugas uji emisi | Input data kendaraan & emisi |

## 4. Fitur Utama

### 4.1 Autentikasi (Auth Page)
**Prioritas:** High

- Form login dengan username dan password
- Session management dengan JWT
- Role-based access control
- Logout functionality

**Acceptance Criteria:**
- [ ] User dapat login dengan username dan password
- [ ] Session expire setelah 24 jam tidak aktif
- [ ] User tidak dapat mengakses halaman yang tidak sesuai role
- [ ] Error message ditampilkan jika login gagal

### 4.2 Home Page
**Prioritas:** High

- Deskripsi aplikasi dengan gambar
- Working instruction dalam format PDF
- Edit, hapus, dan upload konten oleh Admin/Supervisor

**Acceptance Criteria:**
- [ ] Deskripsi dan gambar dapat di-edit oleh Admin/Supervisor
- [ ] Gambar dapat di-upload dan dihapus
- [ ] PDF working instruction dapat di-upload dan dilihat
- [ ] Perubahan terlihat secara real-time

### 4.3 Page 1 - Input Data Kendaraan
**Prioritas:** High

**Form Input:**
| Field | Tipe | Opsi | Validasi |
|-------|------|------|----------|
| Kategori Kendaraan | Select | Dinas/Umum | Required |
| Jenis Kendaraan | Select | Bensin M, Bensin N&O, Solar JBB, Solar GVW | Required |
| Plat Nomor | Text | Contoh: B 1234 CDE | Required, Unique |
| Merek Mobil | Text | Contoh: Toyota | Required |
| Tipe Mobil | Text | Contoh: Avanza | Required |
| Kapasitas Mesin (cc) | Number | Contoh: 1500 | Required, Min: 500 |
| Tahun Pembuatan | Number | Contoh: 2020 | Required, Range: 1900-2026 |

**Fitur Tambahan:**
- Download template CSV
- Upload CSV untuk input batch
- Tabel history input
- Edit dan delete multiple line

**Acceptance Criteria:**
- [ ] Form validasi real-time
- [ ] Template CSV dapat diunduh
- [ ] CSV upload dengan validasi format
- [ ] Tabel menampilkan semua data
- [ ] Edit inline pada tabel
- [ ] Multiple delete dengan checkbox

### 4.4 Page 2 - Input Data Emisi
**Prioritas:** High

**Form Input:**
| Field | Keterangan |
|-------|------------|
| Pilih Kendaraan | Search bar + dropdown |
| CO (%) | Input untuk bensin |
| CO2 (%) | Input untuk bensin |
| HC (ppm) | Input untuk bensin |
| O2 (%) | Input untuk bensin |
| Lambda | Input untuk bensin |
| Opasitas (%) | Input untuk solar |

**Pop-up Notifikasi:**
- Hasil Uji: Lulus/Tidak Lulus
- Valid/Tidak
- Catatan perbaikan

**Acceptance Criteria:**
- [ ] Search bar berfungsi untuk mencari kendaraan
- [ ] Form berbeda untuk bensin dan solar
- [ ] Bensin: tidak perlu input opasitas
- [ ] Solar: hanya input opasitas
- [ ] Pop-up muncul setelah input
- [ ] Hasil uji otomatis berdasarkan parameter

### 4.5 Page 3 - Data Hasil Uji
**Prioritas:** High

**Dashboard (Atas):**
- Total kendaraan masuk
- Total kendaraan dinas
- Total kendaraan umum
- Total lulus uji emisi
- Total tidak lulus uji emisi

**Charts (Tengah):**
- Pie chart: Dinas vs Umum
- Pie chart: Lulus vs Tidak Lulus
- Pie chart: Bensin vs Solar
- Bar chart: Jumlah per kategori

**Tabel Data:**
| Kolom | Keterangan |
|-------|------------|
| No | Nomor urut |
| Waktu Input | Tanggal dan jam |
| Kategori Kendaraan | Dinas/Umum |
| Jenis Kendaraan | Bensin/Solar |
| Plat Nomor | Plat kendaraan |
| Merek Kendaraan | Merek |
| Tipe Kendaraan | Tipe |
| Kapasitas Mesin | cc |
| Tahun Pembuatan | Tahun |
| CO | Persentase |
| CO2 | Persentase |
| HC | ppm |
| O2 | Persentase |
| Lambda | Nilai |
| Opasitas | Persentase |
| Hasil Uji | Lulus/Tidak |
| Valid | Valid/Tidak |
| Catatan | Catatan |

**Fitur:**
- Edit data
- Multiple delete
- Export CSV, Excel, PDF
- Scroll up/down

**Acceptance Criteria:**
- [ ] Data real-time dari database
- [ ] Charts akurat dan interaktif
- [ ] Tabel dengan pagination
- [ ] Export berfungsi dengan benar
- [ ] Edit inline pada tabel

### 4.6 Page 4 - Input Nilai Parameter
**Prioritas:** High

**Hanya dapat diakses oleh Admin dan Supervisor**

**Tabel Parameter:**
| Kategori | Tahun Pembuatan | CO(%) | HC(ppm) | Opasitas(%) |
|----------|-----------------|-------|---------|--------------|
| Bensin M | <2007 | 4 | 1000 | - |
| Bensin M | 2007-2018 | 1 | 150 | - |
| Bensin M | >2018 | 0.5 | 100 | - |
| Bensin N&O | <2007 | 4 | 1000 | - |
| Bensin N&O | 2007-2018 | 1 | 150 | - |
| Bensin N&O | >2018 | 0.5 | 100 | - |
| Solar JBB | <2010 | - | - | 65 |
| Solar JBB | 2010-2021 | - | - | 40 |
| Solar JBB | >2021 | - | - | 30 |
| Solar GVW | <2010 | - | - | 65 |
| Solar GVW | 2010-2021 | - | - | 40 |
| Solar GVW | >2021 | - | - | 35 |

**Pop-up Notes Configuration:**
| Jenis Kendaraan | Tahun | Parameter | Nilai | Note |
|-----------------|-------|-----------|-------|------|
| Bensin M, Bensin N&O | >1990 | O2 | >20 | Indikasi kebocoran exhaust |
| Bensin M | <2018 | HC | >4 | Indikasi kerusakan: busi mati/lemah |
| Solar JBB | >2021 | Opasitas | >70 | Indikasi Boros BBM |

**Acceptance Criteria:**
- [ ] Hanya Admin/Supervisor yang dapat mengakses
- [ ] Parameter dapat ditambah, edit, hapus
- [ ] Pop-up notes dapat dikonfigurasi
- [ ] Perubahan berlaku untuk semua user
- [ ] Validasi input parameter

### 4.7 Page 5 - Halaman Admin
**Prioritas:** High

**Hanya dapat diakses oleh Admin dan Supervisor**

**Tabel User:**
| Kolom | Keterangan |
|-------|------------|
| Nama | Nama lengkap |
| Username | Username |
| Password | Password |
| Role | Anggota/Supervisor/Admin |

**Fitur:**
- Tambah user baru
- Edit user
- Hapus user (kecuali admin)

**Acceptance Criteria:**
- [ ] Admin/Supervisor dapat mengakses
- [ ] Admin tidak dapat dihapus
- [ ] Password di-hash saat disimpan
- [ ] Role hanya dapat dipilih dari opsi yang tersedia
- [ ] Validasi username unik

## 5. Non-Functional Requirements

### 5.1 Performance
- Response time < 2 detik untuk semua operasi
- Support 100 concurrent users
- Database query < 500ms

### 5.2 Security
- Password hashing dengan bcrypt
- JWT token dengan expiry
- CORS configuration
- Input validation dan sanitization
- SQL injection prevention

### 5.3 Usability
- Responsive design (mobile, tablet, desktop)
- Intuitive navigation
- Loading indicators
- Error messages yang jelas
- Form validation real-time

### 5.4 Reliability
- 99.9% uptime
- Graceful error handling
- Data backup capability
- Recovery procedure

### 5.5 Scalability
- Horizontal scaling ready
- Database connection pooling
- Caching strategy

## 6. Tech Stack

### Backend
- Go 1.21+
- Gin Framework
- GORM ORM
- JWT Authentication

### Database
- PostgreSQL 14+

### Frontend
- React 18
- TypeScript
- Vite
- Tailwind CSS
- Recharts

### DevOps
- Docker (optional)
- GitHub Actions (CI/CD)

## 7. Timeline

| Phase | Duration | Deliverables |
|-------|----------|--------------|
| Phase 1 | 2 minggu | Auth, Home Page, DB Setup |
| Phase 2 | 2 minggu | Page 1 & 2 |
| Phase 3 | 2 minggu | Page 3 & 4 |
| Phase 4 | 1 minggu | Page 5, Testing |
| Phase 5 | 1 minggu | Bug fixing, Documentation |

**Total:** 8 minggu

## 8. Success Metrics

- Semua fitur berfungsi sesuai requirements
- Tidak ada critical bug
- User acceptance testing (UAT) passed
- Performance metrics terpenuhi
- Dokumentasi lengkap

## 9. Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Database performance | High | Indexing, query optimization |
| Security vulnerability | High | Regular security audit |
| Scope creep | Medium | Strict requirements management |
| Team availability | Medium | Knowledge sharing, documentation |

## 10. Approval

| Role | Name | Signature | Date |
|------|------|-----------|------|
| Product Owner | | | |
| Tech Lead | | | |
| QA Lead | | | |

---

**Document Version:** 1.0  
**Last Updated:** Agustus 2026
