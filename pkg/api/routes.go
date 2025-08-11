package api

import (
	"bytes"
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
	Id      int
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
// @Success 200 {object} object{isAdmin=bool} "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /login [post]
func (server *Server) LoginHandler(c *gin.Context) {
	login := database.User{}
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	user, err := server.Db.GetSingleUser(login.Name)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	err = auth.AuthenticateUser(user.Password, login.Password)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	isAdmin, err := server.Db.IsUserAdmin(user)
	if err != nil {
		errMsg := fmt.Sprintf("internal server error: %s", err.Error())
		c.String(http.StatusInternalServerError, errMsg)
		return
	}

	newId := server.addNewUUID(user.Id, user.Name, isAdmin)
	c.SetCookie(
		uuidTag,
		newId,
		int(24*time.Hour.Seconds()),
		"/",
		c.Request.Host,
		false,
		true,
	)

	body := map[string]bool{
		"isAdmin": isAdmin,
	}
	c.JSON(http.StatusOK, body)
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
func (server *Server) LogoutHandler(c *gin.Context) {
	id, err := c.Cookie(uuidTag)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	_, exists := server.ActiveUsers[id]
	if !exists {
		c.Status(http.StatusNoContent)
		return
	}

	delete(server.ActiveUsers, id)
	c.Status(http.StatusOK)
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
func (server *Server) AuthHandler(c *gin.Context) {
	userUUID, err := c.Cookie(uuidTag)
	if err != nil {
		errMsg := fmt.Sprintf("expected a uuid cookie: %s", err.Error())
		c.String(http.StatusBadRequest, errMsg)
		c.Abort()
		return
	}

	_, exists := server.ActiveUsers[userUUID]
	if !exists {
		c.String(http.StatusUnauthorized, "uuid not valid")
		c.Abort()
		return
	}

	c.Next()
}

// Check if a user is an admin
//
// @Summary      Check admin status
// @Description  Returns a boolean indicating if the specified user is an admin.
// @Tags         Authentication
// @Produce      json
// @Param        uname path string true "Username"
// @Success      200 {object} map[string]bool "isAdmin"
// @Failure      400 {string} string "Bad Request"
// @Failure      404 {string} string "User Not Found"
// @Router       /api/isadmin/{uname} [get]
func (server *Server) IsAdminHandler(c *gin.Context) {
	userName := c.Param("uname")
	if userName == "" {
		c.String(http.StatusBadRequest, "expected a `:uname` route parameter")
	}

	for _, user := range server.ActiveUsers {
		if user.Name == userName {
			respBody := gin.H{"isAdmin": user.IsAdmin}
			c.JSON(http.StatusOK, respBody)
		}
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
func (server *Server) PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (server *Server) StrataBotHandler(c *gin.Context) {
	// 1. Decode incoming request body into messages slice
	var incoming struct {
		Messages []Message `json:"messages"`
	}
	if err := c.BindJSON(&incoming); err != nil {
		log.Printf("invalid request payload: %v", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ollama_model, err := server.Db.GetSetting("ollama_model")
	if err != nil {
		log.Printf("model lookup error: %v", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 2. Build chat payload
	payload := map[string]interface{}{
		"model":    ollama_model,
		"messages": incoming.Messages,
		"stream":   false,
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		log.Printf("marshal error: %v", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 3. Get host
	ollama_host, err := server.Db.GetSetting("ollama_host")
	if err != nil {
		log.Printf("host lookup error: %v", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 4. Send chat request
	resp, err := http.Post("http://"+ollama_host+"/api/chat",
		"application/json", bytes.NewReader(jsonBody))
	if err != nil {
		log.Printf("request error: %v", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer resp.Body.Close()

	// 5. Decode response
	var respData struct {
		Message Message `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		log.Printf("decode error: %v", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 6. Forward response message object(s)
	c.JSON(http.StatusOK, respData.Message)
}

// Generate UUID if user does not already have one
func (server *Server) addNewUUID(userId int, username string, isAdmin bool) string {
	for sessionUUID, session := range server.ActiveUsers {
		if userId == session.Id {
			return sessionUUID
		}
	}

	newUUId := uuid.NewString()
	for _, ok := server.ActiveUsers[newUUId]; ok; {
		newUUId = uuid.NewString()
	}

	server.ActiveUsers[newUUId] = UserSessionData{
		Id:      userId,
		Name:    username,
		IsAdmin: isAdmin,
	}

	return newUUId
}

func (server *Server) GetAllSettingsHandler(c *gin.Context) {
	settings, err := server.Db.ExecuteCustomQuery("SELECT skey, svalue FROM settings;")
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

func (server *Server) GetSettingHandler(c *gin.Context) {
	key := c.Param("key")

	settings, err := server.Db.GetSetting(key)
	if err != nil {
		errMsg := fmt.Sprintf("unable to get settings: %s", err.Error())
		c.String(http.StatusInternalServerError, errMsg)
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (server *Server) UpdateSettingHandler(c *gin.Context) {
	key := c.Param("key")

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		errMsg := fmt.Sprintf("unable to read request body: %s", err.Error())
		c.String(http.StatusInternalServerError, errMsg)
		return
	}

	err = server.Db.SetSetting(key, string(bodyBytes))
	if err != nil {
		errMsg := fmt.Sprintf("unable to update settings: %s", err.Error())
		c.String(http.StatusInternalServerError, errMsg)
		return
	}

	if key == "cdb" {
		// Disconnect from the client database, then reconnect to this new database
		if server.ClientDb != nil {
			server.ClientDb.Connection.Close()
			server.ClientDb = nil
		}

		// Attempt to connect to new database so you don't have to restart the server
		server.ClientDb, err = database.NewDbManager(string(bodyBytes), context.Background())
		if err != nil {
			errMsg := fmt.Sprintf("Unable to connect to new database: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}
	}

	c.Status(http.StatusOK)
}
