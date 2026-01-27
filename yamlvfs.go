// Package yamlvfs provides a virtual file system backed by YAML content.
//
// It returns a standard fs.FS interface, allowing seamless swapping between
// yamlvfs and os.DirFS for testing or mocking file system operations.
//
// Example YAML structure:
//
//	config.yml: |
//	  name: myapp
//	src/:
//	  main.go: |
//	    package main
//	  utils/:
//	    helper.go: |
//	      package utils
//
// Keys ending with "/" are directories, others are files.
// File values are their content as strings.
package yamlvfs

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	"testing/fstest"

	"gopkg.in/yaml.v3"
)

// Load parses YAML content, validates it against the schema, and returns an fs.FS implementation.
//
// The YAML structure represents a filesystem where:
//   - Keys ending with "/" are directories
//   - Keys without "/" are files
//   - Values are file contents (strings) or nested directory structures
//
// Example:
//
//	fsys, err := yamlvfs.Load(`
//	config.yml: |
//	  name: myapp
//	src/:
//	  main.go: |
//	    package main
//	`)
//	data, _ := fs.ReadFile(fsys, "config.yml")
func Load(content string) (fs.FS, error) {
	if err := Validate(content); err != nil {
		return nil, err
	}

	var tree map[string]any
	if err := yaml.Unmarshal([]byte(content), &tree); err != nil {
		return nil, fmt.Errorf("yamlvfs: failed to parse YAML: %w", err)
	}

	fsys := make(fstest.MapFS)
	flatten(fsys, "", tree)
	return fsys, nil
}

// LoadFile reads a YAML file, validates it, and returns an fs.FS.
func LoadFile(path string) (fs.FS, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("yamlvfs: failed to read file: %w", err)
	}
	return Load(string(data))
}

// flatten recursively converts a nested YAML map into a flat fstest.MapFS.
func flatten(fsys fstest.MapFS, prefix string, tree map[string]any) {
	for key, value := range tree {
		isDir := strings.HasSuffix(key, "/")
		name := strings.TrimSuffix(key, "/")
		path := prefix + name

		if isDir {
			// Add directory entry
			fsys[path] = &fstest.MapFile{Mode: fs.ModeDir}

			// Recurse into children
			if children, ok := value.(map[string]any); ok {
				flatten(fsys, path+"/", children)
			}
		} else {
			// Add file entry
			var data []byte
			if s, ok := value.(string); ok {
				data = []byte(s)
			}
			fsys[path] = &fstest.MapFile{Data: data}
		}
	}
}
