package config

import (
	"fmt"
	"log"
	"time"

	"web-app-uji-emisi/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		AppConfig.DBHost,
		AppConfig.DBPort,
		AppConfig.DBUser,
		AppConfig.DBPassword,
		AppConfig.DBName,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Kendaraan{},
		&models.HasilUji{},
		&models.Parameter{},
		&models.PopupNote{},
		&models.HomeContent{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Batasi pool koneksi. Tanpa ini, tiap request yang bersamaan akan membuka
	// koneksi baru ke PostgreSQL sampai kuota max_connections habis dan database
	// mulai menolak koneksi.
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get underlying SQL connection pool:", err)
	}
	sqlDB.SetMaxOpenConns(AppConfig.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(AppConfig.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	log.Printf(
		"Database connected and migrated successfully (max open conns: %d, max idle: %d)",
		AppConfig.DBMaxOpenConns,
		AppConfig.DBMaxIdleConns,
	)
}
