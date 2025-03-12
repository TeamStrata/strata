package main

import (
	"github.com/TeamStrata/backend/api"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	// Cors setup
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{
		"*",
	}
	server.Use(cors.New(config))

	// Endpoints
	server.GET("/ping", api.PingHandler)

	server.Run(":8080")
}
