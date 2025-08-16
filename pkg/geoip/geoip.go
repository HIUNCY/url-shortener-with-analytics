package geoip

import (
	"log"
	"net"

	"github.com/HIUNCY/url-shortener-with-analytics/configs"
	"github.com/oschwald/geoip2-golang"
)

// LocationData menampung hasil lookup.
type LocationData struct {
	Country string
	Region  string
	City    string
}

type GeoIPService interface {
	Lookup(ipAddress string) (*LocationData, error)
}

type geoIPService struct {
	db *geoip2.Reader
}

func NewGeoIPService(cfg configs.GeoIPConfig) GeoIPService {
	db, err := geoip2.Open(cfg.DBPath)
	if err != nil {
		log.Printf("WARNING: Could not open GeoIP database file: %v. Geolocation will be disabled.", err)
		return &geoIPService{db: nil}
	}
	return &geoIPService{db: db}
}

func (s *geoIPService) Lookup(ipAddress string) (*LocationData, error) {
	if s.db == nil {
		return &LocationData{}, nil // Return empty if DB is not loaded
	}

	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return &LocationData{}, nil
	}

	record, err := s.db.City(ip)
	if err != nil {
		return &LocationData{}, err
	}

	return &LocationData{
		Country: record.Country.IsoCode,
		Region:  getFirstSubdivision(record),
		City:    record.City.Names["en"],
	}, nil
}

func getFirstSubdivision(record *geoip2.City) string {
	if len(record.Subdivisions) > 0 {
		return record.Subdivisions[0].Names["en"]
	}
	return ""
}
