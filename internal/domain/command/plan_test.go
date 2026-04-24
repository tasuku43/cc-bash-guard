package command

import "testing"

func TestParseCommandPlanSimpleGitStatus(t *testing.T) {
	plan := Parse("git status")

	if plan.Shape.Kind != ShellShapeSimple {
		t.Fatalf("Shape.Kind = %q, want %q", plan.Shape.Kind, ShellShapeSimple)
	}
	if !plan.SafeForStructuredAllow {
		t.Fatal("SafeForStructuredAllow = false, want true")
	}
	if len(plan.Commands) != 1 {
		t.Fatalf("len(Commands) = %d, want 1", len(plan.Commands))
	}
	cmd := plan.Commands[0]
	if cmd.Program != "git" || cmd.ProgramToken != "git" {
		t.Fatalf("command program = (%q, %q), want git/git", cmd.Program, cmd.ProgramToken)
	}
	if len(cmd.ActionPath) != 1 || cmd.ActionPath[0] != "status" {
		t.Fatalf("ActionPath = %#v, want [status]", cmd.ActionPath)
	}
}

func TestParseCommandPlanAndListExtractsCommandsButFailsClosed(t *testing.T) {
	plan := Parse("git status && git diff")

	if plan.Shape.Kind != ShellShapeAndList {
		t.Fatalf("Shape.Kind = %q, want %q", plan.Shape.Kind, ShellShapeAndList)
	}
	if !plan.Shape.HasConditional {
		t.Fatal("HasConditional = false, want true")
	}
	if plan.SafeForStructuredAllow {
		t.Fatal("SafeForStructuredAllow = true, want false")
	}
	if len(plan.Commands) != 2 {
		t.Fatalf("len(Commands) = %d, want 2", len(plan.Commands))
	}
	if plan.Commands[0].Raw != "git status" || plan.Commands[1].Raw != "git diff" {
		t.Fatalf("command raws = %#v", []string{plan.Commands[0].Raw, plan.Commands[1].Raw})
	}
	assertNoShellConnectorMetadata(t, plan.Commands)
}

func TestParseCommandPlanPipelineExtractsCommandsButFailsClosed(t *testing.T) {
	plan := Parse("git status | sh")

	if plan.Shape.Kind != ShellShapePipeline {
		t.Fatalf("Shape.Kind = %q, want %q", plan.Shape.Kind, ShellShapePipeline)
	}
	if !plan.Shape.HasPipeline {
		t.Fatal("HasPipeline = false, want true")
	}
	if plan.SafeForStructuredAllow {
		t.Fatal("SafeForStructuredAllow = true, want false")
	}
	if len(plan.Commands) != 2 {
		t.Fatalf("len(Commands) = %d, want 2", len(plan.Commands))
	}
	if plan.Commands[0].Program != "git" || plan.Commands[1].Program != "sh" {
		t.Fatalf("command programs = %#v", []string{plan.Commands[0].Program, plan.Commands[1].Program})
	}
	assertNoShellConnectorMetadata(t, plan.Commands)
}

func TestParseCommandPlanUnsafeShellShapesFailClosed(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		kind ShellShapeKind
	}{
		{name: "sequence", raw: "git status; git diff", kind: ShellShapeSequence},
		{name: "background", raw: "git status &", kind: ShellShapeBackground},
		{name: "redirect", raw: "git status > /tmp/out", kind: ShellShapeRedirect},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := Parse(tt.raw)
			if plan.Shape.Kind != tt.kind {
				t.Fatalf("Shape.Kind = %q, want %q", plan.Shape.Kind, tt.kind)
			}
			if plan.SafeForStructuredAllow {
				t.Fatal("SafeForStructuredAllow = true, want false")
			}
		})
	}
}

func assertNoShellConnectorMetadata(t *testing.T, commands []Command) {
	t.Helper()
	for _, cmd := range commands {
		if cmd.Raw == "&&" || cmd.Raw == "||" || cmd.Raw == ";" || cmd.Raw == "|" || cmd.Raw == "&" {
			t.Fatalf("command contains connector as raw command: %+v", cmd)
		}
		for _, token := range append(append([]string{}, cmd.ActionPath...), cmd.Args...) {
			switch token {
			case "&&", "||", ";", "|", "&":
				t.Fatalf("command contains connector token %q: %+v", token, cmd)
			}
		}
	}
}
