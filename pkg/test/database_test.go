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

// Passes if the database can return a user with the name 'admin'
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
	if err != nil {
		log.Printf("expected user:	%+v", expectedUser)
		log.Printf("actual user:	%+v", actualUser)
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

	if len(users) < 1 {
		t.Fatalf("expected at least 1 user, received: %d", len(users))
	}

	if users[0].Name != expectedUser.Name {
		log.Printf("expected user:	%+v", expectedUser)
		log.Printf("actual user:	%+v", users[0])
		t.Fatalf("expected username \"admin\", received: %s", users[0].Name)
	}
}

func Test_DeleteUser(t *testing.T) {
	db, err := SetupDbManager()
	if err != nil {
		t.Fatalf("error initializing DB manager: %s", err.Error())
	}

	testUser := database.User{
		Name:     "testDelete",
		Password: "testDelete",
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
		Name:     "testAddRole",
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

	_ = db.DeleteUser(testUser.Name)
}

func Test_DeleteUserRole(t *testing.T) {
	db, err := SetupDbManager()
	if err != nil {
		t.Fatalf("error initializing DB manager: %s", err.Error())
	}
	password := auth.HashPassword("test")
	testUser := database.User{
		Name:     "testDeleteRole",
		Password: password,
	}

	err = db.InsertUser(testUser.Name, testUser.Password)
	if err != nil {
		t.Fatalf("error inserting test user into database: %s", err.Error())
	}

	role := "delete"
	err = db.AddUserRole(testUser.Name, role)
	if err != nil {
		t.Fatalf("error adding admin role to user: %s", err.Error())
	}

	err = db.DeleteUserRole(testUser.Name, role)
	if err != nil {
		t.Fatalf("error deleting user role: %s", err.Error())
	}

	_ = db.DeleteUser(testUser.Name)
}
