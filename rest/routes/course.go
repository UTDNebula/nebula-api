package routes

import (
	"github.com/UTDNebula/nebula-api/rest/controllers"
	"github.com/UTDNebula/nebula-api/shared"
	"github.com/gin-gonic/gin"
)

func CourseRoute(router *gin.Engine, repo shared.CourseRepo) {
	handler := controllers.NewCourseController(repo)

	// All routes related to courses come here
	courseGroup := router.Group("/course")

	courseGroup.OPTIONS("", controllers.Preflight)
	courseGroup.GET("", handler.CourseSearch)
	courseGroup.GET(":id", handler.CourseById)
	courseGroup.GET("all", handler.CourseAll)

	// Endpoint to get the list of sections of the queried courses
	courseGroup.GET("/sections", handler.CourseSectionSearch)
	courseGroup.GET("/:id/sections", handler.CourseSectionById)
	courseGroup.GET("/sections/trends", handler.TrendsCourseSectionSearch)

	// Endpoint to get grades for a course by its course id
	courseGroup.GET("/:id/grades", handler.GradesByCourseID)

	// Endpoint to get the list of professors of the queried courses
	courseGroup.GET("/professors", handler.CourseProfessorSearch)
	courseGroup.GET("/:id/professors", handler.CourseProfessorById)
}
