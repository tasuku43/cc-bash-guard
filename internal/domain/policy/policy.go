package policy

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/tasuku43/cmdproxy/internal/domain/directive"
	"github.com/tasuku43/cmdproxy/internal/domain/invocation"
)

var ruleIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

type RuleSpec struct {
	ID            string      `yaml:"id"`
	Pattern       string      `yaml:"pattern"`
	Matcher       MatchSpec   `yaml:"match"`
	Message       string      `yaml:"message"`
	Reject        RejectSpec  `yaml:"reject"`
	Rewrite       RewriteSpec `yaml:"rewrite"`
	BlockExamples []string    `yaml:"block_examples"`
	AllowExamples []string    `yaml:"allow_examples"`
}

type RejectSpec struct {
	Message string `yaml:"message" json:"message,omitempty"`
}

type RewriteSpec struct {
	UnwrapShellDashC bool              `yaml:"unwrap_shell_dash_c" json:"unwrap_shell_dash_c,omitempty"`
	UnwrapWrapper    UnwrapWrapperSpec `yaml:"unwrap_wrapper" json:"unwrap_wrapper,omitempty"`
	MoveFlagToEnv    MoveFlagToEnvSpec `yaml:"move_flag_to_env" json:"move_flag_to_env,omitempty"`
	MoveEnvToFlag    MoveEnvToFlagSpec `yaml:"move_env_to_flag" json:"move_env_to_flag,omitempty"`
}

type MoveFlagToEnvSpec struct {
	Flag string `yaml:"flag" json:"flag,omitempty"`
	Env  string `yaml:"env" json:"env,omitempty"`
}

type MoveEnvToFlagSpec struct {
	Env  string `yaml:"env" json:"env,omitempty"`
	Flag string `yaml:"flag" json:"flag,omitempty"`
}

type UnwrapWrapperSpec struct {
	Wrappers []string `yaml:"wrappers" json:"wrappers,omitempty"`
}

type MatchSpec struct {
	Command      string   `yaml:"command" json:"command,omitempty"`
	CommandIn    []string `yaml:"command_in" json:"command_in,omitempty"`
	Subcommand   string   `yaml:"subcommand" json:"subcommand,omitempty"`
	ArgsContains []string `yaml:"args_contains" json:"args_contains,omitempty"`
	ArgsPrefixes []string `yaml:"args_prefixes" json:"args_prefixes,omitempty"`
	EnvRequires  []string `yaml:"env_requires" json:"env_requires,omitempty"`
	EnvMissing   []string `yaml:"env_missing" json:"env_missing,omitempty"`
}

type Source struct {
	Layer string `json:"layer"`
	Path  string `json:"path"`
}

type Rule struct {
	RuleSpec
	Source Source `json:"source"`
	re     *regexp.Regexp
}

type ValidationError struct {
	Issues []string
}

func (e *ValidationError) Error() string {
	return strings.Join(e.Issues, "; ")
}

type Decision struct {
	Outcome         string
	Rule            *Rule
	Command         string
	OriginalCommand string
}

func NewRule(spec RuleSpec, src Source) Rule {
	r := Rule{RuleSpec: spec, Source: src}
	if strings.TrimSpace(spec.Pattern) != "" {
		r.re, _ = regexp.Compile(spec.Pattern)
	}
	return r
}

func Evaluate(rules []Rule, command string) (Decision, error) {
	for i := range rules {
		matched, err := rules[i].Match(command)
		if err != nil {
			return Decision{}, err
		}
		if matched {
			if rewritten, ok := rules[i].RewriteCommand(command); ok {
				return Decision{Outcome: "rewrite", Rule: &rules[i], Command: rewritten, OriginalCommand: command}, nil
			}
			return Decision{Outcome: "reject", Rule: &rules[i], Command: command, OriginalCommand: command}, nil
		}
	}
	return Decision{Outcome: "pass", Command: command, OriginalCommand: command}, nil
}

func (r Rule) Match(command string) (bool, error) {
	if !IsZeroMatchSpec(r.Matcher) {
		return r.Matcher.matches(invocation.Parse(command)), nil
	}
	if r.re != nil {
		return r.re.MatchString(command), nil
	}
	compiled, err := regexp.Compile(r.Pattern)
	if err != nil {
		return false, err
	}
	return compiled.MatchString(command), nil
}

