package yamlvfs

import (
	"bufio"
	"io/fs"
	"path"
	"strings"
)

// gitignore holds patterns from .gitignore files.
type gitignore struct {
	fsys     fs.FS
	patterns []ignorePattern
	parent   *gitignore
	dir      string
}

type ignorePattern struct {
	pattern  string
	negation bool
	dirOnly  bool
}

// newGitignore creates a gitignore matcher for a directory.
func newGitignore(fsys fs.FS, dir string) *gitignore {
	gi := &gitignore{fsys: fsys, dir: dir}
	gi.load(dir)
	return gi
}

// child returns a gitignore for a child directory.
func (g *gitignore) child(dir string) *gitignore {
	child := &gitignore{fsys: g.fsys, parent: g, dir: dir}
	child.load(dir)
	return child
}

// load reads .gitignore from the given directory.
func (g *gitignore) load(dir string) {
	gitignorePath := path.Join(dir, ".gitignore")
	f, err := g.fsys.Open(gitignorePath)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		p := ignorePattern{}

		if strings.HasPrefix(line, "!") {
			p.negation = true
			line = line[1:]
		}

		if strings.HasSuffix(line, "/") {
			p.dirOnly = true
			line = strings.TrimSuffix(line, "/")
		}

		p.pattern = line
		g.patterns = append(g.patterns, p)
	}
}

// matches returns true if the path should be ignored.
func (g *gitignore) matches(filePath string, isDir bool) bool {
	if g == nil {
		return false
	}

	// Check parent first
	if g.parent != nil && g.parent.matches(filePath, isDir) {
		return true
	}

	// Get path relative to this gitignore's directory
	relPath := filePath
	if g.dir != "." && g.dir != "" {
		relPath = strings.TrimPrefix(filePath, g.dir+"/")
	}
	name := path.Base(filePath)

	ignored := false
	for _, p := range g.patterns {
		if p.dirOnly && !isDir {
			continue
		}

		// Match against basename or full relative path
		matched := matchPattern(p.pattern, name) || matchPattern(p.pattern, relPath)

		if matched {
			ignored = !p.negation
		}
	}

	return ignored
}

// matchPattern matches a gitignore pattern against a path.
func matchPattern(pattern, name string) bool {
	// Handle ** patterns
	if strings.Contains(pattern, "**") {
		pattern = strings.ReplaceAll(pattern, "**", "*")
	}

	matched, _ := path.Match(pattern, name)
	return matched
}
