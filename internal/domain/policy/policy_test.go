package policy

import "testing"

func TestEvaluateFirstMatchWins(t *testing.T) {
	rules := []Rule{
		NewRule(RuleSpec{ID: "first", Pattern: "^git", Message: "first"}, Source{}),
		NewRule(RuleSpec{ID: "second", Pattern: "status$", Message: "second"}, Source{}),
	}

	got, err := Evaluate(rules, "git status")
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got.Outcome != "reject" || got.Rule == nil || got.Rule.ID != "first" {
		t.Fatalf("got %+v", got)
	}
}

func TestEvaluatePredicateRule(t *testing.T) {
	rules := []Rule{
		NewRule(RuleSpec{
			ID: "no-shell-dash-c",
			Matcher: MatchSpec{
				CommandIn:    []string{"bash", "sh"},
				ArgsContains: []string{"-c"},
			},
			Message: "blocked",
		}, Source{}),
	}

	got, err := Evaluate(rules, "/usr/bin/env bash -c 'echo hi'")
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got.Outcome != "reject" || got.Rule == nil || got.Rule.ID != "no-shell-dash-c" {
		t.Fatalf("got %+v", got)
	}
}

func TestEvaluateRewriteRule(t *testing.T) {
	rules := []Rule{
		NewRule(RuleSpec{
			ID: "unwrap-shell-dash-c",
			Matcher: MatchSpec{
				CommandIn:    []string{"bash", "sh"},
				ArgsContains: []string{"-c"},
			},
			Rewrite: RewriteSpec{
				UnwrapShellDashC: true,
			},
			BlockExamples: []string{"bash -c 'git status'"},
			AllowExamples: []string{"bash script.sh"},
		}, Source{}),
	}

	got, err := Evaluate(rules, "bash -c 'git status'")
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got.Outcome != "rewrite" || got.Command != "git status" {
		t.Fatalf("got %+v", got)
	}
}

func TestEvaluateMoveFlagToEnvRewriteRule(t *testing.T) {
	rules := []Rule{
		NewRule(RuleSpec{
			ID: "aws-profile-to-env",
			Matcher: MatchSpec{
				Command:      "aws",
				ArgsContains: []string{"--profile"},
			},
			Rewrite: RewriteSpec{
				MoveFlagToEnv: MoveFlagToEnvSpec{
					Flag: "--profile",
					Env:  "AWS_PROFILE",
				},
			},
			BlockExamples: []string{"aws --profile read-only-profile s3 ls"},
			AllowExamples: []string{"AWS_PROFILE=read-only-profile aws s3 ls"},
		}, Source{}),
	}

	got, err := Evaluate(rules, "aws --profile read-only-profile s3 ls")
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got.Outcome != "rewrite" || got.Command != "AWS_PROFILE=read-only-profile aws s3 ls" {
		t.Fatalf("got %+v", got)
	}
}

func TestEvaluateMoveEnvToFlagRewriteRule(t *testing.T) {
	rules := []Rule{
		NewRule(RuleSpec{
			ID: "aws-env-to-profile",
			Matcher: MatchSpec{
				Command:     "aws",
				EnvRequires: []string{"AWS_PROFILE"},
			},
			Rewrite: RewriteSpec{
				MoveEnvToFlag: MoveEnvToFlagSpec{
					Env:  "AWS_PROFILE",
					Flag: "--profile",
				},
			},
			BlockExamples: []string{"AWS_PROFILE=read-only-profile aws s3 ls"},
			AllowExamples: []string{"aws --profile read-only-profile s3 ls"},
		}, Source{}),
	}

	got, err := Evaluate(rules, "AWS_PROFILE=read-only-profile aws s3 ls")
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got.Outcome != "rewrite" || got.Command != "aws --profile read-only-profile s3 ls" {
		t.Fatalf("got %+v", got)
	}
}

func TestEvaluateUnwrapWrapperRewriteRule(t *testing.T) {
	rules := []Rule{
		NewRule(RuleSpec{
			ID: "unwrap-safe-wrappers",
			Pattern: `^\s*(env|command|exec)\b`,
			Rewrite: RewriteSpec{
				UnwrapWrapper: UnwrapWrapperSpec{
					Wrappers: []string{"env", "command", "exec"},
				},
			},
			BlockExamples: []string{"env AWS_PROFILE=dev command exec aws s3 ls"},
			AllowExamples: []string{"AWS_PROFILE=dev aws s3 ls"},
		}, Source{}),
	}

	got, err := Evaluate(rules, "env AWS_PROFILE=dev command exec aws s3 ls")
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got.Outcome != "rewrite" || got.Command != "AWS_PROFILE=dev aws s3 ls" {
		t.Fatalf("got %+v", got)
	}
}

func TestValidateDirectiveKinds(t *testing.T) {
	issues := ValidateDirective("rules[0]", "legacy", RejectSpec{Message: "new"}, RewriteSpec{})
	if len(issues) != 1 {
		t.Fatalf("issues = %#v", issues)
	}
}

func TestValidateRewriteRejectsMultiplePrimitives(t *testing.T) {
	issues := ValidateRewrite("rules[0].rewrite", RewriteSpec{
		UnwrapShellDashC: true,
		MoveFlagToEnv: MoveFlagToEnvSpec{
			Flag: "--profile",
			Env:  "AWS_PROFILE",
		},
	})
	if len(issues) != 1 || issues[0] != "rules[0].rewrite must set exactly one rewrite primitive" {
		t.Fatalf("issues = %#v", issues)
	}
}
