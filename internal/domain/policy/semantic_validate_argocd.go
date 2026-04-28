package policy

import "strings"

func ValidateArgoCDSemanticMatchSpec(prefix string, semantic SemanticMatchSpec) []string {
	return ValidateArgoCDSemanticSpec(prefix, semantic.ArgoCD())
}

func ValidateArgoCDSemanticSpec(prefix string, semantic ArgoCDSemanticSpec) []string {
	var issues []string
	if IsZeroArgoCDSemanticSpec(semantic) {
		issues = append(issues, prefix+" must not be empty")
	}
	for name, value := range map[string]string{
		"verb":     semantic.Verb,
		"app_name": semantic.AppName,
		"project":  semantic.Project,
		"revision": semantic.Revision,
	} {
		if strings.TrimSpace(value) == "" && value != "" {
			issues = append(issues, prefix+"."+name+" must be non-empty")
		}
	}
	issues = append(issues, validateNonEmptyStrings(prefix+".verb_in", semantic.VerbIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".app_name_in", semantic.AppNameIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".project_in", semantic.ProjectIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".flags_contains", semantic.FlagsContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".flags_prefixes", semantic.FlagsPrefixes)...)
	return issues
}
