package yamlvfs

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Options configures FromFS behavior.
type Options struct {
	Depth              int      // Max traversal depth (-1 = unlimited)
	IncludeFileContent []string // Glob patterns for files to read content
	IncludeDirs        []string // Glob patterns for directories to include
	ExcludeDirs        []string // Glob patterns for directories to exclude
	RespectGitignore   bool     // Honor .gitignore files
}

// DefaultOptions returns default options for FromFS.
func DefaultOptions() *Options {
	return &Options{
		Depth:              -1,
		IncludeFileContent: []string{"*"},
		IncludeDirs:        []string{"*"},
		RespectGitignore:   true,
	}
}

// FromFS creates a yaml.Node from an fs.FS.
func FromFS(fsys fs.FS, opts *Options) (*yaml.Node, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	tree := make(map[string]any)
	gitignores := make(map[string]*gitignore)

	if opts.RespectGitignore {
		gitignores["."] = newGitignore(fsys, ".")
	}

	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == "." {
			return nil
		}

		name := filepath.Base(path)

		// Always skip .git
		if name == ".git" && d.IsDir() {
			return fs.SkipDir
		}

		// Check depth
		depth := strings.Count(path, "/")
		if opts.Depth >= 0 && depth > opts.Depth {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		// Check gitignore
		if opts.RespectGitignore {
			dir := filepath.Dir(path)
			if dir == "." {
				dir = "."
			}
			if gi := gitignores[filepath.ToSlash(dir)]; gi != nil && gi.matches(path, d.IsDir()) {
				if d.IsDir() {
					return fs.SkipDir
				}
				return nil
			}
		}

		if d.IsDir() {
			if !matchesAny(name, opts.IncludeDirs) || matchesAny(name, opts.ExcludeDirs) {
				return fs.SkipDir
			}
			if opts.RespectGitignore {
				parent := gitignores[filepath.ToSlash(filepath.Dir(path))]
				if parent == nil {
					parent = gitignores["."]
				}
				gitignores[filepath.ToSlash(path)] = parent.child(path)
			}
			setPath(tree, path+"/", nil)
		} else {
			var content *string
			if matchesAny(name, opts.IncludeFileContent) {
				data, err := fs.ReadFile(fsys, path)
				if err != nil {
					return fmt.Errorf("failed to read %s: %w", path, err)
				}
				s := strings.ReplaceAll(string(data), "\r\n", "\n")
				content = &s
			}
			setPath(tree, path, content)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	var node yaml.Node
	if err := node.Encode(tree); err != nil {
		return nil, err
	}
	return &node, nil
}

// setPath sets a value in the nested map at the specified path.
func setPath(tree map[string]any, path string, value any) {
	parts := strings.Split(filepath.ToSlash(path), "/")
	current := tree

	for i, part := range parts {
		if part == "" {
			continue
		}
		isLast := i == len(parts)-1 || (i == len(parts)-2 && parts[len(parts)-1] == "")
		isDir := strings.HasSuffix(path, "/")

		if isLast {
			if isDir {
				current[part+"/"] = value
			} else {
				current[part] = value
			}
		} else {
			key := part + "/"
			if v, ok := current[key]; !ok || v == nil {
				current[key] = make(map[string]any)
			}
			current = current[key].(map[string]any)
		}
	}
}

// matchesAny returns true if name matches any glob pattern.
func matchesAny(name string, patterns []string) bool {
	for _, p := range patterns {
		if matched, _ := filepath.Match(p, name); matched {
			return true
		}
	}
	return false
}
