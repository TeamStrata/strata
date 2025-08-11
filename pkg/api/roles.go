package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TeamStrata/strata/pkg/database"
	"github.com/gin-gonic/gin"
)

// RoleUpdateRequest represents the request body for updating a role
// @Description Request body for updating role information
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
// @Description Retrieves a list of roles with their permissions and user counts.
// @Tags Roles
// @Produce json
// @Success 200 {object} []database.Role "Successfully retrieved roles with user counts"
// @Failure 500 {string} string "Internal Server Error"
// @Router /roles [get]
func (server *Server) GetRolesHandler(c *gin.Context) {
	roles, err := server.Db.GetRoles()
	if err != nil {
		errMsg := fmt.Sprintf("unable to get roles: %s", err.Error())
		c.String(http.StatusInternalServerError, errMsg)
		return
	}

	c.JSON(http.StatusOK, roles)
}

// Add a new role to the database using the 'rname' route parameter.
//
// @Summary Add a new role
// @Description Creates a new role in the database with optional permissions.
// @Tags Roles
// @Accept json
// @Produce json
// @Param role body database.Role true "Role to be created"
// @Success 200 {string} string "Role created successfully"
// @Failure 400 {string} string "Bad Request - Invalid JSON format"
// @Failure 500 {string} string "Internal Server Error"
// @Router /role [post]
func (server *Server) AddRoleHandler(c *gin.Context) {
	var newRole database.Role
	err := c.ShouldBindJSON(&newRole)
	if err != nil {
		errMsg := fmt.Sprintf("unable to bind role to JSON: %s", err.Error())
		c.String(http.StatusBadRequest, errMsg)
		return
	}

	_, err = server.Db.AddRole(newRole)
	if err != nil {
		errMsg := fmt.Sprintf("unable to add role: %s", err.Error())
		c.String(http.StatusInternalServerError, errMsg)
		return
	}

	c.Status(http.StatusOK)
}

// Update an existing role using the 'rid' route parameter.
//
// @Summary Update a role
// @Description Updates an existing role's name, color, and/or permissions. Only provided fields will be updated.
// @Tags Roles
// @Accept json
// @Produce json
// @Param rid path int true "Role ID"
// @Param role body RoleUpdateRequest true "Role update data"
// @Success 200 {string} string "Role updated successfully"
// @Failure 400 {string} string "Bad Request - Invalid role ID or JSON format"
// @Failure 500 {string} string "Internal Server Error"
// @Router /role/{rid} [patch]
func (server *Server) UpdateRoleHandler(c *gin.Context) {
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

	err = server.Db.UpdateRole(roleUpdate.ToRole())
	if err != nil {
		errMsg := fmt.Sprintf("unable to update role '%d': %s", roleId, err.Error())
		c.String(http.StatusInternalServerError, errMsg)
		return
	}

	c.Status(http.StatusOK)
}

// Delete a role, using `rid` route parameter.
//
// @Summary Delete a role
// @Description Deletes a role from the database and removes all associated permissions.
// @Tags Roles
// @Produce json
// @Param rid path int true "Role ID"
// @Success 200 {string} string "Role deleted successfully"
// @Failure 400 {string} string "Bad Request - Invalid role ID"
// @Failure 500 {string} string "Internal Server Error"
// @Router /role/{rid} [delete]
func (server *Server) DeleteRoleHandler(c *gin.Context) {
	roleId, err := strconv.Atoi(c.Param("rid"))
	if err != nil {
		c.String(http.StatusBadRequest, "unable to convert role id to int")
		return
	}

	err = server.Db.DeleteRole(roleId)
	if err != nil {
		errMsg := fmt.Sprintf("unable to delete role: %s", err.Error())
		c.String(http.StatusInternalServerError, errMsg)
		return
	}

	c.Status(http.StatusOK)
}
