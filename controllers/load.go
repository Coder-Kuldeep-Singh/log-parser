package controllers

import (
	"log"
	"log-parser/models"
	"log-parser/service"
	"strconv"
)

var (
	updateQueue        []models.Logs
	list               []string
	countries          map[string]int
	uniqueurlqueue     map[string]int
	methodQueue        map[string]int
	ipQueue            map[string]int
	httpQueue          map[string]int
	httpErrorCodeQueue map[string]int
)

// Nmaximum returns the Nth maximum numbers
func Nmaximum(values map[string]int, N int) map[string]int {
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
	return parsed
}

// UniqueIP eliminates the duplicate data and returns back the unique ips
func UniqueIP() []string {
	keys := make(map[string]bool)
	log.Println("Processing: Generating unique IP")
	for _, entry := range updateQueue {
		if _, value := keys[entry.IP]; !value {
			keys[entry.IP] = true
			list = append(list, entry.IP)
		}
	}
	log.Println("Processing: End")
	return list
}

// GetCountries eliminates the duplicate data and returns back the unique ips
func GetCountries(queue []service.Location) {
	log.Println("Processing: Generating unique Countries")
	for _, item := range queue {
		_, exist := countries[item.Country]

		if exist {
			countries[item.Country]++ // increase counter by 1 if already in the map
		} else {
			countries[item.Country] = 1 // else start counting from 1
		}
	}
	log.Println("Processing: End")
}

// GetUniqueURLQueue checks for the every url and returns back the map of url
// which holds the count of the each url how many times called
func GetUniqueURLQueue() {
	log.Printf("Processing: Generating Unique URL Counts Map")
	for _, item := range updateQueue {
		_, exist := uniqueurlqueue[item.URL]
		if exist {
			uniqueurlqueue[item.URL]++ // increase counter by 1 if already in the map
		} else {
			uniqueurlqueue[item.URL] = 1 // else start counting from 1
		}
	}
	log.Println("Processing: End")
}

// GetUniqueMethodQueue returns map
func GetUniqueMethodQueue() {
	log.Printf("Processing: Generating Unique Method Used Counts Map")
	for _, item := range updateQueue {
		_, exist := methodQueue[item.Method]
		if exist {
			methodQueue[item.Method]++ // increase counter by 1 if already in the map
		} else {
			methodQueue[item.Method] = 1 // else start counting from 1
		}
	}
	log.Println("Processing: End")
}

func getIPCounts() {
	log.Printf("Processing: Generating IP's Count")
	for _, item := range updateQueue {
		_, exist := ipQueue[item.IP]
		if exist {
			ipQueue[item.IP]++ // increase counter by 1 if already in the map
		} else {
			ipQueue[item.IP] = 1 // else start counting from 1
		}
	}
	log.Println("Processing: End")
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

func GetHTTPCode() {
	log.Printf("Processing: Generating HTTP Code Count")
	for _, item := range updateQueue {
		_, exist := httpQueue[item.ServerResponse]
		if exist {
			httpQueue[item.ServerResponse]++ // increase counter by 1 if already in the map
		} else {
			httpQueue[item.ServerResponse] = 1 // else start counting from 1
		}
	}
	log.Println("Processing: End")
}

func GetTheErrorStatus() {
	log.Printf("Processing: Generating HTTP Code Count")
	for _, i := range updateQueue {
		if i.ServerResponse == "" {
			continue
		}
		num, err := strconv.Atoi(i.ServerResponse)
		if err != nil {
			log.Printf("error to convert str to int %s", err.Error())
			continue
		}
		if num >= 299 {
			httpErrorCodeQueue[i.ServerResponse]++
		}
	}

	log.Println("Processing: End")
}

// func getDuplicate
