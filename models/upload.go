package models

import (
	"fmt"
	"log"
	"log-parser/config"
)

func insertStatement(logs []Logs) []string {
	queries := []string{}
	for _, keys := range logs {
		queries = append(queries, Query(keys))
	}
	return queries
}

// Query generates the query
func Query(logs Logs) string {
	return fmt.Sprintf(`
	INSERT INTO log(
		ip,
		visited_date,
		status_code,
		visited_url,
		protocol_status,
		server_response,
		send_bytes,
		user_agent
	) VALUES(
		'%s','%s','%s','%s','%s','%s','%s','%s'
	)
	`, logs.IP, logs.Timestamp, logs.Method, logs.URL, logs.ProtocolMethod, logs.ServerResponse, logs.SendBytes, logs.UserAgent)
}

// UploadLogs inserts the logs data into log table
func UploadLogs(logs []Logs) error {
	db, err := config.Connect()
	if err != nil {
		log.Printf("db connection failed %s", err)
		return err
	}
	defer db.Closed()
	statements := insertStatement(logs)
	for _, query := range statements {
		result, err := db.Exec(query)
		if err != nil {
			log.Printf("error to insert %s", err.Error())
			continue
		}
		log.Println(result.LastInsertId())
	}
	return nil
}
