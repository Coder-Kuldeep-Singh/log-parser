package controllers

import (
	"log"
	"log-parser/models"
	"log-parser/service"
	"strconv"
)

var (
	// LogSize holds the total size of the loaded log files
	LogSize int64
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
	f := []string{"./html/static/nginx/access.log", "./html/static/nginx/access.log.1", "./html/static/nginx/access.log.2", "./html/static/nginx/access.log.3", "./html/static/nginx/access.log.4", "./html/static/nginx/access.log.5", "./html/static/nginx/access.log.6", "./html/static/nginx/access.log.7", "./html/static/nginx/access.log.8", "./html/static/nginx/access.log.9", "./html/static/nginx/access.log.10", "./html/static/nginx/access.log.11", "./html/static/nginx/access.log.12"}
	queue := [][]models.Logs{}
	for _, path := range f {
		log.Printf("File Processing Start : [%s]", path)
		file, err := models.OpenFile(path)
		err = models.ErrorHandling(err, "error to open file", models.WARNING)
		if err != nil {
			log.Println(err)
			continue
		}
		Filestate, err := file.Stat()
		if err != nil {
			log.Printf("error to get the file stat %s", err.Error())
			continue
		}
		LogSize += Filestate.Size()
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

// NotFoundPages return the all page urls' whose status is not found(404)
func NotFoundPages(queue []models.Logs) map[string]int {
	URLQueue := make(map[string]int)
	for _, newQueue := range queue {
		if newQueue.URL == "" {
			continue
		}
		if newQueue.ServerResponse == "404" {
			_, exists := URLQueue[newQueue.URL]
			if exists {
				URLQueue[newQueue.URL]++
			} else {
				URLQueue[newQueue.URL] = 1
			}
		}

	}
	return URLQueue
}

// TopVisitedURL return the most visited url counts
func TopVisitedURL(queue []models.Logs) map[string]int {
	URLQueue := make(map[string]int)
	for _, newQueue := range queue {
		if newQueue.URL == "" {
			continue
		}
		_, exists := URLQueue[newQueue.URL]
		if exists {
			URLQueue[newQueue.URL]++
		} else {
			URLQueue[newQueue.URL] = 1
		}
	}
	return URLQueue
}

// GetTotalBytes returns the total bytes used to serve users
func GetTotalBytes(queue []models.Logs) int64 {
	var sum int64
	for _, n := range queue {
		sum += n.SendBytes
	}
	return sum
}

// ErrorCodeCounts returns the count of the all error responses
func ErrorCodeCounts(queue []models.Logs) int {
	Count := 0
	for _, newQueue := range queue {
		if newQueue.ServerResponse > "399" && newQueue.ServerResponse != "404" {
			Count++
		}
	}
	return Count
}

// Error404NotFound returns the count of the all error responses
func Error404NotFound(queue []models.Logs) int {
	Count := 0
	for _, newQueue := range queue {
		if newQueue.ServerResponse == "404" {
			Count++
		}
	}
	return Count
}
