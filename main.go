package main

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"

	"github.com/botr99/gds/controllers"
	"github.com/botr99/gds/db"
)

var server *controllers.AdminServer

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/api/register", server.HandleRegister)
	r.GET("/api/commonstudents", server.HandleCommonStudents)
	r.POST("/api/suspend", server.HandleSuspend)
	r.POST("/api/retrievefornotifications", server.HandleRetrieveForNotifications)
	return r
}

func main() {
	err := godotenv.Load(".env");
	if err != nil {
		panic("Error loading .env file")
	}

	server = controllers.NewAdminServer(db.NewDbAdminService())
	
	r := SetupRouter()
	r.Run(":" + os.Getenv("PORT"))
}
