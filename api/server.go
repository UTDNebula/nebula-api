package main

import (
	"github.com/UTDNebula/nebula-api/api/common/log"
	"github.com/UTDNebula/nebula-api/api/configs"
	_ "github.com/UTDNebula/nebula-api/api/docs"
	"github.com/UTDNebula/nebula-api/api/routes"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Unauthenticated placeholder endpoint for the built-in ginSwagger swagger documentation endpoint
// @Id swagger
// @Router /swagger/index.html [get]
// @Description Returns the OpenAPI/swagger spec for the API
// @Produce text/html
// @Security
// @Success 200
func swagger_controller_placeholder() {}

// @title dev-nebula-api
// @description The developer version of the Nebula Labs API for testing purposes
// @version 0.1.0
// @schemes http https
// @x-google-backend {"address": "REDACTED"}
// @x-google-endpoints [{"name": "dev-nebula-api-2wy9quu2ri5uq.apigateway.nebula-api-368223.cloud.goog", "allowCors": true}]
// @x-google-management {"metrics": [{"name": "read-requests", "displayName": "Read Requests CUSTOM", "valueType": "INT64", "metricKind": "DELTA"}], "quota": {"limits": [{"name": "read-limit", "metric": "read-requests", "unit": "1/min/{project}", "values": {"STANDARD": 1000}}]}}
// @security api_key
// @securitydefinitions.apikey api_key
// @name x-api-key
// @in header
func main() {

	// To avoid unused error on swagger_controller_placeholder
	swagger_controller_placeholder()

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// Establish the connection to the database
	configs.ConnectDB()

	// Configure Gin Router
	router := gin.New()

	// Enable CORS
	router.Use(CORS())

	// Enable Logging
	router.Use(LogRequest)

	// Setup swagger-ui hosted
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Connect Routes
	routes.CourseRoute(router)
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
