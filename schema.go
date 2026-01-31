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

// Schema returns the embedded JSON schema for yamlvfs documents.
func Schema() Document {
	return schemaJSON
}

// Validate checks a Document against the yamlvfs JSON schema.
func Validate(doc Document) error {
	var data any
	if err := yaml.Unmarshal([]byte(doc), &data); err != nil {
		return err
	}
	return schema.Validate(data)
}
