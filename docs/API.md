# API Documentation
## Aplikasi Web Uji Emisi

**Versi:** 1.0  
**Base URL:** `http://localhost:8080/api`

---

## 1. Authentication

### 1.1 Login

**Endpoint:** `POST /api/auth/login`

**Request Body:**
```json
{
    "username": "adminmas",
    "password": "adminmas"
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Login successful",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "user": {
            "id": "uuid",
            "nama": "Superuser Admin",
            "username": "adminmas",
            "role": "admin"
        }
    }
}
```

**Response (401 Unauthorized):**
```json
{
    "success": false,
    "message": "Invalid username or password"
}
```

### 1.2 Logout

**Endpoint:** `POST /api/auth/logout`

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Logout successful"
}
```

### 1.3 Get Current User

**Endpoint:** `GET /api/auth/me`

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
    "success": true,
    "data": {
        "id": "uuid",
        "nama": "Superuser Admin",
        "username": "adminmas",
        "role": "admin"
    }
}
```

---

## 2. Kendaraan

### 2.1 List Kendaraan

**Endpoint:** `GET /api/kendaraan`

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `limit` (integer, optional): Items per page (default: 10)
- `search` (string, optional): Search by plat_nomor, merek, or tipe
- `kategori` (string, optional): Filter by kategori (Dinas/Umum)
- `jenis` (string, optional): Filter by jenis

**Response (200 OK):**
```json
{
    "success": true,
    "data": {
        "items": [
            {
                "id": "uuid",
                "kategori": "Dinas",
                "jenis": "Bensin M",
                "plat_nomor": "B 1234 CDE",
                "merek": "Toyota",
                "tipe": "Avanza",
                "kapasitas_mesin": 1500,
                "tahun_pembuatan": 2020,
                "created_by": "uuid",
                "created_at": "2026-08-16T10:00:00Z"
            }
        ],
        "pagination": {
            "total": 100,
            "page": 1,
            "limit": 10,
            "totalPages": 10
        }
    }
}
```

### 2.2 Tambah Kendaraan

**Endpoint:** `POST /api/kendaraan`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "kategori": "Dinas",
    "jenis": "Bensin M",
    "plat_nomor": "B 1234 CDE",
    "merek": "Toyota",
    "tipe": "Avanza",
    "kapasitas_mesin": 1500,
    "tahun_pembuatan": 2020
}
```

**Response (201 Created):**
```json
{
    "success": true,
    "message": "Kendaraan added successfully",
    "data": {
        "id": "uuid",
        "kategori": "Dinas",
        "jenis": "Bensin M",
        "plat_nomor": "B 1234 CDE",
        "merek": "Toyota",
        "tipe": "Avanza",
        "kapasitas_mesin": 1500,
        "tahun_pembuatan": 2020,
        "created_by": "uuid",
        "created_at": "2026-08-16T10:00:00Z"
    }
}
```

**Response (400 Bad Request):**
```json
{
    "success": false,
    "message": "Validation error",
    "errors": {
        "plat_nomor": "Plat nomor already exists"
    }
}
```

### 2.3 Detail Kendaraan

**Endpoint:** `GET /api/kendaraan/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
    "success": true,
    "data": {
        "id": "uuid",
        "kategori": "Dinas",
        "jenis": "Bensin M",
        "plat_nomor": "B 1234 CDE",
        "merek": "Toyota",
        "tipe": "Avanza",
        "kapasitas_mesin": 1500,
        "tahun_pembuatan": 2020,
        "created_by": "uuid",
        "created_at": "2026-08-16T10:00:00Z",
        "updated_at": "2026-08-16T10:00:00Z"
    }
}
```

### 2.4 Update Kendaraan

**Endpoint:** `PUT /api/kendaraan/:id`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "kategori": "Umum",
    "jenis": "Solar JBB",
    "plat_nomor": "B 1234 CDE",
    "merek": "Toyota",
    "tipe": "Avanza",
    "kapasitas_mesin": 1500,
    "tahun_pembuatan": 2020
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Kendaraan updated successfully",
    "data": {
        "id": "uuid",
        "kategori": "Umum",
        "jenis": "Solar JBB",
        "plat_nomor": "B 1234 CDE",
        "merek": "Toyota",
        "tipe": "Avanza",
        "kapasitas_mesin": 1500,
        "tahun_pembuatan": 2020,
        "created_by": "uuid",
        "created_at": "2026-08-16T10:00:00Z",
        "updated_at": "2026-08-16T11:00:00Z"
    }
}
```

