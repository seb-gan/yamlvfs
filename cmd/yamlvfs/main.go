package main

import (
	"os"

	helpall "github.com/seb-gan/cobra-helpall"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "yamlvfs",
		Short: "Work with yamlvfs YAML filesystems",
		Long: `yamlvfs is a CLI for working with YAML-defined virtual filesystems.

See https://github.com/seb-gan/yamlvfs for more information.`,
	}

	root.AddCommand(
		fromDirCmd,
		toDirCmd,
		treeCmd,
		validateCmd,
		schemaCmd,
	)

	helpall.Install(root)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
