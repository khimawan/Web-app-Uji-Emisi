package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"web-app-uji-emisi/config"
	"web-app-uji-emisi/controllers"
	"web-app-uji-emisi/routes"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to database
	config.ConnectDatabase()

	// Seed initial data
	controllers.SeedInitialData()

	// Setup routes
	r := routes.SetupRoutes()

	server := &http.Server{
		Addr:    ":" + config.AppConfig.Port,
		Handler: r,
		// Tanpa timeout, satu koneksi yang menggantung (slow-loris) bisa
		// menahan satu goroutine dan sebagian kecil RAM selamanya.
		ReadHeaderTimeout: time.Duration(config.AppConfig.ReadTimeout) * time.Second,
		ReadTimeout:       time.Duration(config.AppConfig.ReadTimeout) * time.Second,
		// WriteTimeout diberi longgar karena export/import bisa memproses
		// ribuan baris dan butuh waktu lebih.
		WriteTimeout: time.Duration(config.AppConfig.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.AppConfig.IdleTimeout) * time.Second,
	}

	// Shutdown rapi: selesaikan request yang sedang berjalan lalu tutup dengan
	// halus. Tanpa ini, `docker compose restart` bisa memutus request aktif.
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals

		timeout := time.Duration(config.AppConfig.ShutdownTimeout) * time.Second
		log.Printf("Shutdown signal diterima, menunggu %s...", timeout)

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Shutdown paksa: %v", err)
		}
	}()

	log.Printf(
		"Server starting on port %s (rate limit %d/menit, maks %d request bersamaan)",
		config.AppConfig.Port,
		config.AppConfig.RateLimitPerMinute,
		config.AppConfig.MaxConcurrent,
	)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("Failed to start server:", err)
	}

	log.Println("Server stopped")
}
