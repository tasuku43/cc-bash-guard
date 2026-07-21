package policy

import "testing"

func TestEvaluatePupSemanticRules(t *testing.T) {
	trueValue := true
	p := NewPipeline(PipelineSpec{
		Permission: PermissionSpec{
			Deny: []PermissionRuleSpec{{
				Command: PermissionCommandSpec{Name: "pup", Semantic: &SemanticMatchSpec{
					VerbIn: []string{"create", "delete", "update", "enable", "disable", "run"},
					Yes:    &trueValue,
				}},
			}},
			Allow: []PermissionRuleSpec{{
				Command: PermissionCommandSpec{Name: "pup", Semantic: &SemanticMatchSpec{
					VerbIn: []string{"list", "get", "status", "search", "query", "aggregate"},
				}},
			}},
		},
	}, Source{})

	tests := []struct {
		command string
		want    string
	}{
		{command: "pup monitors get 123", want: "allow"},
		{command: "pup logs archives list", want: "allow"},
		{command: "pup integrations aws cloud-auth persona-mappings get example", want: "allow"},
		{command: "pup logs metrics delete abc", want: "abstain"},
		{command: "pup -y logs metrics delete abc", want: "deny"},
		{command: "pup feature-flags enable feature-id", want: "abstain"},
		{command: "pup --yes feature-flags enable feature-id", want: "deny"},
		{command: "pup future query delete thing", want: "abstain"},
	}
	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			got, err := Evaluate(p, tt.command)
			if err != nil {
				t.Fatalf("Evaluate() error = %v", err)
			}
			if got.Outcome != tt.want {
				t.Fatalf("Outcome = %q, want %q; decision=%+v", got.Outcome, tt.want, got)
			}
		})
	}
}
