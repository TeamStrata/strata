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

func (d *DbManager) ListAllChartTitles() ([]Chart, error) {
	var list []Chart
	query := "SELECT chart_title FROM chart;"

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
		err := rows.Scan(&chart.Title)
		if err != nil {
			return nil, err
		}

		list = append(list, chart)
	}

	return list, nil
}

func (d *DbManager) InsertChart(chart Chart) (int, error) {
	chart_id := 0

	query := "INSERT INTO chart (chart_title, chart_type) VALUES ($1, $2) RETURNING chart_id;"
	err := d.Connection.QueryRow(d.context, query, chart.Title, chart.Type).Scan(&chart_id)

	return chart_id, err
}

func (d *DbManager) GetChart(chart_id int) (Chart, error) {
	var chart Chart

	query := "SELECT * FROM chart WHERE chart_id = $1;"
	err := d.Connection.QueryRow(d.context, query, chart_id).Scan(&chart)
	if err != nil {
		return Chart{}, err
	}

	return chart, nil
}
