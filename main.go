package main

import (
	"log"
	"log-parser/routers"
	"os"

	"github.com/joho/godotenv"
)

// // error levels
// const (
// 	CRITICAL = 1
// 	INFO     = 2
// 	WARNING  = 3
// )

// // Logs handle the structure of the log files
// type Logs struct {
// 	ip             string
// 	timestamp      string
// 	method         string
// 	url            string
// 	ProtocolMethod string
// 	ServerResponse string
// 	SendBytes      string
// 	userAgent      string
// 	// browser        string
// 	// system         string
// }

// // ErrorHandling handles the error and return formated error
// func ErrorHandling(err error, msg string, code int) error {
// 	if err != nil {
// 		log.Printf("%d:%s", code, msg)
// 		return err
// 	}
// 	return nil
// }

// // OpenFile func opens a file and return back the file info as *os.File
// func OpenFile(path string) (*os.File, error) {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return file, nil
// }

// // ReadFile reads the file
// func ReadFile(outcome *os.File) []Logs {
// 	defer func() {
// 		err := outcome.Close()
// 		if err != nil {
// 			log.Printf("error to open the file %s", err.Error())
// 			return
// 		}
// 	}()
// 	rows := []Logs{}
// 	scanner := bufio.NewScanner(outcome)
// 	scanner.Split(bufio.ScanLines)
// 	for scanner.Scan() {
// 		matched := regularExpression(`^(\S+) (\S+) (\S+) \[([\w:/]+\s[+\-]\d{4})\] "(\S+) (\S+)\s*(\S*)" (\d{3}) (\S+) (\S+) "(.*)"`, scanner.Text())
// 		if len(matched) == 0 {
// 			continue
// 		}
// 		rows = append(rows, Logs{
// 			ip:             matched[0][1],
// 			timestamp:      matched[0][4],
// 			method:         matched[0][5],
// 			url:            matched[0][6],
// 			ProtocolMethod: matched[0][7],

// 			ServerResponse: matched[0][8],
// 			SendBytes:      matched[0][9],
// 			userAgent:      matched[0][11],
// 			// browser        :,
// 			// system         :,
// 		})
// 	}
// 	return rows
// }

// // chanToSlice() converts chan []string to []string
// func chanToSlice(records chan []string) []string {
// 	get := []string{}
// 	for record := range records {
// 		for _, found := range record {
// 			get = append(get, found)
// 		}
// 	}
// 	return get
// }

// // printrecords prints the chan data
// func printrecords(records chan []string) {
// 	for record := range records {
// 		fmt.Println(record)
// 	}
// }

// // regularExpression runs the expression on the given data and return back the parse data
// func regularExpression(MatchedCase string, resp string) [][]string {
// 	re := regexp.MustCompile(MatchedCase)
// 	if re.MatchString(resp) {
// 		return re.FindAllStringSubmatch(resp, -1)
// 	}
// 	return nil
// }

// // uniqueIP eliminates the duplicate data and returns back the unique ips
// func uniqueIP(queue []Logs) []string {
// 	keys := make(map[string]bool)
// 	list := []string{}
// 	log.Println("Processing: Generating unique IP")
// 	for _, entry := range queue {
// 		if _, value := keys[entry.ip]; !value {
// 			keys[entry.ip] = true
// 			list = append(list, entry.ip)
// 		}
// 	}
// 	log.Println("Processing: End")
// 	return list
// }

// // getUniqueURLQueue checks for the every url and returns back the map of url
// // which holds the count of the each url how many times called
// func getUniqueURLQueue(record []Logs) map[string]int {
// 	queue := make(map[string]int)
// 	log.Printf("Processing: Generating Unique URL Counts Map")
// 	for _, i := range record {
// 		queue[i.url]++
// 	}
// 	log.Println("Processing: End")
// 	return queue
// }

// func getUniqueMethodQueue(record []Logs) map[string]int {
// 	queue := make(map[string]int)
// 	log.Printf("Processing: Generating Unique Method Used Counts Map")
// 	for _, i := range record {
// 		queue[i.method]++
// 	}
// 	log.Println("Processing: End")
// 	return queue
// }

// // UpdatedQueue return back the []Logs
// func UpdatedQueue(queue [][]Logs) []Logs {
// 	updateQueue := []Logs{}
// 	log.Printf("Processing: Updating Queue")
// 	for _, updated := range queue {
// 		for _, nextQueue := range updated {
// 			updateQueue = append(updateQueue, nextQueue)
// 		}
// 	}
// 	log.Println("Processing: End")
// 	return updateQueue
// }

func main() {
	// files := flag.String("f", "", "location of the nginx log file comma seperated filenames")
	// flag.Parse()
	// if *files == "" {
	// 	fmt.Println("Usage: ./logs -f file1,[file2,...]")
	// 	flag.Usage()
	// 	return
	// }
	// f := strings.Split(*files, ",")
	// queue := [][]Logs{}
	// for _, path := range f {
	// 	log.Printf("File Processing Start : [%s]", path)
	// 	file, err := OpenFile(path)
	// 	err = ErrorHandling(err, "error to open file", WARNING)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// 	queue = append(queue, ReadFile(file))
	// 	log.Printf("File Processing End : [%s]", path)
	// }
	// updateQueue := UpdatedQueue(queue)
	// // removed the all data from the  [][]Logs
	// log.Println("Processing: Deleting the Old queue so we can use the Memory for other queues")
	// queue = nil
	// log.Println("Processing: End")
	// // CountURLQueue := getUniqueURLQueue(updateQueue)
	// _ = getUniqueURLQueue(updateQueue)
	// // for key, value := range CountURLQueue {
	// // 	log.Printf("URL %s >>> %d\n", key, value)
	// // }
	// log.Println(len(uniqueIP(updateQueue)))
	// log.Println(getUniqueMethodQueue(updateQueue))
	// log.Printf("Total Hits >>> %d", len(updateQueue))
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file -> ", err)
		return
	}
	//setup routes
	r := routers.SetupRouter()
	// running
	r.Run(":" + os.Getenv("PORT"))
}
