package api

import (
	"context"
	"log"
	"time"

	"github.com/TeamStrata/strata/pkg/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine      *gin.Engine
	Db          *database.DbManager
	ClientDb    *database.DbManager
	ActiveUsers map[string]UserSessionData
}

// Initialize server engine and database manager.
func InitServer() (*Server, error) {
	engine := gin.Default()

	//cors config (redux)
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "PATCH", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	engine.Use(cors.New(config))

	// Get database connection string
	conStr, err := database.GetConnectionString(database.DefaultPath)
	if err != nil {
		log.Fatalf("error loading .env file: %s", err.Error())
		return nil, err
	}

	// Initialize database manager
	db, err := database.NewDbManager(conStr, context.Background())
	if err != nil {
		log.Fatalf("error initializing DB manager: %s", err.Error())
		return nil, err
	}

	// Create the client database which will be used for custom queries
	cdb_str, err := db.GetSetting("cdb")
	if err != nil {
		log.Printf("error reading DB settings: %s", err.Error())
		return nil, err
	}

	clientDb, err := database.NewDbManager(cdb_str, context.Background())
	if err != nil {
		log.Printf("error initializing client DB manager: %s", err.Error())
		return nil, err
	}

	ActiveUsers := make(map[string]UserSessionData)
	server := Server{
		engine,
		db,
		clientDb,
		ActiveUsers,
	}

	return &server, nil
}