### 2.5 Hapus Kendaraan

**Endpoint:** `DELETE /api/kendaraan/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Kendaraan deleted successfully"
}
```

### 2.6 Upload Data Kendaraan (CSV / Excel)

**Endpoint:** `POST /api/kendaraan/upload`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Request Body:**
```
file: <csv | xlsx | xlsm>   (maks. 5 MB)
```

**CSV Format:**
```csv
kategori,jenis,plat_nomor,merek,tipe,kapasitas_mesin,tahun_pembuatan
Dinas,Bensin M,B 1234 CDE,Toyota,Avanza,1500,2020
Umum,Solar JBB,D 5678 FGH,Honda,Civic,2000,2018
```

Template lengkap (panduan + contoh baris untuk setiap nilai kategori/jenis) dapat diunduh
dalam dua format, isinya sama persis:

| Endpoint | File | Isi |
|---|---|---|
| `GET /api/kendaraan/template/csv` | `template_kendaraan.csv` | baris panduan `#`, header, contoh baris, daftar plat nomor terdaftar |
| `GET /api/kendaraan/template/excel` | `template_kendaraan.xlsx` | 3 sheet: `Panduan`, `Data Kendaraan` (header baris 6, **dropdown** kategori/jenis, baris contoh siap diedit), `Daftar Nilai Valid` (daftar nilai + plat nomor yang sudah terdaftar) |

> Format `.xls` (Excel 97-2003) **tidak didukung** — buka di Excel lalu simpan ulang
> sebagai `.xlsx` atau `.csv`.

**Ketentuan format yang diterima import:**

| Kolom | Nilai yang diterima |
|---|---|
| `kategori` | `Dinas` / `Umum` (huruf besar-kecil bebas) |
| `jenis` | `Bensin M` / `Bensin N&O` / `Solar JBB` / `Solar GVW` |
| `plat_nomor` | Huruf kapital, maks. 20 karakter, tanpa tanda baca |
| `merek`, `tipe` | Teks bebas, maks. 100 karakter |
| `kapasitas_mesin` | Angka bulat cc (1 - 100.000). `1496`, `1.496`, `1496,00` diterima |
| `tahun_pembuatan` | Angka bulat 4 digit (1900 - tahun berjalan + 2) |

Tambahan:
- Pemisah kolom `,`, `;`, tab, atau `|` dideteksi otomatis; BOM UTF-8 dibuang.
- Baris kosong dan baris berawalan `#` dilewati.
- **Header boleh berada di baris mana pun pada 20 baris pertama**, jadi file hasil
  copy-paste dari Excel (dengan baris judul) tetap terbaca. Urutan kolom bebas bila
  nama header dikenali; bila tidak, urutan template digunakan.
- Nama kolom alternatif dikenali: `Plat Nomor`, `No. Polisi`, `Kapasitas (cc)`,
  `Tahun`, `Type`, dan sejenisnya.
