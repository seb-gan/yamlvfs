// Package yamlvfs provides a YAML-based virtual filesystem.
//
// A yamlvfs document is a YAML mapping where each key represents a file or directory.
// Keys ending with "/" are directories; values are either null or nested mappings.
// Other keys are files; values are either null or string content.
package yamlvfs

import (
	"io/fs"
	"os"
	"strings"
	"testing/fstest"

	"gopkg.in/yaml.v3"
)

// Parse parses YAML text into a yaml.Node.
func Parse(s string) (*yaml.Node, error) {
	var node yaml.Node
	if err := yaml.Unmarshal([]byte(s), &node); err != nil {
		return nil, err
	}
	return &node, nil
}

// ParseFile parses a YAML file into a yaml.Node.
func ParseFile(path string) (*yaml.Node, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(string(data))
}

// Open converts a yaml.Node to an fs.FS.
func Open(node *yaml.Node) (fs.FS, error) {
	if err := Validate(node); err != nil {
		return nil, err
	}

	var tree map[string]any
	if err := node.Decode(&tree); err != nil {
		return nil, err
	}

	fsys := make(fstest.MapFS)
	flatten(fsys, "", tree)
	return fsys, nil
}

// Format converts a yaml.Node to YAML text.
// Null values are formatted as implicit (key:) rather than explicit (key: null).
func Format(node *yaml.Node) string {
	out, err := yaml.Marshal(node)
	if err != nil {
		return ""
	}
	return strings.ReplaceAll(string(out), ": null\n", ":\n")
}

// WriteDir writes an fs.FS to disk at destDir.
func WriteDir(fsys fs.FS, destDir string) error {
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		dest := destDir + "/" + path
		if d.IsDir() {
			return os.MkdirAll(dest, 0755)
		}
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(destDir+"/"+dirOf(path), 0755); err != nil {
			return err
		}
		return os.WriteFile(dest, data, 0644)
	})
}

func dirOf(path string) string {
	if i := strings.LastIndex(path, "/"); i >= 0 {
		return path[:i]
	}
	return "."
}

// flatten recursively converts a nested map to fstest.MapFS.
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
