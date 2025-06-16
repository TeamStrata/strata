package api

import (
	"fmt"
	"net/http"
	"strconv"
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

// Create a new user, hash the password, store user in database
//
// @Summary User signup
// @Description Creates a new user, hashes their password, stores them in the database, and sets a session cookie.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body database.User true "New user details"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /signup [post]
func SignUpHandler(d *database.DbManager, activeUsers map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := database.User{}
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		hash, err := auth.HashPassword(user.Password)
		if err != nil {
			errMsg := fmt.Sprintf("unable to hash provided password: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		err = d.InsertUser(user.Name, hash)
		if err != nil {
			errMsg := fmt.Sprintf("unable to add user to database: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
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

// Respond with JSON representation of all users
//
// @Summary Get all users
// @Description Retrieves a list of all users from the database.
// @Tags Users
// @Produce json
// @Success 200 {array} database.User "Successfully retrieved users"
// @Failure 500 {string} string "Internal Server Error"
// @Failure 204 {string} string "No Content"
// @Router /users [get]
func GetUsersHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := d.GetAllUsers()
		if err != nil {
			errMsg := fmt.Sprintf("unable to get users: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		if len(users) <= 0 {
			c.Status(http.StatusNoContent)
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

// Delete a user that matches the name parameter
//
// @Summary Delete user by name
// @Description Deletes a user from the database and removes their session from active users.
// @Tags Users
// @Produce json
// @Param name path string true "User name"
// @Success 200 {string} string "OK"
// @Failure 500 {string} string "Internal Server Error"
// @Router /user/{uname} [delete]
func DeleteUserHandler(d *database.DbManager, activeUsers map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("uname")

		for id, tmpName := range activeUsers {
			if tmpName == name {
				delete(activeUsers, id)
				break
			}
		}

		err := d.DeleteUser(name)
		if err != nil {
			errMsg := fmt.Sprintf("unable to delete user: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		c.Status(http.StatusOK)
	}
}

// Add a role to a user
//
// @Summary Add user role
// @Description Assigns a specific role to a user.
// @Tags User Roles
// @Produce json
// @Param uname path string true "User name"
// @Param rname path string true "Role name"
// @Success 200 {string} string "OK"
// @Failure 500 {string} string "Internal Server Error"
// @Router /user/{uname}/role/{rname} [post]
func AddUserRoleHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userName := c.Param("uname")
		roleName := c.Param("rname")

		err := d.AddUserRole(userName, roleName)
		if err != nil {
			// This errMsg needs to be more descriptive based on the actual error.
			// For now, I'll leave it as a generic message.
			errMsg := fmt.Sprintf("unable to add user role: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		c.Status(http.StatusOK)
	}
}

// Delete a role from a user
//
// @Summary Delete user role
// @Description Removes a specific role from a user.
// @Tags User Roles
// @Produce json
// @Param uname path string true "User name"
// @Param rname path string true "Role name"
// @Success 200 {string} string "OK"
// @Failure 500 {string} string "Internal Server Error"
// @Router /user/{uname}/role/{rname} [delete]
func DeleteUserRoleHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userName := c.Param("uname")
		roleName := c.Param("rname")

		err := d.DeleteUserRole(userName, roleName)
		if err != nil {
			// This errMsg needs to be more descriptive based on the actual error.
			// For now, I'll leave it as a generic message.
			errMsg := fmt.Sprintf("unable to delete user role: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		c.Status(http.StatusOK)
	}
}

// @Summary Get roles
// @Description Retrieves a list of roles and the count of users assigned to each role.
// @Tags Roles
// @Produce json
// @Success 200 {object} []database.Role "Successfully retrieved roles with user counts"
// @Failure 500 {string} string "Internal Server Error"
// @Router /roles [get]
func GetRolesHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, err := d.GetRoles()
		if err != nil {
			errMsg := fmt.Sprintf("unable to get roles: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		c.JSON(http.StatusOK, roles)
	}
}

// Add a new role to the database using the 'rname' route parameter.
//
// @Summary Add a new role
// @Description Creates a new role in the database.
// @Tags Roles
// @Accept json
// @Produce json
// @Param role body database.Role true "Role to be added"
// @Success 200 {string} string "OK"
// @Failure 500 {string} string "Internal Server Error"
// @Router /role [post]
func AddRoleHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newRole database.Role
		err := c.ShouldBindJSON(&newRole)
		if err != nil {
			errMsg := fmt.Sprintf("unable to bind role to JSON: %s", err.Error())
			c.String(http.StatusBadRequest, errMsg)
			return
		}

		_, err = d.AddRole(newRole)
		if err != nil {
			errMsg := fmt.Sprintf("unable to add role: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		c.Status(http.StatusOK)
	}
}

// Update an existing role using the 'rid' route parameter.
//
// @Summary Update a role
// @Description Creates a new role in the database.
// @Tags Roles
// @Produce json
// @Param rid path string true "Role Id"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /role/{rid} [post]
func UpdateRoleHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleIdStr := c.Param("rid")
		roleId, err := strconv.Atoi(roleIdStr)
		if err != nil {
			errMsg := fmt.Sprintf("role id not provided as route parameter: %s", err.Error())
			c.String(http.StatusBadRequest, errMsg)
		}

		// Bind JSON request body
		var roleUpdate database.Role
		err = c.ShouldBindJSON(&roleUpdate)
		if err != nil {
			errMsg := fmt.Sprintf("unable to bind request JSON body: %s", err.Error())
			c.String(http.StatusBadRequest, errMsg)
		}

		err = d.UpdateRole(roleId, roleUpdate.Name, roleUpdate.Color)
		if err != nil {
			errMsg := fmt.Sprintf("unable to update role '%d': %s", roleId, err.Error())
			c.String(http.StatusInternalServerError, errMsg)
		}
	}
}

// Delete a role, using `rid` route parameter.
//
// @Summary Delete a role
// @Description Deletes a role from the database.
// @Tags Roles
// @Produce json
// @Param rid path string true "Id of role to delete"
// @Success 200 {string} string "OK"
// @Failure 500 {string} string "Internal Server Error"
// @Router /role/{rid} [delete]
func DeleteRoleHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleId, err := strconv.Atoi(c.Param("rid"))
		if err != nil {
			c.String(http.StatusBadRequest, "unable to convert role id to int")
			return
		}

		err = d.DeleteRole(roleId)
		if err != nil {
			errMsg := fmt.Sprintf("unable to delete role: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		c.Status(http.StatusOK)
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
