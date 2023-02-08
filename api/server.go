package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	loadMongo()
	serveTraffic()
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Print("There was no .env file found!")
	}
}

func loadMongo() {
	configs.ConnectDB()
}

func serveTraffic() {
	router := gin.Default()
	router.Use(CORS())

	routes.CourseRoute(router)
	routes.DegreeRoute(router)
	routes.ExamRoute(router)
	routes.SectionRoute(router)
	routes.ProfessorRoute(router)

	port, exist := os.LookupEnv("PORT")
	if !exist {
		port = "8080"
	}

	var portString = fmt.Sprintf(":%s", port)

	// Can we use router.run() for this?
	err := http.ListenAndServe(portString, router)
	if err != nil {
		log.Fatal(err)
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")

		if c.Request.Method == "OPTIONS" {
			c.IndentedJSON(204, "")
			return
		}

		c.Next()
	}
}
