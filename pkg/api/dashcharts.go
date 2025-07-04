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
		chart_id_str := c.Param("cid")
		chart_id, err := strconv.Atoi(chart_id_str)
		if err != nil {
			c.Data(400, "text/plain", []byte("The chart ID must be an integer!"))
			c.Done()
			return
		}

		chart, err := d.GetChart(chart_id)
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		data, err := json.Marshal(chart)
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		c.Data(200, "application/json", data)
		c.Done()
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
		chartIDStr := c.Param("cid")
		chartID, err := strconv.Atoi(chartIDStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Chart ID must be an integer"})
			return
		}

		series, err := d.GetChartSeries(chartID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, series)
	}
}

func GetChartSingleSeriesHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func AddChartSeriesHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var seriesList []database.ChartSeries

		if err := c.ShouldBindJSON(&seriesList); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		if len(seriesList) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Series list is empty"})
			return
		}

		for _, series := range seriesList {
			if _, err := d.InsertChartSeries(series); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Series added successfully"})
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
