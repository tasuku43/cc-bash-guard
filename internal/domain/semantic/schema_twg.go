package semantic

var twgSchema = Schema{
	Command:     "twg",
	order:       52,
	Description: "Atlassian Teamwork Graph CLI namespaces, verbs, and conservative read/write classification.",
	Parser:      "twg",
	Fields: []Field{
		stringField("namespace", "Top-level TWG namespace, with aliases normalized (for example bb to bitbucket)."),
		stringListField("namespace_in", "Allowed top-level TWG namespaces."),
		stringField("verb", "Effective TWG verb from a help-backed action path; read shorthands use their effective get or query verb."),
		stringListField("verb_in", "Allowed effective TWG verbs."),
		boolField("read_only", "True only for help/version, documented read-only namespaces, and help-backed read actions."),
		boolField("mutating", "True only for help-backed create/update/delete and other write actions."),
	},
	Examples: []Example{
		{Title: "Allow TWG read-only commands", YAML: `permission:
  allow:
    - command:
        name: twg
        semantic:
          read_only: true`},
	},
	Notes: []string{
		"Unknown actions and authentication/control-plane commands have both read_only and mutating set to false, so a read-only allow rule abstains.",
		"Classification is based on the TWG 1.0.25 help surface and exact action-path matching; positional words are never searched for a read verb.",
	},
}

func init() {
	RegisterSchema(twgSchema)
}
