package policy

import "testing"

func TestTerraformSemanticMatches(t *testing.T) {
	tests := []struct {
		name    string
		command string
		match   MatchSpec
		want    bool
	}{
		{name: "allow read only subcommands", command: "terraform plan -out=tfplan", match: MatchSpec{Command: "terraform", Semantic: &SemanticMatchSpec{SubcommandIn: []string{"init", "validate", "plan", "show", "output", "providers", "graph", "fmt"}}}, want: true},
		{name: "read only does not match apply", command: "terraform apply -auto-approve", match: MatchSpec{Command: "terraform", Semantic: &SemanticMatchSpec{SubcommandIn: []string{"validate", "plan", "show", "output"}}}, want: false},
		{name: "ask apply", command: "terraform apply tfplan", match: MatchSpec{Command: "terraform", Semantic: &SemanticMatchSpec{Subcommand: "apply"}}, want: true},
		{name: "deny destroy auto approve", command: "terraform destroy -auto-approve", match: MatchSpec{Command: "terraform", Semantic: &SemanticMatchSpec{Subcommand: "destroy", AutoApprove: boolPtr(true)}}, want: true},
		{name: "destroy auto approve does not match plan destroy", command: "terraform plan -destroy", match: MatchSpec{Command: "terraform", Semantic: &SemanticMatchSpec{Subcommand: "destroy", AutoApprove: boolPtr(true)}}, want: false},
		{name: "plan destroy", command: "terraform plan -destroy", match: MatchSpec{Command: "terraform", Semantic: &SemanticMatchSpec{Subcommand: "plan", Destroy: boolPtr(true)}}, want: true},
		{name: "state rm", command: "terraform state rm aws_instance.bad", match: MatchSpec{Command: "terraform", Semantic: &SemanticMatchSpec{Subcommand: "state", StateSubcommand: "rm"}}, want: true},
		{name: "global chdir", command: "terraform -chdir infra plan", match: MatchSpec{Command: "terraform", Semantic: &SemanticMatchSpec{Subcommand: "plan", GlobalChdir: "infra"}}, want: true},
		{name: "workspace delete", command: "terraform workspace delete old", match: MatchSpec{Command: "terraform", Semantic: &SemanticMatchSpec{Subcommand: "workspace", WorkspaceSubcommand: "delete"}}, want: true},
		{name: "target", command: "terraform plan -target=module.foo", match: MatchSpec{Command: "terraform", Semantic: &SemanticMatchSpec{Subcommand: "plan", Target: boolPtr(true), TargetsContains: []string{"module.foo"}}}, want: true},
		{name: "replace", command: "terraform apply -replace=aws_instance.bad", match: MatchSpec{Command: "terraform", Semantic: &SemanticMatchSpec{Subcommand: "apply", Replace: boolPtr(true), ReplacesContains: []string{"aws_instance.bad"}}}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.match.MatchMatches(tt.command); got != tt.want {
				t.Fatalf("MatchMatches(%q) = %v, want %v", tt.command, got, tt.want)
			}
		})
	}
}
