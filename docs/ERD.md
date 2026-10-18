# Entity Relationship Diagram (ERD)
## Aplikasi Web Uji Emisi

**Versi:** 1.0  
**Tanggal:** Agustus 2026

---

## 1. Overview Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                              ENTITY RELATIONSHIP DIAGRAM                                │
│                           Aplikasi Web Uji Emisi v1.0                                   │
└─────────────────────────────────────────────────────────────────────────────────────────┘

                                    ┌───────────────┐
                                    │     users     │
                                    └───────┬───────┘
                                            │
                    ┌───────────────────────┼───────────────────────┐
                    │                       │                       │
                    ▼                       ▼                       ▼
            ┌───────────────┐       ┌───────────────┐       ┌───────────────┐
            │   kendaraan   │       │  parameters   │       │  popup_notes  │
            └───────┬───────┘       └───────────────┘       └───────────────┘
                    │
                    ▼
            ┌───────────────┐
            │   hasil_uji   │
            └───────────────┘
```

## 2. Entity Definitions

### 2.1 users

Tabel untuk menyimpan data pengguna sistem.

```
┌─────────────────────────────────────────────────────────────┐
│                          users                              │
├─────────────────────────────────────────────────────────────┤
│ id              UUID            PRIMARY KEY, NOT NULL       │
│ nama            VARCHAR(255)    NOT NULL                    │
│ username        VARCHAR(100)    UNIQUE, NOT NULL            │
│ password        VARCHAR(255)    NOT NULL                    │
│ role            ENUM            NOT NULL                    │
│ created_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP   │
│ updated_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP   │
└─────────────────────────────────────────────────────────────┘

Enum Values for role:
- 'admin'
- 'supervisor'
- 'anggota'
```

**Constraints:**
- Primary Key: `id`
- Unique: `username`
- Check: `role` IN ('admin', 'supervisor', 'anggota')

### 2.2 kendaraan

Tabel untuk menyimpan data kendaraan.

```
┌─────────────────────────────────────────────────────────────┐
│                        kendaraan                            │
├─────────────────────────────────────────────────────────────┤
│ id              UUID            PRIMARY KEY, NOT NULL       │
│ kategori        VARCHAR(20)     NOT NULL                    │
│ jenis           VARCHAR(50)     NOT NULL                    │
│ plat_nomor      VARCHAR(20)     UNIQUE, NOT NULL            │
│ merek           VARCHAR(100)    NOT NULL                    │
│ tipe            VARCHAR(100)    NOT NULL                    │
│ kapasitas_mesin INTEGER         NOT NULL                    │
│ tahun_pembuatan INTEGER         NOT NULL                    │
│ created_by      UUID            FK -> users(id)             │
│ created_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP   │
│ updated_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP   │
└─────────────────────────────────────────────────────────────┘

Enum Values for kategori:
- 'Dinas'
- 'Umum'

Enum Values for jenis:
- 'Bensin M'
- 'Bensin N&O'
- 'Solar JBB'
- 'Solar GVW'
```

**Constraints:**
- Primary Key: `id`
- Foreign Key: `created_by` REFERENCES `users(id)`
- Unique: `plat_nomor`
- Check: `kategori` IN ('Dinas', 'Umum')
- Check: `jenis` IN ('Bensin M', 'Bensin N&O', 'Solar JBB', 'Solar GVW')
- Check: `kapasitas_mesin` > 0
- Check: `tahun_pembuatan` BETWEEN 1900 AND 2026

### 2.3 hasil_uji

Tabel untuk menyimpan hasil uji emisi kendaraan.

```
┌─────────────────────────────────────────────────────────────┐
│                       hasil_uji                             │
├─────────────────────────────────────────────────────────────┤
│ id              UUID            PRIMARY KEY, NOT NULL       │
│ kendaraan_id    UUID            FK -> kendaraan(id)         │
│ co              DECIMAL(10,4)   NULL                        │
│ co2             DECIMAL(10,4)   NULL                        │
│ hc              DECIMAL(10,2)   NULL                        │
│ o2              DECIMAL(10,4)   NULL                        │
│ lambda          DECIMAL(10,4)   NULL                        │
│ opasitas        DECIMAL(10,2)   NULL                        │
│ hasil_uji       VARCHAR(20)     NOT NULL                    │
│ valid           BOOLEAN         NOT NULL                    │
│ catatan         TEXT            NULL                        │
│ tested_by       UUID            FK -> users(id)             │
│ tested_at       TIMESTAMP       DEFAULT CURRENT_TIMESTAMP   │
│ created_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP   │
│ updated_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP   │
└─────────────────────────────────────────────────────────────┘

