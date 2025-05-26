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
func AuthHandler(activeUsers map[string]string) gin.HandlerFunc {
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

		c.Status(http.StatusOK)
	}
}

// Respond with JSON representation of all users
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
func DeleteUserHandler(d *database.DbManager, activeUsers map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")

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

func AddUserRoleHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userName := c.Param("uname")
		roleName := c.Param("rname")

		err := d.AddUserRole(userName, roleName)
		if err != nil {
			errMsg := fmt.Sprintf("unable to add user role: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		c.Status(http.StatusOK)
	}
}

func DeleteUserRoleHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userName := c.Param("uname")
		roleName := c.Param("rname")

		err := d.DeleteUserRole(userName, roleName)
		if err != nil {
			errMsg := fmt.Sprintf("unable to delete user role: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		c.Status(http.StatusOK)
	}
}

// Get list of roles, and number of users per role
func GetUsersPerRoleHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, err := d.UsersPerRole()
		if err != nil {
			errMsg := fmt.Sprintf("unable to get roles: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		c.JSON(http.StatusOK, roles)
	}
}

// Add a new role to the database using the 'rname' route parameter.
func AddRoleHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleName := c.Param("rname")

		err := d.AddRole(roleName)
		if err != nil {
			errMsg := fmt.Sprintf("unable to add role: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		c.Status(http.StatusOK)
	}
}

// Update an existing role, user `rname` and `newname“ route parameters.
func UpdateRoleHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		oldRole := c.Param("rname")
		newRole := c.Param("newname")

		err := d.UpdateRoleName(oldRole, newRole)
		if err != nil {
			errMsg := fmt.Sprintf("unable to update role: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		c.Status(http.StatusOK)
	}
}

// Delete a role, using `rname` route parameter.
func DeleteRoleHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleName := c.Param("rname")

		err := d.DeleteRole(roleName)
		if err != nil {
			errMsg := fmt.Sprintf("unable to delete role: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
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
