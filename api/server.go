package main

import (
	"github.com/UTDNebula/nebula-api/api/common/log"
	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/routes"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// Establish the connection to the database
	configs.ConnectDB()

	// Configure Gin Router
	router := gin.New()

	// Enable CORS
	router.Use(CORS())

	// Enable Logging
	router.Use(LogRequest)

	// Connect Routes
	routes.CourseRoute(router)
	routes.DegreeRoute(router)
	routes.ExamRoute(router)
	routes.SectionRoute(router)
	routes.ProfessorRoute(router)
	routes.GradesRoute(router)
	routes.AutocompleteRoute(router)
	routes.StorageRoute(router)

	// Retrieve the port string to serve traffic on
	portString := configs.GetPortString()

	// Serve Traffic
	router.Run(portString)
	log.Logger.Debug().Str("port", portString).Msg("Listening to port")
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

func LogRequest(c *gin.Context) {
	log.Logger.Info().
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Str("host", c.Request.Host).
		Send()

	c.Next()
}
