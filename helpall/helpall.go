// Package helpall provides a --help-all flag for Cobra CLI applications.
package helpall

import (
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Install adds --help-all to the root command.
func Install(root *cobra.Command) {
	root.PersistentFlags().Bool("help-all", false, "Show full command reference")

	orig := root.HelpFunc()
	root.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if v, _ := cmd.Flags().GetBool("help-all"); v {
			printAll(root)
			return
		}
		orig(cmd, args)
	})
}

func printAll(root *cobra.Command) {
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

const tmpl = `{{ .Name }}{{ with .Short }} - {{ . }}{{ end }}
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
