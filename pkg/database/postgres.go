package database

import (
	"fmt"
	"log"

	"github.com/HIUNCY/url-shortener-with-analytics/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnection(config *configs.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
		config.SSLMode,
	)
	if config.ChannelBinding != "" {
		dsn += " channel_binding=" + config.ChannelBinding
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
		return nil, err
	}

	log.Println("Koneksi ke database berhasil dibuat.")
	return db, nil
}
