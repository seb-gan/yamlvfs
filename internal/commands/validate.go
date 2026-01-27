package commands

import (
	"fmt"

	"github.com/seb-gan/yamlvfs"
	"github.com/spf13/cobra"
)

func NewValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a YAML VFS file",
		RunE: func(cmd *cobra.Command, args []string) error {
			file, _ := cmd.Flags().GetString("file")

			// LoadFile validates automatically
			if _, err := yamlvfs.LoadFile(file); err != nil {
				return err
			}

			fmt.Println("valid")
			return nil
		},
	}

	cmd.Flags().String("file", "", "YAML file to validate")
	cmd.MarkFlagRequired("file")

	return cmd
}
