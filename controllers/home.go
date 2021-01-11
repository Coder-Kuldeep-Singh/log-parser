package controllers

import (
	"log"
	"log-parser/models"
	"log-parser/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	geo "github.com/oschwald/geoip2-golang"
)

func OpenLocationDB(dbPath string) (*geo.Reader, error) {
	db, err := geo.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// MainDashboard holds the combination of all the analysis on the same page
func MainDashboard(c *gin.Context) {
	f := []string{"./html/static/nginx/access.log", "./html/static/nginx/access.log.1"}
	queue := [][]models.Logs{}
	for _, path := range f {
		log.Printf("File Processing Start : [%s]", path)
		file, err := models.OpenFile(path)
		err = models.ErrorHandling(err, "error to open file", models.WARNING)
		if err != nil {
			log.Println(err)
			// return
			continue
		}
		queue = append(queue, models.ReadFile(file))
		log.Printf("File Processing End : [%s]", path)
	}
	updateQueue := UpdatedQueue(queue)
	// removed the all data from the  [][]Logs
	log.Println("Processing: Deleting the Old queue so we can use the Memory for other queues")
	queue = nil
	log.Println("Processing: End")

	ips := getIPCounts(updateQueue)
	Location := []service.Location{}
	db, err := OpenLocationDB("./db/GeoLite2-City.mmdb")
	err = models.ErrorHandling(err, "error to open the location database", models.WARNING)
	if err != nil {
		log.Println(err)
		return
	}
	for key := range ips {
		Location = append(Location, service.GetLocationFromIP(db, key))
	}
	db.Close()
	lastElement := Location[len(Location)-1]
	c.HTML(http.StatusOK, "dashboard.tmpl.html", gin.H{
		"TotalHits": len(updateQueue),
		// "VisitedIP": UniqueIP(updateQueue),
		"IPS":          Nmaximum(ips, 1),
		"Location":     Location[0 : len(Location)-2],
		"LastLocation": lastElement,
		"URLS":         Nmaximum(GetUniqueURLQueue(updateQueue), 1),
		"Methods":      Nmaximum(GetUniqueMethodQueue(updateQueue), 1),
		"Countries":    Nmaximum(GetCountries(Location), 1),
		"HTTPCODE":     Nmaximum(getHTTPCode(updateQueue), 1),
		"HTTPERROR":    Nmaximum(getTheErrorStatus(updateQueue), 1),
	})
}

// ReportIP generates the report of ip
func ReportIP(c *gin.Context) {
	f := []string{"./html/static/nginx/access.log", "./html/static/nginx/access.log.1", "./html/static/nginx/access.log.2"}
	queue := [][]models.Logs{}
	for _, path := range f {
		log.Printf("File Processing Start : [%s]", path)
		file, err := models.OpenFile(path)
		err = models.ErrorHandling(err, "error to open file", models.WARNING)
		if err != nil {
			log.Println(err)
			// return
			continue
		}
		queue = append(queue, models.ReadFile(file))
		log.Printf("File Processing End : [%s]", path)
	}
	updateQueue := UpdatedQueue(queue)
	// removed the all data from the  [][]Logs
	log.Println("Processing: Deleting the Old queue so we can use the Memory for other queues")
	queue = nil
	log.Println("Processing: End")
	c.HTML(http.StatusOK, "reportip.tmpl.html", gin.H{
		"IP":   getIPCounts(updateQueue),
		"Logs": updateQueue,
	})
}

// Nmaximum returns the Nth maximum numbers
func Nmaximum(values map[string]int, N int) map[string]int {
	// if len(values) < N {
	// 	N = len(values)
	// }
	parsed := make(map[string]int)
	for i := 0; i < N; i++ {
		max := 0
		kk := ""
		for key, value := range values {
			if value > max {
				max = value
				kk = key
			}
			delete(values, kk)
			parsed[kk] = max
		}
	}
	// log.Println(parsed)
	return parsed
}

// UniqueIP eliminates the duplicate data and returns back the unique ips
func UniqueIP(queue []models.Logs) []string {
	keys := make(map[string]bool)
	list := []string{}
	log.Println("Processing: Generating unique IP")
	for _, entry := range queue {
		if _, value := keys[entry.IP]; !value {
			keys[entry.IP] = true
			list = append(list, entry.IP)
		}
	}
	log.Println("Processing: End")
	return list
}

// GetCountries eliminates the duplicate data and returns back the unique ips
func GetCountries(queue []service.Location) map[string]int {
	list := make(map[string]int)
	log.Println("Processing: Generating unique Countries")
	for _, entry := range queue {
		if len(entry.Country) == 0 {
			continue
		}
		list[entry.Country]++
	}
	log.Println("Processing: End")
	return list
}

// GetUniqueURLQueue checks for the every url and returns back the map of url
// which holds the count of the each url how many times called
func GetUniqueURLQueue(record []models.Logs) map[string]int {
	queue := make(map[string]int)
	log.Printf("Processing: Generating Unique URL Counts Map")
	for _, i := range record {
		queue[i.URL]++
	}
	log.Println("Processing: End")
	return queue
}

// GetUniqueMethodQueue returns map
func GetUniqueMethodQueue(record []models.Logs) map[string]int {
	queue := make(map[string]int)
	log.Printf("Processing: Generating Unique Method Used Counts Map")
	for _, i := range record {
		queue[i.Method]++
	}
	log.Println("Processing: End")
	return queue
}

func getIPCounts(record []models.Logs) map[string]int {
	queue := make(map[string]int)
	log.Printf("Processing: Generating IP's Count")
	for _, i := range record {
		queue[i.IP]++
	}
	log.Println("Processing: End")
	return queue
}

// UpdatedQueue return back the []Logs
func UpdatedQueue(queue [][]models.Logs) []models.Logs {
	updateQueue := []models.Logs{}
	log.Printf("Processing: Updating Queue")
	for _, updated := range queue {
		for _, nextQueue := range updated {
			updateQueue = append(updateQueue, nextQueue)
		}
	}
	log.Println("Processing: End")
	return updateQueue
}

func getHTTPCode(record []models.Logs) map[string]int {
	queue := make(map[string]int)
	log.Printf("Processing: Generating HTTP Code Count")
	for _, i := range record {
		queue[i.ServerResponse]++
	}
	log.Println("Processing: End")
	return queue
}

func getTheErrorStatus(record []models.Logs) map[string]int {
	queue := make(map[string]int)
	log.Printf("Processing: Generating HTTP Code Count")
	for _, i := range record {
		if i.ServerResponse == "" {
			continue
		}
		num, err := strconv.Atoi(i.ServerResponse)
		if err != nil {
			log.Printf("error to convert str to int %s", err.Error())
			continue
		}
		if num >= 299 {
			queue[i.ServerResponse]++
		}
	}
	log.Println("Processing: End")
	return queue
}
