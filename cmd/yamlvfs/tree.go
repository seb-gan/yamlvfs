package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/seb-gan/yamlvfs"
	"github.com/spf13/cobra"
)

var treeCmd = &cobra.Command{
	Use:     "tree",
	Short:   "Print tree structure of yamlvfs file",
	Example: "  yamlvfs tree --file input.yml",
	RunE:    runTree,
}

func init() {
	f := treeCmd.Flags()
	f.StringP("file", "f", "", "yamlvfs file (required)")
	treeCmd.MarkFlagRequired("file")
}

func runTree(cmd *cobra.Command, args []string) error {
	file, _ := cmd.Flags().GetString("file")

	node, err := yamlvfs.ParseFile(file)
	if err != nil {
		return err
	}

	fsys, err := yamlvfs.Open(node)
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