- Baris yang tidak valid ditolak per baris (baris lain tetap tersimpan).
- Plat nomor yang sudah terdaftar ditolak dengan pesan ramah, bukan error database.
- Nomor baris pada pesan error mengacu ke nomor baris fisik di file.

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Data kendaraan berhasil diunggah: 8 tersimpan, 2 ditolak",
    "data": {
        "total": 10,
        "success": 8,
        "failed": 2,
        "errors": [
            {
                "row": 3,
                "message": "Plat nomor B 1234 CDE sudah terdaftar. Hapus dari file ini atau edit di halaman Input Data Kendaraan."
            },
            {
                "row": 7,
                "message": "jenis \"Benzin\" tidak dikenal, pilih salah satu: Bensin M / Bensin N&O / Solar JBB / Solar GVW"
            }
        ]
    }
}
```

---

## 3. Hasil Uji

### 3.1 List Hasil Uji

**Endpoint:** `GET /api/hasil-uji`

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `limit` (integer, optional): Items per page (default: 10)
- `search` (string, optional): Search by plat_nomor
- `hasil_uji` (string, optional): Filter by hasil (Lulus/Tidak Lulus)
- `kategori` (string, optional): Filter by kategori
- `tanggal_mulai` (string, optional): Filter by start date (YYYY-MM-DD)
- `tanggal_akhir` (string, optional): Filter by end date (YYYY-MM-DD)

**Response (200 OK):**
```json
{
    "success": true,
    "data": {
        "items": [
            {
                "id": "uuid",
                "kendaraan": {
                    "id": "uuid",
                    "plat_nomor": "B 1234 CDE",
                    "merek": "Toyota",
                    "tipe": "Avanza",
                    "kategori": "Dinas",
                    "jenis": "Bensin M"
                },
                "co": 0.5,
                "co2": 12.5,
                "hc": 100,
                "o2": 2.5,
                "lambda": 1.0,
                "opasitas": null,
                "hasil_uji": "Lulus",
                "valid": true,
                "catatan": null,
                "tested_by": "uuid",
                "tested_at": "2026-08-16T10:00:00Z"
            }
        ],
        "pagination": {
            "total": 100,
            "page": 1,
            "limit": 10,
            "totalPages": 10
        }
    }
}
```

### 3.2 Tambah Hasil Uji

**Endpoint:** `POST /api/hasil-uji`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body (Bensin):**
```json
{
    "kendaraan_id": "uuid",
    "co": 0.5,
    "co2": 12.5,
    "hc": 100,
    "o2": 2.5,
    "lambda": 1.0
}
```

**Request Body (Solar):**
```json
{
    "kendaraan_id": "uuid",
    "opasitas": 35.0
}
```

**Response (201 Created):**
```json
{
    "success": true,
    "message": "Hasil uji added successfully",
    "data": {
        "id": "uuid",
        "kendaraan_id": "uuid",
        "co": 0.5,
        "co2": 12.5,
        "hc": 100,
        "o2": 2.5,
        "lambda": 1.0,
        "opasitas": null,
        "hasil_uji": "Lulus",
        "valid": true,
        "catatan": null,
        "tested_by": "uuid",
        "tested_at": "2026-08-16T10:00:00Z"
    }
}
```

### 3.3 Update Hasil Uji

**Endpoint:** `PUT /api/hasil-uji/:id`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "co": 0.6,
    "co2": 13.0,
    "hc": 110,
    "o2": 2.8,
    "lambda": 1.05
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Hasil uji updated successfully",
    "data": {
        "id": "uuid",
        "kendaraan_id": "uuid",
        "co": 0.6,
        "co2": 13.0,
        "hc": 110,
        "o2": 2.8,
        "lambda": 1.05,
        "opasitas": null,
        "hasil_uji": "Lulus",
        "valid": true,
        "catatan": null,
        "tested_by": "uuid",
        "tested_at": "2026-08-16T10:00:00Z"
    }
}
```

### 3.4 Hapus Hasil Uji

**Endpoint:** `DELETE /api/hasil-uji/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Hasil uji deleted successfully"
}
```

### 3.5 Statistik

