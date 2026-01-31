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

yamlvfs
    import-dir # creates yamlvfs document matching directory
        --src-dir <directory>
            # source directory to recursively scan
        [--depth <int>]
            # max depth to scan, default unlimited
        [--out-file <string>]
            # if specified, put yamlvfs in file, else print to stdout
        [--include-file-content <comma separated globs>] default "*"
            # include content for files matching these globs
        [--include-dirs <comma separated globs>] default "*"
            # include directories matching these globs
        [--exclude-dirs <comma separated globs>] default ""
            # exclude directories matching these globs
        [--no-gitignore <bool>] default false
            # When true, read .gitignore in each dir, add to exclusion list going deeper

    write-dir # creates directory and file structure from yamlvfs file
        --src-file <yamlvfs file>
            # source yamlvfs file
        --dest-dir <directory>
            # destination directory to create structure in

    print-tree # prints tree structure of yamlvfs file
        --src-file <string>
            # source yamlvfs file

    validate # validates yamlvfs file structure against embedded schema
        --src-file <string>
            # source yamlvfs file

    schema
        export # Exports the embedded schema to file or stdout
            --dest-dir <string> Default ""
                # if specified, write schema with its default filename to this directory
            [--dest-file <string> ]
                # if specified, write schema to this file
        print # Prints schema to stdout
```

## License

MIT
