package configs

import (
	"github.com/spf13/viper"
)

// Config menampung semua konfigurasi untuk aplikasi.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

// ServerConfig menampung semua konfigurasi untuk server.
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Env  string `mapstructure:"env"`
}

// DatabaseConfig menampung semua konfigurasi untuk database.
type DatabaseConfig struct {
	Host           string `mapstructure:"host"`
	Port           string `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	DBName         string `mapstructure:"dbname"`
	SSLMode        string `mapstructure:"sslmode"`
	ChannelBinding string `mapstructure:"channelbinding"`
}

// LoadConfig membaca konfigurasi dari file .env.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// Viper akan otomatis membaca env vars yang cocok
	viper.AutomaticEnv()

	// Membaca file .env
	err = viper.ReadInConfig()
	if err != nil {
		// Abaikan jika file tidak ditemukan, mungkin variabel diatur di sistem
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	// Unmarshal semua konfigurasi ke dalam struct Config
	err = viper.Unmarshal(&config)
	return
}
