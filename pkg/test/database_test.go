package test

import (
	"context"
	"log"
	"testing"

	"github.com/TeamStrata/strata/pkg/database"
)

func SetupDbManager() (*database.DbManager, error) {
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
		Name:     "admin",
		Password: "$2a$10$GEqRMhGYBay/4uXY50eyP.heui16Vs9WC//cwxt9mHijfJ.4xvi9.",
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

func Test_GetAllUsers(t *testing.T) {
	db, err := SetupDbManager()
	if err != nil {
		t.Fatalf("error initializing DB manager: %s", err.Error())
	}

	expectedUser := database.User{
		Name: "admin",
	}

	users, err := db.GetAllUsers()
	if err != nil {
		t.Fatalf("error getting all user names: %s", err.Error())
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user, received: %d", len(users))
	}

	log.Print(users[0])
	log.Print(expectedUser)
	if users[0].Name != expectedUser.Name {
		t.Fatalf("expected username \"admin\", received: %s", users[0].Name)
	}
}

func Test_GetConnectionString(t *testing.T) {
	connectionString, err := database.GetConnectionString(true)
	if err != nil {
		t.Error(err.Error())
	}

	expectedString := "postgresql://strata:atarts@strata_psql:5432/strata"

	if connectionString != expectedString {
		t.Errorf("expected string != connection string...\n")
		t.Errorf("expected: %s\nconnection string: %s\n", expectedString, connectionString)
		t.Fail()
	}
}
