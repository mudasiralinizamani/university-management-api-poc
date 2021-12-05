package routers

import (
	"university-management-api/src/controllers"

	"github.com/gin-gonic/gin"
)

func UserRouter(app *gin.Engine) {
	app.GET("users/getusers", controllers.GetUsers())
	app.GET("users/getuser/:user_id", controllers.GetUser())
	app.GET("users/getusersbyrole/:user_role", controllers.GetUsersByRole())
}
