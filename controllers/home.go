package controllers

import (
	"log"
	"log-parser/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// MainDashboard holds the combination of all the analysis on the same page
func MainDashboard(c *gin.Context) {
	f := []string{"./html/static/nginx/access.log"}
	queue := [][]models.Logs{}
	for _, path := range f {
		log.Printf("File Processing Start : [%s]", path)
		file, err := models.OpenFile(path)
		err = models.ErrorHandling(err, "error to open file", models.WARNING)
		if err != nil {
			log.Println(err)
			return
		}
		queue = append(queue, models.ReadFile(file))
		log.Printf("File Processing End : [%s]", path)
	}
	updateQueue := UpdatedQueue(queue)
	// removed the all data from the  [][]Logs
	log.Println("Processing: Deleting the Old queue so we can use the Memory for other queues")
	queue = nil
	log.Println("Processing: End")
	// CountURLQueue := getUniqueURLQueue(updateQueue)
	// _ = GetUniqueURLQueue(updateQueue)
	// for key, value := range CountURLQueue {
	// 	log.Printf("URL %s >>> %d\n", key, value)
	// }
	// log.Println(len(UniqueIP(updateQueue)))
	// log.Println(GetUniqueMethodQueue(updateQueue))
	// log.Printf("Total Hits >>> %d", len(updateQueue))
	c.HTML(http.StatusOK, "dashboard.tmpl.html", gin.H{
		"TotalHits": len(updateQueue),
		// "VisitedIP": UniqueIP(updateQueue),
	})
}

// ReportIP generates the report of ip
func ReportIP(c *gin.Context) {
	f := []string{"./html/static/nginx/access.log"}
	queue := [][]models.Logs{}
	for _, path := range f {
		log.Printf("File Processing Start : [%s]", path)
		file, err := models.OpenFile(path)
		err = models.ErrorHandling(err, "error to open file", models.WARNING)
		if err != nil {
			log.Println(err)
			return
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
