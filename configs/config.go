package configs

import (
	"github.com/spf13/viper"
)

type JWTConfig struct {
	SecretKey        string `mapstructure:"secretkey"`
	ExpiresIn        string `mapstructure:"expiresin"`
	RefreshSecretKey string `mapstructure:"refreshsecretkey"`
	RefreshExpiresIn string `mapstructure:"refreshexpiresin"`
}

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	GeoIP    GeoIPConfig    `mapstructure:"geoip"`
}

type ServerConfig struct {
	BaseURL string `mapstructure:"baseurl"`
	Port    string `mapstructure:"port"`
	Env     string `mapstructure:"env"`
}

type DatabaseConfig struct {
	Host           string `mapstructure:"host"`
	Port           string `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	DBName         string `mapstructure:"dbname"`
	SSLMode        string `mapstructure:"sslmode"`
	ChannelBinding string `mapstructure:"channelbinding"`
}

type GeoIPConfig struct {
	DBPath string `mapstructure:"dbpath"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}
