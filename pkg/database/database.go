package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type DbManager struct {
	conStr     string
	Connection *pgxpool.Pool
	context    context.Context
}

type Query struct {
	Id      int    `json:"id"`
	Literal string `json:"literal"`
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

// Return a user based on username. Return error if no user found.
func (d *DbManager) GetUserByUserName(name string) (User, error) {
	user := User{}
	var err error
	if err = d.Connection.Ping(context.Background()); err != nil {
		return User{}, err
	}

	query := "SELECT user_name, password_hash FROM users WHERE user_name = $1"
	err = d.Connection.QueryRow(d.context, query, name).Scan(&user.Name, &user.Password)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// Insert user into the database. Expects the password to be hashed using the auth module.
func (d *DbManager) InsertUser(username string, password string) error {
	user := User{}
	query := "INSERT INTO users (user_name, password_hash) VALUES ($1, $2) RETURNING user_name, password_hash"
	err := d.Connection.QueryRow(d.context, query, username, password).Scan(&user.Name, &user.Password)
	return err
}

// Insert a custom query into the database
func (d *DbManager) InsertCustomQuery(query_string string) (int, error) {
	query_id := 0

	query := "INSERT INTO queries (query_string) VALUES ($1) RETURNING query_id"
	err := d.Connection.QueryRow(d.context, query, query_string).Scan(&query_id)

	return query_id, err
}

// Get the custom query string saved as some ID
func (d *DbManager) GetCustomQuery(query_id int) (string, error) {
	query_string := ""

	query := "SELECT query_string FROM queries WHERE query_id = $1"
	err := d.Connection.QueryRow(d.context, query, query_id).Scan(&query_string)

	return query_string, err
}

// Delete a custom query based on ID
func (d *DbManager) DeleteCustomQuery(query_id int) error {
	query := "DELETE FROM queries WHERE query_id = $1"
	_, err := d.Connection.Exec(d.context, query, query_id)

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
		err := rows.Scan(&query.Id, &query.Literal)
		if err != nil {
			return nil, err
		}

		list = append(list, query)
	}

	// Success!
	return list, nil
}
