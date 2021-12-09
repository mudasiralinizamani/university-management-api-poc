package routers

import (
	"university-management-api/src/controllers"

	"github.com/gin-gonic/gin"
)

func FacultyRouter(app *gin.Engine) {
	app.POST("faculties/create", controllers.CreateFaculty())
	app.GET("faculties/getfaculties", controllers.GetFaculties())
	app.GET("faculties/getfaculty/:faculty_id", controllers.GetFaculty())
}
