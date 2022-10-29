package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/routes"
)

func main() {
	router := gin.Default()

	// enable cors
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowAllOrigins = true

	router.Use(cors.New(corsConfig))

	// connect to database
	configs.ConnectDB()

	// routes
	routes.CourseRoute(router)
	routes.DegreeRoute(router)
	routes.ExamRoute(router)
	routes.ProfessorRoute(router)
	routes.SectionRoute(router)

	// @DEBUG
	// router.GET("/", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "Hello World!")
	// })

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Fatal(err)
	}
}
