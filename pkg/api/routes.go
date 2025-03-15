package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Example endpoint that returns "pong" in the response body JSON
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
