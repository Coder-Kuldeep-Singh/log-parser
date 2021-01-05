package service

import (
	"log"
	"log-parser/models"
	"net"

	geo "github.com/oschwald/geoip2-golang"
)

// Location holds the ip location where it belong
type Location struct {
	Latitude    float64 `json:"Latitude"`
	Longitude   float64 `json:"Longitude"`
	CountryCode string  `json:"country_code"`
	Country     string  `json:"country"`
	PostalCode  string  `json:"postal_code"`
	City        string  `json:"city"`
}

// GetLocationFromIP return the location of the ip address
func GetLocationFromIP(db *geo.Reader, ip string) Location {
	if len(ip) == 0 {
		log.Println("ip is blank")
		return Location{}
	}
	record, err := getCity(db, ip)
	err = models.ErrorHandling(err, "error to while parsing ip into city function", models.WARNING)
	if err != nil {
		log.Println(err)
		return Location{}
	}
	return DecodeDB(record)
}

// DecodeDB helps us to decode db data structure into json table
func DecodeDB(record *geo.City) Location {
	return Location{
		Latitude:    record.Location.Latitude,
		Longitude:   record.Location.Longitude,
		CountryCode: record.Country.IsoCode,
		Country:     record.Country.Names["en"],
		PostalCode:  record.Postal.Code,
		City:        record.City.Names["en"],
	}
}

func getCity(db *geo.Reader, strip string) (*geo.City, error) {
	ip := net.ParseIP(strip)
	record, err := db.City(ip)
	if err != nil {
		return nil, err
	}
	return record, nil
}
