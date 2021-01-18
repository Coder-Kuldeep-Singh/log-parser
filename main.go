package main

import (
	"log"
	"log-parser/controllers"
	"log-parser/routers"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file -> ", err)
		return
	}
	controllers.LoadGlobally()
	controllers.GetUniqueURLQueue()
	controllers.GetUniqueMethodQueue()
	controllers.GetHTTPCode()
	// controllers.GetTheErrorStatus()
	// update ip2location data every day
	go controllers.IP2Location()
	//setup routes
	r := routers.SetupRouter()
	// running
	r.Run(":" + os.Getenv("PORT"))
}
