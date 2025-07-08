package database

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const nilId = -1

// Parse the 'key'. If it parses to an integer, use the 'table'.'intColumn' to compare against 'key'. Otherwise use the 'table'.'stringColumn' to compare against 'key'
func GetSearchSuffix(key string, table string, stringColumn string, intColumn string) string {
	query := ""

	// Attempt to parse the ID into a separate variable
	_, err := strconv.Atoi(key)
	if err != nil {
		query = table + "." + stringColumn + "= $1;"
	} else {
		query = table + "." + intColumn + "=$1;"
	}

	// Return the query
	return query
}

func GetUserSearchSuffix(user_name string) string {
	return GetSearchSuffix(user_name, "users", "user_name", "user_id")
}

func GetQuerySearchSuffix(query_name string) string {
	return GetSearchSuffix(query_name, "queries", "query_name", "query_id")
}

type DbManager struct {
	conStr     string
	Connection *pgxpool.Pool
	context    context.Context
}

type User struct {
	Id       int    `json:"id,omitempty"`
	Name     string `json:"username"`
	Password string `json:"password,omitempty"`
	Roles    []int  `json:"role"`
}

type Role struct {
	Id        int    `json:"id,omitempty"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	UserCount int    `json:"usercount,omitempty"`
}

type Query struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Literal string `json:"literal"`
}

type Settings struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewDbManager(connectionString string, ctx context.Context) (*DbManager, error) {
	if len(connectionString) == 0 {
		msg := "error creating new DbManager: empty connection string"
		return nil, errors.New(msg)
	}

	dbManager := DbManager{
		conStr:     connectionString,
		Connection: nil,
		context:    ctx,
	}

	err := dbManager.ConnectToDatabase()
	if err != nil {
		return nil, err
	}

	return &dbManager, nil
}

// Set connection string
func (d *DbManager) SetConnectionString(conStr string) {
	d.conStr = conStr
}

// Connect to database using set connection string
func (d *DbManager) ConnectToDatabase() error {
	if d.conStr == "" {
		errMsg := "connection string is not set"
		return errors.New(errMsg)
	}

	var err error
	d.Connection, err = pgxpool.New(context.Background(), d.conStr)

	return err
}

func (d *DbManager) GetAllUsers() ([]User, error) {
	query :=
		`SELECT
			users.user_id,
			users.user_name,
			roles.role_id
		FROM
			users
			LEFT JOIN userroles ON users.user_id = userroles.user_id
			LEFT JOIN roles ON userroles.role_id = roles.role_id
		ORDER BY
			users.user_name;`

	rows, err := d.Connection.Query(d.context, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userMap := map[int]User{}
	for rows.Next() {
		user := User{}
		var roleId *int

		err = rows.Scan(&user.Id, &user.Name, &roleId)
		if err != nil {
			return nil, err
		}

		if _, exists := userMap[user.Id]; !exists {
			userMap[user.Id] = User{
				Id:    user.Id,
				Name:  user.Name,
				Roles: []int{},
			}
		}

		u := userMap[user.Id]
		if roleId != nil {
			u.Roles = append(u.Roles, *roleId)
		} else {
			u.Roles = []int{}
		}
		userMap[user.Id] = u
	}

	// Convert map to array
	users := make([]User, 0, len(userMap))
	for _, user := range userMap {
		users = append(users, user)
	}

	return users, nil
}

// Return a user based on username. Return error if no user found.
func (d *DbManager) GetSingleUser(name string) (User, error) {
	var err error
	query :=
		`SELECT
			users.user_id, users.password_hash, roles.role_id, users.user_name
		FROM
			userroles
			RIGHT JOIN users ON userroles.user_id = users.user_id
			RIGHT JOIN roles ON userroles.role_id = roles.role_id
		WHERE ` + GetUserSearchSuffix(name)
	args := []interface{}{name}

	rows, err := d.Connection.Query(d.context, query, args...)
	if err != nil {
		return User{}, err
	}
	defer rows.Close()

	user := User{}
	roles := []int{}
	for rows.Next() {
		var userName string
		var role int

		err = rows.Scan(&user.Id, &user.Password, &role, &userName)
		if err != nil {
			return User{}, err
		}

		user.Name = userName
		roles = append(roles, role)
	}

	slices.Sort(roles)
	user.Roles = roles

	return user, nil
}

// Insert user into the database. Expects the password to be hashed using the auth module.
func (d *DbManager) InsertUser(username string, password string) error {
	user := User{}
	query :=
		`INSERT INTO users (user_name, password_hash)
		VALUES ($1, $2)
		RETURNING user_id, user_name, password_hash`

	var userId int
	err := d.Connection.QueryRow(d.context, query, username, password).Scan(&userId, &user.Name, &user.Password)
	if err != nil {
		return err
	}

	query =
		`WITH default_role AS (
			SELECT role_id
			FROM roles
			WHERE role_name = 'default'
		)
		INSERT INTO
			userroles (user_id, role_id)
			SELECT $1, default_role.role_id
			FROM default_role`

	_, err = d.Connection.Query(d.context, query, userId)
	if err != nil {
		return err
	}

	return nil
}

// Delete a user from the database by username
func (d *DbManager) DeleteUser(username string) error {
	query :=
		`DELETE FROM users
		WHERE user_name = $1`
	_, err := d.Connection.Query(d.context, query, username)
	return err
}

// Update a user, expects pre-hashed password (do not give this plaintext pls)
func (d *DbManager) UpdateUser(userid int, name string, password string) error {

	if name == "" && password == "" {
		return nil
	}

	query := "UPDATE users SET"
	argCount := 1
	var args []any

	if name != "" {
		query += " user_name = $" + strconv.Itoa(argCount)
		args = append(args, name)
		argCount++
	}

	if password != "" {
		if argCount > 1 {
			query += ","
		}
		query += " password_hash = $" + strconv.Itoa(argCount)
		args = append(args, password)
		argCount++
	}

	query += " WHERE user_id = $" + strconv.Itoa(argCount)
	args = append(args, userid)

	_, err := d.Connection.Exec(d.context, query, args...)
	if err != nil {
		errMsg := fmt.Sprintf("unable to update user '%d': %s", userid, err.Error())
		return errors.New(errMsg)
	}

	return nil
}

// Add a role to a user
func (d *DbManager) AddUserRole(userID string, roleID string) error {
	query :=
		`INSERT INTO userroles (user_id, role_id)
		SELECT users.user_id, roles.role_id
		FROM users, roles
		WHERE users.user_id = $1
		AND roles.role_id = $2`

	_, err := d.Connection.Exec(d.context, query, userID, roleID)
	if err != nil {
		errMsg := fmt.Sprintf("unable to add new role to: %s", err.Error())
		return errors.New(errMsg)
	}

	return nil
}

// Delete a role from a user
func (d *DbManager) DeleteUserRole(userID string, roleID string) error {
	query :=
		`DELETE FROM userroles
		WHERE user_id = $1
		AND role_id = $2`

	_, err := d.Connection.Exec(d.context, query, userID, roleID)
	if err != nil {
		errMsg := fmt.Sprintf("unable to delete '%s' role from '%s': %s",
			roleID, userID, err.Error())
		return errors.New(errMsg)
	}

	return nil
}

// Get count of users per role
func (d *DbManager) GetRoles() ([]Role, error) {
	query :=
		`SELECT r.role_id, r.role_name, r.role_color, COUNT(ur.role_id)
		FROM roles r
		LEFT JOIN userroles ur ON r.role_id = ur.role_id
		WHERE r.role_name != 'default'
		GROUP BY r.role_id, r.role_name, r.role_color`

	rows, err := d.Connection.Query(d.context, query)
	if err != nil {
		errMsg := fmt.Sprintf("error getting roles: %s", err.Error())
		return nil, errors.New(errMsg)
	}
	defer rows.Close()

	roles := []Role{}
	for rows.Next() {
		newRole := Role{}
		colorField := pgtype.Text{}
		err = rows.Scan(&newRole.Id, &newRole.Name, &colorField, &newRole.UserCount)
		if err != nil {
			errMsg := fmt.Sprintf("error scanning row: %s", err.Error())
			return nil, errors.New(errMsg)
		}
		if colorField.Valid {
			newRole.Color = colorField.String
		}

		roles = append(roles, newRole)
	}

	return roles, nil
}

// Add a new role to the database
func (d *DbManager) AddRole(role Role) (int, error) {
	query :=
		`INSERT INTO roles (role_name, role_color)
		VALUES ($1, $2)
		ON CONFLICT (role_name) DO NOTHING
		RETURNING role_id`

	id := nilId
	err := d.Connection.QueryRow(d.context, query, role.Name, role.Color).Scan(&id)
	if err != nil {
		errMsg := fmt.Sprintf("unable to add new role: %s", err.Error())
		return nilId, errors.New(errMsg)
	}

	return id, nil
}

// Update a role.
func (d *DbManager) UpdateRole(roleId int, name string, color string) error {
	if name == "" && color == "" {
		return nil
	}

	query := "UPDATE roles SET"
	argCount := 1
	var args []any

	if name != "" {
		query += " role_name = $" + strconv.Itoa(argCount)
		args = append(args, name)
		argCount++
	}

	if color != "" {
		if argCount > 1 {
			query += ","
		}
		query += " role_color = $" + strconv.Itoa(argCount)
		args = append(args, color)
		argCount++
	}

	query += " WHERE role_id = $" + strconv.Itoa(argCount)
	args = append(args, roleId)

	_, err := d.Connection.Exec(d.context, query, args...)
	if err != nil {
		errMsg := fmt.Sprintf("unable to update role '%d': %s", roleId, err.Error())
		return errors.New(errMsg)
	}

	return nil
}

// Update a role name
func (d *DbManager) UpdateRoleName(roleId int, newName string) error {
	query :=
		`UPDATE roles
		SET role_name = $1
		WHERE role_id = $2`

	_, err := d.Connection.Exec(d.context, query, newName, roleId)
	if err != nil {
		errMsg := fmt.Sprintf("unable to update '%d' role: %s", roleId, err.Error())
		return errors.New(errMsg)
	}

	return nil
}

// Update a role color
func (d *DbManager) UpdateRoleColor(roleId int, newColor string) error {
	query :=
		`UPDATE roles
		SET role_color = $1
		WHERE role_id = $2
		`

	_, err := d.Connection.Exec(d.context, query, newColor, roleId)
	if err != nil {
		return err
	}

	return nil
}

// Delete a role
func (d *DbManager) DeleteRole(roleId int) error {
	query :=
		`DELETE FROM roles
		WHERE role_id = $1`

	_, err := d.Connection.Exec(d.context, query, roleId)
	if err != nil {
		errMsg := fmt.Sprintf("unable to delete '%d' role: %s", roleId, err.Error())
		return errors.New(errMsg)
	}

	return nil
}

// Insert a custom query into the database
func (d *DbManager) InsertCustomQuery(query_name string, query_string string) (int, error) {
	query_id := 0
	if len(query_name) == 0 {
		return query_id, fmt.Errorf("Query name was empty")
	}

	query := "INSERT INTO queries (query_name, query_string) VALUES ($1, $2) RETURNING query_id"
	err := d.Connection.QueryRow(d.context, query, query_name, query_string).Scan(&query_id)

	return query_id, err
}

// Get the custom query string saved as some ID
func (d *DbManager) GetCustomQuery(idName string) (string, string, error) {
	query_string := ""
	query_name := ""
	query := "SELECT query_name,query_string FROM queries WHERE " + GetQuerySearchSuffix(idName)

	// Query the Db
	err := d.Connection.QueryRow(d.context, query, idName).Scan(&query_name, &query_string)

	return query_name, query_string, err
}

// Delete a custom query based on ID
func (d *DbManager) DeleteCustomQuery(query_name string) error {
	query := "DELETE FROM queries WHERE " + GetQuerySearchSuffix(query_name)
	_, err := d.Connection.Exec(d.context, query, query_name)

	return err
}

// List the available queries
func (d *DbManager) ListCustomQueries() ([]Query, error) {
	var list []Query
	query := "SELECT * FROM queries"

	// Query the databse for all queries
	rows, err := d.Connection.Query(d.context, query)
	if err != nil {
		return nil, err
	}
	// Ensure the rows are closed properly
	defer rows.Close()

	// Iterate through rows
	for rows.Next() {
		var query Query

		// Scan the row for the query ID and String Literal
		err := rows.Scan(&query.Id, &query.Name, &query.Literal)
		if err != nil {
			return nil, err
		}

		list = append(list, query)
	}

	// Success!
	return list, nil
}

func (d *DbManager) ExecuteCustomQuery(query string) ([]map[string]string, error) {
	if d == nil {
		return nil, errors.New("This service has not been connected to a data-only database.\nUse POST /settings/cdb to set the postgres connection string.")
	}

	var retRows []map[string]string

	// Get Rows
	rows, err := d.Connection.Query(d.context, query)
	if err != nil {
		return nil, err
	}

	// Ensure rows close
	defer rows.Close()

	// Get columns
	fields := rows.FieldDescriptions()

	// Get values
	values := make([]interface{}, len(fields))

	row := 0
	for rows.Next() {
		for i := range values {
			values[i] = new(interface{})
		}

		// Scan for values in the row to populate the interface
		err := rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		// Append a map to the list
		retRows = append(retRows, map[string]string{})

		// Dereference
		for i, valPtr := range values {
			v := *(valPtr.(*interface{}))

			// Add to retRows
			retRows[row][fields[i].Name] = fmt.Sprintf("%v", v)
		}

		// Increment the row counter
		row++
	}

	// Success
	return retRows, nil
}

// Get key-value pair from settings table in database
func (d *DbManager) GetSetting(key string) (string, error) {
	query := "SELECT svalue FROM settings WHERE skey = $1;"
	var value string

	err := d.Connection.QueryRow(d.context, query, key).Scan(&value)
	if err != nil {
		return "", fmt.Errorf("error getting setting '%s': %w", key, err)
	}

	return value, nil
}

// Set or update a key-value pair in the settings table in the database
func (d *DbManager) SetSetting(key string, value string) error {
	query := "INSERT INTO settings (skey, svalue) VALUES ($1, $2) ON CONFLICT (skey) DO UPDATE SET svalue = $2;"

	_, err := d.Connection.Exec(d.context, query, key, value)
	if err != nil {
		return fmt.Errorf("error setting key '%s' to value '%s': %w", key, value, err)
	}

	return nil
}