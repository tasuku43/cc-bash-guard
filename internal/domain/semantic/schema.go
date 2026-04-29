package semantic

type Schema struct {
	Command      string    `json:"command"`
	SemanticPath string    `json:"semantic_path"`
	Description  string    `json:"description"`
	Parser       string    `json:"parser"`
	Fields       []Field   `json:"fields"`
	Examples     []Example `json:"examples"`
	Notes        []string  `json:"notes,omitempty"`
}

type Field struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Values      []string `json:"values,omitempty"`
	Since       string   `json:"since,omitempty"`
}

type Example struct {
	Title string `json:"title"`
	YAML  string `json:"yaml"`
}

var schemas = []Schema{
	gitSchema,
	awsSchema,
	kubectlSchema,
	ghSchema,
	argoCDSchema,
	gwsSchema,
	helmfileSchema,
	terraformSchema,
}

func AllSchemas() []Schema {
	out := make([]Schema, 0, len(schemas))
	for _, schema := range schemas {
		out = append(out, withSemanticPath(schema))
	}
	return out
}

func Lookup(command string) (Schema, bool) {
	for _, schema := range schemas {
		if schema.Command == command {
			return withSemanticPath(schema), true
		}
	}
	return Schema{}, false
}

func SchemasByCommand() map[string]Schema {
	byCommand := map[string]Schema{}
	for _, schema := range AllSchemas() {
		byCommand[schema.Command] = schema
	}
	return byCommand
}

func SupportedCommands() []string {
	commands := make([]string, 0, len(schemas))
	for _, schema := range schemas {
		commands = append(commands, schema.Command)
	}
	return commands
}

func FieldNames(command string) []string {
	schema, ok := Lookup(command)
	if !ok {
		return nil
	}
	names := make([]string, 0, len(schema.Fields))
	for _, field := range schema.Fields {
		names = append(names, field.Name)
	}
	return names
}

func IsFieldSupported(command, field string) bool {
	for _, supported := range FieldNames(command) {
		if supported == field {
			return true
		}
	}
	return false
}

func withSemanticPath(schema Schema) Schema {
	schema.SemanticPath = "command.semantic"
	return schema
}

func stringField(name, description string) Field {
	return Field{Name: name, Type: "string", Description: description}
}

func stringListField(name, description string) Field {
	return Field{Name: name, Type: "[]string", Description: description}
}

func boolField(name, description string) Field {
	return Field{Name: name, Type: "bool", Description: description}
}
