package api

import (
	"encoding/json"
	"io"
	"strconv"

	"github.com/TeamStrata/strata/pkg/database"
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

		println(len(queries))
		data, err := json.Marshal(queries)
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		c.Data(200, "application/json", data)
	}
}

// Return the query SQL string
func ReadQueryLiteralHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		query_id_str := c.Param("qid")

		query_id, err := strconv.Atoi(query_id_str)
		if err != nil {
			c.Data(400, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		query_string, err := d.GetCustomQuery(query_id)
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		c.Data(200, "text/plain", []byte(query_string))
		c.Done()
	}
}

// Execute a saved query (custom or standard saved queries)
func ExecuteQueryHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		query_id_str := c.Param("qid")

		// Parse QID to a number
		query_id, ierr := strconv.Atoi(query_id_str)
		if ierr != nil {
			c.Data(400, "text/plain", []byte(ierr.Error()))
			c.Done()
			return
		}

		query_string, serr := d.GetCustomQuery(query_id)
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

		c.Data(418, "text/plain", data)
		c.Done()
	}
}

func SaveQueryHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		// Inserting custom query
		id, err := d.InsertCustomQuery(string(data))
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		c.Data(200, "text/plain", []byte(strconv.Itoa(id)))

		c.Done()

	}
}

func DeleteQueryHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		query_id_str := c.Param("qid")

		query_id, err := strconv.Atoi(query_id_str)
		if err != nil {
			c.Data(400, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		q_err := d.DeleteCustomQuery(query_id)
		if q_err != nil {
			c.Data(500, "text/plain", []byte(q_err.Error()))
			c.Done()
			return
		}

		c.Status(200)
		c.Done()
	}
}
