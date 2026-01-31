package main

import (
	"fmt"
	"os"

	"github.com/seb-gan/yamlvfs"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a yamlvfs document",
	RunE:  runValidate,
}

func init() {
	f := validateCmd.Flags()
	f.String("src-file", "", "source yamlvfs file (required)")
	validateCmd.MarkFlagRequired("src-file")
}

func runValidate(cmd *cobra.Command, args []string) error {
	srcFile, _ := cmd.Flags().GetString("src-file")

	data, err := os.ReadFile(srcFile)
	if err != nil {
		return err
	}

	if err := yamlvfs.Validate(yamlvfs.Document(data)); err != nil {
		return err
	}

	fmt.Println("valid")
	return nil
}
