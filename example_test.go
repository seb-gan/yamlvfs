package yamlvfs_test

import (
	"fmt"
	"io/fs"
	"testing/fstest"

	"github.com/seb-gan/yamlvfs"
)

func ExampleParse() {
	node, err := yamlvfs.Parse(`
src/:
  main.go: |
    package main
config.yml: |
  port: 8080
`)
	if err != nil {
		panic(err)
	}

	fsys, _ := yamlvfs.Open(node)
	data, _ := fs.ReadFile(fsys, "config.yml")
	fmt.Print(string(data))

	// Output:
	// port: 8080
}

func ExampleOpen() {
	node, _ := yamlvfs.Parse(`
src/:
  main.go:
  util.go:
`)

	fsys, _ := yamlvfs.Open(node)
	entries, _ := fs.ReadDir(fsys, "src")
	for _, e := range entries {
		fmt.Println(e.Name())
	}

	// Output:
	// main.go
	// util.go
}

func ExampleFormat() {
	fsys := fstest.MapFS{
		"main.go": {Data: []byte("package main")},
		"go.mod":  {Data: []byte("module example")},
	}

	node, _ := yamlvfs.FromFS(fsys, nil)
	fmt.Print(yamlvfs.Format(node))

	// Output:
	// go.mod: module example
	// main.go: package main
}

func ExampleFromFS() {
	fsys := fstest.MapFS{
		"src/main.go": {Data: []byte("package main")},
		"README.md":   {Data: []byte("# Hello")},
	}

	node, _ := yamlvfs.FromFS(fsys, nil)
	fmt.Print(yamlvfs.Format(node))

	// Output:
	// README.md: '# Hello'
	// src/:
	//     main.go: package main
}

func ExampleFromFS_options() {
	fsys := fstest.MapFS{
		"main.go":    {Data: []byte("package main")},
		"secret.key": {Data: []byte("supersecret")},
	}

	// Only include content from .go files
	opts := &yamlvfs.Options{
		Depth:              -1,
		IncludeFileContent: []string{"*.go"},
		IncludeDirs:        []string{"*"},
		RespectGitignore:   false,
	}

	node, _ := yamlvfs.FromFS(fsys, opts)
	fmt.Print(yamlvfs.Format(node))

	// Output:
	// main.go: package main
	// secret.key:
}

func ExampleValidate() {
	// Valid
	node, _ := yamlvfs.Parse(`src/:`)
	err := yamlvfs.Validate(node)
	fmt.Println("valid:", err == nil)

	// Invalid: directory cannot have string content
	node, _ = yamlvfs.Parse(`src/: invalid`)
	err = yamlvfs.Validate(node)
	fmt.Println("invalid:", err != nil)

	// Output:
	// valid: true
	// invalid: true
}
