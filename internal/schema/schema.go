package schema

type Table struct {
	ID        string   `json:"id"`
	Name      string   `json:"table"`
	ProjectID string   `json:"project_id"`
	Columns   []Column `json:"columns"`
}

type Column struct {
	Name     string `json:"column"`
	TableID  string `json:"table_id"`
	DataType string `json:"data_type"`
}
