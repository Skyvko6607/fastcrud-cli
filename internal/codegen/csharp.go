package codegen

import (
	"os"
	"path/filepath"

	"github.com/Skyvko6607/fastcrud/cli/internal/schema"
)

type CSharpGenerator struct{}

func (g *CSharpGenerator) Language() string { return "csharp" }

func (g *CSharpGenerator) Generate(tables []schema.Table, outputDir string) error {
	typeMapper := func(dt string) string {
		switch MapType(dt) {
		case "int":
			return "int"
		case "float":
			return "double"
		case "bool":
			return "bool"
		case "datetime":
			return "DateTime"
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

	modelsPath := filepath.Join(outputDir, "Models.cs")
	if err := writeTemplate(modelsPath, csModelsTemplate, data); err != nil {
		return err
	}

	clientPath := filepath.Join(outputDir, "FastCrudClient.cs")
	return writeTemplate(clientPath, csClientTemplate, data)
}

const csModelsTemplate = `using System.Text.Json.Serialization;

namespace FastCrud;
{{range .Tables}}
public class {{.PascalName}}
{
{{- range .Columns}}
    [JsonPropertyName("{{.RawName}}")]
    public {{.LangType}} {{.PascalName}} { get; set; }
{{- end}}
}
{{end}}`

const csClientTemplate = `using System.Net.Http.Headers;
using System.Text;
using System.Text.Json;

namespace FastCrud;

public class FastCrudClient
{
    private readonly HttpClient _http;
    private readonly string _baseUrl;
    private readonly string _accessKeyId;
    private string _token = "";

    public FastCrudClient(string baseUrl, string accessKeyId)
    {
        _baseUrl = baseUrl.TrimEnd('/');
        _accessKeyId = accessKeyId;
        _http = new HttpClient();
    }

    public async Task AuthenticateAsync()
    {
        var resp = await _http.PostAsync(
            $"{_baseUrl}/authenticate/crud/{_accessKeyId}", null);
        resp.EnsureSuccessStatusCode();
        var json = await resp.Content.ReadAsStringAsync();
        var doc = JsonDocument.Parse(json);
        _token = doc.RootElement.GetProperty("access_token").GetString()!;
    }

    private async Task<string> RequestAsync(HttpMethod method, string path, object? body = null)
    {
        var req = new HttpRequestMessage(method, _baseUrl + path);
        req.Headers.Add("Authorization", _token);
        if (body != null)
        {
            req.Content = new StringContent(
                JsonSerializer.Serialize(body), Encoding.UTF8, "application/json");
        }
        var resp = await _http.SendAsync(req);
        var text = await resp.Content.ReadAsStringAsync();
        if (!resp.IsSuccessStatusCode)
            throw new HttpRequestException($"API error ({(int)resp.StatusCode}): {text}");
        return text;
    }

    private static string BuildQuery(string? filter = null, int limit = 0, int offset = 0)
    {
        var parts = new List<string>();
        if (!string.IsNullOrEmpty(filter)) parts.Add($"filter={Uri.EscapeDataString(filter)}");
        if (limit > 0) parts.Add($"limit={limit}");
        if (offset > 0) parts.Add($"offset={offset}");
        return parts.Count > 0 ? "?" + string.Join("&", parts) : "";
    }
{{range .Tables}}
    public async Task<List<{{.PascalName}}>> Query{{.PascalName}}Async(string? filter = null, int limit = 1000, int offset = 0)
    {
        var json = await RequestAsync(HttpMethod.Get, "/crud/{{.RawName}}" + BuildQuery(filter, limit, offset));
        return JsonSerializer.Deserialize<List<{{.PascalName}}>>(json)!;
    }

    public async Task<int> Insert{{.PascalName}}Async(List<{{.PascalName}}> rows)
    {
        var json = await RequestAsync(HttpMethod.Post, "/crud/{{.RawName}}", rows);
        return JsonDocument.Parse(json).RootElement.GetProperty("rowsInserted").GetInt32();
    }

    public async Task<int> Update{{.PascalName}}Async({{.PascalName}} data, string? filter = null)
    {
        var json = await RequestAsync(HttpMethod.Put, "/crud/{{.RawName}}" + BuildQuery(filter), data);
        return JsonDocument.Parse(json).RootElement.GetProperty("rowsAffected").GetInt32();
    }

    public async Task<int> Delete{{.PascalName}}Async(string? filter = null)
    {
        var json = await RequestAsync(HttpMethod.Delete, "/crud/{{.RawName}}" + BuildQuery(filter));
        return JsonDocument.Parse(json).RootElement.GetProperty("rowsAffected").GetInt32();
    }
{{end}}}`
