package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

//// Middleware function to log request information
//func RouteMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// Log request information
//		fmt.Printf("Request: %s %s\n", r.Method, r.URL.Path)
//
//		// Call the next handler in the chain
//		next.ServeHTTP(w, r)
//	})
//}

func RouteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("C:", c)
		// Your middleware logic here
		c.Next() // Call the next middleware/handler
	}
}
