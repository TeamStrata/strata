package main

import (
	"context"
	"log"
	"os"

	"github.com/TeamStrata/strata/pkg/api"
	"github.com/TeamStrata/strata/pkg/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	// Get database connection string
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %s", err.Error())
	}
	conStr := os.Getenv("CONNECTION_STRING")

	// Initialize database manager
	db, err := database.NewDbManager(conStr, context.Background())
	if err != nil {
		log.Fatalf("error initializing DB manager: %s", err.Error())
	}

	// Initialize map for users and uuids
	users := make(map[string]string)

	// Endpoints
	server.POST("/login", api.LoginHandler(db, users))
	server.POST("/signup", api.SignUpHandler(db, users))
	server.POST("/logout", api.LogoutHandler(users))
	server.POST("/auth", api.AuthHandler(users))
	server.GET("/ping", api.PingHandler)
	// Query Endpoints
	server.GET("/queries", api.GetQueryList(db))
	server.GET("/query/:qid", api.ReadQueryLiteralHandler(db))     // Return the query SQL string
	server.GET("/query/:qid/execute", api.ExecuteQueryHandler(db)) // Execute a saved query (custom or standard saved queries)
	server.POST("/query/save", api.SaveQueryHandler(db))
	server.DELETE("/query/:qid", api.DeleteQueryHandler(db))

	err = server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
