package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingController(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "pong!",
	})
}
