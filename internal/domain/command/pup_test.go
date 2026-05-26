package command

import "testing"

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
	}{
		{name: "2-level", raw: `pup logs list --query="*" --limit=3`, wantArea: "logs", wantVerb: "list"},
		{name: "3-level", raw: "pup logs archives list", wantArea: "logs", wantSubArea: "archives", wantVerb: "list"},
		{name: "global flags before area", raw: "pup --org acme --output yaml dashboards create --title=foo", wantArea: "dashboards", wantVerb: "create", wantOrg: "acme", wantOutput: "yaml"},
		{name: "yes and agent", raw: "pup --yes --agent logs metrics delete abc", wantArea: "logs", wantSubArea: "metrics", wantVerb: "delete", wantYes: true, wantAgent: true},
		{name: "no-agent", raw: "pup --no-agent auth status", wantArea: "auth", wantVerb: "status", wantNoAgent: true},
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
		})
	}
}
