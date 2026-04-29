package policy

func IsZeroSemanticMatchSpec(semantic SemanticMatchSpec) bool {
	return !hasGitSemanticFields(semantic) && !hasAWSSemanticFields(semantic) && !hasKubectlSemanticFields(semantic) && !hasGhSemanticFields(semantic) && !hasGwsSemanticFields(semantic) && !hasHelmfileSemanticFields(semantic) && !hasArgoCDSemanticFields(semantic) &&
		len(semantic.FlagsContains) == 0 &&
		len(semantic.FlagsPrefixes) == 0
}

func IsZeroGitSemanticSpec(semantic GitSemanticSpec) bool {
	return semantic.Verb == "" && len(semantic.VerbIn) == 0 &&
		semantic.Remote == "" && len(semantic.RemoteIn) == 0 &&
		semantic.Branch == "" && len(semantic.BranchIn) == 0 &&
		semantic.Ref == "" && len(semantic.RefIn) == 0 &&
		semantic.Force == nil && semantic.ForceWithLease == nil && semantic.ForceIfIncludes == nil &&
		semantic.Hard == nil && semantic.Recursive == nil && semantic.IncludeIgnored == nil &&
		semantic.Cached == nil && semantic.Staged == nil &&
		len(semantic.FlagsContains) == 0 && len(semantic.FlagsPrefixes) == 0
}

func IsZeroAWSSemanticSpec(semantic AWSSemanticSpec) bool {
	return semantic.Service == "" && len(semantic.ServiceIn) == 0 &&
		semantic.Operation == "" && len(semantic.OperationIn) == 0 &&
		semantic.Profile == "" && len(semantic.ProfileIn) == 0 &&
		semantic.Region == "" && len(semantic.RegionIn) == 0 &&
		semantic.EndpointURL == "" && semantic.EndpointURLPrefix == "" &&
		semantic.DryRun == nil && semantic.NoCLIPager == nil &&
		len(semantic.FlagsContains) == 0 && len(semantic.FlagsPrefixes) == 0
}

func IsZeroKubectlSemanticSpec(semantic KubectlSemanticSpec) bool {
	return semantic.Verb == "" && len(semantic.VerbIn) == 0 &&
		semantic.Subverb == "" && len(semantic.SubverbIn) == 0 &&
		semantic.ResourceType == "" && len(semantic.ResourceTypeIn) == 0 &&
		semantic.ResourceName == "" && len(semantic.ResourceNameIn) == 0 &&
		semantic.Namespace == "" && len(semantic.NamespaceIn) == 0 && semantic.NamespaceMissing == nil &&
		semantic.Context == "" && len(semantic.ContextIn) == 0 &&
		semantic.Kubeconfig == "" && semantic.AllNamespaces == nil &&
		semantic.Filename == "" && len(semantic.FilenameIn) == 0 && semantic.FilenamePrefix == "" &&
		semantic.Selector == "" && len(semantic.SelectorIn) == 0 && len(semantic.SelectorContains) == 0 &&
		semantic.SelectorMissing == nil && semantic.Container == "" &&
		semantic.DryRun == nil && semantic.Force == nil && semantic.Recursive == nil &&
		len(semantic.FlagsContains) == 0 && len(semantic.FlagsPrefixes) == 0
}

