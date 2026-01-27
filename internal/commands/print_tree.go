package commands

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/seb-gan/yamlvfs"
	"github.com/spf13/cobra"
)

var PrintTreeCmd = &cobra.Command{
	Use:   "print-tree",
	Short: "Print tree of YAML VFS file",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")

		fsys, err := yamlvfs.LoadFile(file)
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
	},
}

func init() {
	PrintTreeCmd.Flags().String("file", "", "YAML file to print")
	PrintTreeCmd.MarkFlagRequired("file")
}
