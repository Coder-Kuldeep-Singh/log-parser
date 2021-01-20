package controllers

import (
	"log"
	"log-parser/models"
	"log-parser/service"
	"strconv"
)

var (
	// UpdateQueue holds the all data
	UpdateQueue []models.Logs
	// List going to hold the ips
	List []string
	// Uniqueurlqueue going to holds uniqueURL data
	Uniqueurlqueue map[string]int
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

// GetIPs returns back the ips and the counts
func GetIPs(queue []models.Logs) map[string]int {
	IPQueue := make(map[string]int)
	log.Println("Processing: Counting the ips")
	for _, item := range queue {
		_, exist := IPQueue[item.IP]
		if exist {
			IPQueue[item.IP]++ // increase counter by 1 if already in the map
		} else {
			IPQueue[item.IP] = 1 // else start counting from 1
		}
	}
	log.Println("Processing: End")
	return IPQueue

}

// GetMethods returns back the used method counts
func GetMethods(queue []models.Logs) map[string]int {
	MethodQueue := make(map[string]int)
	log.Println("Processing: Counting the used methods")
	for _, item := range queue {
		_, exist := MethodQueue[item.Method]
		if exist {
			MethodQueue[item.Method]++ // increase counter by 1 if already in the map
		} else {
			MethodQueue[item.Method] = 1 // else start counting from 1
		}
	}
	log.Println("Processing: End")
	return MethodQueue

}

// GetReferrer returns back the used method counts
func GetReferrer(queue []models.Logs) map[string]int {
	ReferrerQueue := make(map[string]int)
	log.Println("Processing: Counting the referrer url")
	for _, item := range queue {
		if item.ReferrerURL == "" || len(item.ReferrerURL) <= 3 {
			continue
		}
		_, exist := ReferrerQueue[item.ReferrerURL]
		if exist {
			ReferrerQueue[item.ReferrerURL]++ // increase counter by 1 if already in the map
		} else {
			ReferrerQueue[item.ReferrerURL] = 1 // else start counting from 1
		}
	}
	log.Println("Processing: End")
	return ReferrerQueue

}

// GetCountries creates the hashmap of countries
func GetCountries(queue []service.Location) map[string]int {
	countries := make(map[string]int)
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
	return countries
}

// type UniqueVisitors struct {

// }

// UniqueVisitorsByCity generates the unique cities data
func UniqueVisitorsByCity(queue []service.Location) map[string]int {
	visitors := make(map[string]int)
	for _, newQueue := range queue {
		if newQueue.City == "" {
			continue
		}
		_, exists := visitors[newQueue.City]
		if exists {
			visitors[newQueue.City]++
		} else {
			visitors[newQueue.City] = 1
		}
	}
	return visitors
}

// UniqueBots return the counts of the bots
func UniqueBots(queue []models.Logs) map[string]int {
	bots := make(map[string]int)
	for _, newQueue := range queue {
		if newQueue.Bots == "" || len(newQueue.Bots) == 2 {
			continue
		}
		_, exists := bots[newQueue.Bots]
		if exists {
			bots[newQueue.Bots]++
		} else {
			bots[newQueue.Bots] = 1
		}
	}
	return bots
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
