package policy

import "strings"

func ValidatePupSemanticMatchSpec(prefix string, semantic SemanticMatchSpec) []string {
	return ValidatePUPSemanticSpec(prefix, semantic.Pup())
}

func ValidatePUPSemanticSpec(prefix string, semantic PupSemanticSpec) []string {
	var issues []string
	if IsZeroPupSemanticSpec(semantic) {
		issues = append(issues, prefix+" must not be empty")
	}
	for name, value := range map[string]string{
		"area": semantic.Area, "sub_area": semantic.SubArea, "verb": semantic.Verb,
		"org": semantic.Org, "output": semantic.Output,
	} {
		if strings.TrimSpace(value) == "" && value != "" {
			issues = append(issues, prefix+"."+name+" must be non-empty")
		}
	}
	issues = append(issues, validateNonEmptyStrings(prefix+".area_in", semantic.AreaIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".sub_area_in", semantic.SubAreaIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".verb_in", semantic.VerbIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".org_in", semantic.OrgIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".output_in", semantic.OutputIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".flags_contains", semantic.FlagsContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".flags_prefixes", semantic.FlagsPrefixes)...)
	return issues
}
