package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/TeamStrata/strata/pkg/auth"
	"github.com/TeamStrata/strata/pkg/database"
	"github.com/google/uuid"

	"encoding/json"

	"github.com/gin-gonic/gin"
)

type UserSessionData struct {
	Name    string
	IsAdmin bool
}

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
func LoginHandler(d *database.DbManager, activeUsers map[string]UserSessionData) gin.HandlerFunc {
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

		isAdmin, err := d.IsUserAdmin(user)
		if err != nil {
			errMsg := fmt.Sprintf("internal server error: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		newId := addNewUUID(user.Name, isAdmin, activeUsers)
		c.SetCookie(
			uuidTag,
			newId,
			int(24*time.Hour.Seconds()),
			"/",
			"localhost",
			false,
			true,
		)

		body := map[string]any{
			"isAdmin": isAdmin,
		}
		c.JSON(http.StatusOK, body)
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
func LogoutHandler(activeUsers map[string]UserSessionData) gin.HandlerFunc {
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
func AuthHandler(activeUsers map[string]UserSessionData) gin.HandlerFunc {
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
func addNewUUID(username string, isAdmin bool, users map[string]UserSessionData) string {
	newId := uuid.NewString()
	for _, ok := users[newId]; ok; {
		newId = uuid.NewString()
	}

	users[newId] = UserSessionData{
		Name:    username,
		IsAdmin: isAdmin,
	}

	return newId
}

func GetAllSettingsHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		settings, err := d.ExecuteCustomQuery("SELECT skey, svalue FROM settings;")
		if err != nil {
			errMsg := fmt.Sprintf("Unable to get all settings: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		// Convert the resulting object to JSON
		data, jerr := json.Marshal(settings)
		if jerr != nil {
			c.Data(500, "text/plain", []byte(jerr.Error()))
			c.Done()
			return
		}

		c.Data(200, "application/json", data)
		c.Done()
	}
}

func GetSettingHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")

		settings, err := d.GetSetting(key)
		if err != nil {
			errMsg := fmt.Sprintf("unable to get settings: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		c.JSON(http.StatusOK, settings)
	}
}

func UpdateSettingHandler(d *database.DbManager, cdb **database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			errMsg := fmt.Sprintf("unable to read request body: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		err = d.SetSetting(key, string(bodyBytes))
		if err != nil {
			errMsg := fmt.Sprintf("unable to update settings: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		log.Printf("Before: '%p'\n", *cdb)

		// Disconnect from the client database, then reconnect to this new database
		if *cdb != nil {
			(*cdb).Connection.Close()
			*cdb = nil
		}

		// Attempt to connect to new database so you don't have to restart the server
		*cdb, err = database.NewDbManager(string(bodyBytes), context.Background())
		if err != nil {
			errMsg := fmt.Sprintf("Unable to connect to new database: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}
		log.Printf("After: '%p'\n", *cdb)

		c.Status(http.StatusOK)
	}
}
