package controllers

import (
	"fmt"
	"log"
	"log-parser/models"
	"log-parser/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	geo "github.com/oschwald/geoip2-golang"
)

// OpenLocationDB open the location database
func OpenLocationDB(dbPath string) (*geo.Reader, error) {
	db, err := geo.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// MainDashboard holds the combination of all the analysis on the same page
// func MainDashboard(c *gin.Context) {
// 	Location := []service.Location{}
// 	db, err := OpenLocationDB("./db/GeoLite2-City.mmdb")
// 	err = models.ErrorHandling(err, "error to open the location database", models.WARNING)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	for key := range ipQueue {
// 		Location = append(Location, service.GetLocationFromIP(db, key))
// 	}
// 	db.Close()
// 	lastElement := Location[len(Location)-1]
// 	c.HTML(http.StatusOK, "dashboard.tmpl.html", gin.H{
// 		"TotalHits": len(updateQueue),
// 		// "VisitedIP": UniqueIP(updateQueue),
// 		"IPS":          ipQueue,
// 		"Location":     Location[0 : len(Location)-2],
// 		"LastLocation": lastElement,
// 		"URLS":         uniqueurlqueue,
// 		"Methods":      methodQueue,
// 		// "Countries":    GetCountries(Location),
// 		// "HTTPCODE":     getHTTPCode(updateQueue),
// 		// "HTTPERROR":    getTheErrorStatus(updateQueue),
// 	})
// }

// MainDashboard holds the combination of all the analysis on the same page
func MainDashboard(c *gin.Context) {
	// log.Println(UpdateQueue)
	code := RequestsCode(UpdateQueue)
	errorCode := GetTheErrorStatus(UpdateQueue)
	topIps := GetIPs(UpdateQueue)
	methods := GetMethods(UpdateQueue)
	referrer := GetReferrer(UpdateQueue)
	Location := []service.Location{}
	db, err := OpenLocationDB("./db/GeoLite2-City.mmdb")
	err = models.ErrorHandling(err, "error to open the location database", models.WARNING)
	if err != nil {
		log.Println(err)
		return
	}
	for key := range topIps {
		Location = append(Location, service.GetLocationFromIP(db, key))
	}

	db.Close()
	Counties := GetCountries(Location)
	uniqueVisitors := UniqueVisitorsByCity(Location)
	bots := UniqueBots(UpdateQueue)
	NotFound := NotFoundPages(UpdateQueue)
	topURL := TopVisitedURL(UpdateQueue)
	lastElement := Location[len(Location)-1]
	// log.Println(UpdateQueue[0].Timestamp, UpdateQueue[len(UpdateQueue)-1].Timestamp)
	c.HTML(http.StatusOK, "dashboard.tmpl.html", gin.H{
		"TotalHits":        len(UpdateQueue),
		"LogSize":          fmt.Sprintf("%.2f", float64(LogSize)/float64(1000000)),
		"TotalBytesServed": fmt.Sprintf("%.2f", float64(GetTotalBytes(UpdateQueue))/float64(1000000000)),
		"UniqueVisitors":   len(topIps),
		"FailedRequests":   ErrorCodeCounts(UpdateQueue),
		"NotFound":         Error404NotFound(UpdateQueue),
		"ValidRequests":    len(UpdateQueue) - ErrorCodeCounts(UpdateQueue) - Error404NotFound(UpdateQueue),
		"HTTPCode":         Nmaximum(code, 5),
		"HTTPError":        Nmaximum(errorCode, 5),
		"TopIPS":           Nmaximum(topIps, 10),
		"Methods":          Nmaximum(methods, 3),
		"Referrer":         Nmaximum(referrer, 10),
		"Location":         Location[0 : len(Location)-2],
		"LastLocation":     lastElement,
		"Countries":        Nmaximum(Counties, 10),
		"VisitorsByCity":   Nmaximum(uniqueVisitors, 10),
		"Bots":             Nmaximum(bots, 5),
		"NotFoundURL":      Nmaximum(NotFound, 5),
		"TopURL":           Nmaximum(topURL, 5),
		"ReferrersUnique":  len(referrer),
		"LastUpdated":      time.Now().Format("2006-01-02 15:04:05.000000000"),
		// "HTTPCODE":     getHTTPCode(updateQueue),
		// "HTTPERROR":    getTheErrorStatus(updateQueue),
	})
}

// GetValidLength returns back the valid string
func GetValidLength(data int) string {
	stringValid := strconv.Itoa(data)
	if data < 10000 {
		return fmt.Sprintf("%sk", stringValid[0:1])
	}
	if data < 100000 {
		return fmt.Sprintf("%sk", stringValid[0:2])
	}
	if data < 1000000 {
		return fmt.Sprintf("%sk", stringValid[0:3])
	}
	if data < 10000000 {
		return fmt.Sprintf("%sm %sk", stringValid[0:1], stringValid[1:3])
	}
	return ""
}

// ReportIP generates the report of ip
// func ReportIP(c *gin.Context) {
// 	c.HTML(http.StatusOK, "reportip.tmpl.html", gin.H{
// 		"IP":   getIPCounts(updateQueue),
// 		"Logs": updateQueue,
// 	})
// }
