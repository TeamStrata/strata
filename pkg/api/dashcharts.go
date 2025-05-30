package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/TeamStrata/strata/pkg/database"
	"github.com/gin-gonic/gin"
)

// / Charts
func GetChartListHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		charts, err := d.ListAllCharts()
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		data, err := json.Marshal(charts)
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		c.Data(200, "application/json", data)
		c.Done()
	}
}

func GetChartHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func CreateChartHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		chart := database.Chart{}
		err := c.ShouldBindJSON(&chart)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		// Inserting custom query
		id, err := d.InsertChart(chart)
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		c.Data(200, "text/plain", []byte(strconv.Itoa(id)))
		c.Done()
	}
}

func DeleteChartHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// / Chart Series
func GetChartSeriesListHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func GetChartSingleSeriesHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func AddChartSeriesHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func DeleteChartSingleSeriesHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func DeleteAllChartSeriesHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

// / Dashboard itself
func GetDashboardListHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func GetDashboardHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func CreateDashboardPageHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func DeleteDashboardHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func ListDashboardChartsHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func AppendChartToDashboardHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func RemoveChartFromDashboardHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
