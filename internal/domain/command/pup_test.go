package command

import (
	"reflect"
	"testing"
)

func TestPupParserExtractsSemanticFields(t *testing.T) {
	tests := []struct {
		name        string
		raw         string
		wantArea    string
		wantSubArea string
		wantVerb    string
		wantOrg     string
		wantOutput  string
		wantYes     bool
		wantAgent   bool
		wantNoAgent bool
		wantAction  []string
		wantArgs    []string
	}{
		{name: "2-level", raw: `pup logs list --query="*" --limit=3`, wantArea: "logs", wantVerb: "list", wantAction: []string{"logs", "list"}},
		{name: "2-level positional argument", raw: "pup monitors get 123", wantArea: "monitors", wantVerb: "get", wantAction: []string{"monitors", "get"}, wantArgs: []string{"123"}},
		{name: "3-level", raw: "pup logs archives list", wantArea: "logs", wantSubArea: "archives", wantVerb: "list", wantAction: []string{"logs", "archives", "list"}},
		{name: "3-level positional argument", raw: "pup logs metrics delete abc", wantArea: "logs", wantSubArea: "metrics", wantVerb: "delete", wantAction: []string{"logs", "metrics", "delete"}, wantArgs: []string{"abc"}},
		{name: "deep action", raw: "pup integrations aws cloud-auth persona-mappings get example", wantArea: "integrations", wantSubArea: "aws", wantVerb: "get", wantAction: []string{"integrations", "aws", "cloud-auth", "persona-mappings", "get"}, wantArgs: []string{"example"}},
		{name: "global flags before area", raw: "pup --org acme -oyaml --jq . dashboards create --title=foo", wantArea: "dashboards", wantVerb: "create", wantOrg: "acme", wantOutput: "yaml", wantAction: []string{"dashboards", "create"}},
		{name: "long output flag", raw: "pup --output yaml dashboards create --title=foo", wantArea: "dashboards", wantVerb: "create", wantOutput: "yaml", wantAction: []string{"dashboards", "create"}},
		{name: "yes alias and agent", raw: "pup -y --agent logs metrics delete abc", wantArea: "logs", wantSubArea: "metrics", wantVerb: "delete", wantYes: true, wantAgent: true, wantAction: []string{"logs", "metrics", "delete"}, wantArgs: []string{"abc"}},
		{name: "no-agent", raw: "pup --no-agent auth status", wantArea: "auth", wantVerb: "status", wantNoAgent: true, wantAction: []string{"auth", "status"}},
		{name: "unknown action is not guessed", raw: "pup future query delete thing", wantArea: "future", wantAction: []string{"future"}, wantArgs: []string{"query", "delete", "thing"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := singleParsedCommand(t, tt.raw)
			if got.SemanticParser != "pup" || got.Pup == nil {
				t.Fatalf("Pup semantic missing: %+v", got)
			}
			p := got.Pup
			if p.Area != tt.wantArea || p.SubArea != tt.wantSubArea || p.Verb != tt.wantVerb || p.Org != tt.wantOrg ||
				p.Output != tt.wantOutput || p.Yes != tt.wantYes || p.Agent != tt.wantAgent || p.NoAgent != tt.wantNoAgent {
				t.Fatalf("Pup=%+v", p)
			}
			if !reflect.DeepEqual(got.ActionPath, tt.wantAction) {
				t.Fatalf("ActionPath=%v, want %v", got.ActionPath, tt.wantAction)
			}
			if !reflect.DeepEqual(got.Args, tt.wantArgs) {
				t.Fatalf("Args=%v, want %v", got.Args, tt.wantArgs)
			}
		})
	}
}

func TestPupGeneratedActionInventory(t *testing.T) {
	if pupCommandSchemaVersion == "" {
		t.Fatal("generated pup command schema version is empty")
	}
	for _, path := range []string{"monitors get", "logs archives list", "integrations aws cloud-auth persona-mappings get"} {
		if !pupKnownActionPaths[path] {
			t.Fatalf("generated pup action inventory missing %q", path)
		}
	}
}