Enum Values for hasil_uji:
- 'Lulus'
- 'Tidak Lulus'
```

**Constraints:**
- Primary Key: `id`
- Foreign Key: `kendaraan_id` REFERENCES `kendaraan(id)`
- Foreign Key: `tested_by` REFERENCES `users(id)`
- Check: `hasil_uji` IN ('Lulus', 'Tidak Lulus')
- Check: `co` >= 0 (if not null)
- Check: `co2` >= 0 (if not null)
- Check: `hc` >= 0 (if not null)
- Check: `o2` >= 0 AND o2 <= 100 (if not null)
- Check: `opasitas` >= 0 AND opasitas <= 100 (if not null)

### 2.4 parameters

Tabel untuk menyimpan parameter standar emisi.

```
┌─────────────────────────────────────────────────────────────┐
│                       parameters                            │
├─────────────────────────────────────────────────────────────┤
│ id              UUID            PRIMARY KEY, NOT NULL       │
│ kategori        VARCHAR(50)     NOT NULL                    │
│ tahun_min       INTEGER         NOT NULL                    │
│ tahun_max       INTEGER         NULL                        │
│ tahun_operator  VARCHAR(10)     NOT NULL                    │
│ co_max          DECIMAL(10,4)   NULL                        │
│ hc_max          DECIMAL(10,2)   NULL                        │
│ opasitas_max    DECIMAL(10,2)   NULL                        │
│ is_active       BOOLEAN         DEFAULT TRUE                │
│ created_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP   │
│ updated_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP   │
└─────────────────────────────────────────────────────────────┘

Enum Values for kategori:
- 'Bensin M'
- 'Bensin N&O'
- 'Solar JBB'
- 'Solar GVW'

Enum Values for tahun_operator:
- '<' (kurang dari)
- 'between' (antara dua tahun)
- '>' (lebih dari)
```

**Constraints:**
- Primary Key: `id`
- Check: `kategori` IN ('Bensin M', 'Bensin N&O', 'Solar JBB', 'Solar GVW')
- Check: `tahun_operator` IN ('<', 'between', '>')
- Check: `co_max` >= 0 (if not null)
- Check: `hc_max` >= 0 (if not null)
- Check: `opasitas_max` >= 0 (if not null)

### 2.5 popup_notes

Tabel untuk menyimpan konfigurasi catatan pop-up notifikasi.

```
┌─────────────────────────────────────────────────────────────┐
│                       popup_notes                           │
├─────────────────────────────────────────────────────────────┤
│ id              UUID            PRIMARY KEY, NOT NULL       │
│ jenis_kendaraan VARCHAR(255)    NOT NULL                    │
│ tahun_operator  VARCHAR(10)     NOT NULL                    │
│ tahun_value     INTEGER         NOT NULL                    │
│ parameter_uji   VARCHAR(20)     NOT NULL                    │
│ nilai_operator  VARCHAR(10)     NOT NULL                    │
│ nilai_value     DECIMAL(10,4)   NOT NULL                    │
│ note            TEXT            NOT NULL                    │
│ is_active       BOOLEAN         DEFAULT TRUE                │
│ created_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP   │
│ updated_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP   │
└─────────────────────────────────────────────────────────────┘

Enum Values for parameter_uji:
- 'CO'
- 'CO2'
- 'HC'
- 'O2'
- 'Lambda'
- 'Opasitas'

Enum Values for tahun_operator & nilai_operator:
- '>'
- '<'
- '>='
- '<='
- '='
```

**Constraints:**
- Primary Key: `id`
- Check: `parameter_uji` IN ('CO', 'CO2', 'HC', 'O2', 'Lambda', 'Opasitas')
- Check: `tahun_operator` IN ('>', '<', '>=', '<=', '=')
- Check: `nilai_operator` IN ('>', '<', '>=', '<=', '=')

### 2.6 home_content

Tabel untuk menyimpan konten halaman home (deskripsi, gambar, working instruction).

```
┌─────────────────────────────────────────────────────────────┐
│                      home_content                           │
├─────────────────────────────────────────────────────────────┤
│ id              UUID            PRIMARY KEY, NOT NULL       │
│ content_type    VARCHAR(50)     NOT NULL                    │
│ title           VARCHAR(255)    NULL                        │
│ description     TEXT            NULL                        │
│ file_path       VARCHAR(500)    NULL                        │
│ file_type       VARCHAR(50)     NULL                        │
│ sort_order      INTEGER         DEFAULT 0                   │
│ is_active       BOOLEAN         DEFAULT TRUE                │
│ created_by      UUID            FK -> users(id)             │
│ created_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP   │
│ updated_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP   │
└─────────────────────────────────────────────────────────────┘

