package helpall_test

import (
	"strings"

	"github.com/seb-gan/yamlvfs/helpall"
	"github.com/spf13/cobra"
)

func Example() {
	root := &cobra.Command{
		Use:   "myapp",
		Short: "My application",
	}
	root.AddCommand(&cobra.Command{
		Use:   "serve",
		Short: "Start the server",
	})

	helpall.Install(root)

	// Users can now run: myapp --help-all
}

func Example_withCustomTemplate() {
	root := &cobra.Command{
		Use:   "myapp",
		Short: "My application",
	}

	// Use DefaultTemplate as a starting point
	custom := strings.Replace(
		helpall.DefaultTemplate(),
		"Usage:",
		"Command Reference:",
		1,
	)

	helpall.Install(root, helpall.WithTemplate(custom))
}

func ExampleDefaultTemplate() {
	// Get the default template to inspect or modify
	tmpl := helpall.DefaultTemplate()

	// Customize it
	custom := strings.Replace(tmpl, "Usage:", "Available Commands:", 1)

	root := &cobra.Command{Use: "myapp"}
	helpall.Install(root, helpall.WithTemplate(custom))
}

func ExampleWithTemplate() {
	root := &cobra.Command{
		Use:   "myapp",
		Short: "My application",
	}

	// Simple custom template showing just command names
	simple := `{{ .Name }}
{{ range $cmd := visibleCommands . }}  {{ $cmd.Name }}
{{ end }}`

	helpall.Install(root, helpall.WithTemplate(simple))
}
