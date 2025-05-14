package main

import (
	"context"
	"log"
	"time"

	"github.com/TeamStrata/strata/pkg/api"
	"github.com/TeamStrata/strata/pkg/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	//cors config (redux)
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // 👈 specify your frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true, // 👈 enable sending cookies
		MaxAge:           12 * time.Hour,
	}
	server.Use(cors.New(config))

	// Get database connection string
	conStr, err := database.GetConnectionString(database.DefaultPath)
	if err != nil {
		log.Fatalf("error loading .env file: %s", err.Error())
	}

	// Initialize database manager
	db, err := database.NewDbManager(conStr, context.Background())
	if err != nil {
		log.Fatalf("error initializing DB manager: %s", err.Error())
	}

	// Initialize map for active users and uuids
	activeUsers := make(map[string]string)

	// Auth endpoints
	server.POST("/login", api.LoginHandler(db, activeUsers))
	server.POST("/signup", api.SignUpHandler(db, activeUsers))
	server.POST("/logout", api.LogoutHandler(activeUsers))
	server.POST("/auth", api.AuthHandler(activeUsers))
	server.GET("/users", api.GetUsersHandler(db))
	server.DELETE("/user/:name", api.DeleteUserHandler(db, activeUsers))
	server.PUT("/user/:name/role", api.UpdateUserRoleHandler(db))

	// Frontend
	server.Static("/assets", "ui/dist/assets")
	server.StaticFile("/", "ui/dist/index.html")
	server.StaticFile("/favicon.ico", "ui/dist/favicon.ico")

	// Account / Roles

	// Dashboard Components

	// Query Endpoints
	server.GET("/queries", api.GetQueryList(db))
	server.GET("/query/:qid", api.ReadQueryLiteralHandler(db))     // Return the query SQL string
	server.GET("/query/:qid/execute", api.ExecuteQueryHandler(db)) // Execute a saved query (custom or standard saved queries)
	server.POST("/query/:qid", api.SaveQueryHandler(db))
	server.DELETE("/query/:qid", api.DeleteQueryHandler(db))

	// Misc Endpoints
	server.GET("/ping", api.PingHandler)

	err = server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
