package main

import (
	"github.com/UTDNebula/nebula-api/api/common/log"
	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	// Establish the connection to the database
	configs.ConnectDB()

	// Configure Gin Router
	router := gin.New()

	// Enable CORS
	router.Use(CORS())

	// Enable Logging
	router.Use(Log())

	// Connect Routes
	routes.CourseRoute(router)
	routes.DegreeRoute(router)
	routes.ExamRoute(router)
	routes.SectionRoute(router)
	routes.ProfessorRoute(router)

	// Retrieve the port string to serve traffic on
	portString := configs.GetPortString()

	// Serve Traffic
	router.Run(portString)

}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, x-api-key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")

		if c.Request.Method == "OPTIONS" {
			c.IndentedJSON(204, "")
			return
		}

		c.Next()
	}
}

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info.Printf("%s %s from %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.Host,
		)
		c.Next()
	}
}
