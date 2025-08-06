package database

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx"
)

type ScopeType string

const (
	GlobalScope    ScopeType = "global"
	DashboardScope ScopeType = "dashboard"
	EmptyScope     ScopeType = ""
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

	rows, err := d.Connection.Query(d.Context, query, scope)
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
		`SELECT permission_id, permission_name, permission_scope
         FROM permissions`

	rows, err := d.Connection.Query(d.Context, query)
	if err != nil {
		return nil, err
	}

	permissions := []Permission{}
	for rows.Next() {
		var p Permission
		err = rows.Scan(&p.Id, &p.Name, &p.Scope)
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
	var scope string
	query := `SELECT permission_scope FROM permissions WHERE permission_id = $1;`
	err := d.Connection.QueryRow(d.Context, query, permissionId).Scan(&scope)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			errMsg := fmt.Sprintf("permission with ID %d not found", permissionId)
			return errors.New(errMsg)
		}
		errMsg := fmt.Sprintf("unable to check permission scope: %s", err.Error())
		return errors.New(errMsg)
	}

	if scope != "dashboard" {
		// Explicitly return an error if the scope is not 'dashboard'
		errMsg := fmt.Sprintf("permission with ID %d has scope '%s', but only 'dashboard' scope is allowed", permissionId, scope)
		return errors.New(errMsg)
	}

	query =
		`INSERT INTO
			dashboardRolePermissions (dash_id, role_id, permission_id)
		SELECT
			$1, $2, $3
		FROM
			permissions
		WHERE
			permission_id = $3
		AND
			permission_scope = 'dashboard';`

	_, err = d.Connection.Exec(d.Context, query, dashId, roleId, permissionId)
	if err != nil {
		errMsg := fmt.Sprintf("unable to add dashboard role permission: %s", err.Error())
		return errors.New(errMsg)
	}

	return nil
}

// Delete an association between a dashboard, role, and permission
func (d *DbManager) DeleteDashboardRolePermission(dashId int, roleId int, permissionId int) error {
	query :=
		`DELETE FROM dashboardRolePermissions
		WHERE
			dash_id = $1
		AND
			role_id = $2
		AND
			permission_id = $3`

	_, err := d.Connection.Exec(d.Context, query, dashId, roleId, permissionId)
	if err != nil {
		errMsg := fmt.Sprintf("unable to delete dashboard role permission: %s", err.Error())
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
