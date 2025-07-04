package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/TeamStrata/strata/pkg/auth"
	"github.com/TeamStrata/strata/pkg/database"
	"github.com/gin-gonic/gin"
)

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
func GetUserHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id_str := c.Param("uid")
		user, err := d.GetSingleUser(user_id_str)
		if err != nil {
			errMsg := fmt.Sprintf("unable to get user: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		c.JSON(http.StatusOK, user)
		c.Done()
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

// Update an existing user using the 'uid' route parameter.
//
// @Summary Update a user
// @Description Alters the user data in the database.
// @Tags Users
// @Produce json
// @Param uid path string true "User Id"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /user/{uid} [patch]
func UpdateUserHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {

		// handle route params
		uid, err := strconv.Atoi(c.Param("uid"))

		if err != nil {
			c.String(http.StatusBadRequest, "uid not specified in request")
		}

		// handle body
		var newUserData database.User

		err = c.ShouldBindJSON(&newUserData)
		if err != nil {
			msg := fmt.Sprintf("Error reading body: %s", err.Error())
			c.String(http.StatusBadRequest, msg)
		}

		password_hash := ""
		if newUserData.Password != "" {
			password_hash, err = auth.HashPassword(newUserData.Password)
		}

		if err != nil {
			errMsg := fmt.Sprintf("unable to hash provided password: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		err = d.UpdateUser(uid, newUserData.Name, password_hash)
		if err != nil {
			errMsg := fmt.Sprintf("unable to update user '%d': %s", uid, err.Error())
			c.String(http.StatusInternalServerError, errMsg)
		}
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
