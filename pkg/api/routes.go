package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/TeamStrata/strata/pkg/auth"
	"github.com/TeamStrata/strata/pkg/database"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

const uuidTag = "uuid"

// Login, create and set UUID cookie, add user to the hash map
//
// @Summary User login
// @Description Authenticates a user and sets a session cookie.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body database.User true "User credentials"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /login [post]
func LoginHandler(d *database.DbManager, activeUsers map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		login := database.User{}
		err := c.ShouldBindJSON(&login)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		user, err := d.GetSingleUser(login.Name)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		err = auth.AuthenticateUser(user.Password, login.Password)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		newId := addNewUUID(user.Name, activeUsers)
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
//
// @Summary User logout
// @Description Deletes the user's session UUID from the active users map.
// @Tags Authentication
// @Produce json
// @Success 200 {string} string "OK"
// @Failure 401 {string} string "Unauthorized"
// @Failure 204 {string} string "No Content"
// @Router /logout [post]
func LogoutHandler(activeUsers map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := c.Cookie(uuidTag)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		_, exists := activeUsers[id]
		if !exists {
			c.Status(http.StatusNoContent)
			return
		}

		delete(activeUsers, id)
		c.Status(http.StatusOK)
	}
}

// Check if the UUID cookie is set and valid
//
// @Summary Authenticate session
// @Description Checks if the session UUID cookie is set and corresponds to an active user.
// @Tags Authentication
// @Produce json
// @Success 200 {string} string "OK"
// @Failure 401 {string} string "Unauthorized"
// @Failure 204 {string} string "No Content"
// @Router /auth [get]
func AuthHandler(activeUsers map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid, err := c.Cookie(uuidTag)
		if err != nil {
			errMsg := fmt.Sprintf("expected a uuid cookie: %s", err.Error())
			c.String(http.StatusBadRequest, errMsg)
			c.Abort()
			return
		}

		_, exists := activeUsers[uuid]
		if !exists {
			c.String(http.StatusUnauthorized, "uuid not valid")
			c.Abort()
			return
		}

		c.Next()
	}
}

// Example endpoint that returns "pong" in the response body JSON
//
// @Summary Ping test
// @Description Responds with "pong" to check API status.
// @Tags Utility
// @Produce json
// @Success 200 {object} object{message=string} "Returns message: pong"
// @Router /ping [get]
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