**Endpoint:** `GET /api/hasil-uji/statistics`

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
    "success": true,
    "data": {
        "total_kendaraan": 100,
        "total_dinas": 45,
        "total_umum": 55,
        "total_lulus": 80,
        "total_tidak_lulus": 20,
        "total_bensin": 60,
        "total_solar": 40,
        "by_kategori": [
            {
                "kategori": "Dinas",
                "count": 45,
                "percentage": 45.0
            },
            {
                "kategori": "Umum",
                "count": 55,
                "percentage": 55.0
            }
        ],
        "by_hasil": [
            {
                "hasil": "Lulus",
                "count": 80,
                "percentage": 80.0
            },
            {
                "hasil": "Tidak Lulus",
                "count": 20,
                "percentage": 20.0
            }
        ],
        "by_jenis": [
            {
                "jenis": "Bensin",
                "count": 60,
                "percentage": 60.0
            },
            {
                "jenis": "Solar",
                "count": 40,
                "percentage": 40.0
            }
        ]
    }
}
```

### 3.6 Download Template Import (CSV)

**Endpoint:** `GET /api/hasil-uji/template/csv`

**Headers:**
```
Authorization: Bearer <token>
```

Menghasilkan file `template_import_hasil_uji.csv` berisi baris panduan (diawali `#`), baris
header, baris contoh, dan daftar plat nomor kendaraan yang sudah terdaftar.

Header (kolom boleh diurutkan ulang, boleh dipisah `,` `;` tab atau `|`, boleh diberi satuan
dalam kurung seperti `CO (%)`):

| Kolom | Wajib | Rentang nilai | Contoh |
|---|---|---|---|
| `plat_nomor` | Ya | harus terdaftar di Input Kendaraan | `B 1234 CDE` |
| `co` | minimal 1 | 0 – 100 (%) | `0.85` |
| `co2` | | 0 – 100 (%) | `12.5` |
| `hc` | | 0 – 10000 (ppm) | `120` |
| `o2` | | 0 – 100 (%) | `1.2` |
| `lambda` | | 0.1 – 10 | `0.99` |
| `opasitas` | | 0 – 100 (%) | `12.5` |
| `tanggal_uji` | | `YYYY-MM-DD`, `YYYY/MM/DD` atau `DD/MM/YYYY` | `2024-01-15` |

> Nilai desimal boleh memakai `.` maupun `,` (otomatis dideteksi dari pola). Thousand
> separator juga diterima, contoh `1.234` dan `1,234` dibaca sebagai `1234`.
> Kolom `hasil_uji`, `valid`, dan `catatan` **tidak** diimpor — sistem menghitungnya
> otomatis memakai Parameter aktif yang berlaku sesuai tahun pembuatan kendaraan.

### 3.7 Download Template Import (Excel)

**Endpoint:** `GET /api/hasil-uji/template/excel`

Sama seperti 3.6 tetapi berformat `.xlsx` (3 sheet: `Panduan`, `Data Hasil Uji`,
`Daftar Kendaraan`). Format `.xls` (Excel 97-2003) **tidak didukung** — buka di Excel
lalu simpan ulang sebagai `.xlsx` atau `.csv`.

### 3.8 Upload Hasil Uji (CSV / XLSX)

**Endpoint:** `POST /api/hasil-uji/upload`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Body:** field `file` (`.csv`, `.txt`, `.xlsx`, `.xlsm`; maksimal 5 MB).

Alur pemrosesan:

