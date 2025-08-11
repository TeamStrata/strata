package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/TeamStrata/strata/pkg/database"
	"github.com/gin-gonic/gin"
)

// / Charts
func (server *Server) GetChartListHandler(c *gin.Context) {
	charts, err := server.Db.ListAllCharts()
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

func (server *Server) GetChartHandler(c *gin.Context) {
	chart_id_str := c.Param("cid")
	chart_id, err := strconv.Atoi(chart_id_str)
	if err != nil {
		c.Data(400, "text/plain", []byte("The chart ID must be an integer!"))
		c.Done()
		return
	}

	chart, err := server.Db.GetChart(chart_id)
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

func (server *Server) CreateChartHandler(c *gin.Context) {
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
	id, err := server.Db.InsertChart(chart)
	if err != nil {
		c.Data(500, "text/plain", []byte(err.Error()))
		c.Done()
		return
	}

	c.Data(200, "text/plain", []byte(strconv.Itoa(id)))
	c.Done()
}

func (server *Server) DeleteChartHandler(c *gin.Context) {
	chart_id_str := c.Param("cid")
	chart_id, err := strconv.Atoi(chart_id_str)
	if err != nil {
		c.Data(400, "text/plain", []byte("The chart ID must be an integer!"))
		c.Done()
		return
	}

	err = server.Db.DeleteChart(chart_id)
	if err != nil {
		c.Data(500, "text/plain", []byte(err.Error()))
		c.Done()
		return
	}

	c.Status(200)
	c.Done()
}

func (server *Server) UpdateChartHandler(c *gin.Context) {
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

	err = server.Db.UpdateChart(chart)
	if err != nil {
		c.Data(500, "text/plain", []byte(err.Error()))
		c.Done()
		return
	}

	c.Status(200)
	c.Done()
}

// / Chart Series
func (server *Server) GetChartSeriesListHandler(c *gin.Context) {
	chart_id_str := c.Param("cid")
	chart_id, err := strconv.Atoi(chart_id_str)
	if err != nil {
		c.Data(500, "text/plain", []byte(err.Error()))
		c.Done()
		return
	}

	series, err := server.Db.ListChartSeries(chart_id)
	if err != nil {
		c.Data(500, "text/plain", []byte(err.Error()))
		c.Done()
		return
	}

	c.JSON(200, series)
}

func (server *Server) GetChartSingleSeriesHandler(c *gin.Context) {
	series_id_str := c.Param("sid")
	series_id, err := strconv.Atoi(series_id_str)
	if err != nil {
		c.Data(400, "text/plain", []byte("The series ID must be an integer!"))
		c.Done()
		return
	}

	series, err := server.Db.GetChartSeries(series_id)
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

func (server *Server) AddChartSeriesHandler(c *gin.Context) {
	var series []database.ChartSeries
	if err := c.ShouldBindJSON(&series); err != nil {
		c.Data(http.StatusBadRequest, "text/plain", []byte("ShouldBindJSON Error"))
		return
	}

	id, err := server.Db.InsertChartSeries(series)
	if err != nil {
		c.Data(500, "text/plain", []byte(err.Error()))
		c.Done()
		return
	}

	c.Data(200, "text/plain", []byte(strconv.Itoa(id)))
	c.Done()
}

func (server *Server) DeleteChartSingleSeriesHandler(c *gin.Context) {
	series_id_str := c.Param("sid")
	series_id, err := strconv.Atoi(series_id_str)
	if err != nil {
		c.Data(400, "text/plain", []byte("The series ID must be an integer!"))
		c.Done()
		return
	}

	err = server.Db.DeleteChartSeries(series_id)
	if err != nil {
		c.Data(500, "text/plain", []byte(err.Error()))
		c.Done()
		return
	}

	c.Status(200)
	c.Done()
}

func (server *Server) UpdateChartSeriesHandler(c *gin.Context) {
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

	err = server.Db.UpdateChartSeries(series)
	if err != nil {
		c.Data(500, "text/plain", []byte(err.Error()))
		c.Done()
		return
	}

	c.Status(200)
	c.Done()
}

func (server *Server) DeleteAllChartSeriesHandler(c *gin.Context) {
	chart_id_str := c.Param("cid")
	chart_id, err := strconv.Atoi(chart_id_str)
	if err != nil {
		c.Data(400, "text/plain", []byte("The chart ID must be an integer!"))
		c.Done()
		return
	}

	err = server.Db.DeleteAllChartSeries(chart_id)
	if err != nil {
		c.Data(500, "text/plain", []byte(err.Error()))
		c.Done()
		return
	}

	c.Status(200)
	c.Done()
}

// Dashboard itself
func (server *Server) GetDashboardListHandler(c *gin.Context) {
	userId, err := c.Cookie(uuidTag)
	if err != nil {
		c.String(http.StatusBadRequest, "missing uuid cookie")
		return
	}

	user, exists := server.ActiveUsers[userId]
	if !exists {
		c.String(http.StatusForbidden, "invalid uuid cookie")
	}

	dashboards, err := server.Db.GetUserDashboardPermissions(user.Id)
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

func (server *Server) GetDashboardHandler(c *gin.Context) {
	dash_id_str := c.Param("did")
	dash_id, err := strconv.Atoi(dash_id_str)
	if err != nil {
		c.Data(400, "text/plain", []byte("The dashboard ID must be an integer!"))
		c.Done()
		return
	}
	dash, err := server.Db.GetDashboard(dash_id)
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

func (server *Server) CreateDashboardPageHandler(c *gin.Context) {
	var dash database.Dashboard
	if err := c.ShouldBindJSON(&dash); err != nil {
		c.Data(400, "text/plain", []byte("Invalid dashboard data"))
		return
	}
	id, err := server.Db.InsertDashboard(dash)
	if err != nil {
		c.Data(500, "text/plain", []byte("Internal server error"))
		c.Done()
		return
	}
	c.Data(200, "text/plain", []byte(strconv.Itoa(id)))
	c.Done()
}

func (server *Server) DeleteDashboardHandler(c *gin.Context) {
	dash_id_str := c.Param("did")
	dash_id, err := strconv.Atoi(dash_id_str)
	if err != nil {
		c.Data(400, "text/plain", []byte("The dashboard ID must be an integer!"))
		c.Done()
		return
	}
	err = server.Db.DeleteDashboard(dash_id)
	if err != nil {
		c.Data(500, "text/plain", []byte("Internal server error"))
		c.Done()
		return
	}
	c.Status(200)
	c.Done()
}

func (server *Server) ListDashboardChartsHandler(c *gin.Context) {
	dash_id_str := c.Param("did")
	dash_id, err := strconv.Atoi(dash_id_str)
	if err != nil {
		c.Data(400, "text/plain", []byte("The dashboard ID must be an integer!"))
		c.Done()
		return
	}
	charts, err := server.Db.ListDashboardCharts(dash_id)
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

func (server *Server) AppendChartToDashboardHandler(c *gin.Context) {
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

	err = server.Db.AppendChartToDashboard(dash_id, chart_id, body.SizeX, body.SizeY)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.Status(200)
}

func (server *Server) RemoveChartFromDashboardHandler(c *gin.Context) {
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
	err = server.Db.RemoveChartFromDashboard(dash_id, chart_id)
	if err != nil {
		c.Data(500, "text/plain", []byte("Internal server error"))
		c.Done()
		return
	}
	c.Status(200)
	c.Done()
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
