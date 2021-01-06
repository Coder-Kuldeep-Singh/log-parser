package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log-parser/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetLogsFromClientSide handles upload the logs into server so user can see the logs's visualization
func GetLogsFromClientSide(c *gin.Context) {
	if c.Request.Body == nil {
		log.Printf("request body is missing {%d}\n", http.StatusBadRequest)
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "status_code": http.StatusBadRequest, "error": "response body is nil"})
		return
	}
	if !strings.Contains(c.Request.Header.Get("Content-type"), "application/json") {
		log.Printf("content type not accepted {%s}", c.Request.Header.Get("Content-type"))
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": "not valid content-type"})
		return
	}
	defer c.Request.Body.Close()
	logs, err := ResponseDecode(c.Request.Body)
	if err != nil {
		log.Printf("client request failed because json isn't valid {%v}", err)
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "status_code": http.StatusBadRequest, "error": err})
		return
	}
	err = models.UploadLogs(logs)
	if err != nil {
		log.Printf("client request failed error to upload data into database{%v}", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "status_code": http.StatusBadRequest, "error": "Internal Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "we're processing the logs you have uploaded"})

}

// ResponseDecode returns []models.Logs
func ResponseDecode(buf io.Reader) ([]models.Logs, error) {
	var logs []models.Logs
	err := json.NewDecoder(buf).Decode(&logs)
	if err != nil {
		return nil, fmt.Errorf("json decoding failed {stautsCode:%d} {Error:%v}", http.StatusBadRequest, err)
	}
	// defer c.Request.Body.Close()
	return logs, nil
}

// PostRequest sending post request
// func PostRequest(domain string, leads interface{}) (resp *http.Response, err error) {
// 	URL := fmt.Sprintf("https://%s/api/captureleads", domain)
// 	b := new(bytes.Buffer)
// 	err = json.NewEncoder(b).Encode(leads)
// 	if err != nil {
// 		return nil, fmt.Errorf("error to encode data into json {data:%v} {Error:%v}", leads, err)
// 	}
// 	resp, err = http.Post(URL, "application/json; charset=utf-8", b)
// 	if resp.StatusCode == http.StatusNotFound {
// 		return nil, fmt.Errorf("Targeted Server isn't running {url:%s} {statusCode:%s}", URL, resp.Status)
// 	}
// 	if err != nil {
// 		return nil, fmt.Errorf("error to post {data:%v} {Error:%v}", leads, err)
// 	}
// 	return resp, nil
// }
