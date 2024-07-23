package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"week1/routes"
)

func main() {
	r := gin.Default()

	r.Static("/static", "./static")

	r.GET("/index", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./static/index.html")
	})

	routes.SetupRoutes(r)

	r.Run("0.0.0.0:8000")
}
