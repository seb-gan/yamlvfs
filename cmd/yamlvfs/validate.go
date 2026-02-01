package main

import (
	"fmt"

	"github.com/seb-gan/yamlvfs"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:     "validate",
	Short:   "Validate yamlvfs file structure",
	Example: "  yamlvfs validate --src-file fs.yml",
	RunE:    runValidate,
}

func init() {
	f := validateCmd.Flags()
	f.String("src-file", "", "source yamlvfs file (required)")
	validateCmd.MarkFlagRequired("src-file")
}

func runValidate(cmd *cobra.Command, args []string) error {
	srcFile, _ := cmd.Flags().GetString("src-file")

	node, err := yamlvfs.ParseFile(srcFile)
	if err != nil {
		return err
	}

	if err := yamlvfs.Validate(node); err != nil {
		return err
	}

	fmt.Println("valid")
	return nil
}
