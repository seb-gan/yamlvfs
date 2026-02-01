# yamlvfs

[![Go Reference](https://pkg.go.dev/badge/github.com/seb-gan/yamlvfs.svg)](https://pkg.go.dev/github.com/seb-gan/yamlvfs)
[![Go Report Card](https://goreportcard.com/badge/github.com/seb-gan/yamlvfs)](https://goreportcard.com/report/github.com/seb-gan/yamlvfs)

Define virtual filesystem layouts using YAML.

`yamlvfs` is three things built around one idea:

1. **A YAML format** for representing filesystem trees
2. **A Go library** that loads yamlvfs documents as standard `fs.FS` implementation
3. **A CLI** to generate, inspect, and validate yamlvfs documents

## YAML format

The format is defined in schema  `yamlvfs.schema.json`, which specifies only 2 types of nodes:

- **Directories**: Must end with `/`. The value can either none or one or more directories and/or files.
- **Files**: The value can either be none or a (multiline) string value.

```yaml
data/:
src/:
  version.txt: '1'
  main.go:
  utils.go: |
    // util functions
```

## Go Library

**Usage:**

```bash
go get github.com/seb-gan/yamlvfs
```

```go
// load from file
fsys, _ := yamlvfs.LoadFile("data/sample.yml")

// load from string
fsys, _ := yamlvfs.Load(`
config.yml: |
  name: myapp
src/:
  main.go: |
    package main
`)

// use standard fs.FS implementation
data, _ := fs.ReadFile(fsys, "config.yml")
entries, _ := fs.ReadDir(fsys, "src")
```

## CLI

**Usage:**

```bash
go install github.com/seb-gan/yamlvfs/cmd/yamlvfs@latest

yamlvfs --help-all
```

```
yamlvfs - Work with yamlvfs YAML filesystems

yamlvfs is a CLI for working with YAML-defined virtual filesystems.

See https://github.com/seb-gan/yamlvfs for more information.

Usage:

yamlvfs

    completion  # Generate the autocompletion script for the specified shell

        bash  # Generate the autocompletion script for bash
            [--no-descriptions <bool>]  # disable completion descriptions  (default: false)

        fish  # Generate the autocompletion script for fish
            [--no-descriptions <bool>]  # disable completion descriptions  (default: false)

        powershell  # Generate the autocompletion script for powershell
            [--no-descriptions <bool>]  # disable completion descriptions  (default: false)

        zsh  # Generate the autocompletion script for zsh
            [--no-descriptions <bool>]  # disable completion descriptions  (default: false)

    import-dir  # Create yamlvfs document from directory
        [--depth <int>]  # max traversal depth (-1 = unlimited)  (default: -1)
        [--exclude-dirs <string>]  # glob patterns for directories to exclude (comma-separated)
        [--include-dirs <string>]  # glob patterns for directories to include (comma-separated)  (default: *)
        [--include-file-content <string>]  # glob patterns for files to read content (comma-separated)  (default: *)
        [--no-gitignore <bool>]  # ignore .gitignore files  (default: false)
        [--out-file <string>]  # output file (default: stdout)
        --src-dir <string>  # source directory to scan (required)

        Example:
          yamlvfs import-dir --src-dir mydir --out-file fs.yml
          yamlvfs import-dir --src-dir . --include-file-content=*.go --exclude-dirs=test
          yamlvfs import-dir --src-dir . --no-gitignore

    print-tree  # Print tree structure of yamlvfs file
        --src-file <string>  # source yamlvfs file (required)

        Example:
          yamlvfs print-tree --src-file fs.yml

    schema  # Export or print the embedded schema

        export  # Export JSON schema to file
            [--dest-dir <string>]  # destination directory
            [--dest-file <string>]  # destination file path

            Example:
              yamlvfs schema export --dest-dir .
              yamlvfs schema export --dest-file my-schema.json

        print  # Print JSON schema to stdout

            Example:
              yamlvfs schema print

    validate  # Validate yamlvfs file structure
        --src-file <string>  # source yamlvfs file (required)

        Example:
          yamlvfs validate --src-file fs.yml

    write-dir  # Create directory and file structure from yamlvfs file
        --dest-dir <string>  # destination directory (required)
        --src-file <string>  # source yamlvfs file (required)

        Example:
          yamlvfs write-dir --src-file fs.yml --dest-dir out

```

## License

MIT
