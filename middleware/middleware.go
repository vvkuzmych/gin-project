package middleware

import (
	"github.com/gin-gonic/gin"
)

func RouteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		// Your middleware logic here
		c.Next() // Call the next middleware/handler
	}
}
