package routes

import (
	"github.com/MixhizoR/go-crud/controllers" // Assuming you need this import here
	"github.com/gin-gonic/gin"
)

func UserRoutes(engine *gin.Engine) {
	UserGroup := engine.Group("/users")
	UserGroup.GET("", controllers.GetUsers)
	UserGroup.GET("/:id", controllers.GetUser)
	UserGroup.POST("", controllers.CreateUser)
	UserGroup.PUT("/:id", controllers.UpdateUser)
	UserGroup.DELETE("/:id", controllers.DeleteUser)
}
