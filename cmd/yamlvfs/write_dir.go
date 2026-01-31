package main

import (
	"github.com/seb-gan/yamlvfs"
	"github.com/spf13/cobra"
)

var writeDirCmd = &cobra.Command{
	Use:   "write-dir",
	Short: "Create directory and file structure from yamlvfs file",
	Long:  `Create a directory and file structure from a yamlvfs document.`,

	Example: "  yamlvfs write-dir --src-file fs.yml --dest-dir out",
	RunE:    runWriteDir,
}

func init() {
	f := writeDirCmd.Flags()
	f.String("src-file", "", "source yamlvfs file (required)")
	f.String("dest-dir", "", "destination directory (required)")
	writeDirCmd.MarkFlagRequired("src-file")
	writeDirCmd.MarkFlagRequired("dest-dir")
}

func runWriteDir(cmd *cobra.Command, args []string) error {
	srcFile, _ := cmd.Flags().GetString("src-file")
	destDir, _ := cmd.Flags().GetString("dest-dir")

	fsys, err := yamlvfs.LoadFile(srcFile)
	if err != nil {
		return err
	}

	return yamlvfs.WriteDir(fsys, destDir)
}
