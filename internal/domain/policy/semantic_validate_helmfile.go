package policy

import "strings"

func ValidateHelmfileSemanticMatchSpec(prefix string, semantic SemanticMatchSpec) []string {
	return ValidateHelmfileSemanticSpec(prefix, semantic.Helmfile())
}

func ValidateHelmfileSemanticSpec(prefix string, semantic HelmfileSemanticSpec) []string {
	var issues []string
	if IsZeroHelmfileSemanticSpec(semantic) {
		issues = append(issues, prefix+" must not be empty")
	}
	for name, value := range map[string]string{
		"verb":              semantic.Verb,
		"environment":       semantic.Environment,
		"file":              semantic.File,
		"file_prefix":       semantic.FilePrefix,
		"namespace":         semantic.Namespace,
		"kube_context":      semantic.KubeContext,
		"selector":          semantic.Selector,
		"cascade":           semantic.Cascade,
		"state_values_file": semantic.StateValuesFile,
	} {
		if strings.TrimSpace(value) == "" && value != "" {
			issues = append(issues, prefix+"."+name+" must be non-empty")
		}
	}
	issues = append(issues, validateNonEmptyStrings(prefix+".verb_in", semantic.VerbIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".environment_in", semantic.EnvironmentIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".file_in", semantic.FileIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".namespace_in", semantic.NamespaceIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".kube_context_in", semantic.KubeContextIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".selector_in", semantic.SelectorIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".selector_contains", semantic.SelectorContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".cascade_in", semantic.CascadeIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".state_values_file_in", semantic.StateValuesFileIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".state_values_set_keys_contains", semantic.StateValuesSetKeysContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".state_values_set_string_keys_contains", semantic.StateValuesSetStringKeysContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".flags_contains", semantic.FlagsContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".flags_prefixes", semantic.FlagsPrefixes)...)
	return issues
}
