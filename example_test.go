package yamlvfs_test

import (
	"fmt"
	"io/fs"
	"testing/fstest"

	"github.com/seb-gan/yamlvfs"
)

func ExampleParse() {
	node, err := yamlvfs.Parse(`
cmd/:
  main.go: |
    package main

    func main() {
        println("hello")
    }
internal/:
  util.go: |
    package internal

    func Add(a, b int) int {
        return a + b
    }
go.mod: |
  module example.com/myapp

  go 1.23
`)
	if err != nil {
		panic(err)
	}

	fsys, _ := yamlvfs.Open(node)

	data, _ := fs.ReadFile(fsys, "go.mod")
	fmt.Print(string(data))

	// Output:
	// module example.com/myapp
	//
	// go 1.23
}

func ExampleOpen() {
	node, _ := yamlvfs.Parse(`
cmd/:
  main.go: |
    package main
  cli/:
    root.go: |
      package cli
    help.go: |
      package cli
internal/:
  util.go: |
    package internal
go.mod: module example
`)

	fsys, _ := yamlvfs.Open(node)

	entries, _ := fs.ReadDir(fsys, "cmd/cli")
	for _, e := range entries {
		fmt.Println(e.Name())
	}

	// Output:
	// help.go
	// root.go
}

func ExampleFormat() {
	fsys := fstest.MapFS{
		"cmd/main.go":      {Data: []byte("package main")},
		"internal/util.go": {Data: []byte("package internal")},
		"go.mod":           {Data: []byte("module example")},
		"README.md":        {Data: []byte("# Example")},
	}

	node, _ := yamlvfs.FromFS(fsys, nil)
	fmt.Print(yamlvfs.Format(node))

	// Output:
	// README.md: '# Example'
	// cmd/:
	//     main.go: package main
	// go.mod: module example
	// internal/:
	//     util.go: package internal
}

func ExampleFromFS() {
	fsys := fstest.MapFS{
		"src/main.go":    {Data: []byte("package main")},
		"src/handler.go": {Data: []byte("package main")},
		"docs/api.md":    {Data: []byte("# API")},
		"README.md":      {Data: []byte("# Project")},
	}

	node, _ := yamlvfs.FromFS(fsys, nil)
	fmt.Print(yamlvfs.Format(node))

	// Output:
	// README.md: '# Project'
	// docs/:
	//     api.md: '# API'
	// src/:
	//     handler.go: package main
	//     main.go: package main
}

func ExampleFromFS_options() {
	fsys := fstest.MapFS{
		"cmd/main.go":      {Data: []byte("package main")},
		"internal/util.go": {Data: []byte("package internal")},
		"assets/logo.png":  {Data: []byte{0x89, 0x50, 0x4E, 0x47}},
		"go.mod":           {Data: []byte("module example")},
	}

	opts := &yamlvfs.Options{
		Depth:              -1,
		IncludeFileContent: []string{"*.go", "go.mod"},
		IncludeDirs:        []string{"*"},
		RespectGitignore:   false,
	}

	node, _ := yamlvfs.FromFS(fsys, opts)
	fmt.Print(yamlvfs.Format(node))

	// Output:
	// assets/:
	//     logo.png:
	// cmd/:
	//     main.go: package main
	// go.mod: module example
	// internal/:
	//     util.go: package internal
}

func ExampleValidate() {
	valid, _ := yamlvfs.Parse(`
src/:
  main.go: |
    package main
  util.go:
data/:
config.yml: |
  port: 8080
`)
	err := yamlvfs.Validate(valid)
	fmt.Println("valid:", err == nil)

	invalid, _ := yamlvfs.Parse(`
src/: should be null or mapping
`)
	err = yamlvfs.Validate(invalid)
	fmt.Println("invalid:", err != nil)

	// Output:
	// valid: true
	// invalid: true
}
