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

**Usage:**

```bash
go install github.com/seb-gan/yamlvfs/cmd/yamlvfs@latest
yamlvfs --help-all
```

```txt
yamlvfs - Work with yamlvfs YAML filesystems

yamlvfs is a CLI for working with YAML-defined virtual filesystems.

See https://github.com/seb-gan/yamlvfs for more information.

Usage:

yamlvfs

    import-dir # Create yamlvfs document from directory
        [--depth <int>]  # max traversal depth (-1 = unlimited) (default: -1)
        [--exclude-dirs <string>]  # glob patterns for directories to exclude (comma-separated)
        [--include-dirs <string>]  # glob patterns for directories to include (comma-separated) (default: *)
        [--include-file-content <string>]  # glob patterns for files to read content (comma-separated) (default: *)
        [--no-gitignore <bool>]  # ignore .gitignore files
        [--out-file <string>]  # output file (default: stdout)
        --src-dir <string>  # source directory to scan (required)

    print-tree # Print tree structure of yamlvfs file
        --src-file <string>  # source yamlvfs file (required)

    schema # Export or print the embedded schema

        export # Export JSON schema to file
            [--dest-dir <string>]  # destination directory
            [--dest-file <string>]  # destination file path

        print # Print JSON schema to stdout

    validate # Validate yamlvfs file structure
        --src-file <string>  # source yamlvfs file (required)

    write-dir # Create directory and file structure from yamlvfs file
        --dest-dir <string>  # destination directory (required)
        --src-file <string>  # source yamlvfs file (required)
```

## License

MIT
