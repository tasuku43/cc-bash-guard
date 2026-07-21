package semantic

var pupSchema = Schema{
	Command:     "pup",
	order:       45,
	Description: "Datadog pup CLI area, nested sub-area, leaf verb, and global option semantics.",
	Parser:      "pup",
	Fields: []Field{
		stringField("area", "Top-level pup area, such as logs, monitors, dashboards, metrics, or auth."),
		stringListField("area_in", "Allowed pup areas."),
		stringField("sub_area", "First nested pup area, such as archives in pup logs archives list."),
		stringListField("sub_area_in", "Allowed pup sub-areas."),
		stringField("verb", "Leaf pup verb such as list, get, aggregate, create, or delete."),
		stringListField("verb_in", "Allowed pup verbs."),
		stringField("org", "Organization selected by --org."),
		stringListField("org_in", "Allowed organizations selected by --org."),
		stringField("output", "Output format selected by -o or --output."),
		stringListField("output_in", "Allowed output formats."),
		boolField("yes", "True when --yes is present."),
		boolField("agent", "True when --agent is present."),
		boolField("no_agent", "True when --no-agent is present."),
		stringListField("flags_contains", "Parser-recognized pup option tokens that must be present; this does not scan raw argv words."),
		stringListField("flags_prefixes", "Parser-recognized pup option tokens that must start with these prefixes; this depends on the pup parser."),
	},
	Notes: []string{
		"Known pup action paths are generated from pup's command schema; paths with no known leaf prefix do not infer a verb.",
		"For action paths deeper than three words, sub_area is the first nested area and verb is the final leaf action.",
	},
}

func init() { RegisterSchema(pupSchema) }
