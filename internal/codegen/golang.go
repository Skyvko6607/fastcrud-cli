package codegen

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/Skyvko6607/fastcrud/cli/internal/schema"
)

type GoGenerator struct{}

func (g *GoGenerator) Language() string { return "go" }

func (g *GoGenerator) Generate(tables []schema.Table, outputDir string) error {
	typeMapper := func(dt string) string {
		switch MapType(dt) {
		case "int":
			return "int"
		case "float":
			return "float64"
		case "bool":
			return "bool"
		default:
			return "string"
		}
	}

	data := struct {
		Tables []TableData
	}{
		Tables: BuildTableData(tables, typeMapper),
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	modelsPath := filepath.Join(outputDir, "models.go")
	if err := writeTemplate(modelsPath, goModelsTemplate, data); err != nil {
		return err
	}

	clientPath := filepath.Join(outputDir, "client.go")
	return writeTemplate(clientPath, goClientTemplate, data)
}

func writeTemplate(path, tmplStr string, data any) error {
	t, err := template.New("").Parse(tmplStr)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, data)
}

const goModelsTemplate = `package fastcrud

{{range .Tables}}
type {{.PascalName}} struct {
{{- range .Columns}}
	{{.PascalName}} {{.LangType}} ` + "`" + `json:"{{.RawName}}"` + "`" + `
{{- end}}
}
{{end}}`

const goClientTemplate = `package fastcrud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	BaseURL     string
	AccessKeyID string
	Token       string
	HTTP        *http.Client
}

func NewClient(baseURL, accessKeyID string) *Client {
	return &Client{
		BaseURL:     baseURL,
		AccessKeyID: accessKeyID,
		HTTP:        &http.Client{},
	}
}

func (c *Client) Authenticate() error {
	resp, err := c.HTTP.Post(
		fmt.Sprintf("%s/authenticate/crud/%s", c.BaseURL, c.AccessKeyID),
		"application/json", nil,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string ` + "`" + `json:"access_token"` + "`" + `
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	c.Token = result.AccessToken
	return nil
}

func (c *Client) doRequest(method, path string, body any) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.Token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(data))
	}
	return data, nil
}

func buildQuery(filter string, limit, offset int) string {
	params := url.Values{}
	if filter != "" {
		params.Set("filter", filter)
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		params.Set("offset", strconv.Itoa(offset))
	}
	if len(params) == 0 {
		return ""
	}
	return "?" + params.Encode()
}
{{range .Tables}}
func (c *Client) Query{{.PascalName}}(filter string, limit, offset int) ([]{{.PascalName}}, error) {
	data, err := c.doRequest("GET", "/crud/{{.RawName}}"+buildQuery(filter, limit, offset), nil)
	if err != nil {
		return nil, err
	}
	var rows []{{.PascalName}}
	return rows, json.Unmarshal(data, &rows)
}

func (c *Client) Insert{{.PascalName}}(rows []{{.PascalName}}) (int, error) {
	data, err := c.doRequest("POST", "/crud/{{.RawName}}", rows)
	if err != nil {
		return 0, err
	}
	var result struct{ RowsInserted int ` + "`" + `json:"rowsInserted"` + "`" + ` }
	return result.RowsInserted, json.Unmarshal(data, &result)
}

func (c *Client) Update{{.PascalName}}(row {{.PascalName}}, filter string) (int, error) {
	data, err := c.doRequest("PUT", "/crud/{{.RawName}}"+buildQuery(filter, 0, 0), row)
	if err != nil {
		return 0, err
	}
	var result struct{ RowsAffected int ` + "`" + `json:"rowsAffected"` + "`" + ` }
	return result.RowsAffected, json.Unmarshal(data, &result)
}

func (c *Client) Delete{{.PascalName}}(filter string) (int, error) {
	data, err := c.doRequest("DELETE", "/crud/{{.RawName}}"+buildQuery(filter, 0, 0), nil)
	if err != nil {
		return 0, err
	}
	var result struct{ RowsAffected int ` + "`" + `json:"rowsAffected"` + "`" + ` }
	return result.RowsAffected, json.Unmarshal(data, &result)
}
{{end}}`
