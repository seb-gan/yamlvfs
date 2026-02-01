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
	Long: `Create a yamlvfs document from a directory tree.

Flag notes:
  --include-dirs and --exclude-dirs can be combined; a directory must match include and not match exclude.
  --include-file-content controls which files have their content included; others are listed with empty content.
  By default, .gitignore files are respected and cumulative; use --no-gitignore to disable. The .git directory is always excluded.
`,
	Example: "  yamlvfs import-dir --src-dir mydir --out-file fs.yml",
	RunE:    runImportDir,
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

	opts := &yamlvfs.Options{
		Depth:              depth,
		IncludeFileContent: parseGlobs(includeContent),
		IncludeDirs:        parseGlobs(includeDirs),
		ExcludeDirs:        parseGlobs(excludeDirs),
		RespectGitignore:   !noGitignore,
	}

	node, err := yamlvfs.FromFS(os.DirFS(srcDir), opts)
	if err != nil {
		return err
	}

	output := yamlvfs.Format(node)

	if outFile == "" {
		fmt.Print(output)
	} else {
		if err := os.WriteFile(outFile, []byte(output), 0644); err != nil {
			return err
		}
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
