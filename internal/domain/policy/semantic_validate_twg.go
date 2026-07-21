package policy

import "strings"

func ValidateTWGSemanticMatchSpec(prefix string, semantic SemanticMatchSpec) []string {
	return ValidateTWGSemanticSpec(prefix, semantic.TWG())
}

func ValidateTWGSemanticSpec(prefix string, semantic TWGSemanticSpec) []string {
	var issues []string
	if IsZeroTWGSemanticSpec(semantic) {
		issues = append(issues, prefix+" must not be empty")
	}
	if strings.TrimSpace(semantic.Namespace) == "" && semantic.Namespace != "" {
		issues = append(issues, prefix+".namespace must be non-empty")
	}
	if strings.TrimSpace(semantic.Verb) == "" && semantic.Verb != "" {
		issues = append(issues, prefix+".verb must be non-empty")
	}
	issues = append(issues, validateNonEmptyStrings(prefix+".namespace_in", semantic.NamespaceIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".verb_in", semantic.VerbIn)...)
	return issues
}