Enum Values for content_type:
- 'description'
- 'image'
- 'working_instruction'
```

**Constraints:**
- Primary Key: `id`
- Foreign Key: `created_by` REFERENCES `users(id)`
- Check: `content_type` IN ('description', 'image', 'working_instruction')

## 3. Relationships

### 3.1 Relationship Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                            RELATIONSHIPS                                            │
└─────────────────────────────────────────────────────────────────────────────────────┘

    ┌─────────────┐         ┌─────────────┐         ┌─────────────┐
    │    users    │ 1     N │  kendaraan   │ 1     1 │  hasil_uji  │
    │             │─────────│             │─────────│             │
    │ id (PK)     │         │ id (PK)     │         │ id (PK)     │
    │             │         │ created_by  │         │ kendaraan_id│
    └──────┬──────┘         └─────────────┘         └─────────────┘
           │
           │
    ┌──────┴──────┐
    │             │
    │             │
    │  1        N │
    │─────────────│
    │             │
    │  1        N │
    │─────────────│
    │             │
    │  1        N │
    └─────────────┘
           │
           │
    ┌──────┴──────┐
    │  hasil_uji  │
    │             │
    │ tested_by   │
    └─────────────┘
```

### 3.2 Relationship Details

| Relationship | Type | Description |
|-------------|------|-------------|
| users → kendaraan | 1:N | Satu user dapat membuat banyak data kendaraan |
| users → hasil_uji | 1:N | Satu user dapat melakukan banyak uji emisi |
| users → home_content | 1:N | Satu user dapat membuat banyak konten home |
| kendaraan → hasil_uji | 1:1 | Satu kendaraan memiliki satu hasil uji |

## 4. Database Schema (SQL)

### 4.1 Users Table

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama VARCHAR(255) NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'supervisor', 'anggota')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 4.2 Kendaraan Table

