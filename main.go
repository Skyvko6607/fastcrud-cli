package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Skyvko6607/fastcrud/cli/internal/client"
	"github.com/Skyvko6607/fastcrud/cli/internal/codegen"
)

func main() {
	key := flag.String("key", "", "Access key ID (UUID) for authentication")
	lang := flag.String("lang", "", "Target language: go, csharp, typescript, java")
	output := flag.String("output", "./generated", "Output directory for generated code")
	baseURL := flag.String("url", "https://crud.fastcrud.dev", "Base URL of the FastCRUD API")
	flag.Parse()

	if *key == "" || *lang == "" {
		fmt.Fprintf(os.Stderr, "Usage: fastcrud-cli --key <access-key-id> --lang <language> [--output <dir>] [--url <base-url>]\n\n")
		fmt.Fprintf(os.Stderr, "Supported languages: %s\n", strings.Join(codegen.SupportedLanguages(), ", "))
		os.Exit(1)
	}

	gen := codegen.GetGenerator(*lang)
	if gen == nil {
		fmt.Fprintf(os.Stderr, "Unsupported language: %s\nSupported: %s\n", *lang, strings.Join(codegen.SupportedLanguages(), ", "))
		os.Exit(1)
	}

	fmt.Printf("Authenticating with key %s...\n", *key)
	c := client.New(*baseURL)
	token, err := c.Authenticate(*key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Authenticated.")

	fmt.Println("Fetching schema...")
	tables, err := c.FetchSchema(token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d table(s).\n", len(tables))

	if len(tables) == 0 {
		fmt.Println("No tables found. Make sure your database has been introspected.")
		os.Exit(0)
	}

	fmt.Printf("Generating %s code in %s...\n", gen.Language(), *output)
	if err := gen.Generate(tables, *output); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating code: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Done.")
}
