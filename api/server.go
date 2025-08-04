package main

import (
	"log"

	"github.com/UTDNebula/nebula-api/api/configs"
	_ "github.com/UTDNebula/nebula-api/api/docs"
	"github.com/UTDNebula/nebula-api/api/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Unauthenticated placeholder endpoint for the built-in ginSwagger swagger documentation endpoint
//
//	@Id				swagger
//	@Param			file	path	string	true	"The swagger file to retrieve"
//	@Router			/swagger/{file} [get]
//	@Tags			Other
//	@Description	Returns the OpenAPI/swagger spec for the API
//	@Security
//	@Success	200
func swagger_controller_placeholder() {}

// @title						dev-nebula-api
// @description				The developer Nebula Labs API for access to pertinent UT Dallas data
// @version					1.1.0
// @host						api.utdnebula.com
// @schemes					https http
// @x-google-backend			{"address": "https://dev-nebula-api-1062216541483.us-south1.run.app"}
// @x-google-endpoints			[{"name": "dev-nebula-api-2wy9quu2ri5uq.apigateway.nebula-api-368223.cloud.goog", "allowCors": true}]
// @x-google-management		{"metrics": [{"name": "read-requests", "displayName": "Read Requests CUSTOM", "valueType": "INT64", "metricKind": "DELTA"}], "quota": {"limits": [{"name": "read-limit", "metric": "read-requests", "unit": "1/min/{project}", "values": {"STANDARD": 1000}}]}}
// @security					api_key
// @securitydefinitions.apikey	api_key
// @name						x-api-key
// @in							header

func main() {

	// To avoid unused error on swagger_controller_placeholder
	swagger_controller_placeholder()

	// Set up logging flags
	log.Default().SetFlags(log.Ltime | log.Llongfile)

	// Establish the connection to the database
	configs.ConnectDB()

	// Configure Gin Router
	router := gin.New()
	// Get rid of "trusted all proxies" warning -- we don't care
	router.SetTrustedProxies(nil)

	// Enable CORS
	router.Use(CORS)

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
	routes.RoomsRoute(router)
	routes.EventsRoute(router)
	routes.AstraRoute(router)
	routes.MazevoRoute(router)

	// Retrieve the port string to serve traffic on
	portString := configs.GetPortString()

	// Serve Traffic
	router.Run(portString)
}

func CORS(c *gin.Context) {
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

func LogRequest(c *gin.Context) {
	log.Printf("%s %s %s", c.Request.Method, c.Request.URL.Path, c.Request.Host)
	c.Next()
}
