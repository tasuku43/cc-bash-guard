package policy

func (s SemanticMatchSpec) fieldsUsed() []string {
	var fields []string
	if s.Verb != "" {
		fields = append(fields, "verb")
	}
	if len(s.VerbIn) > 0 {
		fields = append(fields, "verb_in")
	}
	if s.Remote != "" {
		fields = append(fields, "remote")
	}
	if len(s.RemoteIn) > 0 {
		fields = append(fields, "remote_in")
	}
	if s.Branch != "" {
		fields = append(fields, "branch")
	}
	if len(s.BranchIn) > 0 {
		fields = append(fields, "branch_in")
	}
	if s.Ref != "" {
		fields = append(fields, "ref")
	}
	if len(s.RefIn) > 0 {
		fields = append(fields, "ref_in")
	}
	if s.Force != nil {
		fields = append(fields, "force")
	}
	if s.ForceWithLease != nil {
		fields = append(fields, "force_with_lease")
	}
	if s.ForceIfIncludes != nil {
		fields = append(fields, "force_if_includes")
	}
	if s.Hard != nil {
		fields = append(fields, "hard")
	}
	if s.Recursive != nil {
		fields = append(fields, "recursive")
	}
	if s.IncludeIgnored != nil {
		fields = append(fields, "include_ignored")
	}
	if s.Cached != nil {
		fields = append(fields, "cached")
	}
	if s.Staged != nil {
		fields = append(fields, "staged")
	}
	if len(s.FlagsContains) > 0 {
		fields = append(fields, "flags_contains")
	}
	if len(s.FlagsPrefixes) > 0 {
		fields = append(fields, "flags_prefixes")
	}
	if s.Service != "" {
		fields = append(fields, "service")
	}
	if len(s.ServiceIn) > 0 {
		fields = append(fields, "service_in")
	}
	if s.Operation != "" {
		fields = append(fields, "operation")
	}
	if len(s.OperationIn) > 0 {
		fields = append(fields, "operation_in")
	}
	if s.Profile != "" {
		fields = append(fields, "profile")
	}
	if len(s.ProfileIn) > 0 {
		fields = append(fields, "profile_in")
	}
	if s.Region != "" {
		fields = append(fields, "region")
	}
	if len(s.RegionIn) > 0 {
		fields = append(fields, "region_in")
	}
	if s.EndpointURL != "" {
		fields = append(fields, "endpoint_url")
	}
	if s.EndpointURLPrefix != "" {
		fields = append(fields, "endpoint_url_prefix")
	}
	if s.DryRun != nil {
		fields = append(fields, "dry_run")
	}
	if s.NoCLIPager != nil {
		fields = append(fields, "no_cli_pager")
	}
	if s.Subverb != "" {
		fields = append(fields, "subverb")
	}
	if len(s.SubverbIn) > 0 {
		fields = append(fields, "subverb_in")
	}
	if s.ResourceType != "" {
		fields = append(fields, "resource_type")
	}
	if len(s.ResourceTypeIn) > 0 {
		fields = append(fields, "resource_type_in")
	}
	if s.ResourceName != "" {
		fields = append(fields, "resource_name")
	}
	if len(s.ResourceNameIn) > 0 {
		fields = append(fields, "resource_name_in")
	}
	if s.Namespace != "" {
		fields = append(fields, "namespace")
	}
	if len(s.NamespaceIn) > 0 {
		fields = append(fields, "namespace_in")
	}
	if s.NamespaceMissing != nil {
		fields = append(fields, "namespace_missing")
	}
	if s.Context != "" {
		fields = append(fields, "context")
	}
	if len(s.ContextIn) > 0 {
		fields = append(fields, "context_in")
	}
	if s.Kubeconfig != "" {
		fields = append(fields, "kubeconfig")
	}
	if s.AllNamespaces != nil {
		fields = append(fields, "all_namespaces")
	}
	if s.Filename != "" {
		fields = append(fields, "filename")
	}
	if len(s.FilenameIn) > 0 {
		fields = append(fields, "filename_in")
	}
	if s.FilenamePrefix != "" {
		fields = append(fields, "filename_prefix")
	}
	if s.Selector != "" {
		fields = append(fields, "selector")
	}
	if len(s.SelectorIn) > 0 {
		fields = append(fields, "selector_in")
	}
	if len(s.SelectorContains) > 0 {
		fields = append(fields, "selector_contains")
	}
	if s.SelectorMissing != nil {
		fields = append(fields, "selector_missing")
	}
	if s.Container != "" {
		fields = append(fields, "container")
	}
	if s.Environment != "" {
		fields = append(fields, "environment")
	}
	if len(s.EnvironmentIn) > 0 {
		fields = append(fields, "environment_in")
	}
	if s.EnvironmentMissing != nil {
		fields = append(fields, "environment_missing")
	}
	if s.File != "" {
		fields = append(fields, "file")
	}
	if len(s.FileIn) > 0 {
		fields = append(fields, "file_in")
	}
	if s.FilePrefix != "" {
		fields = append(fields, "file_prefix")
	}
	if s.FileMissing != nil {
		fields = append(fields, "file_missing")
	}
	if s.KubeContext != "" {
		fields = append(fields, "kube_context")
	}
	if len(s.KubeContextIn) > 0 {
		fields = append(fields, "kube_context_in")
	}
	if s.KubeContextMissing != nil {
		fields = append(fields, "kube_context_missing")
	}
	if s.Interactive != nil {
		fields = append(fields, "interactive")
	}
	if s.Wait != nil {
		fields = append(fields, "wait")
	}
	if s.WaitForJobs != nil {
		fields = append(fields, "wait_for_jobs")
	}
	if s.SkipDiff != nil {
		fields = append(fields, "skip_diff")
	}
	if s.SkipNeeds != nil {
		fields = append(fields, "skip_needs")
	}
	if s.IncludeNeeds != nil {
		fields = append(fields, "include_needs")
	}
	if s.IncludeTransitiveNeeds != nil {
		fields = append(fields, "include_transitive_needs")
	}
	if s.Purge != nil {
		fields = append(fields, "purge")
	}
	if s.Cascade != "" {
		fields = append(fields, "cascade")
	}
	if len(s.CascadeIn) > 0 {
		fields = append(fields, "cascade_in")
	}
	if s.DeleteWait != nil {
		fields = append(fields, "delete_wait")
	}
	if s.StateValuesFile != "" {
		fields = append(fields, "state_values_file")
	}
	if len(s.StateValuesFileIn) > 0 {
		fields = append(fields, "state_values_file_in")
	}
	if len(s.StateValuesSetKeysContains) > 0 {
		fields = append(fields, "state_values_set_keys_contains")
	}
	if len(s.StateValuesSetStringKeysContains) > 0 {
		fields = append(fields, "state_values_set_string_keys_contains")
	}
	if s.Area != "" {
		fields = append(fields, "area")
	}
	if len(s.AreaIn) > 0 {
		fields = append(fields, "area_in")
	}
	if s.Repo != "" {
		fields = append(fields, "repo")
	}
	if len(s.RepoIn) > 0 {
		fields = append(fields, "repo_in")
	}
	if s.Org != "" {
		fields = append(fields, "org")
	}
	if len(s.OrgIn) > 0 {
		fields = append(fields, "org_in")
	}
	if s.EnvName != "" {
		fields = append(fields, "env")
	}
	if len(s.EnvNameIn) > 0 {
		fields = append(fields, "env_in")
	}
	if s.AppName != "" {
		fields = append(fields, "app_name")
	}
	if len(s.AppNameIn) > 0 {
		fields = append(fields, "app_name_in")
	}
	if s.Project != "" {
		fields = append(fields, "project")
	}
	if len(s.ProjectIn) > 0 {
		fields = append(fields, "project_in")
	}
	if s.Revision != "" {
		fields = append(fields, "revision")
	}
	if s.Hostname != "" {
		fields = append(fields, "hostname")
	}
	if len(s.HostnameIn) > 0 {
		fields = append(fields, "hostname_in")
	}
	if s.Web != nil {
		fields = append(fields, "web")
	}
	if s.Method != "" {
		fields = append(fields, "method")
	}
	if len(s.MethodIn) > 0 {
		fields = append(fields, "method_in")
	}
	if len(s.ResourcePath) > 0 {
		fields = append(fields, "resource_path")
	}
	if len(s.ResourcePathContains) > 0 {
		fields = append(fields, "resource_path_contains")
	}
	if s.Helper != nil {
		fields = append(fields, "helper")
	}
	if s.Mutating != nil {
		fields = append(fields, "mutating")
	}
	if s.Destructive != nil {
		fields = append(fields, "destructive")
	}
	if s.ReadOnly != nil {
		fields = append(fields, "read_only")
	}
	if s.PageAll != nil {
		fields = append(fields, "page_all")
	}
	if s.Upload != nil {
		fields = append(fields, "upload")
	}
	if s.Sanitize != nil {
		fields = append(fields, "sanitize")
	}
	if s.Params != nil {
		fields = append(fields, "params")
	}
	if s.JSONBody != nil {
		fields = append(fields, "json_body")
	}
	if s.Unmasked != nil {
		fields = append(fields, "unmasked")
	}
	if len(s.Scopes) > 0 {
		fields = append(fields, "scopes")
	}
	if s.Endpoint != "" {
		fields = append(fields, "endpoint")
	}
	if s.EndpointPrefix != "" {
		fields = append(fields, "endpoint_prefix")
	}
	if len(s.EndpointContains) > 0 {
		fields = append(fields, "endpoint_contains")
	}
	if s.Paginate != nil {
		fields = append(fields, "paginate")
	}
	if s.Input != nil {
		fields = append(fields, "input")
	}
	if s.Silent != nil {
		fields = append(fields, "silent")
	}
	if s.IncludeHeaders != nil {
		fields = append(fields, "include_headers")
	}
	if len(s.FieldKeysContains) > 0 {
		fields = append(fields, "field_keys_contains")
	}
	if len(s.RawFieldKeysContains) > 0 {
		fields = append(fields, "raw_field_keys_contains")
	}
	if len(s.HeaderKeysContains) > 0 {
		fields = append(fields, "header_keys_contains")
	}
	if s.PRNumber != "" {
		fields = append(fields, "pr_number")
	}
	if s.IssueNumber != "" {
		fields = append(fields, "issue_number")
	}
	if s.SecretName != "" {
		fields = append(fields, "secret_name")
	}
	if len(s.SecretNameIn) > 0 {
		fields = append(fields, "secret_name_in")
	}
	if s.Tag != "" {
		fields = append(fields, "tag")
	}
	if s.WorkflowName != "" {
		fields = append(fields, "workflow_name")
	}
	if s.WorkflowID != "" {
		fields = append(fields, "workflow_id")
	}
	if s.SearchType != "" {
		fields = append(fields, "search_type")
	}
	if len(s.SearchTypeIn) > 0 {
		fields = append(fields, "search_type_in")
	}
	if s.QueryContains != "" {
		fields = append(fields, "query_contains")
	}
	if s.Base != "" {
		fields = append(fields, "base")
	}
	if s.Head != "" {
		fields = append(fields, "head")
	}
	if s.State != "" {
		fields = append(fields, "state")
	}
	if len(s.StateIn) > 0 {
		fields = append(fields, "state_in")
	}
	if len(s.LabelIn) > 0 {
		fields = append(fields, "label_in")
	}
	if len(s.AssigneeIn) > 0 {
		fields = append(fields, "assignee_in")
	}
	if s.TitleContains != "" {
		fields = append(fields, "title_contains")
	}
	if s.BodyContains != "" {
		fields = append(fields, "body_contains")
	}
	if s.Draft != nil {
		fields = append(fields, "draft")
	}
	if s.Prerelease != nil {
		fields = append(fields, "prerelease")
	}
	if s.Latest != nil {
		fields = append(fields, "latest")
	}
	if s.Fill != nil {
		fields = append(fields, "fill")
	}
	if s.Admin != nil {
		fields = append(fields, "admin")
	}
	if s.Auto != nil {
		fields = append(fields, "auto")
	}
	if s.DeleteBranch != nil {
		fields = append(fields, "delete_branch")
	}
	if s.MergeStrategy != "" {
		fields = append(fields, "merge_strategy")
	}
	if len(s.MergeStrategyIn) > 0 {
		fields = append(fields, "merge_strategy_in")
	}
	if s.RunID != "" {
		fields = append(fields, "run_id")
	}
	if s.Failed != nil {
		fields = append(fields, "failed")
	}
	if s.Job != "" {
		fields = append(fields, "job")
	}
	if s.Debug != nil {
		fields = append(fields, "debug")
	}
	if s.ExitStatus != nil {
		fields = append(fields, "exit_status")
	}
	return fields
}
