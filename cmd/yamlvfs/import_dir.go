package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/seb-gan/yamlvfs"
	"github.com/spf13/cobra"
)

var importDirCmd = &cobra.Command{
	Use:   "import-dir",
	Short: "Create yamlvfs document from directory",
	RunE:  runImportDir,
}

func init() {
	f := importDirCmd.Flags()
	f.String("src-dir", "", "source directory to scan (required)")
	f.Int("depth", -1, "max traversal depth (-1 = unlimited)")
	f.String("out-file", "", "output file (default: stdout)")
	f.String("include-file-content", "*", "glob patterns for files to read content (comma-separated)")
	f.String("include-dirs", "*", "glob patterns for directories to include (comma-separated)")
	f.String("exclude-dirs", "", "glob patterns for directories to exclude (comma-separated)")
	f.Bool("no-gitignore", false, "ignore .gitignore files")
	importDirCmd.MarkFlagRequired("src-dir")
}

func runImportDir(cmd *cobra.Command, args []string) error {
	srcDir, _ := cmd.Flags().GetString("src-dir")
	outFile, _ := cmd.Flags().GetString("out-file")
	depth, _ := cmd.Flags().GetInt("depth")
	includeContent, _ := cmd.Flags().GetString("include-file-content")
	includeDirs, _ := cmd.Flags().GetString("include-dirs")
	excludeDirs, _ := cmd.Flags().GetString("exclude-dirs")
	noGitignore, _ := cmd.Flags().GetBool("no-gitignore")

	opts := &yamlvfs.ReadDirOptions{
		Depth:              depth,
		IncludeFileContent: parseGlobs(includeContent),
		IncludeDirs:        parseGlobs(includeDirs),
		ExcludeDirs:        parseGlobs(excludeDirs),
		RespectGitignore:   !noGitignore,
	}

	doc, err := yamlvfs.ReadDir(os.DirFS(srcDir), opts)
	if err != nil {
		return err
	}

	if outFile == "" {
		fmt.Print(doc)
	} else {
		if err := os.WriteFile(outFile, []byte(doc), 0644); err != nil {
			return err
		}
	}

	return nil
}

// parseGlobs splits a comma-separated string into glob patterns.
func parseGlobs(s string) []string {
	if s == "" {
		return nil
	}
	var result []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			result = append(result, p)
		}
	}
	return result
}
