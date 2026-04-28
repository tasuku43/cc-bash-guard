package policy

import "strings"

func ValidateKubectlSemanticMatchSpec(prefix string, semantic SemanticMatchSpec) []string {
	return ValidateKubectlSemanticSpec(prefix, semantic.Kubectl())
}

func ValidateKubectlSemanticSpec(prefix string, semantic KubectlSemanticSpec) []string {
	var issues []string
	if IsZeroKubectlSemanticSpec(semantic) {
		issues = append(issues, prefix+" must not be empty")
	}
	if strings.TrimSpace(semantic.Verb) == "" && semantic.Verb != "" {
		issues = append(issues, prefix+".verb must be non-empty")
	}
	if strings.TrimSpace(semantic.Subverb) == "" && semantic.Subverb != "" {
		issues = append(issues, prefix+".subverb must be non-empty")
	}
	if strings.TrimSpace(semantic.ResourceType) == "" && semantic.ResourceType != "" {
		issues = append(issues, prefix+".resource_type must be non-empty")
	}
	if strings.TrimSpace(semantic.ResourceName) == "" && semantic.ResourceName != "" {
		issues = append(issues, prefix+".resource_name must be non-empty")
	}
	if strings.TrimSpace(semantic.Namespace) == "" && semantic.Namespace != "" {
		issues = append(issues, prefix+".namespace must be non-empty")
	}
	if strings.TrimSpace(semantic.Context) == "" && semantic.Context != "" {
		issues = append(issues, prefix+".context must be non-empty")
	}
	if strings.TrimSpace(semantic.Kubeconfig) == "" && semantic.Kubeconfig != "" {
		issues = append(issues, prefix+".kubeconfig must be non-empty")
	}
	if strings.TrimSpace(semantic.Filename) == "" && semantic.Filename != "" {
		issues = append(issues, prefix+".filename must be non-empty")
	}
	if strings.TrimSpace(semantic.FilenamePrefix) == "" && semantic.FilenamePrefix != "" {
		issues = append(issues, prefix+".filename_prefix must be non-empty")
	}
	if strings.TrimSpace(semantic.Selector) == "" && semantic.Selector != "" {
		issues = append(issues, prefix+".selector must be non-empty")
	}
	if strings.TrimSpace(semantic.Container) == "" && semantic.Container != "" {
		issues = append(issues, prefix+".container must be non-empty")
	}
	issues = append(issues, validateNonEmptyStrings(prefix+".verb_in", semantic.VerbIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".subverb_in", semantic.SubverbIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".resource_type_in", semantic.ResourceTypeIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".resource_name_in", semantic.ResourceNameIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".namespace_in", semantic.NamespaceIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".context_in", semantic.ContextIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".filename_in", semantic.FilenameIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".selector_contains", semantic.SelectorContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".flags_contains", semantic.FlagsContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".flags_prefixes", semantic.FlagsPrefixes)...)
	return issues
}