func IsZeroGHSemanticSpec(semantic GHSemanticSpec) bool {
	return semantic.Area == "" && len(semantic.AreaIn) == 0 &&
		semantic.Verb == "" && len(semantic.VerbIn) == 0 &&
		semantic.Repo == "" && len(semantic.RepoIn) == 0 &&
		semantic.Org == "" && len(semantic.OrgIn) == 0 &&
		semantic.EnvName == "" && len(semantic.EnvNameIn) == 0 &&
		semantic.Hostname == "" && len(semantic.HostnameIn) == 0 &&
		semantic.Web == nil && semantic.Method == "" && len(semantic.MethodIn) == 0 &&
		semantic.Endpoint == "" && semantic.EndpointPrefix == "" && len(semantic.EndpointContains) == 0 &&
		semantic.Paginate == nil && semantic.Input == nil && semantic.Silent == nil && semantic.IncludeHeaders == nil &&
		len(semantic.FieldKeysContains) == 0 && len(semantic.RawFieldKeysContains) == 0 && len(semantic.HeaderKeysContains) == 0 &&
		semantic.PRNumber == "" && semantic.IssueNumber == "" && semantic.SecretName == "" && len(semantic.SecretNameIn) == 0 &&
		semantic.Tag == "" && semantic.WorkflowName == "" && semantic.WorkflowID == "" &&
		semantic.SearchType == "" && len(semantic.SearchTypeIn) == 0 && semantic.QueryContains == "" &&
		semantic.Base == "" && semantic.Head == "" && semantic.Ref == "" && len(semantic.RefIn) == 0 &&
		semantic.State == "" && len(semantic.StateIn) == 0 && len(semantic.LabelIn) == 0 && len(semantic.AssigneeIn) == 0 &&
		semantic.TitleContains == "" && semantic.BodyContains == "" &&
		semantic.Draft == nil && semantic.Prerelease == nil && semantic.Latest == nil && semantic.Fill == nil &&
		semantic.Force == nil && semantic.Admin == nil && semantic.Auto == nil && semantic.DeleteBranch == nil &&
		semantic.MergeStrategy == "" && len(semantic.MergeStrategyIn) == 0 &&
		semantic.RunID == "" && semantic.Failed == nil && semantic.Job == "" && semantic.Debug == nil && semantic.ExitStatus == nil &&
		len(semantic.FlagsContains) == 0 && len(semantic.FlagsPrefixes) == 0
}

func IsZeroGwsSemanticSpec(semantic GwsSemanticSpec) bool {
	return semantic.Service == "" && len(semantic.ServiceIn) == 0 &&
		len(semantic.ResourcePath) == 0 && len(semantic.ResourcePathContains) == 0 &&
		semantic.Method == "" && len(semantic.MethodIn) == 0 &&
		semantic.Helper == nil && semantic.Mutating == nil && semantic.Destructive == nil && semantic.ReadOnly == nil &&
		semantic.DryRun == nil && semantic.PageAll == nil && semantic.Upload == nil && semantic.Sanitize == nil &&
		semantic.Params == nil && semantic.JSONBody == nil && semantic.Unmasked == nil && len(semantic.Scopes) == 0 &&
		len(semantic.FlagsContains) == 0 && len(semantic.FlagsPrefixes) == 0
}

func IsZeroHelmfileSemanticSpec(semantic HelmfileSemanticSpec) bool {
	return semantic.Verb == "" && len(semantic.VerbIn) == 0 &&
		semantic.Environment == "" && len(semantic.EnvironmentIn) == 0 && semantic.EnvironmentMissing == nil &&
		semantic.File == "" && len(semantic.FileIn) == 0 && semantic.FilePrefix == "" && semantic.FileMissing == nil &&
		semantic.Namespace == "" && len(semantic.NamespaceIn) == 0 && semantic.NamespaceMissing == nil &&
		semantic.KubeContext == "" && len(semantic.KubeContextIn) == 0 && semantic.KubeContextMissing == nil &&
		semantic.Selector == "" && len(semantic.SelectorIn) == 0 && len(semantic.SelectorContains) == 0 && semantic.SelectorMissing == nil &&
		semantic.Interactive == nil && semantic.DryRun == nil && semantic.Wait == nil && semantic.WaitForJobs == nil &&
		semantic.SkipDiff == nil && semantic.SkipNeeds == nil && semantic.IncludeNeeds == nil && semantic.IncludeTransitiveNeeds == nil &&
		semantic.Purge == nil && semantic.Cascade == "" && len(semantic.CascadeIn) == 0 && semantic.DeleteWait == nil &&
		semantic.StateValuesFile == "" && len(semantic.StateValuesFileIn) == 0 &&
		len(semantic.StateValuesSetKeysContains) == 0 && len(semantic.StateValuesSetStringKeysContains) == 0 &&
		len(semantic.FlagsContains) == 0 && len(semantic.FlagsPrefixes) == 0
}

