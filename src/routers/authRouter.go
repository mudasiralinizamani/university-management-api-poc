package routers

import (
	"university-management-api/src/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRouter(app *gin.Engine) {
	app.POST("/auth/signup", controllers.Signup())
	app.POST("/auth/signin", controllers.Signin())
}
