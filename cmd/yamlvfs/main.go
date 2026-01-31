package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "yamlvfs",
		Short: "Work with yamlvfs YAML filesystems",
	}

	root.AddCommand(
		importDirCmd,
		writeDirCmd,
		printTreeCmd,
		validateCmd,
		schemaCmd,
	)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
