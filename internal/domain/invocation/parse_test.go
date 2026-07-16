package invocation

import (
	"reflect"
	"testing"
)

func TestParseUnwrapsCommonWrappers(t *testing.T) {
	parsed := Parse("sudo -u root /usr/bin/env bash -c 'echo hi'")
	if parsed.Command != "bash" {
		t.Fatalf("Command = %q", parsed.Command)
	}
	if len(parsed.Args) < 2 || parsed.Args[0] != "-c" || parsed.Args[1] != "echo hi" {
		t.Fatalf("Args = %#v", parsed.Args)
	}
}

func TestParseUnwrapsTimeWrapper(t *testing.T) {
	tests := []string{
		"time git status",
		"/usr/bin/time git status",
		"time -p git status",
		"time -pv git status",
		"time -f %E git status",
		"time --format=%E git status",
		"time -o /tmp/time.txt git status",
		"time --output=/tmp/time.txt git status",
		"time -- git status",
	}

	for _, command := range tests {
		t.Run(command, func(t *testing.T) {
			parsed := Parse(command)
			if parsed.Command != "git" {
				t.Fatalf("Command = %q, want git; parsed=%+v", parsed.Command, parsed)
			}
			if !reflect.DeepEqual(parsed.Args, []string{"status"}) {
				t.Fatalf("Args = %#v, want [status]", parsed.Args)
			}
		})
	}
}

func TestParseKeepsUnsupportedOrIncompleteTimeWrapper(t *testing.T) {
	tests := []string{
		"time --unknown git status",
		"time --format",
		"time -o",
		"time --help",
	}

	for _, command := range tests {
		t.Run(command, func(t *testing.T) {
			parsed := Parse(command)
			if parsed.Command != "time" {
				t.Fatalf("Command = %q, want time; parsed=%+v", parsed.Command, parsed)
			}
		})
	}
}

func TestTokensPreserveQuotedPayload(t *testing.T) {
	got := Tokens("bash -c 'git status'")
	if len(got) != 3 || got[2] != "git status" {
		t.Fatalf("Tokens() = %#v", got)
	}
}

func TestJoinRoundTripPreservesQuotedArgs(t *testing.T) {
	command := `aws s3 cp "hello world" s3://bucket/key`

	want := Tokens(command)
	got := Tokens(Join(want))
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Tokens(Join(Tokens(command))) = %#v, want %#v", got, want)
	}
}

func TestJoinRoundTripPreservesEnvAssignmentWithSpaces(t *testing.T) {
	command := `FOO="hello world" env`

	want := Tokens(command)
	got := Tokens(Join(want))
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Tokens(Join(Tokens(command))) = %#v, want %#v", got, want)
	}
}

func TestJoinRoundTripPreservesSingleQuotesInToken(t *testing.T) {
	command := `printf "%s\n" "it's fine"`

	want := Tokens(command)
	got := Tokens(Join(want))
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Tokens(Join(Tokens(command))) = %#v, want %#v", got, want)
	}
}

func TestIsEnvAssignmentAcceptsEmptyValue(t *testing.T) {
	if !IsEnvAssignment("FOO=") {
		t.Fatal("expected empty env assignment to be treated as env assignment")
	}
}

func TestIsSafeSingleCommandRejectsCompoundPayload(t *testing.T) {
	if IsSafeSingleCommand("git status && git diff") {
		t.Fatal("expected safe single command check to fail")
	}
}

func TestIsASTSafeSimpleCommandRejectsUnsafeShellForms(t *testing.T) {
	tests := []struct {
		name    string
		command string
	}{
		{name: "and list", command: "git status && git diff"},
		{name: "redirect", command: "git status > /tmp/out"},
		{name: "command substitution", command: "git status $(whoami)"},
		{name: "process substitution", command: "cat <(whoami)"},
		{name: "heredoc", command: "cat <<EOF\nhi\nEOF"},
		{name: "background", command: "git status &"},
		{name: "comment", command: "git status # harmless"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if IsASTSafeSimpleCommand(tt.command) {
				t.Fatalf("IsASTSafeSimpleCommand(%q) = true, want false", tt.command)
			}
		})
	}
}

func TestClassify(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    CommandClass
	}{
		{name: "simple", command: "git status", want: CommandClassSimple},
		{name: "env prefixed", command: "AWS_PROFILE=dev git status", want: CommandClassEnvPrefixedSimple},
		{name: "wrapper prefixed env", command: "env AWS_PROFILE=dev git status", want: CommandClassWrapperPrefixed},
		{name: "wrapper prefixed sudo", command: "sudo -u root git status", want: CommandClassWrapperPrefixed},
		{name: "wrapper prefixed time", command: "time git status", want: CommandClassWrapperPrefixed},
		{name: "compound and", command: "git status && rm -rf /tmp/x", want: CommandClassUnsafeCompound},
		{name: "compound semicolon", command: "git status; rm -rf /tmp/x", want: CommandClassUnsafeCompound},
		{name: "pipeline", command: "git status | sh", want: CommandClassUnsafeCompound},
		{name: "redirect", command: "git status > /tmp/out", want: CommandClassUnsafeCompound},
		{name: "command substitution", command: "git status $(whoami)", want: CommandClassUnsafeCompound},
		{name: "comment", command: "git status # harmless", want: CommandClassUnsafeCompound},
		{name: "bash c unsafe", command: "bash -c 'git status && rm -rf /tmp/x'", want: CommandClassUnsafeCompound},
		{name: "bash c redirect", command: "bash -c 'git status > /tmp/out'", want: CommandClassUnsafeCompound},
		{name: "bash c command substitution", command: "bash -c 'git status $(whoami)'", want: CommandClassUnsafeCompound},
		{name: "bash c process substitution", command: "bash -c 'cat <(whoami)'", want: CommandClassUnsafeCompound},
		{name: "bash c heredoc", command: "bash -c 'cat <<EOF\nhi\nEOF'", want: CommandClassUnsafeCompound},
		{name: "bash c background", command: "bash -c 'git status &'", want: CommandClassUnsafeCompound},
		{name: "bash c comment", command: "bash -c 'git status # harmless'", want: CommandClassUnsafeCompound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Classify(tt.command); got != tt.want {
				t.Fatalf("Classify(%q) = %q, want %q", tt.command, got, tt.want)
			}
		})
	}
}
