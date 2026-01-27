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

yamlvfs generate --dir ./my-dir --out-file my-dir.yml
yamlvfs print-tree --file my-dir.yml
yamlvfs validate --file my-dir.yml
```

## License

MIT
