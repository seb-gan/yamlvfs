package commands

import (
	"fmt"

	"github.com/seb-gan/yamlvfs"
	"github.com/spf13/cobra"
)

var ValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a YAML VFS file",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")

		if _, err := yamlvfs.LoadFile(file); err != nil {
			return err
		}

		fmt.Println("valid")
		return nil
	},
}

func init() {
	ValidateCmd.Flags().String("file", "", "YAML file to validate")
	ValidateCmd.MarkFlagRequired("file")
}
