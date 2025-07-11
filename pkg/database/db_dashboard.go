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
}

type ChartSeries struct {
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
	rows, err := d.Connection.Query(d.context, query)
	if err != nil {
		return nil, err
	}
	// Ensure the rows are closed properly
	defer rows.Close()

	// Iterate through rows
	for rows.Next() {
		var chart Chart

		// Scan the row for the query ID and String Literal
		err := rows.Scan(&chart.Id, &chart.Title, &chart.Type)
		if err != nil {
			return nil, err
		}

		list = append(list, chart)
	}

	return list, nil
}

func (d *DbManager) InsertChart(chart Chart) (int, error) {
	var chartID int

	//Check if the chart already exists
	queryCheck := "SELECT chart_id FROM chart WHERE chart_title = $1;"
	err := d.Connection.QueryRow(d.context, queryCheck, chart.Title).Scan(&chartID)

	if err == nil {
		queryUpdate := "UPDATE chart SET chart_type = $1 WHERE chart_id = $2;"
		_, updateErr := d.Connection.Exec(d.context, queryUpdate, chart.Type, chartID)
		return chartID, updateErr
	}

	//If no chart found (no rows), insert a new one
	if err.Error() == "no rows in result set" {
		queryInsert := "INSERT INTO chart (chart_title, chart_type) VALUES ($1, $2) RETURNING chart_id;"
		insertErr := d.Connection.QueryRow(d.context, queryInsert, chart.Title, chart.Type).Scan(&chartID)
		return chartID, insertErr
	}
	return 0, err
}

func (d *DbManager) GetChart(chart_id int) (Chart, error) {
	var chart Chart

	query := "SELECT chart_id, chart_title, chart_type FROM chart WHERE chart_id = $1;"
	err := d.Connection.QueryRow(d.context, query, chart_id).Scan(&chart.Id, &chart.Title, &chart.Type)
	if err != nil {
		return Chart{}, err
	}

	return chart, nil
}

func (d *DbManager) DeleteChart(chart_id int) error {
	query := "DELETE FROM chart WHERE chart_id = $1;"
	_, err := d.Connection.Exec(d.context, query, chart_id)
	return err
}

func (d *DbManager) UpdateChart(chart Chart) error {
    query := "UPDATE chart SET chart_title = $1, chart_type = $2 WHERE chart_id = $3;"
    _, err := d.Connection.Exec(d.context, query, chart.Title, chart.Type, chart.Id)
    return err
}

func (d *DbManager) ListChartSeries(chart_id int) ([]ChartSeries, error) {
	var list []ChartSeries
	query := "SELECT series_id, chart_id, query_id, x_column, y_column FROM chartSeries WHERE chart_id = $1;"
	rows, err := d.Connection.Query(d.context, query, chart_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s ChartSeries
		err := rows.Scan(&s.Id, &s.ChartID, &s.QueryID, &s.XCol, &s.YCol)
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
	err := d.Connection.QueryRow(d.context, query, series_id).Scan(&s.Id, &s.ChartID, &s.QueryID, &s.XCol, &s.YCol)
	return s, err
}

func (d *DbManager) InsertChartSeries(series ChartSeries) (int, error) {
	var id int
	query := "INSERT INTO chartSeries (chart_id, query_id, x_column, y_column) VALUES ($1, $2, $3, $4) RETURNING series_id;"
	err := d.Connection.QueryRow(d.context, query, series.ChartID, series.QueryID, series.XCol, series.YCol).Scan(&id)
	return id, err
}

func (d *DbManager) DeleteChartSeries(series_id int) error {
	query := "DELETE FROM chartSeries WHERE series_id = $1;"
	_, err := d.Connection.Exec(d.context, query, series_id)
	return err
}

func (d *DbManager) DeleteAllChartSeries(chart_id int) error {
	query := "DELETE FROM chartSeries WHERE chart_id = $1;"
	_, err := d.Connection.Exec(d.context, query, chart_id)
	return err
}

func (d* DbManager) UpdateChartSeries(series ChartSeries) error {
    query := "UPDATE chartSeries SET chart_id = $1, query_id = $2, x_column = $3, y_column = $4 WHERE series_id = $5;"
    _, err := d.Connection.Exec(d.context, query, series.ChartID, series.QueryID, series.XCol, series.YCol, series.Id)
    return err
}

// List all dashboards
func (d *DbManager) ListAllDashboards() ([]Dashboard, error) {
	var dashboards []Dashboard
	query := "SELECT dashboard_id, dashboard_title, dashboard_content FROM dashboards;"
	rows, err := d.Connection.Query(d.context, query)
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
	err := d.Connection.QueryRow(d.context, query, dashID).Scan(&dash.Id, &dash.Title, &dash.Content)
	return dash, err
}

// Insert a new dashboard, returns new dashboard ID
func (d *DbManager) InsertDashboard(dash Dashboard) (int, error) {
	var id int
	query := "INSERT INTO dashboards (dashboard_title, dashboard_content) VALUES ($1, $2) RETURNING dashboard_id;"
	err := d.Connection.QueryRow(d.context, query, dash.Title, dash.Content).Scan(&id)
	return id, err
}

// Delete a dashboard by ID
func (d *DbManager) DeleteDashboard(dashID int) error {
	query := "DELETE FROM dashboards WHERE dashboard_id = $1;"
	_, err := d.Connection.Exec(d.context, query, dashID)
	return err
}

// List all charts for a dashboard
func (d *DbManager) ListDashboardCharts(dashID int) ([]DashboardGraphs, error) {
	var graphs []DashboardGraphs
	query := "SELECT dash_id, chart_id, size_x, size_y, \"order\" FROM dashboard_graphs WHERE dash_id = $1 ORDER BY \"order\";"
	rows, err := d.Connection.Query(d.context, query, dashID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var graph DashboardGraphs
		if err := rows.Scan(&graph.DashId, &graph.ChartId, &graph.SizeX, &graph.SizeY, &graph.Order); err != nil {
			return nil, err
		}
		graphs = append(graphs, graph)
	}
	return graphs, nil
}

// Append a chart to a dashboard
func (d *DbManager) AppendChartToDashboard(dashID, chartID int) error {
	query := "INSERT INTO dashboard_graphs (dash_id, chart_id, size_x, size_y, \"order\") VALUES ($1, $2, 1, 1, (SELECT COALESCE(MAX(\"order\"), 0) + 1 FROM dashboard_graphs WHERE dash_id = $1));"
	_, err := d.Connection.Exec(d.context, query, dashID, chartID)
	return err
}

// Remove a chart from a dashboard
func (d *DbManager) RemoveChartFromDashboard(dashID, chartID int) error {
	query := "DELETE FROM dashboard_graphs WHERE dash_id = $1 AND chart_id = $2;"
	_, err := d.Connection.Exec(d.context, query, dashID, chartID)
	return err
}
