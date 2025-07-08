package database

import (
	"errors"
	"fmt"
)

type ScopeType string

const (
	GlobalScope    ScopeType = "global"
	DashboardScope ScopeType = "dashboard"
	EmptyScope ScopeType = ""
)

type Permission struct {
	Id    int       `json:"id"`
	Name  string    `json:"name"`
	Scope ScopeType `json:"scope,omitempty"`
}

// Get all permissions with a specific scope (e.g. 'global', 'dashboard')
func (d *DbManager) GetScopedPermissions(scope ScopeType) ([]Permission, error) {
	query :=
		`SELECT permission_id, permission_name
         FROM permissions
         WHERE permissions.permission_scope = $1`

	rows, err := d.Connection.Query(d.context, query, scope)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := []Permission{}
	for rows.Next() {
		var p Permission
		err = rows.Scan(&p.Id, &p.Name)
		if err != nil {
			errMsg := fmt.Sprintf("unable to scan row: %s", err.Error())
			return nil, errors.New(errMsg)
		}

		permissions = append(permissions, p)
	}

	return permissions, nil
}

// Get all permissions
func (d *DbManager) GetPermissions() ([]Permission, error) {
	query :=
		`SELECT permission_id, permission_name
         FROM permissions`

	rows, err := d.Connection.Query(d.context, query)
	if err != nil {
		return nil, err
	}

	permissions := []Permission{}
	for rows.Next() {
		var p Permission
		err = rows.Scan(&p.Id, &p.Name)
		if err != nil {
			errMsg := fmt.Sprintf("unable to scan row: %s", err.Error())
			return nil, errors.New(errMsg)
		}

		permissions = append(permissions, p)
	}

	return permissions, nil	
}

// Associate a dashboard, role, and permission
func (d *DbManager) AddDashboardRolePermission(dashId int, roleId int, permissionId int) error {
	query :=
		`INSERT INTO dashboardRolePermissions (dash_id, role_id, permission_id)
         VALUES ($1, $2, $3)
         WHERE permission_scope`

	_, err := d.Connection.Exec(d.context, query, dashId, roleId, permissionId)
	if err != nil {
		errMsg := fmt.Sprintf("unable to add dashboard role permission: %s", err.Error())
		return errors.New(errMsg)
	}

	return nil
}

func StringToScopeType(s string) ScopeType {
	switch s {
	case string(GlobalScope):
		return GlobalScope
	case string(DashboardScope):
		return DashboardScope
	default:
		return EmptyScope
	}
}
