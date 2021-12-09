package middlewares

import (
	"net/http"
	"university-management-api/src/services"

	"github.com/gin-gonic/gin"
)

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "NoTokenProvided", "error": "No Authorization Header Provided"})
			c.Abort()
			return
		}

		claims, err := services.ValidateToken(clientToken)

		if err != "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "InvalidToken", "error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("firstName", claims.FirstName)
		c.Set("lastName", claims.LastName)
		c.Set("uid", claims.Uid)
		c.Set("role", claims.Role)
		c.Next()
	}
}
