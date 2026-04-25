package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"

	semanticpkg "github.com/tasuku43/cc-bash-guard/internal/domain/semantic"
)

func writeUsage(w io.Writer) {
	fmt.Fprint(w, `cc-bash-guard

Declarative, testable command policy for AI-agent shell commands.

Typical workflow:
  1. Edit ~/.config/cc-bash-guard/cc-bash-guard.yml
  2. Optionally add .cc-bash-guard/cc-bash-guard.yaml in the project
  3. Add permission and E2E tests
  4. Run cc-bash-guard verify
  5. Let Claude Code call cc-bash-guard hook from PreToolUse

Usage:
  cc-bash-guard <command> [flags]

Commands:
  init     create the user config and print the Claude Code hook snippet
  doctor   inspect config quality and installation state
  verify   verify config tests, trust-critical setup, and build metadata
  version  print build and source metadata for the running binary
  hook     Claude Code hook entrypoint
  semantic-schema
           print supported semantic match schemas

Help:
  cc-bash-guard help <command>
  cc-bash-guard <command> --help
  cc-bash-guard help config
  cc-bash-guard help match
  cc-bash-guard help semantic
  cc-bash-guard help semantic git

Examples:
  cc-bash-guard init
  cc-bash-guard verify --format json
  cc-bash-guard semantic-schema --format json
  cc-bash-guard version --format json
  cc-bash-guard hook --rtk
  cc-bash-guard doctor --format json
`)
}

func writeHelp(stdout, stderr io.Writer, args []string) int {
	if len(args) == 0 {
		writeUsage(stdout)
		return exitAllow
	}
	if args[0] == "semantic" {
		if err := writeSemanticHelp(stdout, args[1:]); err != nil {
			writeErr(stderr, err.Error())
			return exitError
		}
		return exitAllow
	}
	writeCommandHelp(stdout, args[0])
	return exitAllow
}

func writeCommandHelp(w io.Writer, command string) {
	switch command {
	case "init":
		fmt.Fprint(w, `cc-bash-guard init

Create ~/.config/cc-bash-guard/cc-bash-guard.yml when it does not exist and print the
Claude Code PreToolUse hook snippet.

Usage:
  cc-bash-guard init

Typical use:
  cc-bash-guard init
`)
	case "doctor":
		fmt.Fprint(w, `cc-bash-guard doctor

Inspect config validity, pipeline quality, and Claude Code hook registration.

Usage:
  cc-bash-guard doctor [--format json]

Examples:
  cc-bash-guard doctor
  cc-bash-guard doctor --format json
`)
	case "verify":
		fmt.Fprint(w, `cc-bash-guard verify

Verify the local trust-critical cc-bash-guard setup.
This command is stricter than doctor: it fails when the config is broken, when
configured tests fail, when the effective global/local tool settings and
cc-bash-guard policy disagree with expected E2E outcomes, or when build metadata is
missing.

Usage:
  cc-bash-guard verify [--format json]

Examples:
  cc-bash-guard verify
  cc-bash-guard verify --format json
`)
	case "hook":
		fmt.Fprint(w, `cc-bash-guard hook

Claude Code hook entrypoint.
Reads stdin JSON, parses the command, evaluates permission policy, and
returns Claude Code hook JSON for allow, ask, deny, or error outcomes.

Usage:
  cc-bash-guard hook [--rtk] [--auto-verify]

Options:
  --rtk          run "rtk rewrite" once after cc-bash-guard policy evaluation
  --auto-verify  regenerate verified hook artifacts when they are missing or stale

Note:
  You usually do not run this manually. Edit rules and use cc-bash-guard verify
  while authoring policy instead. Without --auto-verify, the hook fails closed
  when verified artifacts are missing or stale. --auto-verify is convenient, but
  it lets hook-time config changes become active without a separate review step.

RTK compatibility:
  --rtk applies rtk rewrite after cc-bash-guard permission evaluation in the same
  hook invocation. Permission checks therefore see the command before the rtk
  rename/rewrite is applied. Stacking multiple Bash hooks can make the visible
  renamed command differ from the command that cc-bash-guard checked. The old
  command was cmdproxy hook claude --rtk; use cc-bash-guard hook --rtk now.
`)
	case "version":
		fmt.Fprint(w, `cc-bash-guard version

Print build metadata for the running binary. Use this to inspect the module,
Go toolchain, and VCS information embedded in the installed executable.

Usage:
  cc-bash-guard version [--format json]

Examples:
  cc-bash-guard version
  cc-bash-guard version --format json
`)
	case "semantic-schema":
		fmt.Fprint(w, `cc-bash-guard semantic-schema

Print supported command-specific semantic match schemas.

Usage:
  cc-bash-guard semantic-schema [command] [--format json]

Examples:
  cc-bash-guard semantic-schema --format json
  cc-bash-guard semantic-schema git --format json
`)
	case "config":
		fmt.Fprint(w, `cc-bash-guard help config

Config files live at:
  - ~/.config/cc-bash-guard/cc-bash-guard.yml
  - ./.cc-bash-guard/cc-bash-guard.yaml (project-local, optional)

Top-level sections are:
  - claude_permission_merge_mode: strict / migration_compat / cc_bash_guard_authoritative
  - permission: deny / ask / allow buckets
  - test: end-to-end expect cases

Top-level rewrite is no longer supported. cc-bash-guard never changes the
command string it evaluates or returns to Claude. Parser-backed normalization is
evaluation-only: shell -c wrappers are inspected as inner commands, absolute
paths match by basename, and AWS profile flags are parsed semantically.

Permission merge mode:
  claude_permission_merge_mode controls how Claude settings.json permissions and
  cc-bash-guard rules combine.

  strict:
    Recommended for security-first setups. cc-bash-guard is fail-closed and
    settings.json allow entries do not silently broaden cc-bash-guard policy.

  migration_compat:
    Use while migrating existing Claude permissions. Existing settings can still
    contribute compatibility decisions, but cc-bash-guard deny/ask remains the
    safer override.

  cc_bash_guard_authoritative:
    Use when cc-bash-guard should be the authoritative policy surface.

Decision order:
  deny wins over ask, ask wins over allow, and abstain means no local rule
  matched. When no permission rule matches, the fallback decision is ask.

Permission rule example:
  permission:
    allow:
      - command:
          name: aws
          semantic:
            service: sts
            operation: get-caller-identity
        env:
          requires:
            - "AWS_PROFILE"
        test:
          allow:
            - "AWS_PROFILE=read-only-profile aws sts get-caller-identity"
          pass:
            - "AWS_PROFILE=read-only-profile aws s3 ls"

E2E test example:
  test:
    - in: "AWS_PROFILE=read-only-profile aws sts get-caller-identity"
      decision: allow

For permission predicate fields, run:
  cc-bash-guard help match

For semantic command schemas, run:
  cc-bash-guard help semantic

`)
	case "match":
		fmt.Fprint(w, `cc-bash-guard help match

Permission rules do not use match or pattern. Permission rules use:
  - command: command name plus command-specific semantic matcher
  - env: execution environment matcher with requires and missing
  - patterns: raw command string regex list
Permission command does not support command_in; use multiple patterns for
multi-command raw fallbacks.

permission command.semantic:
  command.semantic is command-specific. The schema is selected by exact
  command.name. Do not write semantic.git or semantic.gh.

  GenericParser fallback never satisfies semantic match.

  semantic.flags_contains and semantic.flags_prefixes inspect tokens recognized
  as options/flags by the command-specific parser. They do not run when a
  semantic parser is unavailable.

Permission predicate combinations:
  command, command + env, command + semantic, command + semantic + env,
  patterns, patterns + env, and env only. command + patterns is invalid.

Discover semantic schemas:
  cc-bash-guard help semantic
  cc-bash-guard help semantic <command>
  cc-bash-guard semantic-schema --format json

Example:
  command:
    name: aws
    semantic:
      service: sts
  env:
    requires:
      - AWS_PROFILE

patterns is the raw regex escape hatch for permission rules.

Example:
  patterns:
    - '^\s*helm\s+upgrade\b'
`)
	default:
		writeUsage(w)
	}
}

