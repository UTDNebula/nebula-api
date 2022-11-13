package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/controllers"
)

func ExamRoute(router *gin.Engine) {
	// All routes related to exams come here
	router.OPTIONS("/exam", controllers.Preflight)
	examGroup := router.Group("/exam")

	examGroup.GET("", controllers.ExamSearch())
	examGroup.GET("all", controllers.ExamAll())
	examGroup.GET(":id", controllers.ExamById())
}