1. Header dideteksi pada 20 baris pertama, jadi file hasil salin-tempel dari Excel tetap aman.
2. Setiap baris divalidasi (plat terdaftar, rentang nilai, format tanggal).
3. `hasil_uji` / `valid` / `catatan` dihitung otomatis oleh mesin evaluasi emisi.
4. Kendaraan yang sudah memiliki hasil uji akan ditolak agar tidak menimpa data.

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Import selesai: 2 tersimpan, 1 ditolak",
    "data": {
        "total": 3,
        "success": 2,
        "failed": 1,
        "lulus": 1,
        "tidak_lulus": 1,
        "errors": [
            { "row": 8, "message": "co 9999 di luar rentang yang diperbolehkan (0 - 100 %)" }
        ]
    }
}
```

> `row` adalah nomor baris fisik pada file (termasuk baris panduan/header), sehingga
> mudah dilacak di Excel. Field `errors_truncated` bernilai `true` bila daftar error
> dipotong pada 50 baris pertama.

Contoh pesan penolakan:

| Pesan | Penyebab |
|---|---|
| `Kendaraan dengan plat nomor B 1234 CDE tidak terdaftar. Daftarkan dulu di halaman 'Input Data Kendaraan'.` | plat nomor tidak ada di master kendaraan |
| `Kendaraan dengan plat nomor B 1234 CDE sudah memiliki data hasil uji.` | hasil uji untuk kendaraan itu sudah ada |
| `co "abc" bukan angka yang valid (contoh: 0.85)` | nilai bukan angka |
| `lambda 15 di luar rentang yang diperbolehkan (0.1 - 10)` | nilai di luar rentang |
| `tanggal_uji "2024-13-01" bukan tanggal yang valid` | tanggal tidak nyata |
| `Baris X plat nomor tidak boleh kosong` | plat nomor kosong |

---

## 4. Parameters

### 4.1 List Parameters

**Endpoint:** `GET /api/parameters`

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `kategori` (string, optional): Filter by kategori

**Response (200 OK):**
```json
{
    "success": true,
    "data": [
        {
            "id": "uuid",
            "kategori": "Bensin M",
            "tahun_min": 0,
            "tahun_max": 2007,
            "tahun_operator": "<",
            "co_max": 4.0,
            "hc_max": 1000.0,
            "opasitas_max": null,
            "is_active": true,
            "created_at": "2026-08-16T10:00:00Z"
        }
    ]
}
```

### 4.2 Tambah Parameter

**Endpoint:** `POST /api/parameters`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "kategori": "Bensin M",
    "tahun_min": 0,
    "tahun_max": 2007,
    "tahun_operator": "<",
    "co_max": 4.0,
    "hc_max": 1000.0,
    "opasitas_max": null
}
```

**Response (201 Created):**
```json
{
    "success": true,
    "message": "Parameter added successfully",
    "data": {
        "id": "uuid",
        "kategori": "Bensin M",
        "tahun_min": 0,
        "tahun_max": 2007,
        "tahun_operator": "<",
        "co_max": 4.0,
        "hc_max": 1000.0,
        "opasitas_max": null,
        "is_active": true,
        "created_at": "2026-08-16T10:00:00Z"
    }
}
```

### 4.3 Update Parameter

**Endpoint:** `PUT /api/parameters/:id`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "co_max": 3.5,
    "hc_max": 900.0
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Parameter updated successfully",
    "data": {
        "id": "uuid",
        "kategori": "Bensin M",
        "tahun_min": 0,
        "tahun_max": 2007,
        "tahun_operator": "<",
        "co_max": 3.5,
        "hc_max": 900.0,
        "opasitas_max": null,
        "is_active": true,
        "created_at": "2026-08-16T10:00:00Z",
        "updated_at": "2026-08-16T11:00:00Z"
    }
}
```

### 4.4 Hapus Parameter

**Endpoint:** `DELETE /api/parameters/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Parameter deleted successfully"
}
```

---

## 5. Pop-up Notes

### 5.1 List Pop-up Notes

**Endpoint:** `GET /api/popup-notes`

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
    "success": true,
    "data": [
        {
            "id": "uuid",
            "jenis_kendaraan": "Bensin M, Bensin N&O",
            "tahun_operator": ">",
            "tahun_value": 1990,
            "parameter_uji": "O2",
            "nilai_operator": ">",
            "nilai_value": 20.0,
            "note": "Indikasi kebocoran exhaust",
            "is_active": true,
            "created_at": "2026-08-16T10:00:00Z"
        }
    ]
}
```

### 5.2 Tambah Pop-up Note

