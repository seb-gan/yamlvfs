package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/seb-gan/yamlvfs"
	"github.com/spf13/cobra"
)

var printTreeCmd = &cobra.Command{
	Use:   "print-tree",
	Short: "Print tree of yamlvfs document",
	RunE:  runPrintTree,
}

func init() {
	f := printTreeCmd.Flags()
	f.String("src-file", "", "source yamlvfs file (required)")
	printTreeCmd.MarkFlagRequired("src-file")
}

func runPrintTree(cmd *cobra.Command, args []string) error {
	srcFile, _ := cmd.Flags().GetString("src-file")

	fsys, err := yamlvfs.LoadFile(srcFile)
	if err != nil {
		return err
	}

	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		indent := strings.Repeat("  ", strings.Count(path, "/"))
		name := filepath.Base(path)
		if d.IsDir() {
			name += "/"
		}
		fmt.Println(indent + name)
		return nil
	})
}
