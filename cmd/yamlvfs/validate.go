package main

import (
	"fmt"

	"github.com/seb-gan/yamlvfs"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:     "validate",
	Short:   "Validate yamlvfs file",
	Example: "  yamlvfs validate --file input.yml",
	RunE:    runValidate,
}

func init() {
	f := validateCmd.Flags()
	f.StringP("file", "f", "", "yamlvfs file (required)")
	validateCmd.MarkFlagRequired("file")
}

func runValidate(cmd *cobra.Command, args []string) error {
	file, _ := cmd.Flags().GetString("file")

	node, err := yamlvfs.ParseFile(file)
	if err != nil {
		return err
	}

	if err := yamlvfs.Validate(node); err != nil {
		return err
	}

	fmt.Println("valid")
	return nil
}
