package codegen

import (
	"strings"
	"unicode"

	"github.com/Skyvko6607/fastcrud/cli/internal/schema"
)

type Generator interface {
	Generate(tables []schema.Table, outputDir string) error
	Language() string
}

type TableData struct {
	RawName    string
	PascalName string
	CamelName  string
	Columns    []ColumnData
}

type ColumnData struct {
	RawName    string
	PascalName string
	CamelName  string
	LangType   string
}

func BuildTableData(tables []schema.Table, typeMapper func(string) string) []TableData {
	var result []TableData
	for _, t := range tables {
		td := TableData{
			RawName:    t.Name,
			PascalName: ToPascalCase(t.Name),
			CamelName:  ToCamelCase(t.Name),
		}
		for _, c := range t.Columns {
			td.Columns = append(td.Columns, ColumnData{
				RawName:    c.Name,
				PascalName: ToPascalCase(c.Name),
				CamelName:  ToCamelCase(c.Name),
				LangType:   typeMapper(c.DataType),
			})
		}
		result = append(result, td)
	}
	return result
}

func ToPascalCase(s string) string {
	var sb strings.Builder
	capitalize := true
	for _, r := range s {
		if r == '_' || r == '-' || r == '.' || r == ' ' {
			capitalize = true
			continue
		}
		if capitalize {
			sb.WriteRune(unicode.ToUpper(r))
			capitalize = false
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func ToCamelCase(s string) string {
	p := ToPascalCase(s)
	if p == "" {
		return ""
	}
	runes := []rune(p)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func MapType(dataType string) string {
	dt := strings.ToLower(dataType)
	switch {
	case strings.Contains(dt, "int"):
		return "int"
	case strings.Contains(dt, "float") || strings.Contains(dt, "numeric") ||
		strings.Contains(dt, "decimal") || strings.Contains(dt, "double") || strings.Contains(dt, "real"):
		return "float"
	case strings.Contains(dt, "bool"):
		return "bool"
	case strings.Contains(dt, "timestamp") || strings.Contains(dt, "date") || strings.Contains(dt, "time"):
		return "datetime"
	default:
		return "string"
	}
}

func GetGenerator(lang string) Generator {
	switch strings.ToLower(lang) {
	case "go", "golang":
		return &GoGenerator{}
	case "csharp", "c#", "cs":
		return &CSharpGenerator{}
	case "typescript", "ts", "node", "nodejs":
		return &TypeScriptGenerator{}
	case "java":
		return &JavaGenerator{}
	default:
		return nil
	}
}

func SupportedLanguages() []string {
	return []string{"go", "csharp", "typescript", "java"}
}
