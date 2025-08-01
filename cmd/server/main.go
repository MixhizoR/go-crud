package main

import (
	"github.com/MixhizoR/go-crud/config"
	"github.com/MixhizoR/go-crud/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	engine := gin.Default()
	routes.UserRoutes(engine)
	engine.Run(":8080")
}
