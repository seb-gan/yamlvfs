package commands

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate YAML VFS from directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			srcDir, _ := cmd.Flags().GetString("src-dir")
			outFile, _ := cmd.Flags().GetString("out-file")

			tree := make(map[string]any)
			err := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				rel, _ := filepath.Rel(srcDir, path)
				if rel == "." {
					return nil
				}

				if d.IsDir() {
					setPath(tree, rel+"/", nil)
				} else {
					content, err := os.ReadFile(path)
					if err != nil {
						return err
					}
					setPath(tree, rel, string(content))
				}

				return nil
			})

			if err != nil {
				return err
			}

			out, err := yaml.Marshal(tree)
			if err != nil {
				return err
			}

			if outFile == "" {
				fmt.Print(string(out))
			} else {
				return os.WriteFile(outFile, out, 0644)
			}

			return nil
		},
	}

	cmd.Flags().String("src-dir", "", "Source directory")
	cmd.Flags().String("out-file", "", "Output file (stdout if not specified)")
	cmd.MarkFlagRequired("src-dir")

	return cmd
}

func setPath(tree map[string]any, path string, value any) {
	parts := strings.Split(filepath.ToSlash(path), "/")
	current := tree
	for i, part := range parts {
		if part == "" {
			continue
		}
		isLast := i == len(parts)-1 || (i == len(parts)-2 && parts[len(parts)-1] == "")
		isDir := strings.HasSuffix(path, "/")

		if isLast {
			if isDir {
				current[part+"/"] = value
			} else {
				current[part] = value
			}
		} else {
			key := part + "/"
			if v, ok := current[key]; !ok || v == nil {
				current[key] = make(map[string]any)
			}
			current = current[key].(map[string]any)
		}
	}
}
