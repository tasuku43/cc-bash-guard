package policy

import "strings"

func ValidateAWSSemanticMatchSpec(prefix string, semantic SemanticMatchSpec) []string {
	return ValidateAWSSemanticSpec(prefix, semantic.AWS())
}

func ValidateAWSSemanticSpec(prefix string, semantic AWSSemanticSpec) []string {
	var issues []string
	if IsZeroAWSSemanticSpec(semantic) {
		issues = append(issues, prefix+" must not be empty")
	}
	if strings.TrimSpace(semantic.Service) == "" && semantic.Service != "" {
		issues = append(issues, prefix+".service must be non-empty")
	}
	if strings.TrimSpace(semantic.Operation) == "" && semantic.Operation != "" {
		issues = append(issues, prefix+".operation must be non-empty")
	}
	if strings.TrimSpace(semantic.Profile) == "" && semantic.Profile != "" {
		issues = append(issues, prefix+".profile must be non-empty")
	}
	if strings.TrimSpace(semantic.Region) == "" && semantic.Region != "" {
		issues = append(issues, prefix+".region must be non-empty")
	}
	if strings.TrimSpace(semantic.EndpointURL) == "" && semantic.EndpointURL != "" {
		issues = append(issues, prefix+".endpoint_url must be non-empty")
	}
	if strings.TrimSpace(semantic.EndpointURLPrefix) == "" && semantic.EndpointURLPrefix != "" {
		issues = append(issues, prefix+".endpoint_url_prefix must be non-empty")
	}
	issues = append(issues, validateNonEmptyStrings(prefix+".service_in", semantic.ServiceIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".operation_in", semantic.OperationIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".profile_in", semantic.ProfileIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".region_in", semantic.RegionIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".flags_contains", semantic.FlagsContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".flags_prefixes", semantic.FlagsPrefixes)...)
	return issues
}
