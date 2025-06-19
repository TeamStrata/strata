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

func GetChartTitlesListHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		charts, err := d.ListAllChartTitles()
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
		if err == nil {
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
		chart_id_str := c.Param("cid")
		chart_id, err := strconv.Atoi(chart_id_str)
		if err != nil {
			c.Data(400, "text/plain", []byte("The chart ID must be an integer!"))
			c.Done()
			return
		}

		err = d.DeleteChart(chart_id)
		if err != nil {
			c.Data(500, "text/plain", []byte("Could not delete this chart!"))
			c.Done()
			return
		}

		c.Status(200)
		c.Done()
	}
}

// / Chart Series
func GetChartSeriesListHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		chart_id_str := c.Param("cid")
		chart_id, err := strconv.Atoi(chart_id_str)
		if err != nil {
			c.Data(400, "text/plain", []byte("The chart ID must be an integer!"))
			c.Done()
			return
		}

		series, err := d.ListChartSeries(chart_id)
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		data, err := json.Marshal(series)
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		c.Data(200, "application/json", data)
		c.Done()
	}
}

func GetChartSingleSeriesHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		series_id_str := c.Param("sid")
		series_id, err := strconv.Atoi(series_id_str)
		if err != nil {
			c.Data(400, "text/plain", []byte("The series ID must be an integer!"))
			c.Done()
			return
		}

		series, err := d.GetChartSeries(series_id)
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		data, err := json.Marshal(series)
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		c.Data(200, "application/json", data)
		c.Done()
	}
}

func AddChartSeriesHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		chart_id_str := c.Param("cid")
		chart_id, err := strconv.Atoi(chart_id_str)
		if err != nil {
			c.Data(400, "text/plain", []byte("The chart ID must be an integer!"))
			c.Done()
			return
		}

		var series database.ChartSeries
		if err := c.ShouldBindJSON(&series); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		series.ChartID = chart_id

		id, err := d.InsertChartSeries(series)
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		c.Data(200, "text/plain", []byte(strconv.Itoa(id)))
		c.Done()
	}
}

func DeleteChartSingleSeriesHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		series_id_str := c.Param("sid")
		series_id, err := strconv.Atoi(series_id_str)
		if err != nil {
			c.Data(400, "text/plain", []byte("The series ID must be an integer!"))
			c.Done()
			return
		}

		err = d.DeleteChartSeries(series_id)
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		c.Status(200)
		c.Done()
	}
}

func DeleteAllChartSeriesHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		chart_id_str := c.Param("cid")
		chart_id, err := strconv.Atoi(chart_id_str)
		if err != nil {
			c.Data(400, "text/plain", []byte("The chart ID must be an integer!"))
			c.Done()
			return
		}

		err = d.DeleteAllChartSeries(chart_id)
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		c.Status(200)
		c.Done()
	}
}

// / Dashboard itself
func GetDashboardListHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		dashboards, err := d.ListAllDashboards()
		if err != nil {
			c.Data(500, "text/plain", []byte("Internal server error"))
			c.Done()
			return
		}
		data, err := json.Marshal(dashboards)
		if err != nil {
			c.Data(500, "text/plain", []byte("Internal server error"))
			c.Done()
			return
		}
		c.Data(200, "application/json", data)
		c.Done()
	}
}

func GetDashboardHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		dash_id_str := c.Param("did")
		dash_id, err := strconv.Atoi(dash_id_str)
		if err != nil {
			c.Data(400, "text/plain", []byte("The dashboard ID must be an integer!"))
			c.Done()
			return
		}
		dash, err := d.GetDashboard(dash_id)
		if err != nil {
			c.Data(404, "text/plain", []byte("Dashboard not found"))
			c.Done()
			return
		}
		data, err := json.Marshal(dash)
		if err != nil {
			c.Data(500, "text/plain", []byte("Internal server error"))
			c.Done()
			return
		}
		c.Data(200, "application/json", data)
		c.Done()
	}
}

func CreateDashboardPageHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dash database.Dashboard
		if err := c.ShouldBindJSON(&dash); err != nil {
			c.Data(400, "text/plain", []byte("Invalid dashboard data"))
			return
		}
		id, err := d.InsertDashboard(dash)
		if err != nil {
			c.Data(500, "text/plain", []byte("Internal server error"))
			c.Done()
			return
		}
		c.Data(200, "text/plain", []byte(strconv.Itoa(id)))
		c.Done()
	}
}

func DeleteDashboardHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		dash_id_str := c.Param("did")
		dash_id, err := strconv.Atoi(dash_id_str)
		if err != nil {
			c.Data(400, "text/plain", []byte("The dashboard ID must be an integer!"))
			c.Done()
			return
		}
		err = d.DeleteDashboard(dash_id)
		if err != nil {
			c.Data(500, "text/plain", []byte("Internal server error"))
			c.Done()
			return
		}
		c.Status(200)
		c.Done()
	}
}

func ListDashboardChartsHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		dash_id_str := c.Param("did")
		dash_id, err := strconv.Atoi(dash_id_str)
		if err != nil {
			c.Data(400, "text/plain", []byte("The dashboard ID must be an integer!"))
			c.Done()
			return
		}
		graphs, err := d.ListDashboardCharts(dash_id)
		if err != nil {
			c.Data(500, "text/plain", []byte("Internal server error"))
			c.Done()
			return
		}
		data, err := json.Marshal(graphs)
		if err != nil {
			c.Data(500, "text/plain", []byte("Internal server error"))
			c.Done()
			return
		}
		c.Data(200, "application/json", data)
		c.Done()
	}
}

func AppendChartToDashboardHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		dash_id_str := c.Param("did")
		dash_id, err := strconv.Atoi(dash_id_str)
		if err != nil {
			c.Data(400, "text/plain", []byte("The dashboard ID must be an integer!"))
			c.Done()
			return
		}
		chart_id_str := c.Param("cid")
		chart_id, err := strconv.Atoi(chart_id_str)
		if err != nil {
			c.Data(400, "text/plain", []byte("The chart ID must be an integer!"))
			c.Done()
			return
		}
		err = d.AppendChartToDashboard(dash_id, chart_id)
		if err != nil {
			c.Data(500, "text/plain", []byte("Internal server error"))
			c.Done()
			return
		}
		c.Status(200)
		c.Done()
	}
}

func RemoveChartFromDashboardHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		dash_id_str := c.Param("did")
		dash_id, err := strconv.Atoi(dash_id_str)
		if err != nil {
			c.Data(400, "text/plain", []byte("The dashboard ID must be an integer!"))
			c.Done()
			return
		}
		chart_id_str := c.Param("cid")
		chart_id, err := strconv.Atoi(chart_id_str)
		if err != nil {
			c.Data(400, "text/plain", []byte("The chart ID must be an integer!"))
			c.Done()
			return
		}
		err = d.RemoveChartFromDashboard(dash_id, chart_id)
		if err != nil {
			c.Data(500, "text/plain", []byte("Internal server error"))
			c.Done()
			return
		}
		c.Status(200)
		c.Done()
	}
}
