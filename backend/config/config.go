package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	Port       string

	// Pembatas beban kerja.
	RateLimitPerMinute  int
	RateLimitBurst      int
	LoginRateLimit      int
	LoginRateLimitBurst int
	MaxConcurrent       int

	// Pool koneksi database. Tanpa batas, GORM membuka koneksi tanpa henti
	// dan menghabiskan seluruh kuota PostgreSQL.
	DBMaxOpenConns int
	DBMaxIdleConns int

	// Proxy yang dipercaya untuk dibaca X-Forwarded-For. Nilai selain ini
	// diabaikan, supaya klien tidak bisa memalsukan IP untuk lolos rate limit.
	TrustedProxies []string

	// Batas waktu server (detik).
	ReadTimeout       int
	WriteTimeout      int
	IdleTimeout       int
	ShutdownTimeout   int
}

var AppConfig Config

func LoadConfig() {
	godotenv.Load()

	AppConfig = Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "uji_emisi"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key"),
		Port:       getEnv("PORT", "8080"),

		RateLimitPerMinute:  getEnvInt("RATE_LIMIT_PER_MINUTE", 120),
		RateLimitBurst:      getEnvInt("RATE_LIMIT_BURST", 40),
		LoginRateLimit:      getEnvInt("LOGIN_RATE_LIMIT_PER_MINUTE", 10),
		LoginRateLimitBurst: getEnvInt("LOGIN_RATE_LIMIT_BURST", 5),
		MaxConcurrent:       getEnvInt("MAX_CONCURRENT_REQUESTS", 64),

		DBMaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 20),
		DBMaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 5),

		// Rentang private 172.16/12 dipakai Docker bridge; 127/8 untuk akses lokal.
		TrustedProxies: strings.Split(getEnv("TRUSTED_PROXIES", "172.16.0.0/12,127.0.0.1/8"), ","),

		ReadTimeout:     getEnvInt("SERVER_READ_TIMEOUT", 15),
		WriteTimeout:    getEnvInt("SERVER_WRITE_TIMEOUT", 120),
		IdleTimeout:     getEnvInt("SERVER_IDLE_TIMEOUT", 60),
		ShutdownTimeout: getEnvInt("SERVER_SHUTDOWN_TIMEOUT", 20),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			return parsed
		}
	}
	return defaultValue
}
