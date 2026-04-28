package policy

import "strings"

func ValidateGhSemanticMatchSpec(prefix string, semantic SemanticMatchSpec) []string {
	return ValidateGHSemanticSpec(prefix, semantic.GH())
}

func ValidateGHSemanticSpec(prefix string, semantic GHSemanticSpec) []string {
	var issues []string
	if IsZeroGHSemanticSpec(semantic) {
		issues = append(issues, prefix+" must not be empty")
	}
	for name, value := range map[string]string{
		"area":            semantic.Area,
		"verb":            semantic.Verb,
		"repo":            semantic.Repo,
		"org":             semantic.Org,
		"env":             semantic.EnvName,
		"hostname":        semantic.Hostname,
		"method":          semantic.Method,
		"endpoint":        semantic.Endpoint,
		"endpoint_prefix": semantic.EndpointPrefix,
		"pr_number":       semantic.PRNumber,
		"issue_number":    semantic.IssueNumber,
		"secret_name":     semantic.SecretName,
		"tag":             semantic.Tag,
		"workflow_name":   semantic.WorkflowName,
		"workflow_id":     semantic.WorkflowID,
		"search_type":     semantic.SearchType,
		"query_contains":  semantic.QueryContains,
		"base":            semantic.Base,
		"head":            semantic.Head,
		"ref":             semantic.Ref,
		"state":           semantic.State,
		"title_contains":  semantic.TitleContains,
		"body_contains":   semantic.BodyContains,
		"merge_strategy":  semantic.MergeStrategy,
		"run_id":          semantic.RunID,
		"job":             semantic.Job,
	} {
		if strings.TrimSpace(value) == "" && value != "" {
			issues = append(issues, prefix+"."+name+" must be non-empty")
		}
	}
	issues = append(issues, validateNonEmptyStrings(prefix+".area_in", semantic.AreaIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".verb_in", semantic.VerbIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".repo_in", semantic.RepoIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".org_in", semantic.OrgIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".env_in", semantic.EnvNameIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".hostname_in", semantic.HostnameIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".method_in", semantic.MethodIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".endpoint_contains", semantic.EndpointContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".ref_in", semantic.RefIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".state_in", semantic.StateIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".label_in", semantic.LabelIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".assignee_in", semantic.AssigneeIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".secret_name_in", semantic.SecretNameIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".search_type_in", semantic.SearchTypeIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".field_keys_contains", semantic.FieldKeysContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".raw_field_keys_contains", semantic.RawFieldKeysContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".header_keys_contains", semantic.HeaderKeysContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".merge_strategy_in", semantic.MergeStrategyIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".flags_contains", semantic.FlagsContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".flags_prefixes", semantic.FlagsPrefixes)...)
	return issues
}
