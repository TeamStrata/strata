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

		// Validate chart data
		if chart.Title == "" || chart.Type == "" || chart.Xname == "" || chart.Yname == "" {
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
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		c.Status(200)
		c.Done()
	}
}

func UpdateChartHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		chart_id_str := c.Param("cid")
		chart_id, err := strconv.Atoi(chart_id_str)
		if err != nil {
			c.Data(400, "text/plain", []byte("The chart ID must be an integer!"))
			c.Done()
			return
		}

		var chart database.Chart
		if err := c.ShouldBindJSON(&chart); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		chart.Id = chart_id

		err = d.UpdateChart(chart)
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
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
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		series, err := d.ListChartSeries(chart_id)
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		c.JSON(200, series)
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

		var series []database.ChartSeries
		if err := c.ShouldBindJSON(&series); err != nil {
			c.Data(http.StatusBadRequest, "text/plain", []byte("ShouldBindJSON Error"))
			return
		}

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

func UpdateChartSeriesHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		series_id_str := c.Param("sid")
		series_id, err := strconv.Atoi(series_id_str)
		if err != nil {
			c.Data(400, "text/plain", []byte("The series ID must be an integer!"))
			c.Done()
			return
		}

		var series database.ChartSeries
		if err := c.ShouldBindJSON(&series); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		series.Id = series_id

		err = d.UpdateChartSeries(series)
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

// Dashboard itself
func GetDashboardListHandler(d *database.DbManager, activeUsers map[string]UserSessionData) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := c.Cookie(uuidTag)
		if err != nil {
			c.String(http.StatusBadRequest, "missing uuid cookie")
			return
		}

		user, exists := activeUsers[userId]
		if !exists {
			c.String(http.StatusForbidden, "invalid uuid cookie")
		}

		dashboards, err := d.GetUserDashboardPermissions(user.Id)
		if err != nil {
			c.Data(500, "text/plain", []byte("Internal server error:"+err.Error()))
			c.Done()
			return
		}
		data, err := json.Marshal(dashboards)
		if err != nil {
			c.Data(500, "text/plain", []byte("Internal server error"+err.Error()))
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
		charts, err := d.ListDashboardCharts(dash_id)
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

func AppendChartToDashboardHandler(d *database.DbManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		dash_id, err := strconv.Atoi(c.Param("did"))
		if err != nil {
			c.String(400, "The dashboard ID must be an integer!")
			return
		}
		chart_id, err := strconv.Atoi(c.Param("cid"))
		if err != nil {
			c.String(400, "The chart ID must be an integer!")
			return
		}

		var body struct {
			SizeX int `json:"size_x"`
			SizeY int `json:"size_y"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.String(400, "Invalid JSON payload")
			return
		}

		err = d.AppendChartToDashboard(dash_id, chart_id, body.SizeX, body.SizeY)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		c.Status(200)
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

func GetFullDashboardData(d *database.DbManager, cdb **database.DbManager) gin.HandlerFunc {
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
		charts, err := d.ListDashboardCharts(dash_id)
		if err != nil {
			c.Data(500, "text/plain", []byte(err.Error()))
			c.Done()
			return
		}

		type SeriesRow = map[string]any

		type SeriesWithData struct {
			database.ChartSeries             // embed existing fields (id, name, query_id, etc.)
			Data                 []SeriesRow `json:"data,omitempty"`  // attach query results here
			Error                string      `json:"error,omitempty"` // optional per-series error
		}

		// For each chart, fetch its series and axis names
		type ChartWithSeries struct {
			ChartId int              `json:"chart_id"`
			SizeX   int              `json:"size_x"`
			SizeY   int              `json:"size_y"`
			Order   int              `json:"order"`
			Xname   string           `json:"xname"`
			Yname   string           `json:"yname"`
			Title   string           `json:"title"`
			Type    string           `json:"type"`
			Series  []SeriesWithData `json:"series"`
		}

		var chartsWithSeries []ChartWithSeries
		for _, chart := range charts {
			series, err := d.ListChartSeries(chart.ChartId)
			if err != nil {
				c.Data(500, "text/plain", []byte(err.Error()))
				c.Done()
				return
			}
			// Get chart details for axis names, title, type
			chartDetails, err := d.GetChart(chart.ChartId)
			if err != nil {
				c.Data(500, "text/plain", []byte(err.Error()))
				c.Done()
				return
			}

			swd := make([]SeriesWithData, len(series))
			for i, s := range series {
				swd[i] = SeriesWithData{ChartSeries: s}
				var queryErr error
				rawData, err := executeQuery(d, cdb, s.QueryID)
				if err != nil {
					c.Data(500, "text/plain", []byte(err.Error()))
					c.Done()
					return
				}
				seriesRows := make([]SeriesRow, len(rawData))
				for j, row := range rawData {
					seriesRows[j] = make(SeriesRow)
					for k, v := range row {
						seriesRows[j][k] = v
					}
				}
				swd[i].Data = seriesRows
				if queryErr != nil {
					c.Data(500, "text/plain", []byte(queryErr.Error()))
					c.Done()
					return
				}
			}

			chartsWithSeries = append(chartsWithSeries, ChartWithSeries{
				ChartId: chart.ChartId,
				SizeX:   chart.SizeX,
				SizeY:   chart.SizeY,
				Order:   chart.Order,
				Xname:   chartDetails.Xname,
				Yname:   chartDetails.Yname,
				Title:   chartDetails.Title,
				Type:    chartDetails.Type,
				Series:  swd,
			})
		}

		result := map[string]interface{}{
			"dashboard": dash,
			"charts":    chartsWithSeries,
		}
		data, err := json.Marshal(result)
		if err != nil {
			c.Data(500, "text/plain", []byte("Internal server error"))
			c.Done()
			return
		}
		c.Data(200, "application/json", data)
		c.Done()
	}
}
