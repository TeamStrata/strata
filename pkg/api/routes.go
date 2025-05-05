package api

import (
	"github.com/TeamStrata/strata/pkg/auth"
	"github.com/TeamStrata/strata/pkg/database"
	"github.com/google/uuid"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const uuidTag = "uuid"

// Login, create and set UUID cookie, add user to the hash map
func LoginHandler(d *database.DbManager, users map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		login := database.User{}
		err := c.ShouldBindJSON(&login)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		user, err := d.GetUserByUserName(login.Name)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		err = auth.AuthenticateUser(user.Password, login.Password)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		newId := addNewUUID(user.Name, users)
		c.SetCookie(
			uuidTag,
			newId,
			int(24*time.Hour.Seconds()),
			"/",
			"localhost",
			false,
			true,
		)
		c.Status(http.StatusOK)
	}
}

// Create a new user, hash the password, store user in database
func SignUpHandler(d *database.DbManager, users map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := database.User{}
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		hash := auth.HashPassword(user.Password)
		if hash == "" {
			c.Status(http.StatusInternalServerError)
			return
		}

		err = d.InsertUser(user.Name, hash)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		newId := addNewUUID(user.Name, users)
		c.SetCookie(
			uuidTag,
			newId,
			int(24*time.Hour.Seconds()),
			"/",
			"localhost",
			false,
			true,
		)
		c.Status(http.StatusOK)
	}
}

// Log out a user, delete their session UUID from the hash map
func LogoutHandler(users map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := c.Cookie(uuidTag)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		_, exists := users[id]
		if !exists {
			c.Status(http.StatusNoContent)
			return
		}

		delete(users, id)
		c.Status(http.StatusOK)
	}
}

// Check if the UUID cookie is set and valid
func AuthHandler(users map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := c.Cookie(uuidTag)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		_, exists := users[id]
		if !exists {
			c.Status(http.StatusNoContent)
			return
		}

		c.Status(http.StatusOK)
	}
}

// Example endpoint that returns "pong" in the response body JSON
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// Generate UUID if user does not already have one
func addNewUUID(username string, users map[string]string) string {
	newId := uuid.NewString()
	for _, ok := users[newId]; ok; {
		newId = uuid.NewString()
	}

	users[newId] = username
	return newId
}
