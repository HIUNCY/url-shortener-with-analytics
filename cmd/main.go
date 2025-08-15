package main

import (
	"fmt"
	"log"

	"github.com/HIUNCY/url-shortener-with-analytics/configs"
	"github.com/HIUNCY/url-shortener-with-analytics/pkg/database"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatalf("Tidak dapat memuat konfigurasi: %v", err)
	}
	fmt.Printf("Konfigurasi berhasil dimuat. Env: %s\n", config.Server.Env)

	db, err := database.NewPostgresConnection(&config.Database)
	if err != nil {
		return
	}

	// Untuk development, kita bisa cek koneksi dengan Ping
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Gagal mendapatkan objek DB: %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("Gagal ping ke database: %v", err)
	}

	fmt.Println("Ping ke database sukses!")
	// Server Gin akan dijalankan di sini nantinya...
}
