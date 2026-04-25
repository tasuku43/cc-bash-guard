package claude

import (
	"strings"

	"github.com/tasuku43/cc-bash-guard/internal/domain/policy"
)

const Tool = "claude"

func Supported(tool string) bool {
	switch strings.TrimSpace(tool) {
	case Tool:
		return true
	default:
		return false
	}
}

func ApplyPermissionBridge(tool string, decision policy.Decision, cwd string, home string) policy.Decision {
	switch strings.TrimSpace(tool) {
	case Tool:
		return applyPermissionBridge(decision, cwd, home)
	default:
		return decision
	}
}

func applyPermissionBridge(decision policy.Decision, cwd string, home string) policy.Decision {
	verdict := CheckCommand(decision.Command, cwd, home)
	claudeOutcome := permissionVerdictOutcome(verdict)
	decision.Trace = append(decision.Trace, ccBashGuardPolicyTrace(decision))
	decision.Trace = append(decision.Trace, policy.TraceStep{
		Action:  "permission",
		Name:    "claude_settings",
		Effect:  claudeOutcome,
		Message: claudeSettingsTraceMessage(claudeOutcome),
	})
	return mergePermissionSources(decision, claudeOutcome)
}

func ccBashGuardPolicyTrace(decision policy.Decision) policy.TraceStep {
	outcome := strings.TrimSpace(decision.Outcome)
	if outcome == "" {
		outcome = "abstain"
	}
	reason := "explicit"
	if outcome == "abstain" {
		reason = "abstain"
	}
	step := policy.TraceStep{
		Action: "permission",
		Name:   "cc_bash_guard_policy",
		Effect: outcome,
		Reason: reason,
	}
	if matched := matchedPolicyTraceName(decision.Trace); matched != "" {
		step.Message = "matched rule: " + matched
	}
	return step
}

func matchedPolicyTraceName(trace []policy.TraceStep) string {
	for i := len(trace) - 1; i >= 0; i-- {
		step := trace[i]
		if step.Action != "permission" {
			continue
		}
		if step.Name == "no_match" || step.Name == "fail_closed" || step.Name == "composition" || step.Name == "composition.command" {
			continue
		}
		if strings.TrimSpace(step.Name) != "" {
			return step.Name
		}
	}
	return ""
}

func permissionVerdictOutcome(verdict PermissionVerdict) string {
	switch verdict {
	case PermissionDeny:
		return "deny"
	case PermissionAsk:
		return "ask"
	case PermissionAllow:
		return "allow"
	default:
		return "abstain"
	}
}

func claudeSettingsTraceMessage(outcome string) string {
	switch outcome {
	case "deny":
		return "Claude settings deny matched"
	case "ask":
		return "Claude settings ask matched"
	case "allow":
		return "Claude settings allow matched"
	default:
		return "Claude settings did not define a matching permission"
	}
}

func mergePermissionSources(decision policy.Decision, claudeOutcome string) policy.Decision {
	ccOutcome := strings.TrimSpace(decision.Outcome)
	if ccOutcome == "" {
		ccOutcome = "abstain"
	}

	switch {
	case ccOutcome == "deny" || claudeOutcome == "deny":
		if claudeOutcome == "deny" {
			decision.Outcome = "deny"
			decision.Explicit = true
			decision.Reason = "claude_settings"
			if strings.TrimSpace(decision.Message) == "" {
				decision.Message = "blocked by Claude settings"
			}
		}
		decision.Trace = append(decision.Trace, finalMergeTrace("deny", "source denied"))
	case ccOutcome == "ask" || claudeOutcome == "ask":
		if ccOutcome != "ask" && claudeOutcome == "ask" {
			decision.Outcome = "ask"
			decision.Explicit = true
			decision.Reason = "claude_settings"
		}
		decision.Trace = append(decision.Trace, finalMergeTrace("ask", "source asked"))
	case ccOutcome == "allow" || claudeOutcome == "allow":
		if ccOutcome != "allow" && claudeOutcome == "allow" {
			decision.Outcome = "allow"
			decision.Explicit = true
			decision.Reason = "claude_settings"
		}
		decision.Trace = append(decision.Trace, finalMergeTrace("allow", "source allowed"))
	default:
		decision.Outcome = "ask"
		decision.Explicit = false
		decision.Reason = "default_fallback"
		decision.Trace = append(decision.Trace, finalMergeTrace("ask", "all sources abstained; fallback ask"))
	}
	return decision
}

func finalMergeTrace(outcome string, reason string) policy.TraceStep {
	return policy.TraceStep{
		Action:  "permission",
		Name:    "permission_sources_merge",
		Effect:  outcome,
		Reason:  reason,
		Message: "permission sources merged using deny > ask > allow > abstain",
	}
}
