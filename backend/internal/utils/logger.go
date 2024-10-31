package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

func SetupLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("%s %s %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
		c.Next()
	}
}
