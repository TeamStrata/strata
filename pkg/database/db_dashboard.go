package database

const (
	CT_Line    string = "line"
	CT_Area    string = "area"
	CT_Column  string = "column"
	CT_Bar     string = "bar"
	CT_Treemap string = "treemap"
	CT_Heatmap string = "heatmap"
	CT_Pie     string = "pie"
	CT_Radar   string = "radar"
	CT_Polar   string = "polar"
	CT_Scatter string = "scatter"
)

type Chart struct {
	Id    int    `json:"id,omitempty"`
	Title string `json:"title"`
	Type  string `json:"type"`
	Xname string `json:"x_axis"`
	Yname string `json:"y_axis"`
}

type ChartSeries struct {
	Name    string `json:"name"`
	Id      int    `json:"id,omitempty"`
	ChartID int    `json:"chart_id"`
	QueryID int    `json:"query_id"`
	XCol    string `json:"x_col_name"`
	YCol    string `json:"y_col_name"`
}

type Dashboard struct {
	Id      int    `json:"id,omitempty"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type DashboardGraphs struct {
	DashId  int `json:"dash_id"`
	ChartId int `json:"chart_id"`
	SizeX   int `json:"size_x"`
	SizeY   int `json:"size_y"`
	Order   int `json:"order"`
}

func (d *DbManager) ListAllCharts() ([]Chart, error) {
	var list []Chart
	query := "SELECT * FROM chart;"

	// Query the databse for all queries
	rows, err := d.Connection.Query(d.Context, query)
	if err != nil {
		return nil, err
	}
	// Ensure the rows are closed properly
	defer rows.Close()

	// Iterate through rows
	for rows.Next() {
		var chart Chart

		// Scan the row for the query ID and String Literal
		err := rows.Scan(&chart.Id, &chart.Title, &chart.Type, &chart.Xname, &chart.Yname)
		if err != nil {
			return nil, err
		}

		list = append(list, chart)
	}

	return list, nil
}

func (d *DbManager) InsertChart(chart Chart) (int, error) {

	//Check if the chart already exists
	queryCheck := "SELECT chart_id FROM chart WHERE chart_title = $1;"
	err := d.Connection.QueryRow(d.Context, queryCheck, chart.Title).Scan(&chart.Id)

	//If no chart found (no rows), insert a new one
	// if err.Error() == "no rows in result set" {
	if err != nil {
		queryInsert := "INSERT INTO chart (chart_title, chart_type, x_axis_name, y_axis_name) VALUES ($1, $2, $3, $4) RETURNING chart_id;"
		insertErr := d.Connection.QueryRow(d.Context, queryInsert, chart.Title, chart.Type, chart.Xname, chart.Yname).Scan(&chart.Id)
		return chart.Id, insertErr
	}

	updateErr := d.UpdateChart(chart)
	if updateErr != nil {
		return -1, updateErr
	}
	// if err == nil {
	// 	queryUpdate := "UPDATE chart SET chart_type = $1 WHERE chart_id = $2;"
	// 	_, updateErr := d.Connection.Exec(d.Context, queryUpdate, chart.Type, chartID)
	// 	return chartID, updateErr
	// }

	return chart.Id, err
}

func (d *DbManager) GetChart(chart_id int) (Chart, error) {
	var chart Chart

	query := "SELECT chart_id, chart_title, chart_type, x_axis_name, y_axis_name FROM chart WHERE chart_id = $1;"
	err := d.Connection.QueryRow(d.Context, query, chart_id).Scan(&chart.Id, &chart.Title, &chart.Type, &chart.Xname, &chart.Yname)
	if err != nil {
		return Chart{}, err
	}

	return chart, nil
}

func (d *DbManager) DeleteChart(chart_id int) error {
	query := "DELETE FROM chart WHERE chart_id = $1;"
	_, err := d.Connection.Exec(d.Context, query, chart_id)
	return err
}

func (d *DbManager) UpdateChart(chart Chart) error {
	query := "UPDATE chart SET chart_title = $1, chart_type = $2, x_axis_name = $3, y_axis_name = $4 WHERE chart_id = $5;"
	_, err := d.Connection.Exec(d.Context, query, chart.Title, chart.Type, chart.Xname, chart.Yname, chart.Id)
	return err
}

func (d *DbManager) ListChartSeries(chart_id int) ([]ChartSeries, error) {
	var list []ChartSeries
	query := "SELECT series_id, chart_id, query_id, x_column, y_column, seriesName FROM chartSeries WHERE chart_id = $1;"
	rows, err := d.Connection.Query(d.Context, query, chart_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s ChartSeries
		err := rows.Scan(&s.Id, &s.ChartID, &s.QueryID, &s.XCol, &s.YCol, &s.Name)
		if err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, nil
}

func (d *DbManager) GetChartSeries(series_id int) (ChartSeries, error) {
	var s ChartSeries
	query := "SELECT series_id, chart_id, query_id, x_column, y_column FROM chartSeries WHERE series_id = $1;"
	err := d.Connection.QueryRow(d.Context, query, series_id).Scan(&s.Id, &s.ChartID, &s.QueryID, &s.XCol, &s.YCol)
	return s, err
}

func (d *DbManager) InsertChartSeries(series []ChartSeries) (int, error) {
	var id int
	if len(series) > 0 {
		err := d.DeleteAllChartSeries(series[0].ChartID)
		if err != nil {
			return -1, err
		}
	}
	for _, s := range series {
		query := "INSERT INTO chartSeries (chart_id, query_id, x_column, y_column, seriesName) SELECT $1, $2, $3, $4, $5 WHERE NOT EXISTS (SELECT 1 FROM chartSeries WHERE series_id = $6) RETURNING series_id;"
		err := d.Connection.QueryRow(d.Context, query, s.ChartID, s.QueryID, s.XCol, s.YCol, s.Name, s.Id).Scan(&id)
		if err != nil {
			return -1, err
		}
	}
	return id, nil
}

func (d *DbManager) DeleteChartSeries(series_id int) error {
	query := "DELETE FROM chartSeries WHERE series_id = $1;"
	_, err := d.Connection.Exec(d.Context, query, series_id)
	return err
}

func (d *DbManager) DeleteAllChartSeries(chart_id int) error {
	query := "DELETE FROM chartSeries WHERE chart_id = $1;"
	_, err := d.Connection.Exec(d.Context, query, chart_id)
	return err
}

func (d *DbManager) UpdateChartSeries(series ChartSeries) error {
	query := "UPDATE chartSeries SET chart_id = $1, query_id = $2, x_column = $3, y_column = $4, seriesName = $5, WHERE series_id = $6;"
	_, err := d.Connection.Exec(d.Context, query, series.ChartID, series.QueryID, series.XCol, series.YCol, series.Name, series.Id)
	return err
}

// List all dashboards
func (d *DbManager) ListAllDashboards() ([]Dashboard, error) {
	var dashboards []Dashboard
	query := "SELECT dashboard_id, dashboard_title, dashboard_content FROM dashboards;"
	rows, err := d.Connection.Query(d.Context, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var dash Dashboard
		if err := rows.Scan(&dash.Id, &dash.Title, &dash.Content); err != nil {
			return nil, err
		}
		dashboards = append(dashboards, dash)
	}
	return dashboards, nil
}

// Get a single dashboard by ID
func (d *DbManager) GetDashboard(dashID int) (Dashboard, error) {
	var dash Dashboard
	query := "SELECT dashboard_id, dashboard_title, dashboard_content FROM dashboards WHERE dashboard_id = $1;"
	err := d.Connection.QueryRow(d.Context, query, dashID).Scan(&dash.Id, &dash.Title, &dash.Content)
	return dash, err
}

// Insert a new dashboard, returns new dashboard ID
func (d *DbManager) InsertDashboard(dash Dashboard) (int, error) {
	var id int
	query := "INSERT INTO dashboards (dashboard_title, dashboard_content) VALUES ($1, $2) RETURNING dashboard_id;"
	err := d.Connection.QueryRow(d.Context, query, dash.Title, dash.Content).Scan(&id)
	return id, err
}

// Delete a dashboard by ID
func (d *DbManager) DeleteDashboard(dashID int) error {
	query := "DELETE FROM dashboards WHERE dashboard_id = $1;"
	_, err := d.Connection.Exec(d.Context, query, dashID)
	return err
}

// List all charts for a dashboard
func (d *DbManager) ListDashboardCharts(dashID int) ([]DashboardGraphs, error) {
	var charts []DashboardGraphs
	query := "SELECT dashboard_id, chart_id, size_x, size_y, chart_order FROM dashboardGraphs WHERE dashboard_id = $1 ORDER BY chart_order;"
	rows, err := d.Connection.Query(d.Context, query, dashID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var chart DashboardGraphs
		if err := rows.Scan(&chart.DashId, &chart.ChartId, &chart.SizeX, &chart.SizeY, &chart.Order); err != nil {
			return nil, err
		}
		charts = append(charts, chart)
	}
	return charts, nil
}

func (d *DbManager) AppendChartToDashboard(dashID, chartID, sizeX, sizeY int) error {
	query := `
		INSERT INTO dashboardGraphs 
    (dashboard_id, chart_id, size_x, size_y, chart_order)
VALUES 
    ($1, $2, $3, $4, 
     COALESCE(
       (SELECT chart_order FROM dashboardGraphs WHERE dashboard_id = $1 AND chart_id = $2),
       (SELECT COALESCE(MAX(chart_order), 0) + 1 FROM dashboardGraphs WHERE dashboard_id = $1)
     )
)
ON CONFLICT (dashboard_id, chart_id)
DO UPDATE SET 
    size_x = EXCLUDED.size_x,
    size_y = EXCLUDED.size_y,
    chart_order = EXCLUDED.chart_order;
	`
	_, err := d.Connection.Exec(d.Context, query, dashID, chartID, sizeX, sizeY)
	return err
}

// Remove a chart from a dashboard
func (d *DbManager) RemoveChartFromDashboard(dashID, chartID int) error {
	query := "DELETE FROM dashboardGraphs WHERE dashboard_id = $1 AND chart_id = $2;"
	_, err := d.Connection.Exec(d.Context, query, dashID, chartID)
	return err
}
