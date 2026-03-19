package codegen

import (
	"os"
	"path/filepath"

	"github.com/Skyvko6607/fastcrud/cli/internal/schema"
)

type TypeScriptGenerator struct{}

func (g *TypeScriptGenerator) Language() string { return "typescript" }

func (g *TypeScriptGenerator) Generate(tables []schema.Table, outputDir string) error {
	typeMapper := func(dt string) string {
		switch MapType(dt) {
		case "int", "float":
			return "number"
		case "bool":
			return "boolean"
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

	modelsPath := filepath.Join(outputDir, "models.ts")
	if err := writeTemplate(modelsPath, tsModelsTemplate, data); err != nil {
		return err
	}

	clientPath := filepath.Join(outputDir, "client.ts")
	return writeTemplate(clientPath, tsClientTemplate, data)
}

const tsModelsTemplate = `{{range .Tables}}export interface {{.PascalName}} {
{{- range .Columns}}
  {{.CamelName}}: {{.LangType}};
{{- end}}
}

{{end}}`

const tsClientTemplate = `{{range .Tables}}import type { {{.PascalName}} } from "./models";
{{end}}
export class FastCrudClient {
  private baseUrl: string;
  private accessKeyId: string;
  private token = "";

  constructor(baseUrl: string, accessKeyId: string) {
    this.baseUrl = baseUrl.replace(/\/+$/, "");
    this.accessKeyId = accessKeyId;
  }

  async authenticate(): Promise<void> {
    const resp = await fetch(
      ` + "`${this.baseUrl}/authenticate/crud/${this.accessKeyId}`" + `,
      { method: "POST" }
    );
    if (!resp.ok) throw new Error(` + "`Auth failed: ${resp.status}`" + `);
    const data = await resp.json();
    this.token = data.access_token;
  }

  private async request<T>(method: string, path: string, body?: unknown): Promise<T> {
    const opts: RequestInit = {
      method,
      headers: {
        Authorization: this.token,
        "Content-Type": "application/json",
      },
    };
    if (body !== undefined) {
      opts.body = JSON.stringify(body);
    }
    const resp = await fetch(this.baseUrl + path, opts);
    const text = await resp.text();
    if (!resp.ok) throw new Error(` + "`API error (${resp.status}): ${text}`" + `);
    return JSON.parse(text);
  }

  private buildQuery(filter?: string, limit?: number, offset?: number): string {
    const params = new URLSearchParams();
    if (filter) params.set("filter", filter);
    if (limit !== undefined && limit > 0) params.set("limit", String(limit));
    if (offset !== undefined && offset > 0) params.set("offset", String(offset));
    const qs = params.toString();
    return qs ? "?" + qs : "";
  }
{{range .Tables}}
  async query{{.PascalName}}(filter?: string, limit = 1000, offset = 0): Promise<{{.PascalName}}[]> {
    return this.request("GET", "/crud/{{.RawName}}" + this.buildQuery(filter, limit, offset));
  }

  async insert{{.PascalName}}(rows: {{.PascalName}}[]): Promise<number> {
    const res = await this.request<{ rowsInserted: number }>("POST", "/crud/{{.RawName}}", rows);
    return res.rowsInserted;
  }

  async update{{.PascalName}}(data: Partial<{{.PascalName}}>, filter?: string): Promise<number> {
    const res = await this.request<{ rowsAffected: number }>(
      "PUT", "/crud/{{.RawName}}" + this.buildQuery(filter), data
    );
    return res.rowsAffected;
  }

  async delete{{.PascalName}}(filter?: string): Promise<number> {
    const res = await this.request<{ rowsAffected: number }>(
      "DELETE", "/crud/{{.RawName}}" + this.buildQuery(filter)
    );
    return res.rowsAffected;
  }
{{end}}}`
