// Package yamlvfs loads YAML documents as [fs.FS] implementations.
//
// Keys ending with "/" are directories; all others are files.
// Directories may optionally contain nested files or directories.
// Files may optionally contain (multi-line) string content.
//
// Example YAML structure:
//
//	src/:
//	  main.go: |
//	    package main
//	config.yml: |
//	  name: myapp
package yamlvfs

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	"testing/fstest"

	"gopkg.in/yaml.v3"
)

// FS is a filesystem created from a yamlvfs Document.
type FS = fs.FS

// Document is a yamlvfs-formatted YAML string.
type Document = string

// Load parses and validates YAML content, returning an [fs.FS].
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

// LoadFile parses and validates YAML content from a file, returning an [fs.FS].
func LoadFile(path string) (fs.FS, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("yamlvfs: failed to read file: %w", err)
	}
	return Load(string(data))
}

// flatten recursively flattens a nested map into fstest.MapFS.
func flatten(fsys fstest.MapFS, prefix string, tree map[string]any) {
	for key, value := range tree {
		isDir := strings.HasSuffix(key, "/")
		name := strings.TrimSuffix(key, "/")
		path := prefix + name

		if isDir {
			fsys[path] = &fstest.MapFile{Mode: fs.ModeDir}
			if children, ok := value.(map[string]any); ok {
				flatten(fsys, path+"/", children)
			}
		} else {
			var data []byte
			if s, ok := value.(string); ok {
				data = []byte(s)
			}
			fsys[path] = &fstest.MapFile{Data: data}
		}
	}
}
