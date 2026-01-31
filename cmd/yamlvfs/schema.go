package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/seb-gan/yamlvfs"
	"github.com/spf13/cobra"
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Export or print the embedded schema",
	Long:  `Export or print the embedded JSON schema for yamlvfs documents.`,
}

var schemaExportCmd = &cobra.Command{
	Use:     "export",
	Short:   "Export JSON schema to file",
	Long:    `Export the embedded JSON schema to a file or directory.`,
	Example: "  yamlvfs schema export --dest-dir .\n  yamlvfs schema export --dest-file my-schema.json",
	RunE:    runSchemaExport,
}

var schemaPrintCmd = &cobra.Command{
	Use:     "print",
	Short:   "Print JSON schema to stdout",
	Long:    `Print the embedded JSON schema to stdout.`,
	Example: "  yamlvfs schema print",
	RunE:    runSchemaPrint,
}

func init() {
	schemaCmd.AddCommand(schemaExportCmd, schemaPrintCmd)

	f := schemaExportCmd.Flags()
	f.String("dest-dir", "", "destination directory")
	f.String("dest-file", "", "destination file path")
}

func runSchemaExport(cmd *cobra.Command, args []string) error {
	destDir, _ := cmd.Flags().GetString("dest-dir")
	destFile, _ := cmd.Flags().GetString("dest-file")

	if destDir == "" && destFile == "" {
		return fmt.Errorf("either --dest-dir or --dest-file is required")
	}

	path := destFile
	if destDir != "" {
		path = filepath.Join(destDir, "yamlvfs.schema.json")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(yamlvfs.Schema()), 0644)
}

func runSchemaPrint(cmd *cobra.Command, args []string) error {
	fmt.Print(yamlvfs.Schema())
	return nil
}
