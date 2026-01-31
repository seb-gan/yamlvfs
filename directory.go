package yamlvfs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ReadDirOptions configures ReadDir behavior.
type ReadDirOptions struct {
	// Depth limits directory traversal depth. -1 means unlimited.
	Depth int

	// IncludeFileContent is a list of glob patterns for files whose content
	// should be read. Files not matching are included with empty content.
	// Default: ["*"] (all files).
	IncludeFileContent []string

	// IncludeDirs is a list of glob patterns for directories to include.
	// Default: ["*"] (all directories).
	IncludeDirs []string

	// ExcludeDirs is a list of glob patterns for directories to exclude.
	// Default: [] (none excluded).
	ExcludeDirs []string

	// RespectGitignore parses .gitignore files and excludes matching paths.
	// Default: true.
	RespectGitignore bool
}

// DefaultReadDirOptions returns sensible defaults for ReadDir.
func DefaultReadDirOptions() *ReadDirOptions {
	return &ReadDirOptions{
		Depth:              -1,
		IncludeFileContent: []string{"*"},
		IncludeDirs:        []string{"*"},
		ExcludeDirs:        []string{},
		RespectGitignore:   true,
	}
}

// ReadDir serializes an fs.FS to a yamlvfs Document.
// If opts is nil, DefaultReadDirOptions is used.
func ReadDir(fsys fs.FS, opts *ReadDirOptions) (Document, error) {
	if opts == nil {
		opts = DefaultReadDirOptions()
	}

	tree := make(map[string]any)
	var gi *gitignore

	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip root
		if path == "." {
			if opts.RespectGitignore {
				gi = newGitignore(fsys, ".")
			}
			return nil
		}

		// Always skip .git directory
		name := filepath.Base(path)
		if name == ".git" && d.IsDir() {
			return fs.SkipDir
		}

		depth := strings.Count(path, "/")
		if opts.Depth >= 0 && depth > opts.Depth {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		// Handle gitignore
		if opts.RespectGitignore && gi != nil && gi.matches(path, d.IsDir()) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			// Check include/exclude patterns
			if !matchesAny(name, opts.IncludeDirs) || matchesAny(name, opts.ExcludeDirs) {
				return fs.SkipDir
			}

			// Load gitignore for this directory
			if opts.RespectGitignore {
				gi = gi.child(path)
			}

			setPath(tree, path+"/", nil)
		} else {
			var content string
			if matchesAny(name, opts.IncludeFileContent) {
				data, err := fs.ReadFile(fsys, path)
				if err != nil {
					return fmt.Errorf("yamlvfs: failed to read %s: %w", path, err)
				}
				// Normalize line endings to LF for clean YAML output
				content = strings.ReplaceAll(string(data), "\r\n", "\n")
			}
			setPath(tree, path, content)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	out, err := yaml.Marshal(tree)
	if err != nil {
		return "", fmt.Errorf("yamlvfs: failed to marshal: %w", err)
	}

	return Document(out), nil
}

// WriteDir writes an fs.FS to disk at destDir.
func WriteDir(fsys fs.FS, destDir string) error {
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		dest := filepath.Join(destDir, path)

		if d.IsDir() {
			return os.MkdirAll(dest, 0755)
		}

		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("yamlvfs: failed to read %s: %w", path, err)
		}

		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}

		return os.WriteFile(dest, data, 0644)
	})
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

// matchesAny returns true if name matches any of the glob patterns.
func matchesAny(name string, patterns []string) bool {
	for _, p := range patterns {
		if matched, _ := filepath.Match(p, name); matched {
			return true
		}
	}
	return false
}
