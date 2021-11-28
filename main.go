package main

import (
	"os"
	"university-management-api/src/routers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	app := gin.New()
	app.Use(gin.Logger())
	app.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))

	// Pipeline - Mudasir Ali
	routers.AuthRouter(app)

	app.Run(":" + port)
}
