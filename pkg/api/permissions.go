package api

import (
	"fmt"
	"net/http"

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