func IsZeroArgoCDSemanticSpec(semantic ArgoCDSemanticSpec) bool {
	return semantic.Verb == "" && len(semantic.VerbIn) == 0 &&
		semantic.AppName == "" && len(semantic.AppNameIn) == 0 &&
		semantic.Project == "" && len(semantic.ProjectIn) == 0 &&
		semantic.Revision == "" &&
		len(semantic.FlagsContains) == 0 && len(semantic.FlagsPrefixes) == 0
}

func hasGitSemanticFields(semantic SemanticMatchSpec) bool {
	return semantic.Verb != "" ||
		len(semantic.VerbIn) > 0 ||
		hasGitOnlySemanticFields(semantic) ||
		semantic.Force != nil ||
		semantic.ForceWithLease != nil ||
		semantic.ForceIfIncludes != nil ||
		semantic.Recursive != nil
}

func hasGitOnlySemanticFields(semantic SemanticMatchSpec) bool {
	return semantic.Remote != "" ||
		len(semantic.RemoteIn) > 0 ||
		semantic.Branch != "" ||
		len(semantic.BranchIn) > 0 ||
		semantic.Ref != "" ||
		len(semantic.RefIn) > 0 ||
		semantic.Hard != nil ||
		semantic.IncludeIgnored != nil ||
		semantic.Cached != nil ||
		semantic.Staged != nil
}

func hasAWSSemanticFields(semantic SemanticMatchSpec) bool {
	return hasAWSOnlySemanticFields(semantic) ||
		semantic.DryRun != nil
}

func hasAWSOnlySemanticFields(semantic SemanticMatchSpec) bool {
	return semantic.Service != "" ||
		len(semantic.ServiceIn) > 0 ||
		semantic.Operation != "" ||
		len(semantic.OperationIn) > 0 ||
		semantic.Profile != "" ||
		len(semantic.ProfileIn) > 0 ||
		semantic.Region != "" ||
		len(semantic.RegionIn) > 0 ||
		semantic.EndpointURL != "" ||
		semantic.EndpointURLPrefix != "" ||
		semantic.NoCLIPager != nil
}

func hasKubectlSemanticFields(semantic SemanticMatchSpec) bool {
	return semantic.Verb != "" ||
		len(semantic.VerbIn) > 0 ||
		semantic.Force != nil ||
		semantic.Recursive != nil ||
		semantic.DryRun != nil ||
		hasKubectlOnlySemanticFields(semantic)
}

func hasKubectlOnlySemanticFields(semantic SemanticMatchSpec) bool {
	return semantic.Subverb != "" ||
		len(semantic.SubverbIn) > 0 ||
		semantic.ResourceType != "" ||
		len(semantic.ResourceTypeIn) > 0 ||
		semantic.ResourceName != "" ||
		len(semantic.ResourceNameIn) > 0 ||
		semantic.Namespace != "" ||
		len(semantic.NamespaceIn) > 0 ||
		semantic.Context != "" ||
		len(semantic.ContextIn) > 0 ||
		semantic.Kubeconfig != "" ||
		semantic.AllNamespaces != nil ||
		semantic.Filename != "" ||
		len(semantic.FilenameIn) > 0 ||
		semantic.FilenamePrefix != "" ||
		semantic.Selector != "" ||
		len(semantic.SelectorContains) > 0 ||
		semantic.Container != ""
}

func hasKubectlOnlySemanticFieldsForHelmfile(semantic SemanticMatchSpec) bool {
	return semantic.Subverb != "" ||
		len(semantic.SubverbIn) > 0 ||
		semantic.ResourceType != "" ||
		len(semantic.ResourceTypeIn) > 0 ||
		semantic.ResourceName != "" ||
		len(semantic.ResourceNameIn) > 0 ||
		semantic.Context != "" ||
		len(semantic.ContextIn) > 0 ||
		semantic.Kubeconfig != "" ||
		semantic.AllNamespaces != nil ||
		semantic.Filename != "" ||
		len(semantic.FilenameIn) > 0 ||
		semantic.FilenamePrefix != "" ||
		semantic.Container != ""
}

