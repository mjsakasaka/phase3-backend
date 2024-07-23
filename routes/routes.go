package routes

import (
	"week1/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/users", controllers.GetUsers)
}