func (m MatchSpec) matches(parsed invocation.Parsed) bool {
	if parsed.Command == "" {
		return false
	}
	if m.Command != "" && parsed.Command != m.Command {
		return false
	}
	if len(m.CommandIn) > 0 && !containsString(m.CommandIn, parsed.Command) {
		return false
	}
	if m.Subcommand != "" && parsed.Subcommand != m.Subcommand {
		return false
	}
	for _, arg := range m.ArgsContains {
		if !containsString(parsed.Args, arg) {
			return false
		}
	}
	for _, prefix := range m.ArgsPrefixes {
		if !containsPrefix(parsed.Args, prefix) {
			return false
		}
	}
	for _, env := range m.EnvRequires {
		if _, ok := parsed.EnvAssignments[env]; !ok {
			return false
		}
	}
	for _, env := range m.EnvMissing {
		if _, ok := parsed.EnvAssignments[env]; ok {
			return false
		}
	}
	return true
}

func ValidateRuleMatcher(prefix string, pattern string, match MatchSpec) []string {
	var issues []string
	hasPattern := strings.TrimSpace(pattern) != ""
	hasMatch := !IsZeroMatchSpec(match)
	switch {
	case hasPattern && hasMatch:
		issues = append(issues, prefix+" must not set both pattern and match")
	case !hasPattern && !hasMatch:
		issues = append(issues, prefix+" must set exactly one of pattern or match")
	case hasPattern:
		if _, err := regexp.Compile(pattern); err != nil {
			issues = append(issues, prefix+".pattern failed to compile: "+err.Error())
		}
	case hasMatch:
		issues = append(issues, ValidateMatchSpec(prefix+".match", match)...)
	}
	return issues
}

func ValidateDirective(prefix string, message string, reject RejectSpec, rewrite RewriteSpec) []string {
	hasMessage := strings.TrimSpace(message) != ""
	hasReject := strings.TrimSpace(reject.Message) != ""
	hasRewrite := !IsZeroRewriteSpec(rewrite)
	switch {
	case countDirectiveKinds(hasMessage, hasReject, hasRewrite) > 1:
		return []string{prefix + " must set exactly one directive kind"}
	case countDirectiveKinds(hasMessage, hasReject, hasRewrite) == 0:
		return []string{prefix + " must set one directive"}
	case hasRewrite:
		return ValidateRewrite(prefix+".rewrite", rewrite)
	default:
		return nil
	}
}

func ValidateRewrite(prefix string, rewrite RewriteSpec) []string {
	var issues []string
	primitiveCount := 0
	if rewrite.UnwrapShellDashC {
		primitiveCount++
	}
	if !IsZeroUnwrapWrapperSpec(rewrite.UnwrapWrapper) {
		primitiveCount++
		issues = append(issues, validateNonEmptyStrings(prefix+".unwrap_wrapper.wrappers", rewrite.UnwrapWrapper.Wrappers)...)
	}
	if !IsZeroMoveFlagToEnvSpec(rewrite.MoveFlagToEnv) {
		primitiveCount++
		if strings.TrimSpace(rewrite.MoveFlagToEnv.Flag) == "" {
			issues = append(issues, prefix+".move_flag_to_env.flag must be non-empty")
		}
		if strings.TrimSpace(rewrite.MoveFlagToEnv.Env) == "" {
			issues = append(issues, prefix+".move_flag_to_env.env must be non-empty")
		}
	}
	if !IsZeroMoveEnvToFlagSpec(rewrite.MoveEnvToFlag) {
		primitiveCount++
		if strings.TrimSpace(rewrite.MoveEnvToFlag.Env) == "" {
			issues = append(issues, prefix+".move_env_to_flag.env must be non-empty")
		}
		if strings.TrimSpace(rewrite.MoveEnvToFlag.Flag) == "" {
			issues = append(issues, prefix+".move_env_to_flag.flag must be non-empty")
		}
	}
	switch {
	case primitiveCount == 0:
		issues = append(issues, prefix+" must not be empty")
	case primitiveCount > 1:
		issues = append(issues, prefix+" must set exactly one rewrite primitive")
	}
	return issues
}

func ValidateMatchSpec(prefix string, match MatchSpec) []string {
	var issues []string
	if IsZeroMatchSpec(match) {
		return []string{prefix + " must not be empty"}
	}
	if strings.TrimSpace(match.Command) == "" && match.Command != "" {
		issues = append(issues, prefix+".command must be non-empty")
	}
	if strings.TrimSpace(match.Subcommand) == "" && match.Subcommand != "" {
		issues = append(issues, prefix+".subcommand must be non-empty")
	}
	issues = append(issues, validateNonEmptyStrings(prefix+".command_in", match.CommandIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".args_contains", match.ArgsContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".args_prefixes", match.ArgsPrefixes)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".env_requires", match.EnvRequires)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".env_missing", match.EnvMissing)...)
	return issues
}

