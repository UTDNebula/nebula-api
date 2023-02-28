package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/controllers"
)

func ExamRoute(router *gin.Engine) {
	// All routes related to exams come here
	examGroup := router.Group("/exam")

	examGroup.OPTIONS("", controllers.Preflight)
	examGroup.GET("", controllers.ExamSearch())
	examGroup.GET("all", controllers.ExamAll())
	examGroup.GET(":id", controllers.ExamById())
}
