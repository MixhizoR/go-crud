package main

import (
	gin "github.com/gin-gonic/gin"
)

func helloWorld(c *gin.Context) {
	c.String(200, "Hello, World!")
}

func helloUsers(c *gin.Context) {
	c.String(200, "Hello, Users!")
}

func main() {
	engine := gin.New()
	engine.GET("/", helloWorld)
	engine.GET("/users", helloUsers)
	engine.Run(":8080")
}
