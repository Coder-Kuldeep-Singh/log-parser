package routers

import (
	"log-parser/controllers"
	"os"

	"github.com/gin-gonic/gin"
)

//SetupRouter sets up routes
func SetupRouter() *gin.Engine {
	router := gin.Default()
	gin.SetMode(os.Getenv("GIN_MODE"))
	router.LoadHTMLGlob("./html/templates/*")
	router.Static("/home", "./html/static")
	router.GET("/analysis", controllers.MainDashboard)
	router.GET("/analysis/ip-report", controllers.ReportIP)
	router.POST("api/upload_logs", controllers.GetLogsFromClientSide)
	return router
}
