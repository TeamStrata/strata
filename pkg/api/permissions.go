package api

import (
	"fmt"
	"net/http"
	"github.com/TeamStrata/strata/pkg/database"
	"github.com/gin-gonic/gin"
)

func GetScopedPermissionsHandler(d *database.DbManager) gin.HandlerFunc {
	return func (c *gin.Context) {
		scope := database.StringToScopeType(c.Param("scope"))
		if scope == database.EmptyScope {
			c.String(http.StatusBadRequest, "expected route parameter: ':scope'")
			return
		}
		
		permissions, err := d.GetScopedPermissions(scope)
		if err != nil {
			errMsg := fmt.Sprintf("unable to get global permissions: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		c.JSON(http.StatusOK, permissions)
	}
}

func GetPermissionsHandler(d *database.DbManager) gin.HandlerFunc {
	return func (c *gin.Context) {
		permissions, err := d.GetPermissions()
		if err != nil {
			errMsg := fmt.Sprintf("unable to get all permissions: %s", err.Error())
			c.String(http.StatusInternalServerError, errMsg)
			return
		}

		c.JSON(http.StatusOK, permissions)
	}
}

