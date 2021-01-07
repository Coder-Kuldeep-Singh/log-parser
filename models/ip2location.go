package models

import (
	"fmt"
	"log"
	"log-parser/config"
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

// OpenLocationDB  opens the mmdb database
func OpenLocationDB(dbPath string) (*geo.Reader, error) {
	db, err := geo.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// GetLocationFromIP return the location of the ip address
func GetLocationFromIP(db *geo.Reader, ip string) Location {
	if len(ip) == 0 {
		log.Println("ip is blank")
		return Location{}
	}
	record, err := getCity(db, ip)
	err = ErrorHandling(err, "error to while parsing ip into city function", WARNING)
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

func getIP(db *config.SQL) ([]string, error) {
	query := `SELECT IP	FROM logs;`
	rows, err := db.Query(query)
	if err != nil {
		log.Println("error to get the rows of ip data from logs table")
		return nil, err
	}
	ipAddress := []string{}
	for rows.Next() {
		var ip string
		err := rows.Scan(&ip)
		if err != nil {
			log.Printf("error to scan ip row %s", err.Error())
			continue
		}
	}
	return ipAddress, nil
}

func ip2locationQuery(ip string, location Location) string {
	return fmt.Sprintf(`
	INSERT INTO ip2location(
		ip
		latitude
		longitude
		country_code
		country
		postalCode
		city
	) VALUES(
		'%s',%.5f,%.5f,'%s','%s','%s','%s'
	)
	`, ip, location.Latitude, location.Longitude, location.CountryCode, location.Country, location.PostalCode, location.City)
}

// UploadDailyIP2Location generates the location from the ip so we can avoid the file readings
// for this
func UploadDailyIP2Location() {
	db, err := config.Connect()
	if err != nil {
		log.Printf("db connection failed %s", err.Error())
		return
	}
	defer db.Closed()
	ip, err := getIP(db)
	if err != nil {
		log.Printf("error to get the ip's %s", err.Error())
		return
	}
	citDB, err := OpenLocationDB("./db/GeoLite2-City.mmdb")
	if err != nil {
		log.Printf("error to open the city database %s", err.Error())
		return
	}
	for _, next := range ip {
		location := GetLocationFromIP(citDB, next)
		query := ip2locationQuery(next, location)
		result, err := db.Exec(query)
		if err != nil {
			log.Printf("error to insert %s", err.Error())
			continue
		}
		log.Println(result.LastInsertId())
	}
	return
}
