package api

import (
	//	"github.com/TeamStrata/strata/pkg/auth"
	"encoding/json"
	"io"
	"strconv"

	"github.com/TeamStrata/strata/pkg/database"
	//	"github.com/google/uuid"
	//	"net/http"
	//	"time"

	"github.com/gin-gonic/gin"
)

// Get the list of queries
func GetQueryList(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		queries, err := d.ListCustomQueries()
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		data, err2 := json.Marshal(queries)
		if err2 != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		c.Data(200, "application/json", data)
		c.Done()
	}
}

// Return the query SQL string
func ReadQueryLiteralHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		query_id_str := c.Param("qid")

		query_string, err := d.GetCustomQuery(query_id_str)
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		c.Data(200, "application/sql", []byte(query_string))
		c.Done()
	}
}

// Execute a saved query (custom or standard saved queries)
func ExecuteQueryHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		query_id_str := c.Param("qid")

		query_string, serr := d.GetCustomQuery(query_id_str)
		if serr != nil {
			c.Data(500, "text/plain", []byte(serr.Error()))
			c.Done()
			return
		}
		
		// Perform the custom query
		rows, qerr := d.ExecuteCustomQuery(query_string);
		if qerr != nil {
			c.Data(500, "text/plain", []byte(qerr.Error()))
			c.Done()
			return
		}

		// Convert the resulting object to JSON
		data, jerr := json.Marshal(rows)
		if jerr != nil {
			c.Data(500, "text/plain", []byte(jerr.Error()))
			c.Done()
			return
		}

		c.Data(418, "application/json", data)
		c.Done();
	}
}

// Save a custom query to the database
func SaveQueryHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		query_name := c.Param("qid")
		
		// To prevent the query from being inaccessible later on, we need to sanitize
		// If the query name can be parsed to an integer, reject it. It will clash with the ID, the user may accidentally select the wrong thing.
		_, err := strconv.Atoi(query_name)
		// If there is NO error parsing to integer, then we raise our own error
		if err == nil {
			c.Data(400, "text/plain", []byte("The query name must not be an integer!"))
			c.Done()
			return
		}


		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		// Inserting custom query
		id, err := d.InsertCustomQuery(query_name, string(data))
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		c.Data(200, "text/plain", []byte(strconv.Itoa(id)))
		c.Done()
	}
}

// Delete a query 
func DeleteQueryHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		query_id_str := c.Param("qid")

		q_err := d.DeleteCustomQuery(query_id_str)
		if q_err != nil {
			c.Data(500, "text/plain", []byte(q_err.Error()))
			c.Done()
			return
		}

		c.Status(200)
		c.Done()
	}
}
