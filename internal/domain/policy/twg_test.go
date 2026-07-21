package policy

import "testing"

func TestEvaluateTWGSemanticReadOnlyRule(t *testing.T) {
	trueValue := true
	p := NewPipeline(PipelineSpec{
		Permission: PermissionSpec{
			Allow: []PermissionRuleSpec{
				{Command: PermissionCommandSpec{Name: "cd"}},
				{Command: PermissionCommandSpec{Name: "twg", Semantic: &SemanticMatchSpec{ReadOnly: &trueValue}}},
			},
		},
	}, Source{})

	tests := []struct {
		command string
		want    string
	}{
		{command: "twg confluence content get 12345", want: "allow"},
		{command: "twg jira workitem PROJ-123", want: "allow"},
		{command: "twg -o json jira workitem get PROJ-123", want: "allow"},
		{command: "twg search 'deployment failures'", want: "allow"},
		{command: "twg doctor", want: "allow"},
		{command: "twg confluence content delete 12345", want: "abstain"},
		{command: "twg jira workitem create --space PROJ --summary get", want: "abstain"},
		{command: "twg api /rest/api/3/myself", want: "abstain"},
		{command: "twg login", want: "abstain"},
		{command: "cd repo && twg jira workitem PROJ-123", want: "allow"},
		{command: "bash -c 'twg confluence content get 12345'", want: "allow"},
		{command: "bash -c 'twg confluence content delete 12345'", want: "abstain"},
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

func TestEvaluateTWGNamespaceAndVerbRules(t *testing.T) {
	trueValue := true
	p := NewPipeline(PipelineSpec{
		Permission: PermissionSpec{
			Ask: []PermissionRuleSpec{{
				Command: PermissionCommandSpec{Name: "twg", Semantic: &SemanticMatchSpec{Mutating: &trueValue}},
			}},
			Allow: []PermissionRuleSpec{{
				Command: PermissionCommandSpec{Name: "twg", Semantic: &SemanticMatchSpec{
					Namespace: "jira", VerbIn: []string{"get", "query", "search"}, ReadOnly: &trueValue,
				}},
			}},
		},
	}, Source{})

	tests := []struct {
		command string
		want    string
	}{
		{command: "twg jira workitem get PROJ-123", want: "allow"},
		{command: "twg confluence content get 12345", want: "abstain"},
		{command: "twg jira workitem create --space PROJ --summary list", want: "ask"},
		{command: "twg login", want: "abstain"},
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
