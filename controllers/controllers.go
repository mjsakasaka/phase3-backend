package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	users := "bob"

	c.JSON(http.StatusOK, users)
}
