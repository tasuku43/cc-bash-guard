package policy

import "strings"

func ValidateGitSemanticMatchSpec(prefix string, semantic SemanticMatchSpec) []string {
	return ValidateGitSemanticSpec(prefix, semantic.Git())
}

func ValidateGitSemanticSpec(prefix string, semantic GitSemanticSpec) []string {
	var issues []string
	if IsZeroGitSemanticSpec(semantic) {
		issues = append(issues, prefix+" must not be empty")
	}
	if strings.TrimSpace(semantic.Verb) == "" && semantic.Verb != "" {
		issues = append(issues, prefix+".verb must be non-empty")
	}
	if strings.TrimSpace(semantic.Remote) == "" && semantic.Remote != "" {
		issues = append(issues, prefix+".remote must be non-empty")
	}
	if strings.TrimSpace(semantic.Branch) == "" && semantic.Branch != "" {
		issues = append(issues, prefix+".branch must be non-empty")
	}
	if strings.TrimSpace(semantic.Ref) == "" && semantic.Ref != "" {
		issues = append(issues, prefix+".ref must be non-empty")
	}
	issues = append(issues, validateNonEmptyStrings(prefix+".verb_in", semantic.VerbIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".remote_in", semantic.RemoteIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".branch_in", semantic.BranchIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".ref_in", semantic.RefIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".flags_contains", semantic.FlagsContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".flags_prefixes", semantic.FlagsPrefixes)...)
	return issues
}
