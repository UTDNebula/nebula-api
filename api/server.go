package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/routes"
)

func main() {
	router := gin.Default()

	// enable cors
	/*router.Use(cors.New(cors.Config{
	    AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT"},
	    AllowHeaders:     []string{"*"},
	    ExposeHeaders:    []string{"Content-Length"},
	    AllowCredentials: false,
	    AllowAllOrigins:  true,
	    AllowOriginFunc: func(origin string) bool { return true },
	    MaxAge:          86400,
	}))*/

	router.Use(CORS())

	// connect to database
	configs.ConnectDB()

	// routes
	routes.CourseRoute(router)
	routes.DegreeRoute(router)
	routes.ExamRoute(router)
	routes.SectionRoute(router)
	routes.ProfessorRoute(router)
	routes.GradesRoute(router)

	//router.OPTIONS("*", CORSOptionsHandler())

	// @DEBUG
	// router.GET("/", func(c *gin.Context) {
	//     c.String(http.StatusOK, "Hello World!")
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

func CORSOptionsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatus(200)
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
