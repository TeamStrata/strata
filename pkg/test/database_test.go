package test

import (
	"context"
	"github.com/TeamStrata/strata/pkg/database"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func SetupDbManager() (*database.DbManager, error) {
	// Get database connection string
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("error loading .env file: %s", err.Error())
		return nil, err
	}
	conStr := os.Getenv("CONNECTION_STRING")

	// Initialize database manager
	db, err := database.NewDbManager(conStr, context.Background())
	if err != nil {
		log.Fatalf("error initializing DB manager: %s", err.Error())
		return nil, err
	}

	err = db.ConnectToDatabase()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Passes if the database can return a user with the name 'gopher'
func Test_GetUserByUsername(t *testing.T) {
	db, err := SetupDbManager()
	if err != nil {
		t.Fatalf("error initializing DB manager: %s", err.Error())
	}

	realUser := database.User{
		"gopher",
		"123",
	}

	actualUser, err := db.GetUserByUserName(realUser.Name)
	log.Print(realUser.Name)
	log.Print(actualUser)
	if err != nil {
		t.Fatalf("error getting user by username: %s", err.Error())
	}

	if actualUser != realUser {
		t.Fatalf("actual user does not match real user")
	}
}
