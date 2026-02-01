# yamlvfs

[![Go Reference](https://pkg.go.dev/badge/github.com/seb-gan/yamlvfs.svg)](https://pkg.go.dev/github.com/seb-gan/yamlvfs)
[![Go Report Card](https://goreportcard.com/badge/github.com/seb-gan/yamlvfs)](https://goreportcard.com/report/github.com/seb-gan/yamlvfs)

`yamlvfs` is a schema, a Go library, and a CLI.

The schema specifies a format that allows describing a virtual filesystem using YAML.
The library provides functions for loading a yamlvfs document into a standard `fs.FS` implementation.
The CLI provides commands for importing and exporting yamlvfs documents.

## YAMLVFS Format

The yamlvfs format must comply with [`yamlvfs.schema.json`](yamlvfs.schema.json). The document is a YAML mapping where each key-value pair represents a filesystem entry:

- **Directory**: Key ends with `/`. Value is either null or a mapping containing nested entries.
- **File**: Key does not end with `/`. Value is either null or a string scalar.

```yaml
data/:
src/:
  version.txt: '1'
  main.go:
  utils.go: |
    // util functions
```

## Go Library

```bash
go get github.com/seb-gan/yamlvfs
```

```go
// Parse YAML, open as filesystem
node, _ := yamlvfs.ParseFile("data/sample.yml")
fsys, _ := yamlvfs.Open(node)

// Use standard fs.FS
data, _ := fs.ReadFile(fsys, "config.yml")
entries, _ := fs.ReadDir(fsys, "src")

// Capture filesystem to YAML
node, _ := yamlvfs.FromFS(os.DirFS("."), nil)
yaml := yamlvfs.Format(node)
```

## CLI

```bash
go install github.com/seb-gan/yamlvfs/cmd/yamlvfs@latest
```

```txt
yamlvfs - Work with yamlvfs YAML filesystems

yamlvfs is a CLI for working with YAML-defined virtual filesystems.

See https://github.com/seb-gan/yamlvfs for more information.

Usage:

yamlvfs

    from-dir # Create yamlvfs document from directory
        [--depth <int>]  # max depth (-1 = unlimited) (default: -1)
        --dir <string>  # directory to scan (required)
        [--exclude-dirs <string>]  # glob patterns to exclude (comma-separated)
        [--include-content <string>]  # glob patterns for file content (comma-separated) (default: *)
        [--include-dirs <string>]  # glob patterns for directories (comma-separated) (default: *)
        [--no-gitignore <bool>]  # do not skip .gitignore paths
        [--out <string>]  # output file (default: stdout)

    schema # Export or print the embedded schema

        export # Export JSON schema to file
            [--dest-dir <string>]  # destination directory
            [--dest-file <string>]  # destination file path

        print # Print JSON schema to stdout

    to-dir # Write yamlvfs document to directory
        --file <string>  # yamlvfs file (required)
        --out <string>  # output directory (required)

    tree # Print tree structure of yamlvfs file
        --file <string>  # yamlvfs file (required)

    validate # Validate yamlvfs file
        --file <string>  # yamlvfs file (required)


```

## License

MIT
