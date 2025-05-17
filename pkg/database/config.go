package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type PathType bool

const (
	DefaultPath     PathType = false
	TestPath        PathType = true
	dockerBuildFlag          = "STRATA_RELEASE"
)

func GetConnectionString(testPath PathType) (string, error) {
	path := ".env"
	if testPath {
		path = fmt.Sprintf("../../%s", path)
	}

	var suffix string
	if os.Getenv(dockerBuildFlag) == "true" {
		suffix = "" // Any release version is expected to populate .env with the desired environment variables
	} else {
		suffix = ".local"
	}
	path = fmt.Sprintf("%s%s", path, suffix)

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