**Endpoint:** `POST /api/popup-notes`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "jenis_kendaraan": "Bensin M, Bensin N&O",
    "tahun_operator": ">",
    "tahun_value": 1990,
    "parameter_uji": "O2",
    "nilai_operator": ">",
    "nilai_value": 20.0,
    "note": "Indikasi kebocoran exhaust"
}
```

**Response (201 Created):**
```json
{
    "success": true,
    "message": "Pop-up note added successfully",
    "data": {
        "id": "uuid",
        "jenis_kendaraan": "Bensin M, Bensin N&O",
        "tahun_operator": ">",
        "tahun_value": 1990,
        "parameter_uji": "O2",
        "nilai_operator": ">",
        "nilai_value": 20.0,
        "note": "Indikasi kebocoran exhaust",
        "is_active": true,
        "created_at": "2026-08-16T10:00:00Z"
    }
}
```

### 5.3 Update Pop-up Note

**Endpoint:** `PUT /api/popup-notes/:id`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "note": "Indikasi kebocoran exhaust (updated)"
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Pop-up note updated successfully",
    "data": {
        "id": "uuid",
        "jenis_kendaraan": "Bensin M, Bensin N&O",
        "tahun_operator": ">",
        "tahun_value": 1990,
        "parameter_uji": "O2",
        "nilai_operator": ">",
        "nilai_value": 20.0,
        "note": "Indikasi kebocoran exhaust (updated)",
        "is_active": true,
        "created_at": "2026-08-16T10:00:00Z",
        "updated_at": "2026-08-16T11:00:00Z"
    }
}
```

### 5.4 Hapus Pop-up Note

**Endpoint:** `DELETE /api/popup-notes/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Pop-up note deleted successfully"
}
```

---

## 6. Home Content

### 6.1 Get Home Content

**Endpoint:** `GET /api/home`

**Response (200 OK):**
```json
{
    "success": true,
    "data": {
        "descriptions": [
            {
                "id": "uuid",
                "title": "Tentang Aplikasi",
                "description": "Aplikasi Uji Emisi adalah solusi web-based...",
                "sort_order": 0
            }
        ],
        "images": [
            {
                "id": "uuid",
                "file_path": "/uploads/home/image1.jpg",
                "file_type": "image/jpeg",
                "sort_order": 0
            }
        ],
        "working_instruction": {
            "id": "uuid",
            "file_path": "uploads/home/instruction.pdf",
            "file_type": "application/pdf"
        }
    }
}
```

> `working_instruction` bernilai `null` jika belum ada file aktif. Endpoint ini publik (tanpa token) sehingga seluruh pengguna bisa melihat dan mengunduh; hanya `admin`/`supervisor` yang boleh mengunggah.

> **`file_path` harus berupa path absolut** (diawali `/`), contoh `/uploads/home/xxx.pdf`. Path relatif akan resolve ke origin frontend dan gagal dimuat. Di mode development, `vite.config.ts` memproxy `/uploads` ke `http://localhost:8080`; di production backend dan SPA harus berada pada origin yang sama.

### 6.2 Update Description

**Endpoint:** `PUT /api/home/description/:id`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "title": "Tentang Aplikasi",
    "description": "Aplikasi Uji Emisi adalah solusi web-based..."
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Description updated successfully"
}
```

### 6.3 Upload Image

**Endpoint:** `POST /api/home/image`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Request Body:**
```
file: <image_file>
title: "Gambar Aplikasi"
```

**Response (201 Created):**
```json
{
    "success": true,
    "message": "Image uploaded successfully",
    "data": {
        "id": "uuid",
        "file_path": "/uploads/home/image1.jpg",
        "file_type": "image/jpeg",
        "title": "Gambar Aplikasi"
    }
}
```

### 6.4 Delete Image

**Endpoint:** `DELETE /api/home/image/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Image deleted successfully"
}
```

### 6.5 Upload Working Instruction

**Endpoint:** `POST /api/home/working-instruction`

**Akses:** `admin` atau `supervisor` (403 `Access denied` untuk role lain).

**Headers:**
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Request Body:**
```
file: <pdf_file>
```

**Aturan validasi:**

| Aturan | Nilai |
| --- | --- |
| Ekstensi | `.pdf` (case-insensitive) |
| Magic bytes | harus diawali `%PDF-` |
| Ukuran maksimal | 10 MB |
| Jumlah file aktif | 1 — upload baru otomatis menonaktifkan file sebelumnya |

**Response (201 Created):**
```json
{
    "success": true,
    "message": "Working instruction berhasil diunggah",
    "data": {
        "id": "uuid",
        "content_type": "working_instruction",
        "title": "Instruksi Kerja 2026.pdf",
        "file_path": "uploads/home/uuid.pdf",
        "file_type": "application/pdf",
        "is_active": true
    }
}
```

**Response (400 Bad Request):**
```json
{
    "success": false,
    "message": "File bukan PDF yang valid"
}
```

### 6.6 Delete Working Instruction

**Endpoint:** `DELETE /api/home/working-instruction/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Working instruction deleted successfully"
}
```

---

## 7. Users (Admin)

### 7.1 List Users

**Endpoint:** `GET /api/users`

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `limit` (integer, optional): Items per page (default: 10)
- `role` (string, optional): Filter by role

**Response (200 OK):**
```json
{
    "success": true,
    "data": {
        "items": [
            {
                "id": "uuid",
                "nama": "Superuser Admin",
                "username": "adminmas",
                "role": "admin",
                "created_at": "2026-08-16T10:00:00Z"
            }
        ],
        "pagination": {
            "total": 3,
            "page": 1,
            "limit": 10,
            "totalPages": 1
        }
    }
}
```

### 7.2 Tambah User

**Endpoint:** `POST /api/users`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "nama": "New User",
    "username": "newuser",
    "password": "password123",
    "role": "anggota"
}
```

