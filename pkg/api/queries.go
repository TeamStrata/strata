package api

import (
	"encoding/json"
	"io"
	"strconv"

	"github.com/TeamStrata/strata/pkg/database"
	"github.com/gin-gonic/gin"
)

// Get the list of queries
//
// @Summary List Queries
// @Description Returns the list of custom queries created on the system.
// @Tags Queries
// @Produce json
// @Success 200 {array} database.Query "Successfully retrieved users"
// @Failure 500 {string} string "Internal Server Error"
// @Router /queries [get]
func GetQueryList(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		queries, err := d.ListCustomQueries()
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		data, err := json.Marshal(queries)
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		c.Data(200, "application/json", data)
		c.Done()
	}
}

// Return the query SQL string
//
// @Summary Get Single Query
// @Description Returns information about a single query
// @Tags Queries
// @Produce json
// @Success 200 {object} database.Query
// @Failure 500 {string} string "Internal Server Error"
// @Param qid query string true "Query name or ID to retrieve"
// @Router /query/{qid} [get]
func ReadQueryLiteralHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		query_id_str := c.Param("qid")

		query_name, query_string, err := d.GetCustomQuery(query_id_str)
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		c.JSON(200, gin.H{
			"name":    query_name,
			"literal": query_string,
		})
		c.Done()
	}
}

// Execute a saved query (custom or standard saved queries)
//
// @Summary Execute Query
// @Description Execute a custom query and return the custom rows and columns
// @Tags Queries
// @Produce json
// @Success 200 {string} string "OK"
// @Failure 500 {string} string "Internal Server Error"
// @Param qid query string type "Query name or ID"
// @Router /query/{qid}/execute [get]
func ExecuteQueryHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		query_id_str := c.Param("qid")

		_, query_string, serr := d.GetCustomQuery(query_id_str)
		if serr != nil {
			c.Data(500, "text/plain", []byte(serr.Error()))
			c.Done()
			return
		}

		// Perform the custom query
		rows, qerr := d.ExecuteCustomQuery(query_string)
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

		c.Data(200, "application/json", data)
		c.Done()
	}
}

// Save a custom query to the database
//
// @Summary Save Query
// @Description Save a custom query to the database. Responds with the Query ID.
// @Tags Queries
// @Accept json
// @Param query body database.Query true "New query details"
// @Success 200 {integer} 3
// @Failure 400 {string} string "Query Name Invalid"
// @Failure 500 {string} string "Internal Server Error" 
// @Param qid query string type "Query name or ID"
// @Router /query/{qid} [post]
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
//
// @Summary Delete Query
// @Description Delete a custom query from the database.
// @Tags Queries
// @Success 200 {string} string "OK"
// @Failure 500 {string} string "Internal Server Error"
// @Param qid query string type "Query name or ID"
// @Router /query/{qid} [delete]
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
