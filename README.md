# FastCRUD CLI

FastCRUD CLI is a command-line tool that automates the generation of strongly-typed CRUD clients for your database schema. It works seamlessly with the FastCRUD API to introspect your database and generate code in your preferred programming language.

## Features

- **Multi-language Support**: Generate code for Go, C#, Java, and TypeScript.
- **Automated Schema Mapping**: Automatically maps your database tables and columns to language-specific structures.
- **Secure Authentication**: Uses UUID access keys for secure interaction with the FastCRUD API.
- **Customizable Output**: Specify your own output directory and API base URL.

## Installation

To install the FastCRUD CLI, you can use `go install`:

```bash
go install github.com/Skyvko6607/fastcrud/cli@latest
```

Alternatively, you can build it from source:

```bash
cd cli
go build -o fastcrud-cli
```

## Usage

The CLI requires an access key and a target language to generate code.

```bash
fastcrud-cli --key <your-access-key-id> --lang <language> [--output <dir>] [--url <base-url>]
```

### Options

| Flag | Description | Default |
|------|-------------|---------|
| `--key` | **Required.** Access key ID (UUID) for authentication. | |
| `--lang` | **Required.** Target language: `go`, `csharp`, `typescript`, `java`. | |
| `--output` | Output directory for generated code. | `./generated` |
| `--url` | Base URL of the FastCRUD API. | `https://crud.fastcrud.dev` |

### Example

Generate Go code in the `./pkg/client` directory:

```bash
fastcrud-cli --key 123e4567-e89b-12d3-a456-426614174000 --lang go --output ./pkg/client
```

## Supported Languages

- **Go**: Native Go structures and database/sql compatible code.
- **C#**: Strongly-typed classes for .NET environments.
- **Java**: POJOs and repository patterns.
- **TypeScript**: Interfaces and fetch-based clients for web development.

## License

[License Information - Replace with appropriate license if applicable]