func writeSemanticHelp(w io.Writer, args []string) error {
	if len(args) == 0 {
		fmt.Fprint(w, `Semantic match schemas

command.semantic is command-specific. The schema is selected by command.name.
Do not nest another command key under semantic.

Supported commands:
`)
		for _, schema := range semanticpkg.AllSchemas() {
			fmt.Fprintf(w, "  %-10s %s\n", schema.Command, schema.Description)
		}
		fmt.Fprint(w, `
Usage:
  cc-bash-guard help semantic <command>
  cc-bash-guard semantic-schema --format json
  cc-bash-guard semantic-schema <command> --format json

Example:
  permission:
    deny:
      - command:
          name: git
          semantic:
            verb: push
            force: true

Notes:
  semantic.flags_contains / semantic.flags_prefixes inspect options
  recognized by the command-specific parser and never match on GenericParser
  fallback.
`)
		return nil
	}
	if len(args) > 1 {
		return errors.New("usage: cc-bash-guard help semantic [command]")
	}
	schema, ok := semanticpkg.Lookup(args[0])
	if !ok {
		return fmt.Errorf("unknown semantic command %q. Supported commands: %s", args[0], strings.Join(semanticpkg.SupportedCommands(), ", "))
	}
	fmt.Fprintf(w, "Semantic schema: %s\n\n", schema.Command)
	fmt.Fprintf(w, "Description: %s\n", schema.Description)
	fmt.Fprintf(w, "Parser support: %s\n\n", schema.Parser)
	fmt.Fprint(w, "Fields:\n")
	for _, field := range schema.Fields {
		fmt.Fprintf(w, "  %-38s %-9s %s\n", field.Name, field.Type, field.Description)
	}
	if len(schema.Notes) > 0 {
		fmt.Fprint(w, "\nBoolean field definitions and notes:\n")
		for _, note := range schema.Notes {
			fmt.Fprintf(w, "  - %s\n", note)
		}
	} else {
		fmt.Fprint(w, "\nBoolean field definitions are included in the field descriptions above.\n")
	}
	fmt.Fprint(w, "\nValidation rules:\n")
	fmt.Fprint(w, "  - permission command.semantic requires exact command.name.\n")
	fmt.Fprint(w, "  - top-level rewrite is unsupported.\n")
	fmt.Fprint(w, "  - unsupported fields and unsupported value types fail verify.\n")
	fmt.Fprint(w, "  - GenericParser fallback never satisfies semantic match.\n")
	if len(schema.Examples) > 0 {
		fmt.Fprint(w, "\nExamples:\n")
		for _, example := range schema.Examples {
			fmt.Fprintf(w, "  %s:\n%s\n", example.Title, indent(example.YAML, "    "))
		}
	}
	return nil
}

func indent(s, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

func wantsHelp(args []string) bool {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			return true
		}
	}
	return false
}

func writeErr(w io.Writer, msg string) {
	fmt.Fprintln(w, msg)
}
