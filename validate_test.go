package yamlvfs

import (
	"testing"
)

func TestValidate_ValidYAML(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "simple file",
			content: `file.txt: hello`,
		},
		{
			name:    "empty file",
			content: `file.txt:`,
		},
		{
			name: "directory with files",
			content: `
src/:
  main.go: |
    package main
  utils.go: helper
`,
		},
		{
			name:    "empty directory",
			content: `empty/:`,
		},
		{
			name: "nested structure",
			content: `
src/:
  cmd/:
    main.go: package main
  internal/:
    helper.go: package internal
config.yml: |
  name: test
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(tt.content); err != nil {
				t.Errorf("Validate() error = %v", err)
			}
		})
	}
}

func TestValidate_InvalidYAML(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "directory with string value",
			content: `dir/: some content`,
		},
		{
			name: "array value",
			content: `file.txt:
  - item1
  - item2`,
		},
		{
			name:    "invalid filename characters",
			content: `file<name>.txt: content`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(tt.content); err == nil {
				t.Error("Validate() expected error, got nil")
			}
		})
	}
}