```sql
CREATE TABLE kendaraan (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kategori VARCHAR(20) NOT NULL CHECK (kategori IN ('Dinas', 'Umum')),
    jenis VARCHAR(50) NOT NULL CHECK (jenis IN ('Bensin M', 'Bensin N&O', 'Solar JBB', 'Solar GVW')),
    plat_nomor VARCHAR(20) UNIQUE NOT NULL,
    merek VARCHAR(100) NOT NULL,
    tipe VARCHAR(100) NOT NULL,
    kapasitas_mesin INTEGER NOT NULL CHECK (kapasitas_mesin > 0),
    tahun_pembuatan INTEGER NOT NULL CHECK (tahun_pembuatan BETWEEN 1900 AND 2026),
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 4.3 Hasil Uji Table

```sql
CREATE TABLE hasil_uji (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kendaraan_id UUID UNIQUE REFERENCES kendaraan(id),
    co DECIMAL(10,4) CHECK (co >= 0),
    co2 DECIMAL(10,4) CHECK (co2 >= 0),
    hc DECIMAL(10,2) CHECK (hc >= 0),
    o2 DECIMAL(10,4) CHECK (o2 >= 0 AND o2 <= 100),
    lambda DECIMAL(10,4),
    opasitas DECIMAL(10,2) CHECK (opasitas >= 0 AND opasitas <= 100),
    hasil_uji VARCHAR(20) NOT NULL CHECK (hasil_uji IN ('Lulus', 'Tidak Lulus')),
    valid BOOLEAN NOT NULL,
    catatan TEXT,
    tested_by UUID REFERENCES users(id),
    tested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 4.4 Parameters Table

```sql
CREATE TABLE parameters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kategori VARCHAR(50) NOT NULL CHECK (kategori IN ('Bensin M', 'Bensin N&O', 'Solar JBB', 'Solar GVW')),
    tahun_min INTEGER NOT NULL,
    tahun_max INTEGER,
    tahun_operator VARCHAR(10) NOT NULL CHECK (tahun_operator IN ('<', 'between', '>')),
    co_max DECIMAL(10,4) CHECK (co_max >= 0),
    hc_max DECIMAL(10,2) CHECK (hc_max >= 0),
    opasitas_max DECIMAL(10,2) CHECK (opasitas_max >= 0),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 4.5 Popup Notes Table

```sql
CREATE TABLE popup_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    jenis_kendaraan VARCHAR(255) NOT NULL,
    tahun_operator VARCHAR(10) NOT NULL CHECK (tahun_operator IN ('>', '<', '>=', '<=', '=')),
    tahun_value INTEGER NOT NULL,
    parameter_uji VARCHAR(20) NOT NULL CHECK (parameter_uji IN ('CO', 'CO2', 'HC', 'O2', 'Lambda', 'Opasitas')),
    nilai_operator VARCHAR(10) NOT NULL CHECK (nilai_operator IN ('>', '<', '>=', '<=', '=')),
    nilai_value DECIMAL(10,4) NOT NULL,
    note TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 4.6 Home Content Table

```sql
CREATE TABLE home_content (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_type VARCHAR(50) NOT NULL CHECK (content_type IN ('description', 'image', 'working_instruction')),
    title VARCHAR(255),
    description TEXT,
    file_path VARCHAR(500),
    file_type VARCHAR(50),
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 5. Indexes

```sql
-- Performance indexes
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_kendaraan_plat ON kendaraan(plat_nomor);
CREATE INDEX idx_kendaraan_kategori ON kendaraan(kategori);
CREATE INDEX idx_kendaraan_jenis ON kendaraan(jenis);
CREATE INDEX idx_kendaraan_created_by ON kendaraan(created_by);
CREATE INDEX idx_hasil_uji_kendaraan ON hasil_uji(kendaraan_id);
CREATE INDEX idx_hasil_uji_hasil ON hasil_uji(hasil_uji);
CREATE INDEX idx_hasil_uji_tested_by ON hasil_uji(tested_by);
CREATE INDEX idx_parameters_kategori ON parameters(kategori);
CREATE INDEX idx_parameters_tahun ON parameters(tahun_min, tahun_max);
CREATE INDEX idx_popup_notes_parameter ON popup_notes(parameter_uji);
CREATE INDEX idx_home_content_type ON home_content(content_type);
```

## 6. Triggers

### 6.1 Updated At Trigger

```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_kendaraan_updated_at BEFORE UPDATE ON kendaraan
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_hasil_uji_updated_at BEFORE UPDATE ON hasil_uji
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_parameters_updated_at BEFORE UPDATE ON parameters
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_popup_notes_updated_at BEFORE UPDATE ON popup_notes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_home_content_updated_at BEFORE UPDATE ON home_content
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

## 7. Initial Data

### 7.1 Default Users

```sql
INSERT INTO users (nama, username, password, role) VALUES
('Superuser Admin', 'adminmas', '$2a$10$hashed_password_here', 'admin'),
('Kang Supervisor', 'supervisoraja', '$2a$10$hashed_password_here', 'supervisor'),
('Kroco01', 'kroco01', '$2a$10$hashed_password_here', 'anggota');
```

### 7.2 Default Parameters

```sql
-- Bensin M
INSERT INTO parameters (kategori, tahun_min, tahun_max, tahun_operator, co_max, hc_max, opasitas_max) VALUES
('Bensin M', 0, 2007, '<', 4.0, 1000.0, NULL),
('Bensin M', 2007, 2018, 'between', 1.0, 150.0, NULL),
('Bensin M', 2018, NULL, '>', 0.5, 100.0, NULL);

-- Bensin N&O
INSERT INTO parameters (kategori, tahun_min, tahun_max, tahun_operator, co_max, hc_max, opasitas_max) VALUES
('Bensin N&O', 0, 2007, '<', 4.0, 1000.0, NULL),
('Bensin N&O', 2007, 2018, 'between', 1.0, 150.0, NULL),
('Bensin N&O', 2018, NULL, '>', 0.5, 100.0, NULL);

-- Solar JBB
INSERT INTO parameters (kategori, tahun_min, tahun_max, tahun_operator, co_max, hc_max, opasitas_max) VALUES
('Solar JBB', 0, 2010, '<', NULL, NULL, 65.0),
('Solar JBB', 2010, 2021, 'between', NULL, NULL, 40.0),
('Solar JBB', 2021, NULL, '>', NULL, NULL, 30.0);

-- Solar GVW
INSERT INTO parameters (kategori, tahun_min, tahun_max, tahun_operator, co_max, hc_max, opasitas_max) VALUES
('Solar GVW', 0, 2010, '<', NULL, NULL, 65.0),
('Solar GVW', 2010, 2021, 'between', NULL, NULL, 40.0),
('Solar GVW', 2021, NULL, '>', NULL, NULL, 35.0);
```

---

**Document Version:** 1.0  
**Last Updated:** Agustus 2026
