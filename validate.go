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

// Schema returns the embedded JSON schema.
func Schema() string {
	return schemaJSON
}

// Validate checks a yaml.Node against the yamlvfs schema.
func Validate(node *yaml.Node) error {
	var data any
	if err := node.Decode(&data); err != nil {
		return err
	}
	return schema.Validate(data)
}
