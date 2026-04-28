package policy

import "strings"

func ValidateGwsSemanticMatchSpec(prefix string, semantic SemanticMatchSpec) []string {
	return ValidateGwsSemanticSpec(prefix, semantic.Gws())
}

func ValidateGwsSemanticSpec(prefix string, semantic GwsSemanticSpec) []string {
	var issues []string
	if IsZeroGwsSemanticSpec(semantic) {
		issues = append(issues, prefix+" must not be empty")
	}
	for name, value := range map[string]string{
		"service": semantic.Service,
		"method":  semantic.Method,
	} {
		if strings.TrimSpace(value) == "" && value != "" {
			issues = append(issues, prefix+"."+name+" must be non-empty")
		}
	}
	issues = append(issues, validateNonEmptyStrings(prefix+".service_in", semantic.ServiceIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".resource_path", semantic.ResourcePath)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".resource_path_contains", semantic.ResourcePathContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".method_in", semantic.MethodIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".scopes", semantic.Scopes)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".flags_contains", semantic.FlagsContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".flags_prefixes", semantic.FlagsPrefixes)...)
	return issues
}
