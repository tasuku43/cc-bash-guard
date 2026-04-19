package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadEffectiveUsesUserConfig(t *testing.T) {
	home := t.TempDir()
	userPath := filepath.Join(home, ".config", "cmdproxy", "cmdproxy.yml")
	if err := os.MkdirAll(filepath.Dir(userPath), 0o755); err != nil {
		t.Fatal(err)
	}
	body := `version: 1
rules:
  - id: user-rule
    pattern: "^echo"
    message: "user"
    block_examples: ["echo hi"]
    allow_examples: ["git status"]
`
	if err := os.WriteFile(userPath, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	loaded := LoadEffective(home, "")
	if len(loaded.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", loaded.Errors)
	}
	if len(loaded.Rules) != 1 || loaded.Rules[0].ID != "user-rule" {
		t.Fatalf("rules = %#v", loaded.Rules)
	}
}

func TestLoadFileForEvalIfPresentSupportsRewriteDirective(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cmdproxy.yml")
	cachePath := filepath.Join(t.TempDir(), "eval-cache-v1.json")
	body := `version: 1
rules:
  - id: unwrap-shell-dash-c
    match:
      command_in: ["bash", "sh"]
      args_contains: ["-c"]
    rewrite:
      unwrap_shell_dash_c: true
    block_examples: ["bash -c 'git status'"]
    allow_examples: ["bash script.sh"]
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	rules, err := LoadFileForEvalIfPresent(Source{Layer: LayerUser, Path: path}, cachePath)
	if err != nil {
		t.Fatalf("LoadFileForEvalIfPresent() error = %v", err)
	}
	rewritten, ok := rules[0].RewriteCommand("bash -c 'git status'")
	if !ok || rewritten != "git status" {
		t.Fatalf("RewriteCommand() = %q ok=%v", rewritten, ok)
	}
}

func TestLoadFileForEvalIfPresentSupportsMoveFlagToEnv(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cmdproxy.yml")
	cachePath := filepath.Join(t.TempDir(), "eval-cache-v1.json")
	body := `version: 1
rules:
  - id: aws-profile-to-env
    match:
      command: aws
      args_contains: ["--profile"]
    rewrite:
      move_flag_to_env:
        flag: "--profile"
        env: "AWS_PROFILE"
    block_examples: ["aws --profile read-only-profile s3 ls"]
    allow_examples: ["AWS_PROFILE=read-only-profile aws s3 ls"]
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	rules, err := LoadFileForEvalIfPresent(Source{Layer: LayerUser, Path: path}, cachePath)
	if err != nil {
		t.Fatalf("LoadFileForEvalIfPresent() error = %v", err)
	}
	rewritten, ok := rules[0].RewriteCommand("aws --profile read-only-profile s3 ls")
	if !ok || rewritten != "AWS_PROFILE=read-only-profile aws s3 ls" {
		t.Fatalf("RewriteCommand() = %q ok=%v", rewritten, ok)
	}
}

func TestLoadFileForEvalIfPresentSupportsMoveEnvToFlag(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cmdproxy.yml")
	cachePath := filepath.Join(t.TempDir(), "eval-cache-v1.json")
	body := `version: 1
rules:
  - id: aws-env-to-profile
    match:
      command: aws
      env_requires: ["AWS_PROFILE"]
    rewrite:
      move_env_to_flag:
        env: "AWS_PROFILE"
        flag: "--profile"
    block_examples: ["AWS_PROFILE=read-only-profile aws s3 ls"]
    allow_examples: ["aws --profile read-only-profile s3 ls"]
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	rules, err := LoadFileForEvalIfPresent(Source{Layer: LayerUser, Path: path}, cachePath)
	if err != nil {
		t.Fatalf("LoadFileForEvalIfPresent() error = %v", err)
	}
	rewritten, ok := rules[0].RewriteCommand("AWS_PROFILE=read-only-profile aws s3 ls")
	if !ok || rewritten != "aws --profile read-only-profile s3 ls" {
		t.Fatalf("RewriteCommand() = %q ok=%v", rewritten, ok)
	}
}

func TestLoadFileForEvalIfPresentSupportsUnwrapWrapper(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cmdproxy.yml")
	cachePath := filepath.Join(t.TempDir(), "eval-cache-v1.json")
	body := `version: 1
rules:
  - id: unwrap-safe-wrappers
    pattern: '^\s*(env|command|exec)\b'
    rewrite:
      unwrap_wrapper:
        wrappers: ["env", "command", "exec"]
    block_examples: ["env AWS_PROFILE=dev command exec aws s3 ls"]
    allow_examples: ["AWS_PROFILE=dev aws s3 ls"]
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	rules, err := LoadFileForEvalIfPresent(Source{Layer: LayerUser, Path: path}, cachePath)
	if err != nil {
		t.Fatalf("LoadFileForEvalIfPresent() error = %v", err)
	}
	rewritten, ok := rules[0].RewriteCommand("env AWS_PROFILE=dev command exec aws s3 ls")
	if !ok || rewritten != "AWS_PROFILE=dev aws s3 ls" {
		t.Fatalf("RewriteCommand() = %q ok=%v", rewritten, ok)
	}
}

func TestLoadFileIfPresentRejectsPatternAndMatchTogether(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cmdproxy.yml")
	body := `version: 1
rules:
  - id: bad-rule
    pattern: "^git"
    match:
      command: git
    message: "bad"
    block_examples: ["git status"]
    allow_examples: ["echo ok"]
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadFileIfPresent(Source{Layer: LayerUser, Path: path})
	if err == nil || !strings.Contains(err.Error(), "must not set both pattern and match") {
		t.Fatalf("unexpected error: %v", err)
	}
}
