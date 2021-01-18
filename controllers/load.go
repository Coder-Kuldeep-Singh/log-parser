package controllers

import (
	"log"
	"log-parser/models"
	"strconv"
)

var (
	// UpdateQueue holds the all data
	UpdateQueue []models.Logs
	// List going to hold the ips
	List []string
	// Countries going to holds countries data
	Countries map[string]int
	// Uniqueurlqueue going to holds uniqueURL data
	Uniqueurlqueue map[string]int
	// MethodQueue going to holds Methods data
	MethodQueue map[string]int
	// IPQueue going to holds ips data
	IPQueue map[string]int
	// HTTPErrorCodeQueue going to holds http error code data
	HTTPErrorCodeQueue map[string]int
)

// LoadGlobally loads the data globally
func LoadGlobally() {
	f := []string{"./html/static/nginx/access.log", "./html/static/nginx/access.log.1", "./html/static/nginx/access.log.2", "./html/static/nginx/access.log.3", "./html/static/nginx/access.log.4", "./html/static/nginx/access.log.5"}
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
	UpdatedQueue(queue)
}

// UpdatedQueue return back the []Logs
func UpdatedQueue(queue [][]models.Logs) {
	log.Printf("Processing: Updating Queue")
	for _, updated := range queue {
		for _, nextQueue := range updated {
			UpdateQueue = append(UpdateQueue, nextQueue)
		}
	}
	log.Println("Processing: End")
}

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
		}
		delete(values, kk)
		parsed[kk] = max
	}
	return parsed
}

// RequestsCode counts the requests
func RequestsCode(queue []models.Logs) map[string]int {
	HTTPQueue := make(map[string]int)
	log.Println("Processing: Counting the requests status code")
	for _, item := range queue {
		_, exist := HTTPQueue[item.ServerResponse]
		if exist {
			HTTPQueue[item.ServerResponse]++ // increase counter by 1 if already in the map
		} else {
			HTTPQueue[item.ServerResponse] = 1 // else start counting from 1
		}
	}
	log.Println("Processing: End")
	return HTTPQueue
}

// GetTheErrorStatus get the error status codes
func GetTheErrorStatus(queue []models.Logs) map[string]int {
	httpErrorCodeQueue := make(map[string]int)
	log.Printf("Processing: Generating HTTP Code Count")
	for _, i := range queue {
		if i.ServerResponse == "" {
			continue
		}
		num, err := strconv.Atoi(i.ServerResponse)
		if err != nil {
			log.Printf("error to convert str to int %s", err.Error())
			continue
		}
		if num >= 399 {
			httpErrorCodeQueue[i.ServerResponse]++
		}
	}
	log.Println("Processing: End")
	return httpErrorCodeQueue
}

// // UniqueIP eliminates the duplicate data and returns back the unique ips
// func UniqueIP() []string {
// 	keys := make(map[string]bool)
// 	log.Println("Processing: Generating unique IP")
// 	for _, entry := range updateQueue {
// 		if _, value := keys[entry.IP]; !value {
// 			keys[entry.IP] = true
// 			list = append(list, entry.IP)
// 		}
// 	}
// 	log.Println("Processing: End")
// 	return list
// }

// // GetCountries eliminates the duplicate data and returns back the unique ips
// func GetCountries(queue []service.Location) {
// 	log.Println("Processing: Generating unique Countries")
// 	for _, item := range queue {
// 		_, exist := countries[item.Country]

// 		if exist {
// 			countries[item.Country]++ // increase counter by 1 if already in the map
// 		}
// 		// else {
// 		// 	countries[item.Country] = 1 // else start counting from 1
// 		// }
// 	}
// 	log.Println("Processing: End")
// }

// // GetUniqueURLQueue checks for the every url and returns back the map of url
// // which holds the count of the each url how many times called
// func GetUniqueURLQueue() {
// 	log.Printf("Processing: Generating Unique URL Counts Map")
// 	for _, item := range updateQueue {
// 		_, exist := uniqueurlqueue[item.URL]
// 		if exist {
// 			uniqueurlqueue[item.URL]++ // increase counter by 1 if already in the map
// 			// continue
// 		}
// 		// uniqueurlqueue[item.URL] = 1 // else start counting from 1
// 	}
// 	log.Println("Processing: End")
// }

// // GetUniqueMethodQueue returns map
// func GetUniqueMethodQueue() {
// 	log.Printf("Processing: Generating Unique Method Used Counts Map")
// 	for _, item := range updateQueue {
// 		_, exist := methodQueue[item.Method]
// 		if exist {
// 			methodQueue[item.Method]++ // increase counter by 1 if already in the map
// 		}
// 		// else {
// 		// 	methodQueue[item.Method] = 1 // else start counting from 1
// 		// }
// 	}
// 	log.Println("Processing: End")
// }

// func getIPCounts() {
// 	log.Printf("Processing: Generating IP's Count")
// 	for _, item := range updateQueue {
// 		_, exist := ipQueue[item.IP]
// 		if exist {
// 			ipQueue[item.IP]++ // increase counter by 1 if already in the map
// 		}
// 		// else {
// 		// 	ipQueue[item.IP] = 1 // else start counting from 1
// 		// }
// 	}
// 	log.Println("Processing: End")
// }

// func GetHTTPCode() {
// 	log.Printf("Processing: Generating HTTP Code Count")
// 	for _, item := range updateQueue {
// 		_, exist := httpQueue[item.ServerResponse]
// 		if exist {
// 			httpQueue[item.ServerResponse]++ // increase counter by 1 if already in the map
// 		}
// 		// else {
// 		// 	httpQueue[item.ServerResponse] = 1 // else start counting from 1
// 		// }
// 	}
// 	log.Println("Processing: End")
// }

// func GetTheErrorStatus() {
// 	log.Printf("Processing: Generating HTTP Code Count")
// 	for _, i := range updateQueue {
// 		if i.ServerResponse == "" {
// 			continue
// 		}
// 		num, err := strconv.Atoi(i.ServerResponse)
// 		if err != nil {
// 			log.Printf("error to convert str to int %s", err.Error())
// 			continue
// 		}
// 		if num >= 299 {
// 			httpErrorCodeQueue[i.ServerResponse]++
// 		}
// 	}

// 	log.Println("Processing: End")
// }

// func getDuplicate