func ValidateRules(rules []RuleSpec) []string {
	var issues []string
	seen := map[string]struct{}{}
	for i, r := range rules {
		prefix := fmt.Sprintf("rules[%d]", i)
		if !ruleIDPattern.MatchString(r.ID) {
			issues = append(issues, prefix+".id must match [a-z0-9][a-z0-9-]*")
		}
		if _, ok := seen[r.ID]; ok && r.ID != "" {
			issues = append(issues, prefix+".id duplicates another rule in the same file")
		}
		seen[r.ID] = struct{}{}
		issues = append(issues, ValidateRuleMatcher(prefix, r.Pattern, r.Matcher)...)
		issues = append(issues, ValidateDirective(prefix, r.Message, r.Reject, r.Rewrite)...)
		if len(r.BlockExamples) == 0 {
			issues = append(issues, prefix+".block_examples must be non-empty")
		}
		if len(r.AllowExamples) == 0 {
			issues = append(issues, prefix+".allow_examples must be non-empty")
		}
	}
	return issues
}

func ValidateDuplicateIDs(rules []Rule) []error {
	seen := map[string]Source{}
	var errs []error
	for _, r := range rules {
		if prev, ok := seen[r.ID]; ok {
			errs = append(errs, fmt.Errorf("duplicate rule id %q across %s and %s", r.ID, prev.Path, r.Source.Path))
			continue
		}
		seen[r.ID] = r.Source
	}
	return errs
}

func ErrorStrings(errs []error) []string {
	parts := make([]string, 0, len(errs))
	for _, err := range errs {
		if err == nil {
			continue
		}
		var ve *ValidationError
		if errors.As(err, &ve) {
			parts = append(parts, ve.Issues...)
			continue
		}
		parts = append(parts, err.Error())
	}
	slices.Sort(parts)
	return parts
}

func (r Rule) RejectMessage() string {
	if strings.TrimSpace(r.Reject.Message) != "" {
		return r.Reject.Message
	}
	return r.Message
}

func (r Rule) RewriteCommand(command string) (string, bool) {
	if r.Rewrite.UnwrapShellDashC {
		return directive.UnwrapShellDashC(command)
	}
	if !IsZeroUnwrapWrapperSpec(r.Rewrite.UnwrapWrapper) {
		return directive.UnwrapWrapper(command, r.Rewrite.UnwrapWrapper.Wrappers)
	}
	if !IsZeroMoveFlagToEnvSpec(r.Rewrite.MoveFlagToEnv) {
		return directive.MoveFlagToEnv(command, r.Rewrite.MoveFlagToEnv.Flag, r.Rewrite.MoveFlagToEnv.Env)
	}
	if !IsZeroMoveEnvToFlagSpec(r.Rewrite.MoveEnvToFlag) {
		return directive.MoveEnvToFlag(command, r.Rewrite.MoveEnvToFlag.Env, r.Rewrite.MoveEnvToFlag.Flag)
	}
	return "", false
}

func IsZeroMatchSpec(match MatchSpec) bool {
	return match.Command == "" &&
		len(match.CommandIn) == 0 &&
		match.Subcommand == "" &&
		len(match.ArgsContains) == 0 &&
		len(match.ArgsPrefixes) == 0 &&
		len(match.EnvRequires) == 0 &&
		len(match.EnvMissing) == 0
}

func IsZeroRewriteSpec(rewrite RewriteSpec) bool {
	return !rewrite.UnwrapShellDashC &&
		IsZeroUnwrapWrapperSpec(rewrite.UnwrapWrapper) &&
		IsZeroMoveFlagToEnvSpec(rewrite.MoveFlagToEnv) &&
		IsZeroMoveEnvToFlagSpec(rewrite.MoveEnvToFlag)
}

func IsZeroMoveFlagToEnvSpec(spec MoveFlagToEnvSpec) bool {
	return strings.TrimSpace(spec.Flag) == "" && strings.TrimSpace(spec.Env) == ""
}

func IsZeroMoveEnvToFlagSpec(spec MoveEnvToFlagSpec) bool {
	return strings.TrimSpace(spec.Env) == "" && strings.TrimSpace(spec.Flag) == ""
}

func IsZeroUnwrapWrapperSpec(spec UnwrapWrapperSpec) bool {
	return len(spec.Wrappers) == 0
}

func countDirectiveKinds(hasMessage bool, hasReject bool, hasRewrite bool) int {
	n := 0
	if hasMessage {
		n++
	}
	if hasReject {
		n++
	}
	if hasRewrite {
		n++
	}
	return n
}

func validateNonEmptyStrings(prefix string, values []string) []string {
	var issues []string
	for i, value := range values {
		if strings.TrimSpace(value) == "" {
			issues = append(issues, fmt.Sprintf("%s[%d] must be non-empty", prefix, i))
		}
	}
	return issues
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func containsPrefix(values []string, prefix string) bool {
	for _, value := range values {
		if strings.HasPrefix(value, prefix) {
			return true
		}
	}
	return false
}
