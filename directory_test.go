package yamlvfs

import (
	"strings"
	"testing"
	"testing/fstest"
)

func TestReadDir_GitignoreScoping(t *testing.T) {
	// Test that .gitignore files only affect their own subtree, not siblings
	//
	// Structure:
	//   dir1/.gitignore: *.txt
	//   dir1/no.txt        <- excluded by dir1/.gitignore
	//   dir1/sub/no.txt    <- excluded (inherited from dir1/.gitignore)
	//   dir2/dir3/yes.txt  <- included (dir1's .gitignore doesn't affect dir2)
	//   dir2/dir3/dir4/.gitignore: *.txt
	//   dir2/dir3/dir4/no.txt <- excluded by dir4/.gitignore

	fsys := fstest.MapFS{
		"dir1/.gitignore":           {Data: []byte("*.txt\n")},
		"dir1/no.txt":               {Data: []byte("")},
		"dir1/sub/no.txt":           {Data: []byte("")},
		"dir2/dir3/yes.txt":         {Data: []byte("")},
		"dir2/dir3/dir4/.gitignore": {Data: []byte("*.txt\n")},
		"dir2/dir3/dir4/no.txt":     {Data: []byte("")},
	}

	doc, err := ReadDir(fsys, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Should be excluded (by dir1/.gitignore)
	if strings.Contains(doc, "dir1/no.txt") || strings.Contains(doc, "no.txt:") && strings.Contains(doc, "dir1/") {
		// Check more carefully - the structure is nested
		if strings.Contains(doc, "dir1/:") && strings.Contains(doc, "no.txt") {
			lines := strings.Split(doc, "\n")
			inDir1 := false
			for _, line := range lines {
				if strings.Contains(line, "dir1/:") {
					inDir1 = true
				} else if inDir1 && !strings.HasPrefix(strings.TrimSpace(line), "") {
					inDir1 = false
				}
				if inDir1 && strings.Contains(line, "no.txt") {
					t.Error("dir1/no.txt should be excluded by dir1/.gitignore")
				}
			}
		}
	}

	// Should be included (dir1's .gitignore doesn't affect dir2)
	if !strings.Contains(doc, "yes.txt") {
		t.Error("dir2/dir3/yes.txt should be included (not affected by dir1/.gitignore)")
	}

	// Should be excluded (by dir4/.gitignore)
	// We check that dir4 exists but no.txt is not in it
	if strings.Contains(doc, "dir4/:") {
		// dir4 should exist, but should not contain no.txt
		lines := strings.Split(doc, "\n")
		inDir4 := false
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "dir4/:") {
				inDir4 = true
			} else if inDir4 && !strings.HasPrefix(line, "                ") && trimmed != "" {
				inDir4 = false
			}
			if inDir4 && strings.Contains(line, "no.txt") {
				t.Error("dir2/dir3/dir4/no.txt should be excluded by dir4/.gitignore")
			}
		}
	}
}

func TestReadDir_GitignoreInheritance(t *testing.T) {
	// Test that .gitignore rules are cumulative within a tree
	//
	// Structure:
	//   .gitignore: *.txt          <- excludes .txt in entire tree
	//   dir1/no.txt                <- excluded by root .gitignore
	//   dir1/yes.js                <- included
	//   dir2/no.txt                <- excluded by root .gitignore (inherited)
	//   dir2/.gitignore: *.js      <- adds .js exclusion for dir2 subtree
	//   dir2/no.js                 <- excluded by dir2/.gitignore
	//   dir2/dir3/no.js            <- excluded by dir2/.gitignore (inherited)
	//   dir2/dir3/no.txt           <- excluded by root .gitignore (inherited)

	fsys := fstest.MapFS{
		".gitignore":       {Data: []byte("*.txt\n")},
		"dir1/no.txt":      {Data: []byte("")},
		"dir1/yes.js":      {Data: []byte("// js")},
		"dir2/no.txt":      {Data: []byte("")},
		"dir2/.gitignore":  {Data: []byte("*.js\n")},
		"dir2/no.js":       {Data: []byte("")},
		"dir2/dir3/no.js":  {Data: []byte("")},
		"dir2/dir3/no.txt": {Data: []byte("")},
	}

	doc, err := ReadDir(fsys, nil)
	if err != nil {
		t.Fatal(err)
	}

	// dir1/yes.js should be included (no .js rule in root or dir1)
	if !strings.Contains(doc, "yes.js") {
		t.Error("dir1/yes.js should be included")
	}

	// All .txt files should be excluded (root .gitignore)
	if containsInDir(doc, "dir1", "no.txt") {
		t.Error("dir1/no.txt should be excluded by root .gitignore")
	}
	if containsInDir(doc, "dir2", "no.txt") {
		t.Error("dir2/no.txt should be excluded by root .gitignore")
	}

	// dir2/*.js and dir2/**/*.js should be excluded
	if containsInDir(doc, "dir2", "no.js") {
		t.Error("dir2/no.js should be excluded by dir2/.gitignore")
	}
}

// containsInDir checks if a file appears under a directory in the YAML output
func containsInDir(doc, dir, file string) bool {
	lines := strings.Split(doc, "\n")
	inDir := false
	dirIndent := -1

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		indent := len(line) - len(strings.TrimLeft(line, " "))

		if strings.HasPrefix(trimmed, dir+"/:") {
			inDir = true
			dirIndent = indent
			continue
		}

		if inDir {
			// Exited the directory (same or lower indent level)
			if indent <= dirIndent && trimmed != "" {
				inDir = false
				dirIndent = -1
				continue
			}

			// Check for the file
			if strings.HasPrefix(trimmed, file+":") || strings.HasPrefix(trimmed, file+" ") {
				return true
			}
		}
	}
	return false
}

func TestReadDir_NoGitignore(t *testing.T) {
	// Test that RespectGitignore=false includes everything
	fsys := fstest.MapFS{
		".gitignore":  {Data: []byte("*.txt\n")},
		"included.go": {Data: []byte("package main")},
		"ignored.txt": {Data: []byte("should be included when gitignore disabled")},
	}

	opts := &ReadDirOptions{
		Depth:              -1,
		IncludeFileContent: []string{"*"},
		IncludeDirs:        []string{"*"},
		ExcludeDirs:        []string{},
		RespectGitignore:   false,
	}

	doc, err := ReadDir(fsys, opts)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(doc, "ignored.txt") {
		t.Error("ignored.txt should be included when RespectGitignore=false")
	}
	if !strings.Contains(doc, "included.go") {
		t.Error("included.go should be included")
	}
}
