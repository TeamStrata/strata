package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TeamStrata/strata/pkg/database"
	"github.com/gin-gonic/gin"
)

// @Summary Get permissions by scope
// @Description Retrieves all permissions filtered by a specific scope (e.g., global, server, channel).
// @Tags Permissions
// @Produce json
// @Param scope path string true "Permission scope" Enums(global,server,channel)
// @Success 200 {object} []database.Permission "Successfully retrieved scoped permissions"
// @Failure 400 {string} string "Bad Request - Invalid or missing scope parameter"
// @Failure 500 {string} string "Internal Server Error"
// @Router /permissions/{scope} [get]
func GetScopedPermissionsHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		scope := database.StringToScopeType(c.Param("scope"))
		if scope == database.EmptyScope {
			c.String(http.StatusBadRequest, "expected route parameter: ':scope'")
			return
		}

		permissions, err := d.GetScopedPermissions(scope)
		if err != nil {
			errMsg := fmt.Sprintf("unable to get scoped permissions: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		c.JSON(http.StatusOK, permissions)
	}
}

// @Summary Get all permissions
// @Description Retrieves all permissions from the database regardless of scope.
// @Tags Permissions
// @Produce json
// @Success 200 {object} []database.Permission "Successfully retrieved all permissions"
// @Failure 500 {string} string "Internal Server Error"
// @Router /permissions [get]
func GetPermissionsHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, err := d.GetPermissions()
		if err != nil {
			errMsg := fmt.Sprintf("unable to get all permissions: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		c.JSON(http.StatusOK, permissions)
	}
}

// AddDashboardRolePermissionHandler
// @Summary Add a permission to a role for a specific dashboard
// @Description Adds a permission to a role associated with a given dashboard.
// @Tags Admin
// @Accept json
// @Produce json
// @Param did path int true "Dashboard ID"
// @Param rid path int true "Role ID"
// @Param pid path int true "Permission ID"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/admin/dashboard/{did}/role/{rid}/permission/{pid} [post]
func AddDashboardRolePermissionHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		dashIdStr := c.Param("did")
		roleIdStr := c.Param("rid")
		permIdStr := c.Param("pid")

		if dashIdStr == "" || roleIdStr == "" || permIdStr == "" {
			c.String(http.StatusBadRequest, "missing route parameters")
			return
		}

		dashId, err := strconv.Atoi(dashIdStr)
		if err != nil {
			c.String(http.StatusBadRequest, "invalid dashboard id param")
		}

		roleId, err := strconv.Atoi(roleIdStr)
		if err != nil {
			c.String(http.StatusBadRequest, "invalid dashboard id param")
		}

		permId, err := strconv.Atoi(permIdStr)
		if err != nil {
			c.String(http.StatusBadRequest, "invalid dashboard id param")
		}

		err = d.AddDashboardRolePermission(dashId, roleId, permId)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.Status(http.StatusOK)
	}
}

// DeleteDashboardRolePermissionHandler
// @Summary Delete a permission from a role for a specific dashboard
// @Description Deletes a permission from a role associated with a given dashboard.
// @Tags Admin
// @Accept json
// @Produce json
// @Param did path int true "Dashboard ID"
// @Param rid path int true "Role ID"
// @Param pid path int true "Permission ID"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/admin/dashboard/{did}/role/{rid}/permission/{pid} [delete]
func DeleteDashboardRolePermissionHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		dashIdStr := c.Param("did")
		roleIdStr := c.Param("rid")
		permIdStr := c.Param("pid")

		if dashIdStr == "" || roleIdStr == "" || permIdStr == "" {
			c.String(http.StatusBadRequest, "missing route parameters")
			return
		}

		dashId, err := strconv.Atoi(dashIdStr)
		if err != nil {
			c.String(http.StatusBadRequest, "invalid dashboard id param")
		}

		roleId, err := strconv.Atoi(roleIdStr)
		if err != nil {
			c.String(http.StatusBadRequest, "invalid dashboard id param")
		}

		permId, err := strconv.Atoi(permIdStr)
		if err != nil {
			c.String(http.StatusBadRequest, "invalid dashboard id param")
		}

		err = d.DeleteDashboardRolePermission(dashId, roleId, permId)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.Status(http.StatusOK)
	}
}
