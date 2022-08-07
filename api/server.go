package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/routes"
)

func main() {
	router := gin.Default()

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

	http.ListenAndServe(":8080", router)
}
