package routers

import (
	"university-management-api/src/controllers"

	"github.com/gin-gonic/gin"
)

func DepartmentRouter(app *gin.Engine) {
	app.POST("/departments/create", controllers.CreateDepartment())
}
