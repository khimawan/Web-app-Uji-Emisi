-- Aplikasi Uji Emisi Database Schema

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(255) NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'supervisor', 'anggota')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Kendaraan table
CREATE TABLE kendaraan (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kategori VARCHAR(20) NOT NULL CHECK (kategori IN ('Dinas', 'Umum')),
    jenis VARCHAR(50) NOT NULL CHECK (jenis IN ('Bensin M', 'Bensin N&O', 'Solar JBB', 'Solar GVW')),
    plat_nomor VARCHAR(20) UNIQUE NOT NULL,
    merek VARCHAR(100) NOT NULL,
    tipe VARCHAR(100) NOT NULL,
    kapasitas_mesin INTEGER NOT NULL CHECK (kapasitas_mesin > 0),
    tahun_pembuatan INTEGER NOT NULL CHECK (tahun_pembuatan BETWEEN 1900 AND 2026),
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Hasil Uji table
CREATE TABLE hasil_uji (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
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
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Parameters table
CREATE TABLE parameters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kategori VARCHAR(50) NOT NULL CHECK (kategori IN ('Bensin M', 'Bensin N&O', 'Solar JBB', 'Solar GVW')),
    tahun_min INTEGER NOT NULL,
    tahun_max INTEGER,
    tahun_operator VARCHAR(10) NOT NULL CHECK (tahun_operator IN ('<', '>', 'between')),
    co_max DECIMAL(10,4) CHECK (co_max >= 0),
    hc_max DECIMAL(10,2) CHECK (hc_max >= 0),
    opasitas_max DECIMAL(10,2) CHECK (opasitas_max >= 0),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Popup Notes table
CREATE TABLE popup_notes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    jenis_kendaraan VARCHAR(255) NOT NULL,
    tahun_operator VARCHAR(10) NOT NULL CHECK (tahun_operator IN ('>', '<', '>=', '<=', '=')),
    tahun_value INTEGER NOT NULL,
    parameter_uji VARCHAR(20) NOT NULL CHECK (parameter_uji IN ('CO', 'CO2', 'HC', 'O2', 'Lambda', 'Opasitas')),
    nilai_operator VARCHAR(10) NOT NULL CHECK (nilai_operator IN ('>', '<', '>=', '<=', '=')),
    nilai_value DECIMAL(10,4) NOT NULL,
    note TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Home Content table
CREATE TABLE home_content (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content_type VARCHAR(50) NOT NULL CHECK (content_type IN ('description', 'image', 'working_instruction')),
    title VARCHAR(255),
    description TEXT,
    file_path VARCHAR(500),
    file_type VARCHAR(50),
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Indexes
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

-- Updated at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers
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

-- Initial data will be seeded by the application
