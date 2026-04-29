package policy

import commandpkg "github.com/tasuku43/cc-bash-guard/internal/domain/command"

func init() {
	registerSemanticHandler(semanticHandler{
		command:  "terraform",
		match:    func(s SemanticMatchSpec, cmd commandpkg.Command) bool { return s.Terraform().matches(cmd) },
		validate: ValidateTerraformSemanticMatchSpec,
	})
}

func (s TerraformSemanticSpec) matches(cmd commandpkg.Command) bool {
	if cmd.SemanticParser != "terraform" || cmd.Terraform == nil {
		return false
	}
	t := cmd.Terraform
	if s.Subcommand != "" && t.Subcommand != s.Subcommand {
		return false
	}
	if len(s.SubcommandIn) > 0 && !containsString(s.SubcommandIn, t.Subcommand) {
		return false
	}
	if s.GlobalChdir != "" && t.GlobalChdir != s.GlobalChdir {
		return false
	}
	if s.WorkspaceSubcommand != "" && t.WorkspaceSubcommand != s.WorkspaceSubcommand {
		return false
	}
	if len(s.WorkspaceSubcommandIn) > 0 && !containsString(s.WorkspaceSubcommandIn, t.WorkspaceSubcommand) {
		return false
	}
	if s.StateSubcommand != "" && t.StateSubcommand != s.StateSubcommand {
		return false
	}
	if len(s.StateSubcommandIn) > 0 && !containsString(s.StateSubcommandIn, t.StateSubcommand) {
		return false
	}
	if s.Target != nil && t.Target != *s.Target {
		return false
	}
	for _, target := range s.TargetsContains {
		if !containsString(t.Targets, target) {
			return false
		}
	}
	if s.Replace != nil && t.Replace != *s.Replace {
		return false
	}
	for _, replace := range s.ReplacesContains {
		if !containsString(t.Replaces, replace) {
			return false
		}
	}
	if s.Destroy != nil && t.Destroy != *s.Destroy {
		return false
	}
	if s.AutoApprove != nil && t.AutoApprove != *s.AutoApprove {
		return false
	}
	if s.Input != nil {
		if t.Input == nil || *t.Input != *s.Input {
			return false
		}
	}
	if s.Lock != nil {
		if t.Lock == nil || *t.Lock != *s.Lock {
			return false
		}
	}
	if s.Refresh != nil {
		if t.Refresh == nil || *t.Refresh != *s.Refresh {
			return false
		}
	}
	if s.RefreshOnly != nil && t.RefreshOnly != *s.RefreshOnly {
		return false
	}
	if s.Out != "" && t.Out != s.Out {
		return false
	}
	if s.PlanFile != "" && t.PlanFile != s.PlanFile {
		return false
	}
	for _, file := range s.VarFilesContains {
		if !containsString(t.VarFiles, file) {
			return false
		}
	}
	if s.Vars != nil && t.Vars != *s.Vars {
		return false
	}
	if s.Backend != nil {
		if t.Backend == nil || *t.Backend != *s.Backend {
			return false
		}
	}
	if s.Upgrade != nil && t.Upgrade != *s.Upgrade {
		return false
	}
	if s.Reconfigure != nil && t.Reconfigure != *s.Reconfigure {
		return false
	}
	if s.MigrateState != nil && t.MigrateState != *s.MigrateState {
		return false
	}
	if s.Recursive != nil && t.Recursive != *s.Recursive {
		return false
	}
	if s.Check != nil && t.Check != *s.Check {
		return false
	}
	if s.JSON != nil && t.JSON != *s.JSON {
		return false
	}
	if s.Force != nil && t.Force != *s.Force {
		return false
	}
	for _, flag := range s.FlagsContains {
		if !containsString(t.Flags, flag) {
			return false
		}
	}
	for _, prefix := range s.FlagsPrefixes {
		if !containsPrefix(t.Flags, prefix) {
			return false
		}
	}
	return true
}
