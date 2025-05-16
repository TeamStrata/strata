package test

import (
	"context"
	"log"
	"slices"
	"testing"

	"github.com/TeamStrata/strata/pkg/auth"
	"github.com/TeamStrata/strata/pkg/database"
)

// Return true if two users have identical field values.
func compareUsers(first database.User, second database.User) bool {
	return first.Name == second.Name &&
		first.Password == second.Password &&
		slices.Compare(first.Roles, second.Roles) == 0
}

func SetupDbManager() (*database.DbManager, error) {
	// Get database connection string
	conStr, err := database.GetConnectionString(database.TestPath)
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

// Passes if the database can return a user with the name 'gopher'
func Test_GetUserByUsername(t *testing.T) {
	db, err := SetupDbManager()
	if err != nil {
		t.Fatalf("error initializing DB manager: %s", err.Error())
	}

	expectedUser := database.User{
		Name:     "admin",
		Password: "$2a$10$GEqRMhGYBay/4uXY50eyP.heui16Vs9WC//cwxt9mHijfJ.4xvi9.",
		Roles:    []string{"admin"},
	}

	actualUser, err := db.GetUserByUserName(expectedUser.Name)
	log.Printf("expected user:	%+v", expectedUser)
	log.Printf("actual user:	%+v", actualUser)
	if err != nil {
		t.Fatalf("error getting user by username: %s", err.Error())
	}

	if compareUsers(actualUser, expectedUser) {
		t.Fatalf("actual user does not match real user")
	}
}

func Test_GetAllUsers(t *testing.T) {
	db, err := SetupDbManager()
	if err != nil {
		t.Fatalf("error initializing DB manager: %s", err.Error())
	}

	expectedUser := database.User{
		Name:  "admin",
		Roles: []string{"admin"},
	}

	users, err := db.GetAllUsers()
	if err != nil {
		t.Fatalf("error getting all user names: %s", err.Error())
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user, received: %d", len(users))
	}

	log.Printf("expected user:	%+v", expectedUser)
	log.Printf("actual user:	%+v", users[0])
	if users[0].Name != expectedUser.Name {
		t.Fatalf("expected username \"admin\", received: %s", users[0].Name)
	}
}

func Test_DeleteUser(t *testing.T) {
	db, err := SetupDbManager()
	if err != nil {
		t.Fatalf("error initializing DB manager: %s", err.Error())
	}

	testUser := database.User{
		Name:     "test",
		Password: "test",
	}
	err = db.InsertUser(testUser.Name, testUser.Password)
	if err != nil {
		t.Fatalf("error inserting test user into database: %s", err.Error())
	}

	err = db.DeleteUser(testUser.Name)
	if err != nil {
		t.Fatalf("error deleting user: %s", err.Error())
	}
}

func Test_AddUserRole(t *testing.T) {
	db, err := SetupDbManager()
	if err != nil {
		t.Fatalf("error initializing DB manager: %s", err.Error())
	}
	password := auth.HashPassword("test")
	testUser := database.User{
		Name:     "test",
		Password: password,
	}

	err = db.InsertUser(testUser.Name, testUser.Password)
	if err != nil {
		t.Fatalf("error inserting test user into database: %s", err.Error())
	}

	err = db.AddUserRole(testUser.Name, "admin")
	if err != nil {
		t.Fatalf("error adding admin role to user: %s", err.Error())
	}
}
