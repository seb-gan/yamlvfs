package yamlvfs

import (
	_ "embed"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"
)

//go:embed yamlvfs.schema.json
var schemaJSON string

var schema *jsonschema.Schema

func init() {
	doc, _ := jsonschema.UnmarshalJSON(strings.NewReader(schemaJSON))
	c := jsonschema.NewCompiler()
	c.AddResource("schema.json", doc)
	schema, _ = c.Compile("schema.json")
}

// Validate checks if the YAML content conforms to the YAML VFS schema.
func Validate(content string) error {
	var data any
	if err := yaml.Unmarshal([]byte(content), &data); err != nil {
		return err
	}

	return schema.Validate(data)
}
