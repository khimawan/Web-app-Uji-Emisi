package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID        string         `json:"id" gorm:"type:uuid;primaryKey"`
	Nama      string         `json:"nama" gorm:"type:varchar(255);not null"`
	Username  string         `json:"username" gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string         `json:"password" gorm:"type:varchar(255);not null"`
	Role      string         `json:"role" gorm:"type:varchar(20);not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Kendaraan struct {
	ID             string         `json:"id" gorm:"type:uuid;primaryKey"`
	Kategori       string         `json:"kategori" gorm:"type:varchar(20);not null"`
	Jenis          string         `json:"jenis" gorm:"type:varchar(50);not null"`
	PlatNomor      string         `json:"plat_nomor" gorm:"type:varchar(20);uniqueIndex;not null"`
	Merek          string         `json:"merek" gorm:"type:varchar(100);not null"`
	Tipe           string         `json:"tipe" gorm:"type:varchar(100);not null"`
	KapasitasMesin int            `json:"kapasitas_mesin" gorm:"not null"`
	TahunPembuatan int            `json:"tahun_pembuatan" gorm:"not null"`
	CreatedBy      string         `json:"created_by" gorm:"type:uuid"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

type HasilUji struct {
	ID          string         `json:"id" gorm:"type:uuid;primaryKey"`
	KendaraanID string         `json:"kendaraan_id" gorm:"type:uuid;uniqueIndex"`
	CO          *float64       `json:"co" gorm:"type:decimal(10,4)"`
	CO2         *float64       `json:"co2" gorm:"type:decimal(10,4)"`
	HC          *float64       `json:"hc" gorm:"type:decimal(10,2)"`
	O2          *float64       `json:"o2" gorm:"type:decimal(10,4)"`
	Lambda      *float64       `json:"lambda" gorm:"type:decimal(10,4)"`
	Opasitas    *float64       `json:"opasitas" gorm:"type:decimal(10,2)"`
	HasilUji    string         `json:"hasil_uji" gorm:"type:varchar(20);not null"`
	Valid       bool           `json:"valid" gorm:"not null"`
	Catatan     *string        `json:"catatan" gorm:"type:text"`
	TestedBy    string         `json:"tested_by" gorm:"type:uuid"`
	TestedAt    time.Time      `json:"tested_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Parameter struct {
	ID            string         `json:"id" gorm:"type:uuid;primaryKey"`
	Kategori      string         `json:"kategori" gorm:"type:varchar(50);not null"`
	TahunMin      int            `json:"tahun_min" gorm:"not null"`
	TahunMax      *int           `json:"tahun_max"`
	TahunOperator string         `json:"tahun_operator" gorm:"type:varchar(10);not null"`
	COMax         *float64       `json:"co_max" gorm:"type:decimal(10,4)"`
	HCMax         *float64       `json:"hc_max" gorm:"type:decimal(10,2)"`
	OpasitasMax   *float64       `json:"opasitas_max" gorm:"type:decimal(10,2)"`
	IsActive      bool           `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type PopupNote struct {
	ID              string         `json:"id" gorm:"type:uuid;primaryKey"`
	JenisKendaraan  string         `json:"jenis_kendaraan" gorm:"type:varchar(255);not null"`
	TahunOperator   string         `json:"tahun_operator" gorm:"type:varchar(10);not null"`
	TahunValue      int            `json:"tahun_value" gorm:"not null"`
	ParameterUji    string         `json:"parameter_uji" gorm:"type:varchar(20);not null"`
	NilaiOperator   string         `json:"nilai_operator" gorm:"type:varchar(10);not null"`
	NilaiValue      float64        `json:"nilai_value" gorm:"type:decimal(10,4);not null"`
	Note            string         `json:"note" gorm:"type:text;not null"`
	IsActive        bool           `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type HomeContent struct {
	ID          string         `json:"id" gorm:"type:uuid;primaryKey"`
	ContentType string         `json:"content_type" gorm:"type:varchar(50);not null"`
	Title       *string        `json:"title" gorm:"type:varchar(255)"`
	Description *string        `json:"description" gorm:"type:text"`
	FilePath    *string        `json:"file_path" gorm:"type:varchar(500)"`
	FileType    *string        `json:"file_type" gorm:"type:varchar(50)"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedBy   string         `json:"created_by" gorm:"type:uuid"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
