package main

import (
	"os"

	"github.com/seb-gan/yamlvfs/internal/commands"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{Use: "yamlvfs"}
	root.AddCommand(
		commands.NewValidateCmd(),
		commands.NewPrintTreeCmd(),
		commands.NewGenerateCmd(),
	)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
