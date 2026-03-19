package codegen

import (
	"os"
	"path/filepath"

	"github.com/Skyvko6607/fastcrud/cli/internal/schema"
)

type JavaGenerator struct{}

func (g *JavaGenerator) Language() string { return "java" }

func (g *JavaGenerator) Generate(tables []schema.Table, outputDir string) error {
	typeMapper := func(dt string) string {
		switch MapType(dt) {
		case "int":
			return "int"
		case "float":
			return "double"
		case "bool":
			return "boolean"
		case "datetime":
			return "String"
		default:
			return "String"
		}
	}

	tableData := BuildTableData(tables, typeMapper)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	for _, td := range tableData {
		modelPath := filepath.Join(outputDir, td.PascalName+".java")
		modelData := struct{ Table TableData }{Table: td}
		if err := writeTemplate(modelPath, javaModelTemplate, modelData); err != nil {
			return err
		}
	}

	clientPath := filepath.Join(outputDir, "FastCrudClient.java")
	data := struct{ Tables []TableData }{Tables: tableData}
	return writeTemplate(clientPath, javaClientTemplate, data)
}

const javaModelTemplate = `package fastcrud;

import com.google.gson.annotations.SerializedName;

public class {{.Table.PascalName}} {
{{- range .Table.Columns}}
    @SerializedName("{{.RawName}}")
    private {{.LangType}} {{.CamelName}};
{{- end}}

    public {{.Table.PascalName}}() {}
{{range .Table.Columns}}
    public {{.LangType}} get{{.PascalName}}() { return {{.CamelName}}; }
    public void set{{.PascalName}}({{.LangType}} {{.CamelName}}) { this.{{.CamelName}} = {{.CamelName}}; }
{{end}}}
`

const javaClientTemplate = `package fastcrud;

import com.google.gson.Gson;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.google.gson.reflect.TypeToken;

import java.io.IOException;
import java.net.URI;
import java.net.URLEncoder;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;

public class FastCrudClient {
    private final HttpClient http;
    private final Gson gson;
    private final String baseUrl;
    private final String accessKeyId;
    private String token = "";

    public FastCrudClient(String baseUrl, String accessKeyId) {
        this.baseUrl = baseUrl.replaceAll("/+$", "");
        this.accessKeyId = accessKeyId;
        this.http = HttpClient.newHttpClient();
        this.gson = new Gson();
    }

    public void authenticate() throws IOException, InterruptedException {
        HttpRequest req = HttpRequest.newBuilder()
            .uri(URI.create(baseUrl + "/authenticate/crud/" + accessKeyId))
            .POST(HttpRequest.BodyPublishers.noBody())
            .build();
        HttpResponse<String> resp = http.send(req, HttpResponse.BodyHandlers.ofString());
        if (resp.statusCode() != 200)
            throw new IOException("Auth failed (" + resp.statusCode() + "): " + resp.body());
        JsonObject json = JsonParser.parseString(resp.body()).getAsJsonObject();
        token = json.get("access_token").getAsString();
    }

    private String request(String method, String path, String body) throws IOException, InterruptedException {
        HttpRequest.Builder builder = HttpRequest.newBuilder()
            .uri(URI.create(baseUrl + path))
            .header("Authorization", token)
            .header("Content-Type", "application/json");

        switch (method) {
            case "GET":    builder.GET(); break;
            case "DELETE": builder.DELETE(); break;
            case "POST":   builder.POST(HttpRequest.BodyPublishers.ofString(body != null ? body : "")); break;
            case "PUT":    builder.PUT(HttpRequest.BodyPublishers.ofString(body != null ? body : "")); break;
        }

        HttpResponse<String> resp = http.send(builder.build(), HttpResponse.BodyHandlers.ofString());
        if (resp.statusCode() >= 400)
            throw new IOException("API error (" + resp.statusCode() + "): " + resp.body());
        return resp.body();
    }

    private String buildQuery(String filter, int limit, int offset) {
        List<String> parts = new ArrayList<>();
        if (filter != null && !filter.isEmpty())
            parts.add("filter=" + URLEncoder.encode(filter, StandardCharsets.UTF_8));
        if (limit > 0) parts.add("limit=" + limit);
        if (offset > 0) parts.add("offset=" + offset);
        return parts.isEmpty() ? "" : "?" + String.join("&", parts);
    }
{{range .Tables}}
    public List<{{.PascalName}}> query{{.PascalName}}(String filter, int limit, int offset)
            throws IOException, InterruptedException {
        String json = request("GET", "/crud/{{.RawName}}" + buildQuery(filter, limit, offset), null);
        return gson.fromJson(json, new TypeToken<List<{{.PascalName}}>>(){}.getType());
    }

    public int insert{{.PascalName}}(List<{{.PascalName}}> rows) throws IOException, InterruptedException {
        String json = request("POST", "/crud/{{.RawName}}", gson.toJson(rows));
        return JsonParser.parseString(json).getAsJsonObject().get("rowsInserted").getAsInt();
    }

    public int update{{.PascalName}}({{.PascalName}} data, String filter) throws IOException, InterruptedException {
        String json = request("PUT", "/crud/{{.RawName}}" + buildQuery(filter, 0, 0), gson.toJson(data));
        return JsonParser.parseString(json).getAsJsonObject().get("rowsAffected").getAsInt();
    }

    public int delete{{.PascalName}}(String filter) throws IOException, InterruptedException {
        String json = request("DELETE", "/crud/{{.RawName}}" + buildQuery(filter, 0, 0), null);
        return JsonParser.parseString(json).getAsJsonObject().get("rowsAffected").getAsInt();
    }
{{end}}}
`
