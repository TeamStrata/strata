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

type UserDashboardPermissions struct {
	DashboardId int    `json:"id"`
	Name        string `json:"name"`
	CanView     bool   `json:"canView"`
	CanEdit     bool   `json:"canEdit"`
	CanDelete   bool   `json:"canDelete"`
}

type RolePermission struct {
	RoleId    int  `json:"roleId"`
	CanView   bool `json:"canView"`
	CanEdit   bool `json:"canEdit"`
	CanDelete bool `json:"canDelete"`
}

type DashboardRolePermissions struct {
	Id          int              `json:"id"`
	Name        string           `json:"name"`
	Content     string           `json:"content"`
	Permissions []RolePermission `json:"permissions"`
}

func (d *DbManager) GetUserDashboardPermissions(userId int) ([]UserDashboardPermissions, error) {
	query :=
		`SELECT
			d.dashboard_id,
			d.dashboard_title,
		BOOL_OR(p.permission_name = 'edit_dashboard')   AS can_edit,
		BOOL_OR(p.permission_name = 'delete_dashboard') AS can_delete,
		BOOL_OR(p.permission_name = 'view_dashboard')   AS can_view,
		FROM dashboards					AS d
		JOIN dashboardRolePermissions	AS drp ON drp.dash_id = d.dashboard_id
		JOIN rolePermissions			AS rp  ON rp.role_id = drp.role_id
											AND rp.permission_id = drp.permission_id
		JOIN userRoles					AS ur  ON ur.role_id = rp.role_id
											AND ur.user_id  = $1
		JOIN permissions				AS p ON p.permission_id = drp.permission_id
		GROUP BY d.dashboard_id, d.dashboard_title
		HAVING BOOL_OR(p.permission_name = 'view_dashboard');
		`

	rows, err := d.Connection.Query(d.Context, query, userId)
	if err != nil {
		return nil, err
	}

	dashboards := []UserDashboardPermissions{}
	for rows.Next() {
		tmp := UserDashboardPermissions{}
		err = rows.Scan(&tmp.DashboardId, &tmp.Name, &tmp.CanEdit, &tmp.CanDelete)
		if err != nil {
			return nil, err
		}
		dashboards = append(dashboards, tmp)
	}

	return dashboards, nil
}

func (d *DbManager) GetDashboardRolePermissions(id int) (DashboardRolePermissions, error) {
	query := `
		SELECT
			d.dash_id AS id,
			d.name,
			d.description,
			drp.role_id,
			BOOL_OR(p.permission_name = 'view_dashboard')   AS can_view,
			BOOL_OR(p.permission_name = 'edit_dashboard')   AS can_edit,
			BOOL_OR(p.permission_name = 'delete_dashboard') AS can_delete
		FROM dashboards d
		JOIN dashboardRolePermissions drp
			ON d.dash_id = drp.dash_id
		JOIN permissions p
			ON drp.permission_id = p.permission_id
		WHERE d.dash_id = $1
		GROUP BY d.dash_id, d.name, d.description, drp.role_id
		ORDER BY drp.role_id;
	`

	rows, err := d.Connection.Query(d.Context, query, id)
	if err != nil {
		return DashboardRolePermissions{}, err
	}
	defer rows.Close()

	var result DashboardRolePermissions
	var permissions []RolePermission
	first := true

	for rows.Next() {
		var rp RolePermission
		var id int
		var name, description string
		var canView, canEdit, canDelete bool

		err := rows.Scan(
			&id,
			&name,
			&description,
			&rp.RoleId,
			&canView,
			&canEdit,
			&canDelete,
		)
		if err != nil {
			return DashboardRolePermissions{}, err
		}

		if first {
			result.Id = id
			result.Name = name
			result.Content = description
			first = false
		}

		rp.CanView = canView
		rp.CanEdit = canEdit
		rp.CanDelete = canDelete
		permissions = append(permissions, rp)
	}

	if rows.Err() != nil {
		return DashboardRolePermissions{}, rows.Err()
	}

	result.Permissions = permissions
	return result, nil
}

// Bulk update role permissions for a dashboard.
func (d *DbManager) UpdateDashboardRolePermissions(dashPermissions DashboardRolePermissions) error {
	// Get view_dashboard, edit_dashboard, and delete_dashboard permission id.
	var readPermID, editPermID, deletePermID int
	query := `SELECT permission_id FROM permissions WHERE permission_name = 'view_dashboard'`
	err := d.Connection.QueryRow(d.Context, query).Scan(&readPermID)
	if err != nil {
		return err
	}

	query = `SELECT permission_id FROM permissions WHERE permission_name = 'edit_dashboard'`
	err = d.Connection.QueryRow(d.Context, query).Scan(&editPermID)
	if err != nil {
		return err
	}

	query = `SELECT permission_id FROM permissions WHERE permission_name = 'delete_dashboard'`
	err = d.Connection.QueryRow(d.Context, query).Scan(&deletePermID)
	if err != nil {
		return err
	}

	tx, err := d.Connection.Begin(d.Context)
	if err != nil {
		return err
	}
	defer tx.Rollback(d.Context)

	// Update name and description.
	query = `
		UPDATE dashboards
		SET dashboard_name = $1,
			dashboard_content = $2
		WHERE dash_id = $3
	`
	_, err = tx.Exec(
		d.Context,
		query,
		dashPermissions.Name,
		dashPermissions.Content,
		dashPermissions.Id,
	)
	if err != nil {
		return err
	}

	// Delete all permissions for dashboard.
	query = `DELETE FROM dashboardRolePermissions WHERE dash_id = $1`
	_, err = tx.Exec(d.Context, query, dashPermissions.Id)
	if err != nil {
		return err
	}

	// Insert given permissions for dashboard.
	query = `INSERT INTO dashboardRolePermissions (dash_id, role_id, permission_id) VALUES ($1, $2, $3)`
	for _, p := range dashPermissions.Permissions {
		if p.CanView {
			_, err := tx.Exec(d.Context, query, dashPermissions.Id, p.RoleId, readPermID)
			if err != nil {
				return err
			}
		}
		if p.CanEdit {
			_, err := tx.Exec(d.Context, query, dashPermissions.Id, p.RoleId, editPermID)
			if err != nil {
				return err
			}
		}
		if p.CanDelete {
			_, err := tx.Exec(d.Context, query, dashPermissions.Id, p.RoleId, deletePermID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(d.Context)
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
