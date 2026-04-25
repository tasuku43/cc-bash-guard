package claude

import (
	"path/filepath"
	"testing"

	"github.com/tasuku43/cc-bash-guard/internal/domain/policy"
)

func TestSecurityRegressionMatrixPermissionSourceMerge(t *testing.T) {
	tests := []struct {
		name         string
		baseOutcome  string
		claudeAllow  []string
		claudeAsk    []string
		claudeDeny   []string
		want         string
		wantExplicit bool
		wantReason   string
		wantMerge    string
	}{
		{
			name:         "cc deny plus Claude allow denies",
			baseOutcome:  "deny",
			claudeAllow:  []string{"git status"},
			want:         "deny",
			wantExplicit: true,
			wantMerge:    "source denied",
		},
		{
			name:         "cc ask plus Claude allow asks",
			baseOutcome:  "ask",
			claudeAllow:  []string{"git status"},
			want:         "ask",
			wantExplicit: true,
			wantMerge:    "source asked",
		},
		{
			name:         "cc allow plus Claude ask asks",
			baseOutcome:  "allow",
			claudeAsk:    []string{"git status"},
			want:         "ask",
			wantExplicit: true,
			wantReason:   "claude_settings",
			wantMerge:    "source asked",
		},
		{
			name:         "cc abstain plus Claude allow allows",
			baseOutcome:  "abstain",
			claudeAllow:  []string{"git status"},
			want:         "allow",
			wantExplicit: true,
			wantReason:   "claude_settings",
			wantMerge:    "source allowed",
		},
		{
			name:         "cc abstain plus Claude ask asks",
			baseOutcome:  "abstain",
			claudeAsk:    []string{"git status"},
			want:         "ask",
			wantExplicit: true,
			wantReason:   "claude_settings",
			wantMerge:    "source asked",
		},
		{
			name:         "cc abstain plus Claude deny denies",
			baseOutcome:  "abstain",
			claudeDeny:   []string{"git status"},
			want:         "deny",
			wantExplicit: true,
			wantReason:   "claude_settings",
			wantMerge:    "source denied",
		},
		{
			name:         "both abstain fallback asks",
			baseOutcome:  "abstain",
			want:         "ask",
			wantExplicit: false,
			wantReason:   "default_fallback",
			wantMerge:    "all sources abstained; fallback ask",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home := t.TempDir()
			cwd := t.TempDir()
			writeSettings(t, filepath.Join(home, ".claude", "settings.json"), claudeSettingsJSON(tt.claudeDeny, tt.claudeAsk, tt.claudeAllow))

			decision := ApplyPermissionBridge(Tool, policy.Decision{
				Outcome:  tt.baseOutcome,
				Explicit: tt.baseOutcome != "abstain",
				Reason:   "rule_match",
				Command:  "git status",
			}, cwd, home)
			if decision.Outcome != tt.want {
				t.Fatalf("Outcome = %q, want %q; decision=%+v", decision.Outcome, tt.want, decision)
			}
			if decision.Explicit != tt.wantExplicit {
				t.Fatalf("Explicit = %v, want %v; decision=%+v", decision.Explicit, tt.wantExplicit, decision)
			}
			if tt.wantReason != "" && decision.Reason != tt.wantReason {
				t.Fatalf("Reason = %q, want %q; decision=%+v", decision.Reason, tt.wantReason, decision)
			}
			if !bridgeTraceContainsReason(decision.Trace, "permission_sources_merge", tt.wantMerge) {
				t.Fatalf("trace missing merge reason %q; trace=%+v", tt.wantMerge, decision.Trace)
			}
			if tt.baseOutcome == "deny" && decision.Outcome == "allow" {
				t.Fatalf("deny widened to allow; decision=%+v", decision)
			}
		})
	}
}

func claudeSettingsJSON(deny []string, ask []string, allow []string) string {
	return `{"permissions":{"deny":` + bashPatternsJSON(deny) + `,"ask":` + bashPatternsJSON(ask) + `,"allow":` + bashPatternsJSON(allow) + `}}`
}

func bashPatternsJSON(patterns []string) string {
	if len(patterns) == 0 {
		return `[]`
	}
	out := `[`
	for i, pattern := range patterns {
		if i > 0 {
			out += `,`
		}
		out += `"Bash(` + pattern + `)"`
	}
	return out + `]`
}

func bridgeTraceContainsReason(trace []policy.TraceStep, name string, reason string) bool {
	for _, step := range trace {
		if step.Name == name && step.Reason == reason {
			return true
		}
	}
	return false
}