func hasGhSemanticFields(semantic SemanticMatchSpec) bool {
	return semantic.Area != "" ||
		len(semantic.AreaIn) > 0 ||
		semantic.Verb != "" ||
		len(semantic.VerbIn) > 0 ||
		semantic.Repo != "" ||
		len(semantic.RepoIn) > 0 ||
		semantic.Hostname != "" ||
		len(semantic.HostnameIn) > 0 ||
		semantic.Web != nil ||
		semantic.Force != nil ||
		hasGhOnlySemanticFields(semantic)
}

func hasGhOnlySemanticFields(semantic SemanticMatchSpec) bool {
	return semantic.Method != "" ||
		len(semantic.MethodIn) > 0 ||
		semantic.Endpoint != "" ||
		semantic.EndpointPrefix != "" ||
		len(semantic.EndpointContains) > 0 ||
		semantic.Paginate != nil ||
		semantic.Input != nil ||
		semantic.Silent != nil ||
		semantic.IncludeHeaders != nil ||
		len(semantic.FieldKeysContains) > 0 ||
		len(semantic.RawFieldKeysContains) > 0 ||
		len(semantic.HeaderKeysContains) > 0 ||
		semantic.Org != "" ||
		len(semantic.OrgIn) > 0 ||
		semantic.EnvName != "" ||
		len(semantic.EnvNameIn) > 0 ||
		semantic.PRNumber != "" ||
		semantic.IssueNumber != "" ||
		semantic.SecretName != "" ||
		len(semantic.SecretNameIn) > 0 ||
		semantic.Tag != "" ||
		semantic.WorkflowName != "" ||
		semantic.WorkflowID != "" ||
		semantic.SearchType != "" ||
		len(semantic.SearchTypeIn) > 0 ||
		semantic.QueryContains != "" ||
		semantic.Base != "" ||
		semantic.Head != "" ||
		semantic.State != "" ||
		len(semantic.StateIn) > 0 ||
		len(semantic.LabelIn) > 0 ||
		len(semantic.AssigneeIn) > 0 ||
		semantic.TitleContains != "" ||
		semantic.BodyContains != "" ||
		semantic.Draft != nil ||
		semantic.Prerelease != nil ||
		semantic.Latest != nil ||
		semantic.Fill != nil ||
		semantic.Admin != nil ||
		semantic.Auto != nil ||
		semantic.DeleteBranch != nil ||
		semantic.MergeStrategy != "" ||
		len(semantic.MergeStrategyIn) > 0 ||
		semantic.RunID != "" ||
		semantic.Failed != nil ||
		semantic.Job != "" ||
		semantic.Debug != nil ||
		semantic.ExitStatus != nil
}

func hasGwsSemanticFields(semantic SemanticMatchSpec) bool {
	return semantic.Service != "" ||
		len(semantic.ServiceIn) > 0 ||
		len(semantic.ResourcePath) > 0 ||
		len(semantic.ResourcePathContains) > 0 ||
		semantic.Method != "" ||
		len(semantic.MethodIn) > 0 ||
		semantic.Helper != nil ||
		semantic.Mutating != nil ||
		semantic.Destructive != nil ||
		semantic.ReadOnly != nil ||
		semantic.DryRun != nil ||
		semantic.PageAll != nil ||
		semantic.Upload != nil ||
		semantic.Sanitize != nil ||
		semantic.Params != nil ||
		semantic.JSONBody != nil ||
		semantic.Unmasked != nil ||
		len(semantic.Scopes) > 0
}

