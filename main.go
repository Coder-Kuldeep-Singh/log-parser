package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

// error levels
const (
	CRITICAL = 1
	INFO     = 2
	WARNING  = 3
)

// Logs handle the structure of the log files
type Logs struct {
	ip             string
	timestamp      string
	method         string
	url            string
	ProtocolMethod string
	ServerResponse string
	SendBytes      string
	userAgent      string
	// browser        string
	// system         string
}

// ErrorHandling handles the error and return formated error
func ErrorHandling(err error, msg string, code int) error {
	if err != nil {
		log.Printf("%d:%s", code, msg)
		return err
	}
	return nil
}

// OpenFile func opens a file and return back the file info as *os.File
func OpenFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// ReadFile reads the file
func ReadFile(outcome *os.File) chan Logs {
	rows := make(chan Logs)
	go func() {
		defer func() {
			close(rows)
		}()
		scanner := bufio.NewScanner(outcome)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			matched := regularExpression(`^(\S+) (\S+) (\S+) \[([\w:/]+\s[+\-]\d{4})\] "(\S+) (\S+)\s*(\S*)" (\d{3}) (\S+) (\S+) "(.*)"`, scanner.Text())
			if len(matched) == 0 {
				continue
			}
			rows <- Logs{
				ip:             matched[0][1],
				timestamp:      matched[0][4],
				method:         matched[0][5],
				url:            matched[0][6],
				ProtocolMethod: matched[0][7],

				ServerResponse: matched[0][8],
				SendBytes:      matched[0][9],
				userAgent:      matched[0][11],
				// browser        :,
				// system         :,
			}
		}
	}()
	return rows
}

// chanToSlice() converts chan []string to []string
func chanToSlice(records chan []string) []string {
	get := []string{}
	for record := range records {
		for _, found := range record {
			get = append(get, found)
		}
	}
	return get
}

// printrecords prints the chan data
func printrecords(records chan []string) {
	for record := range records {
		fmt.Println(record)
	}
}

// regularExpression runs the expression on the given data and return back the parse data
func regularExpression(MatchedCase string, resp string) [][]string {
	re := regexp.MustCompile(MatchedCase)
	if re.MatchString(resp) {
		return re.FindAllStringSubmatch(resp, -1)
	}
	return nil
}

func getIP(logRecord []string) chan []string {
	ips := make(chan []string)
	go func() {
		defer func() {
			close(ips)
		}()
		for _, record := range logRecord {
			expression := `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`
			matched := regularExpression(expression, record)
			if matched == nil {
				continue
			}
			for _, sliceOfIPS := range matched {
				ips <- []string{sliceOfIPS[0]}
			}
		}
	}()
	return ips
}

// \[\d{1,2}\/\w{3}\/\d{1,4}(:[0-9]{1,2}){3} \+([0-9]){1,4}\]

func getTimeStamp(logRecord []string) chan []string {
	stamp := make(chan []string)
	go func() {
		defer func() {
			// c, ok := <-stamp
			// if ok {
			close(stamp)
			// }
			// log.Println("channel still open")

		}()
		for _, record := range logRecord {
			expression := `\[\d{1,2}\/\w{3}\/\d{1,4}(:[0-9]{1,2}){3} \+([0-9]){1,4}\]`
			matched := regularExpression(expression, record)
			if matched == nil {
				continue
			}
			for _, sliceOfStamps := range matched {
				stamp <- []string{sliceOfStamps[0]}
			}
		}
	}()
	return stamp
}

// TopIPS return the top ips from logs
func TopIPS(queue map[string]int) {
	// top := []string{}
	for i := range queue {
		// if queue[i] > queue[i]+1 {
		// 	log.Println(queue[i])
		// }
		log.Printf("%d----->%d", queue[i], queue[i])
	}
}

func getUniqueQueue(record []string, length int) map[string]int {
	if length == 0 || length == 1 {
		return nil
	}
	queue := make(map[string]int)
	for idx := 0; idx < length-1; idx++ {
		if record[idx] != record[idx+1] {
			queue[record[idx]]++
		}
	}
	return queue
}

func main() {
	file, err := OpenFile("./access.log")
	err = ErrorHandling(err, "error to open file", WARNING)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("error to open the file %s", err.Error())
			return
		}
	}()
	rows := ReadFile(file)
	for i := range rows {
		log.Println(i)
	}
	// log.Println(rows)
	// parallel := 2
	// var wg sync.WaitGroup
	// wg.Add(parallel)
	// go func() {
	// 	defer wg.Done()
	// go func() {
	// ips := getIP(chanToSlice(rows))
	// converted := chanToSlice(ips)
	// unique := getUniqueQueue(converted, len(converted))
	// for key, value := range unique {
	// 	log.Printf("IP %s >>> %d\n", key, value)
	// }
	// TopIPS(unique)
	// log.Println(len(unique))
	// defer close(ips)
	// }()

	// go func() {
	// stamp := getTimeStamp(chanToSlice(rows))
	// // defer close(stamp)
	// printrecords(stamp)
	// }()
	// }()
	// wg.Wait()

}