**Response (201 Created):**
```json
{
    "success": true,
    "message": "User added successfully",
    "data": {
        "id": "uuid",
        "nama": "New User",
        "username": "newuser",
        "role": "anggota",
        "created_at": "2026-08-16T10:00:00Z"
    }
}
```

### 7.3 Update User

**Endpoint:** `PUT /api/users/:id`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "nama": "Updated Name",
    "password": "newpassword123",
    "role": "supervisor"
}
```

| Field | Wajib | Catatan |
|---|---|---|
| `nama` | tidak | Diabaikan bila kosong |
| `username` | tidak | Diabaikan bila kosong. Untuk akun `admin`, hanya boleh diisi oleh admin itu sendiri |
| `password` | tidak | Diabaikan bila kosong. Untuk akun `admin`, hanya boleh diisi oleh admin itu sendiri |
| `role` | tidak | Diabaikan bila kosong |

> **Kunci kredensial admin.** Bila target user punya role `admin` dan pemohon bukan admin itu sendiri, pengiriman `username` atau `password` akan ditolak `403`:
> `"Username admin hanya dapat diubah oleh admin itu sendiri"` / `"Password admin hanya dapat diubah oleh admin itu sendiri"`. Field `nama` dan `role` tetap boleh diubah.

**Response (403 Forbidden):**
```json
{
    "success": false,
    "message": "Password admin hanya dapat diubah oleh admin itu sendiri"
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "User updated successfully",
    "data": {
        "id": "uuid",
        "nama": "Updated Name",
        "username": "newuser",
        "role": "supervisor",
        "created_at": "2026-08-16T10:00:00Z",
        "updated_at": "2026-08-16T11:00:00Z"
    }
}
```

### 7.4 Hapus User

**Endpoint:** `DELETE /api/users/:id`

> User dengan role `admin` tidak dapat dihapus — dilindungi di backend dengan `400 Cannot delete admin user`.

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "User deleted successfully"
}
```

**Response (400 Bad Request):**
```json
{
    "success": false,
    "message": "Cannot delete admin user"
}
```

---

## 8. Error Responses

### 400 Bad Request
```json
{
    "success": false,
    "message": "Validation error",
    "errors": {
        "field": "Error message"
    }
}
```

### 401 Unauthorized
```json
{
    "success": false,
    "message": "Unauthorized"
}
```

### 403 Forbidden
```json
{
    "success": false,
    "message": "Forbidden"
}
```

### 404 Not Found
```json
{
    "success": false,
    "message": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
    "success": false,
    "message": "Internal server error"
}
```

---

**Document Version:** 1.0  
**Last Updated:** Agustus 2026