func hasGhNonSharedSemanticFields(semantic SemanticMatchSpec) bool {
	return semantic.Area != "" ||
		len(semantic.AreaIn) > 0 ||
		semantic.Repo != "" ||
		len(semantic.RepoIn) > 0 ||
		semantic.Hostname != "" ||
		len(semantic.HostnameIn) > 0 ||
		semantic.Web != nil ||
		hasGhOnlySemanticFields(semantic)
}

func hasHelmfileSemanticFields(semantic SemanticMatchSpec) bool {
	return semantic.Verb != "" ||
		len(semantic.VerbIn) > 0 ||
		semantic.Namespace != "" ||
		len(semantic.NamespaceIn) > 0 ||
		semantic.NamespaceMissing != nil ||
		semantic.Selector != "" ||
		len(semantic.SelectorIn) > 0 ||
		len(semantic.SelectorContains) > 0 ||
		semantic.SelectorMissing != nil ||
		semantic.DryRun != nil ||
		hasHelmfileOnlySemanticFields(semantic)
}

func hasHelmfileOnlySemanticFields(semantic SemanticMatchSpec) bool {
	return semantic.Environment != "" ||
		len(semantic.EnvironmentIn) > 0 ||
		semantic.EnvironmentMissing != nil ||
		semantic.NamespaceMissing != nil ||
		semantic.File != "" ||
		len(semantic.FileIn) > 0 ||
		semantic.FilePrefix != "" ||
		semantic.FileMissing != nil ||
		semantic.KubeContext != "" ||
		len(semantic.KubeContextIn) > 0 ||
		semantic.KubeContextMissing != nil ||
		len(semantic.SelectorIn) > 0 ||
		semantic.SelectorMissing != nil ||
		semantic.Interactive != nil ||
		semantic.Wait != nil ||
		semantic.WaitForJobs != nil ||
		semantic.SkipDiff != nil ||
		semantic.SkipNeeds != nil ||
		semantic.IncludeNeeds != nil ||
		semantic.IncludeTransitiveNeeds != nil ||
		semantic.Purge != nil ||
		semantic.Cascade != "" ||
		len(semantic.CascadeIn) > 0 ||
		semantic.DeleteWait != nil ||
		semantic.StateValuesFile != "" ||
		len(semantic.StateValuesFileIn) > 0 ||
		len(semantic.StateValuesSetKeysContains) > 0 ||
		len(semantic.StateValuesSetStringKeysContains) > 0
}

func hasArgoCDSemanticFields(semantic SemanticMatchSpec) bool {
	return semantic.Verb != "" ||
		len(semantic.VerbIn) > 0 ||
		semantic.AppName != "" ||
		len(semantic.AppNameIn) > 0 ||
		semantic.Project != "" ||
		len(semantic.ProjectIn) > 0 ||
		semantic.Revision != ""
}

func IsZeroTerraformSemanticSpec(semantic TerraformSemanticSpec) bool {
	return semantic.Subcommand == "" && len(semantic.SubcommandIn) == 0 &&
		semantic.GlobalChdir == "" &&
		semantic.WorkspaceSubcommand == "" && len(semantic.WorkspaceSubcommandIn) == 0 &&
		semantic.StateSubcommand == "" && len(semantic.StateSubcommandIn) == 0 &&
		semantic.Target == nil && len(semantic.TargetsContains) == 0 &&
		semantic.Replace == nil && len(semantic.ReplacesContains) == 0 &&
		semantic.Destroy == nil && semantic.AutoApprove == nil && semantic.Input == nil &&
		semantic.Lock == nil && semantic.Refresh == nil && semantic.RefreshOnly == nil &&
		semantic.Out == "" && semantic.PlanFile == "" && len(semantic.VarFilesContains) == 0 &&
		semantic.Vars == nil && semantic.Backend == nil && semantic.Upgrade == nil &&
		semantic.Reconfigure == nil && semantic.MigrateState == nil && semantic.Recursive == nil &&
		semantic.Check == nil && semantic.JSON == nil && semantic.Force == nil &&
		len(semantic.FlagsContains) == 0 && len(semantic.FlagsPrefixes) == 0
}
