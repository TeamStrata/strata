package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type PathType bool

const (
	DefaultPath PathType = false
	TestPath    PathType = true
)

func GetConnectionString(testPath PathType) (string, error) {
	path := ".env"
	if testPath {
		path = "../../" + path
	}

	err := godotenv.Load(path)
	if err != nil {
		return "", err
	}

	host := os.Getenv("DB_HOSTNAME")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_DATABASE")
	port := os.Getenv("DB_PORT")

	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username, password, host, port, name)

	return connectionString, nil
}
