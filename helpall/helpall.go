// Package helpall adds a --help-all flag to Cobra CLI applications that displays
// a complete command reference showing all commands, subcommands, and flags.
//
// # Basic Usage
//
// Call Install on your root command before Execute:
//
//	func main() {
//	    root := &cobra.Command{Use: "myapp", Short: "My application"}
//	    root.AddCommand(subCmd1, subCmd2)
//	    helpall.Install(root)
//	    root.Execute()
//	}
//
// Users can then run:
//
//	myapp --help-all
//
// # Custom Templates
//
// You can provide a custom template using WithTemplate:
//
//	helpall.Install(root, helpall.WithTemplate(myTemplate))
//
// Use DefaultTemplate as a starting point for customization:
//
//	base := helpall.DefaultTemplate()
//	custom := strings.Replace(base, "Usage:", "Commands:", 1)
//	helpall.Install(root, helpall.WithTemplate(custom))
//
// # Template Functions
//
// The following functions are available in templates:
//
//   - indent(n int) string - returns n spaces
//   - add(a, b int) int - returns a + b
//   - mul(a, b int) int - returns a * b
//   - prefixLines(n int, s string) string - indents each line of s by n spaces
//   - visibleCommands(cmd) []*cobra.Command - returns visible subcommands
//   - visibleFlags(cmd) []*pflag.Flag - returns visible flags (excludes help flags)
//   - isRequired(flag) bool - returns true if flag is marked required
//   - dict(pairs ...any) map[string]any - creates a map from key/value pairs
package helpall

import (
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Option configures [Install] behavior.
type Option func(*options)

type options struct {
	template string
}

// WithTemplate sets a custom Go text/template for the help-all output.
// The template receives the root [*cobra.Command] as its data.
// See package documentation for available template functions.
func WithTemplate(t string) Option {
	return func(o *options) {
		o.template = t
	}
}

// DefaultTemplate returns the default template string.
// Use this as a starting point when creating custom templates.
func DefaultTemplate() string {
	return defaultTmpl
}

// Install adds a --help-all persistent flag to root and configures the help
// system to display a complete command reference when the flag is set.
func Install(root *cobra.Command, opts ...Option) {
	o := &options{template: defaultTmpl}
	for _, opt := range opts {
		opt(o)
	}

	root.PersistentFlags().Bool("help-all", false, "Show full command reference")

	orig := root.HelpFunc()
	root.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if v, _ := cmd.Flags().GetBool("help-all"); v {
			printAll(root, o.template)
			return
		}
		orig(cmd, args)
	})
}

func printAll(root *cobra.Command, tmpl string) {
	t := template.Must(template.New("").Funcs(funcs).Parse(tmpl))
	t.Execute(os.Stdout, root)
}

var funcs = template.FuncMap{
	"indent": func(n int) string { return strings.Repeat(" ", n) },
	"add":    func(a, b int) int { return a + b },
	"mul":    func(a, b int) int { return a * b },
	"prefixLines": func(n int, s string) string {
		prefix := strings.Repeat(" ", n)
		lines := strings.Split(strings.TrimSpace(s), "\n")
		for i, line := range lines {
			lines[i] = prefix + strings.TrimSpace(line)
		}
		return strings.Join(lines, "\n")
	},
	"visibleCommands": func(cmd *cobra.Command) []*cobra.Command {
		var out []*cobra.Command
		for _, c := range cmd.Commands() {
			if c.IsAvailableCommand() && c.Name() != "help" {
				out = append(out, c)
			}
		}
		return out
	},
	"visibleFlags": func(cmd *cobra.Command) []*pflag.Flag {
		var out []*pflag.Flag
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Name != "help" && f.Name != "help-all" {
				out = append(out, f)
			}
		})
		return out
	},
	"isRequired": func(f *pflag.Flag) bool {
		_, ok := f.Annotations[cobra.BashCompOneRequiredFlag]
		return ok
	},
}

const defaultTmpl = `{{ .Name }}{{ with .Short }} - {{ . }}{{ end }}
{{ with .Long }}
{{ . }}
{{ end }}
Usage:

{{ .Name }}
{{ range $cmd := visibleCommands . }}
{{ template "cmd" (dict "Cmd" $cmd "Depth" 1) }}{{ end }}
{{- define "cmd" -}}
{{ $base := mul .Depth 4 -}}
{{ indent $base }}{{ .Cmd.Name }}{{ with .Cmd.Short }}  # {{ . }}{{ end }}
{{ range $f := visibleFlags .Cmd }}{{ template "flag" (dict "Flag" $f "Base" $base) }}
{{ end -}}
{{ with .Cmd.Example }}
{{ indent (add $base 4) }}Example:
{{ prefixLines (add $base 6) . }}
{{ end -}}
{{ range $c := visibleCommands .Cmd }}
{{ template "cmd" (dict "Cmd" $c "Depth" (add $.Depth 1)) }}{{ end -}}
{{ end -}}
{{- define "flag" -}}
{{ indent (add .Base 4) }}{{ if isRequired .Flag }}--{{ .Flag.Name }} <{{ .Flag.Value.Type }}>{{ else }}[--{{ .Flag.Name }} <{{ .Flag.Value.Type }}>]{{ end }}{{ with .Flag.Usage }}  # {{ . }}{{ end }}{{ with .Flag.DefValue }}  (default: {{ . }}){{ end }}
{{- end -}}
`

func init() {
	funcs["dict"] = func(pairs ...any) map[string]any {
		m := make(map[string]any)
		for i := 0; i < len(pairs); i += 2 {
			m[pairs[i].(string)] = pairs[i+1]
		}
		return m
	}
}
