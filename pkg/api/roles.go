package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TeamStrata/strata/pkg/database"
	"github.com/gin-gonic/gin"
)

// Doing this to simplify API,
// much easier to return slice of ints
// than Permissions from front end.
type RoleUpdateRequest struct {
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Color       string `json:"color,omitempty"`
	Permissions []int  `json:"permissions,omitempty"`
}

func (rUpdate *RoleUpdateRequest) ToRole() database.Role {
	var role database.Role
	role.Id = rUpdate.Id
	role.Name = rUpdate.Name
	role.Color = rUpdate.Color
	for _, num := range rUpdate.Permissions {
		role.Permissions = append(role.Permissions, database.Permission{Id: num, Name: ""})
	}
	return role
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
			return
		}

		// Bind JSON request body
		var roleUpdate RoleUpdateRequest
		err = c.ShouldBindJSON(&roleUpdate)
		if err != nil {
			errMsg := fmt.Sprintf("unable to bind request JSON body: %s", err.Error())
			c.String(http.StatusBadRequest, errMsg)
			return
		}

		err = d.UpdateRole(roleUpdate.ToRole())
		if err != nil {
			errMsg := fmt.Sprintf("unable to update role '%d': %s", roleId, err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		c.Status(http.StatusOK)
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
