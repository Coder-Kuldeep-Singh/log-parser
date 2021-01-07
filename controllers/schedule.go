package controllers

import (
	"log"
	"log-parser/models"
	"os"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
)

// IP2Location runs the service on specific time
func IP2Location() {
	schedule := gocron.NewScheduler(time.UTC)
	interval, err := strconv.ParseUint(os.Getenv("INTERVALS"), 10, 64)
	if err != nil {
		log.Printf("error to parse interval for ip2location %s", err.Error())
		return
	}
	_, err = schedule.Every(interval).Day().At(os.Getenv("AUTO_PING_IP2LOCATION")).Do(models.UploadDailyIP2Location)
	if err != nil {
		log.Printf("cron failed %s", err.Error())
		return
	}
	<-schedule.StartAsync()
	schedule.Stop()
}
