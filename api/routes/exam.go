package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/controllers"
)

func ExamRoute(router *gin.Engine) {
	// All routes related to exams come here
	examGroup := router.Group("/exam")

	examGroup.GET("/", controllers.ExamSearch())
	examGroup.GET("/:id", controllers.ExamById())
	examGroup.GET("/all", controllers.ExamAll())
}
