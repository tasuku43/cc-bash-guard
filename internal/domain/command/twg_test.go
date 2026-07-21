package command

import (
	"strings"
	"testing"
)

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
		{name: "alias outside declared parent remains unknown", raw: "twg jira prs query", namespace: "jira"},
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

func TestTWGParserNormalizesHelpSurfaceAliases(t *testing.T) {
	tests := []struct {
		name       string
		raw        string
		actionPath string
		namespace  string
		verb       string
		readOnly   bool
		mutating   bool
	}{
		{name: "canonical bitbucket pull requests", raw: "twg bitbucket pull-requests query", actionPath: "bitbucket pull-requests query", namespace: "bitbucket", verb: "query", readOnly: true},
		{name: "bitbucket pull requests alias", raw: "twg bitbucket prs query", actionPath: "bitbucket pull-requests query", namespace: "bitbucket", verb: "query", readOnly: true},
		{name: "bitbucket namespace and pull requests aliases", raw: "twg bb prs query", actionPath: "bitbucket pull-requests query", namespace: "bitbucket", verb: "query", readOnly: true},
		{name: "assets schema", raw: "twg assets schema get abc", actionPath: "assets objectschema get", namespace: "assets", verb: "get", readOnly: true},
		{name: "assets reference type", raw: "twg assets referencetype create", actionPath: "assets reference-type create", namespace: "assets", verb: "create", mutating: true},
		{name: "assets list attributes", raw: "twg assets type list-attributes query", actionPath: "assets type list-attr query", namespace: "assets", verb: "query", readOnly: true},
		{name: "confluence comment", raw: "twg confluence content comment create 123", actionPath: "confluence content comments create", namespace: "confluence", verb: "create", mutating: true},
		{name: "confluence label", raw: "twg confluence content label list 123", actionPath: "confluence content labels list", namespace: "confluence", verb: "list", readOnly: true},
		{name: "csm organisation", raw: "twg csm organisation", actionPath: "csm organization", namespace: "csm", verb: "organization", readOnly: true},
		{name: "jira custom fields", raw: "twg jira workitem custom-fields create-metadata", actionPath: "jira workitem field create-metadata", namespace: "jira", verb: "create-metadata", readOnly: true},
		{name: "jsm help centre", raw: "twg jsm help-centre query", actionPath: "jsm help-center query", namespace: "jsm", verb: "query", readOnly: true},
		{name: "jsm affected services", raw: "twg jsm incident affected-services query INC-1", actionPath: "jsm incident affected-service query", namespace: "jsm", verb: "query", readOnly: true},
		{name: "jsm pir", raw: "twg jsm pir delete PIR-1", actionPath: "jsm post-incident-review delete", namespace: "jsm", verb: "delete", mutating: true},
		{name: "jsm request type fields", raw: "twg jsm request-type fields get 1", actionPath: "jsm request-type field get", namespace: "jsm", verb: "get", readOnly: true},
		{name: "loom videos", raw: "twg loom videos get 1", actionPath: "loom video get", namespace: "loom", verb: "get", readOnly: true},
		{name: "rovo list connectors", raw: "twg rovo list-connectors", actionPath: "rovo list-apps", namespace: "rovo", verb: "list-apps", readOnly: true},
		{name: "search list connectors", raw: "twg search list-connectors", actionPath: "search list-apps", namespace: "search", verb: "list-apps", readOnly: true},
		{name: "teams search", raw: "twg teams search platform", actionPath: "teams query", namespace: "teams", verb: "query", readOnly: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := singleParsedCommand(t, tt.raw)
			got := cmd.TWG
			if got == nil {
				t.Fatal("TWG semantic is nil")
			}
			if actionPath := strings.Join(cmd.ActionPath, " "); actionPath != tt.actionPath {
				t.Fatalf("ActionPath = %q, want %q", actionPath, tt.actionPath)
			}
			if got.Namespace != tt.namespace || got.Verb != tt.verb || got.ReadOnly != tt.readOnly || got.Mutating != tt.mutating {
				t.Fatalf("TWG=%+v, want namespace=%q verb=%q readOnly=%v mutating=%v", got, tt.namespace, tt.verb, tt.readOnly, tt.mutating)
			}
		})
	}
}
