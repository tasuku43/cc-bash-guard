package policy

import "strings"

func ValidateTerraformSemanticMatchSpec(prefix string, semantic SemanticMatchSpec) []string {
	return ValidateTerraformSemanticSpec(prefix, semantic.Terraform())
}

func ValidateTerraformSemanticSpec(prefix string, semantic TerraformSemanticSpec) []string {
	var issues []string
	if IsZeroTerraformSemanticSpec(semantic) {
		issues = append(issues, prefix+" must not be empty")
	}
	for name, value := range map[string]string{
		"subcommand":           semantic.Subcommand,
		"global_chdir":         semantic.GlobalChdir,
		"workspace_subcommand": semantic.WorkspaceSubcommand,
		"state_subcommand":     semantic.StateSubcommand,
		"out":                  semantic.Out,
		"plan_file":            semantic.PlanFile,
	} {
		if strings.TrimSpace(value) == "" && value != "" {
			issues = append(issues, prefix+"."+name+" must be non-empty")
		}
	}
	issues = append(issues, validateNonEmptyStrings(prefix+".subcommand_in", semantic.SubcommandIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".workspace_subcommand_in", semantic.WorkspaceSubcommandIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".state_subcommand_in", semantic.StateSubcommandIn)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".targets_contains", semantic.TargetsContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".replaces_contains", semantic.ReplacesContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".var_files_contains", semantic.VarFilesContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".flags_contains", semantic.FlagsContains)...)
	issues = append(issues, validateNonEmptyStrings(prefix+".flags_prefixes", semantic.FlagsPrefixes)...)
	return issues
}
