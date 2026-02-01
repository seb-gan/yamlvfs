package main

import (
	"github.com/seb-gan/yamlvfs"
	"github.com/spf13/cobra"
)

var toDirCmd = &cobra.Command{
	Use:     "to-dir",
	Short:   "Write yamlvfs document to directory",
	Example: "  yamlvfs to-dir --file input.yml --out ./output",
	RunE:    runToDir,
}

func init() {
	f := toDirCmd.Flags()
	f.StringP("file", "f", "", "yamlvfs file (required)")
	f.StringP("out", "o", "", "output directory (required)")
	toDirCmd.MarkFlagRequired("file")
	toDirCmd.MarkFlagRequired("out")
}

func runToDir(cmd *cobra.Command, args []string) error {
	file, _ := cmd.Flags().GetString("file")
	out, _ := cmd.Flags().GetString("out")

	node, err := yamlvfs.ParseFile(file)
	if err != nil {
		return err
	}

	fsys, err := yamlvfs.Open(node)
	if err != nil {
		return err
	}

	return yamlvfs.WriteDir(fsys, out)
}
