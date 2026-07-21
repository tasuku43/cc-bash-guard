package command

import "testing"

func TestTWGParserClassifiesCommandsFromHelpSurface(t *testing.T) {
	tests := []struct {
		name      string
		raw       string
		namespace string
		verb      string
		readOnly  bool
		mutating  bool
	}{
		{name: "confluence get", raw: "twg confluence content get 12345", namespace: "confluence", verb: "get", readOnly: true},
		{name: "jira shorthand", raw: "twg jira workitem PROJ-123", namespace: "jira", verb: "get", readOnly: true},
		{name: "leading global output", raw: "twg -o json jira workitem get PROJ-123", namespace: "jira", verb: "get", readOnly: true},
		{name: "optional global output summary", raw: "twg --output-summary jira workitem get PROJ-123", namespace: "jira", verb: "get", readOnly: true},
		{name: "search namespace", raw: "twg search 'deployment failures'", namespace: "search", verb: "search", readOnly: true},
		{name: "doctor", raw: "twg doctor", namespace: "doctor", verb: "doctor", readOnly: true},
		{name: "help namespace", raw: "twg help jira workitem", namespace: "help", verb: "help", readOnly: true},
		{name: "help flag overrides delete", raw: "twg confluence content delete 12345 --help", namespace: "confluence", verb: "help", readOnly: true},
		{name: "version", raw: "twg --version", verb: "version", readOnly: true},
		{name: "short version", raw: "twg -V", verb: "version", readOnly: true},
		{name: "read alias normalized", raw: "twg whoami", namespace: "user", verb: "user", readOnly: true},
		{name: "confluence delete", raw: "twg confluence content delete 12345", namespace: "confluence", verb: "delete", mutating: true},
		{name: "jira create", raw: "twg jira workitem create --space PROJ --summary get", namespace: "jira", verb: "create", mutating: true},
		{name: "jira link target", raw: "twg jira workitem link goal PROJ-123 GOAL-1", namespace: "jira", verb: "link", mutating: true},
		{name: "read verb positional does not override write", raw: "twg confluence content create get", namespace: "confluence", verb: "create", mutating: true},
		{name: "api remains unknown", raw: "twg api /rest/api/3/myself", namespace: "api"},
		{name: "login remains unknown", raw: "twg login", namespace: "login"},
		{name: "rovo auth remains unknown", raw: "twg rovo auth slack", namespace: "rovo"},
		{name: "unknown mixed path remains unknown", raw: "twg confluence unknown get", namespace: "confluence"},
		{name: "option before write does not preserve shorthand read", raw: "twg jira workitem --unknown value create", namespace: "jira"},
		{name: "explicit write remains write before unknown option", raw: "twg jira workitem create --unknown get", namespace: "jira", verb: "create", mutating: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := singleParsedCommand(t, tt.raw)
			if cmd.Parser != "twg" || cmd.SemanticParser != "twg" || cmd.TWG == nil {
				t.Fatalf("parser state = (%q, %q, %v), want twg semantic", cmd.Parser, cmd.SemanticParser, cmd.TWG)
			}
			got := cmd.TWG
			if got.Namespace != tt.namespace || got.Verb != tt.verb || got.ReadOnly != tt.readOnly || got.Mutating != tt.mutating {
				t.Fatalf("TWG=%+v, want namespace=%q verb=%q readOnly=%v mutating=%v", got, tt.namespace, tt.verb, tt.readOnly, tt.mutating)
			}
		})
	}
}
