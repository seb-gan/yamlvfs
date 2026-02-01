package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/seb-gan/yamlvfs"
	"github.com/spf13/cobra"
)

var fromDirCmd = &cobra.Command{
	Use:     "from-dir",
	Short:   "Create yamlvfs document from directory",
	Example: "  yamlvfs from-dir --dir ./src -o output.yml",
	RunE:    runFromDir,
}

func init() {
	f := fromDirCmd.Flags()
	f.String("dir", "", "directory to scan (required)")
	f.StringP("out", "o", "", "output file (default: stdout)")
	f.Int("depth", -1, "max depth (-1 = unlimited)")
	f.String("include-content", "*", "glob patterns for file content (comma-separated)")
	f.String("include-dirs", "*", "glob patterns for directories (comma-separated)")
	f.String("exclude-dirs", "", "glob patterns to exclude (comma-separated)")
	f.Bool("no-gitignore", false, "do not skip .gitignore paths")
	fromDirCmd.MarkFlagRequired("dir")
}

func runFromDir(cmd *cobra.Command, args []string) error {
	dir, _ := cmd.Flags().GetString("dir")
	out, _ := cmd.Flags().GetString("out")
	depth, _ := cmd.Flags().GetInt("depth")
	includeContent, _ := cmd.Flags().GetString("include-content")
	includeDirs, _ := cmd.Flags().GetString("include-dirs")
	excludeDirs, _ := cmd.Flags().GetString("exclude-dirs")
	noGitignore, _ := cmd.Flags().GetBool("no-gitignore")

	opts := &yamlvfs.Options{
		Depth:              depth,
		IncludeFileContent: parseGlobs(includeContent),
		IncludeDirs:        parseGlobs(includeDirs),
		ExcludeDirs:        parseGlobs(excludeDirs),
		RespectGitignore:   !noGitignore,
	}

	node, err := yamlvfs.FromFS(os.DirFS(dir), opts)
	if err != nil {
		return err
	}

	yaml := yamlvfs.Format(node)

	if out == "" {
		fmt.Print(yaml)
	} else {
		return os.WriteFile(out, []byte(yaml), 0644)
	}
	return nil
}

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
