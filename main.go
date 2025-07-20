package main

import (
	gin "github.com/gin-gonic/gin"
)

func main() {
	engine := gin.New()
	engine.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})
	engine.Run(":8080")
}
