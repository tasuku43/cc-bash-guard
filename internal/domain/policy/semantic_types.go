package policy

type SemanticMatchSpec struct {
	Verb                             string   `yaml:"verb" json:"verb,omitempty"`
	VerbIn                           []string `yaml:"verb_in" json:"verb_in,omitempty"`
	Remote                           string   `yaml:"remote" json:"remote,omitempty"`
	RemoteIn                         []string `yaml:"remote_in" json:"remote_in,omitempty"`
	Branch                           string   `yaml:"branch" json:"branch,omitempty"`
	BranchIn                         []string `yaml:"branch_in" json:"branch_in,omitempty"`
	Ref                              string   `yaml:"ref" json:"ref,omitempty"`
	RefIn                            []string `yaml:"ref_in" json:"ref_in,omitempty"`
	Force                            *bool    `yaml:"force" json:"force,omitempty"`
	ForceWithLease                   *bool    `yaml:"force_with_lease" json:"force_with_lease,omitempty"`
	ForceIfIncludes                  *bool    `yaml:"force_if_includes" json:"force_if_includes,omitempty"`
	Hard                             *bool    `yaml:"hard" json:"hard,omitempty"`
	Recursive                        *bool    `yaml:"recursive" json:"recursive,omitempty"`
	IncludeIgnored                   *bool    `yaml:"include_ignored" json:"include_ignored,omitempty"`
	Cached                           *bool    `yaml:"cached" json:"cached,omitempty"`
	Staged                           *bool    `yaml:"staged" json:"staged,omitempty"`
	FlagsContains                    []string `yaml:"flags_contains" json:"flags_contains,omitempty"`
	FlagsPrefixes                    []string `yaml:"flags_prefixes" json:"flags_prefixes,omitempty"`
	Service                          string   `yaml:"service" json:"service,omitempty"`
	ServiceIn                        []string `yaml:"service_in" json:"service_in,omitempty"`
	Operation                        string   `yaml:"operation" json:"operation,omitempty"`
	OperationIn                      []string `yaml:"operation_in" json:"operation_in,omitempty"`
	Profile                          string   `yaml:"profile" json:"profile,omitempty"`
	ProfileIn                        []string `yaml:"profile_in" json:"profile_in,omitempty"`
	Region                           string   `yaml:"region" json:"region,omitempty"`
	RegionIn                         []string `yaml:"region_in" json:"region_in,omitempty"`
	EndpointURL                      string   `yaml:"endpoint_url" json:"endpoint_url,omitempty"`
	EndpointURLPrefix                string   `yaml:"endpoint_url_prefix" json:"endpoint_url_prefix,omitempty"`
	DryRun                           *bool    `yaml:"dry_run" json:"dry_run,omitempty"`
	NoCLIPager                       *bool    `yaml:"no_cli_pager" json:"no_cli_pager,omitempty"`
	Subverb                          string   `yaml:"subverb" json:"subverb,omitempty"`
	SubverbIn                        []string `yaml:"subverb_in" json:"subverb_in,omitempty"`
	ResourceType                     string   `yaml:"resource_type" json:"resource_type,omitempty"`
	ResourceTypeIn                   []string `yaml:"resource_type_in" json:"resource_type_in,omitempty"`
	ResourceName                     string   `yaml:"resource_name" json:"resource_name,omitempty"`
	ResourceNameIn                   []string `yaml:"resource_name_in" json:"resource_name_in,omitempty"`
	Namespace                        string   `yaml:"namespace" json:"namespace,omitempty"`
	NamespaceIn                      []string `yaml:"namespace_in" json:"namespace_in,omitempty"`
	NamespaceMissing                 *bool    `yaml:"namespace_missing" json:"namespace_missing,omitempty"`
	Context                          string   `yaml:"context" json:"context,omitempty"`
	ContextIn                        []string `yaml:"context_in" json:"context_in,omitempty"`
	Kubeconfig                       string   `yaml:"kubeconfig" json:"kubeconfig,omitempty"`
	AllNamespaces                    *bool    `yaml:"all_namespaces" json:"all_namespaces,omitempty"`
	Filename                         string   `yaml:"filename" json:"filename,omitempty"`
	FilenameIn                       []string `yaml:"filename_in" json:"filename_in,omitempty"`
	FilenamePrefix                   string   `yaml:"filename_prefix" json:"filename_prefix,omitempty"`
	Selector                         string   `yaml:"selector" json:"selector,omitempty"`
	SelectorIn                       []string `yaml:"selector_in" json:"selector_in,omitempty"`
	SelectorContains                 []string `yaml:"selector_contains" json:"selector_contains,omitempty"`
	SelectorMissing                  *bool    `yaml:"selector_missing" json:"selector_missing,omitempty"`
	Container                        string   `yaml:"container" json:"container,omitempty"`
	Environment                      string   `yaml:"environment" json:"environment,omitempty"`
	EnvironmentIn                    []string `yaml:"environment_in" json:"environment_in,omitempty"`
	EnvironmentMissing               *bool    `yaml:"environment_missing" json:"environment_missing,omitempty"`
	File                             string   `yaml:"file" json:"file,omitempty"`
	FileIn                           []string `yaml:"file_in" json:"file_in,omitempty"`
	FilePrefix                       string   `yaml:"file_prefix" json:"file_prefix,omitempty"`
	FileMissing                      *bool    `yaml:"file_missing" json:"file_missing,omitempty"`
	KubeContext                      string   `yaml:"kube_context" json:"kube_context,omitempty"`
	KubeContextIn                    []string `yaml:"kube_context_in" json:"kube_context_in,omitempty"`
	KubeContextMissing               *bool    `yaml:"kube_context_missing" json:"kube_context_missing,omitempty"`
	Interactive                      *bool    `yaml:"interactive" json:"interactive,omitempty"`
	Wait                             *bool    `yaml:"wait" json:"wait,omitempty"`
	WaitForJobs                      *bool    `yaml:"wait_for_jobs" json:"wait_for_jobs,omitempty"`
	SkipDiff                         *bool    `yaml:"skip_diff" json:"skip_diff,omitempty"`
	SkipNeeds                        *bool    `yaml:"skip_needs" json:"skip_needs,omitempty"`
	IncludeNeeds                     *bool    `yaml:"include_needs" json:"include_needs,omitempty"`
	IncludeTransitiveNeeds           *bool    `yaml:"include_transitive_needs" json:"include_transitive_needs,omitempty"`
	Purge                            *bool    `yaml:"purge" json:"purge,omitempty"`
	Cascade                          string   `yaml:"cascade" json:"cascade,omitempty"`
	CascadeIn                        []string `yaml:"cascade_in" json:"cascade_in,omitempty"`
	DeleteWait                       *bool    `yaml:"delete_wait" json:"delete_wait,omitempty"`
	StateValuesFile                  string   `yaml:"state_values_file" json:"state_values_file,omitempty"`
	StateValuesFileIn                []string `yaml:"state_values_file_in" json:"state_values_file_in,omitempty"`
	StateValuesSetKeysContains       []string `yaml:"state_values_set_keys_contains" json:"state_values_set_keys_contains,omitempty"`
	StateValuesSetStringKeysContains []string `yaml:"state_values_set_string_keys_contains" json:"state_values_set_string_keys_contains,omitempty"`
	Area                             string   `yaml:"area" json:"area,omitempty"`
	AreaIn                           []string `yaml:"area_in" json:"area_in,omitempty"`
	Repo                             string   `yaml:"repo" json:"repo,omitempty"`
	RepoIn                           []string `yaml:"repo_in" json:"repo_in,omitempty"`
	Org                              string   `yaml:"org" json:"org,omitempty"`
	OrgIn                            []string `yaml:"org_in" json:"org_in,omitempty"`
	EnvName                          string   `yaml:"env" json:"env,omitempty"`
	EnvNameIn                        []string `yaml:"env_in" json:"env_in,omitempty"`
	AppName                          string   `yaml:"app_name" json:"app_name,omitempty"`
	AppNameIn                        []string `yaml:"app_name_in" json:"app_name_in,omitempty"`
	Project                          string   `yaml:"project" json:"project,omitempty"`
	ProjectIn                        []string `yaml:"project_in" json:"project_in,omitempty"`
	Revision                         string   `yaml:"revision" json:"revision,omitempty"`
	Hostname                         string   `yaml:"hostname" json:"hostname,omitempty"`
	HostnameIn                       []string `yaml:"hostname_in" json:"hostname_in,omitempty"`
	Web                              *bool    `yaml:"web" json:"web,omitempty"`
	Method                           string   `yaml:"method" json:"method,omitempty"`
	MethodIn                         []string `yaml:"method_in" json:"method_in,omitempty"`
	ResourcePath                     []string `yaml:"resource_path" json:"resource_path,omitempty"`
	ResourcePathContains             []string `yaml:"resource_path_contains" json:"resource_path_contains,omitempty"`
	Helper                           *bool    `yaml:"helper" json:"helper,omitempty"`
	Mutating                         *bool    `yaml:"mutating" json:"mutating,omitempty"`
	Destructive                      *bool    `yaml:"destructive" json:"destructive,omitempty"`
	ReadOnly                         *bool    `yaml:"read_only" json:"read_only,omitempty"`
	PageAll                          *bool    `yaml:"page_all" json:"page_all,omitempty"`
	Upload                           *bool    `yaml:"upload" json:"upload,omitempty"`
	Sanitize                         *bool    `yaml:"sanitize" json:"sanitize,omitempty"`
	Params                           *bool    `yaml:"params" json:"params,omitempty"`
	JSONBody                         *bool    `yaml:"json_body" json:"json_body,omitempty"`
	Unmasked                         *bool    `yaml:"unmasked" json:"unmasked,omitempty"`
	Scopes                           []string `yaml:"scopes" json:"scopes,omitempty"`
	Endpoint                         string   `yaml:"endpoint" json:"endpoint,omitempty"`
	EndpointPrefix                   string   `yaml:"endpoint_prefix" json:"endpoint_prefix,omitempty"`
	EndpointContains                 []string `yaml:"endpoint_contains" json:"endpoint_contains,omitempty"`
	Paginate                         *bool    `yaml:"paginate" json:"paginate,omitempty"`
	Input                            *bool    `yaml:"input" json:"input,omitempty"`
	Silent                           *bool    `yaml:"silent" json:"silent,omitempty"`
	IncludeHeaders                   *bool    `yaml:"include_headers" json:"include_headers,omitempty"`
	FieldKeysContains                []string `yaml:"field_keys_contains" json:"field_keys_contains,omitempty"`
	RawFieldKeysContains             []string `yaml:"raw_field_keys_contains" json:"raw_field_keys_contains,omitempty"`
	HeaderKeysContains               []string `yaml:"header_keys_contains" json:"header_keys_contains,omitempty"`
	PRNumber                         string   `yaml:"pr_number" json:"pr_number,omitempty"`
	IssueNumber                      string   `yaml:"issue_number" json:"issue_number,omitempty"`
	SecretName                       string   `yaml:"secret_name" json:"secret_name,omitempty"`
	SecretNameIn                     []string `yaml:"secret_name_in" json:"secret_name_in,omitempty"`
	Tag                              string   `yaml:"tag" json:"tag,omitempty"`
	WorkflowName                     string   `yaml:"workflow_name" json:"workflow_name,omitempty"`
	WorkflowID                       string   `yaml:"workflow_id" json:"workflow_id,omitempty"`
	SearchType                       string   `yaml:"search_type" json:"search_type,omitempty"`
	SearchTypeIn                     []string `yaml:"search_type_in" json:"search_type_in,omitempty"`
	QueryContains                    string   `yaml:"query_contains" json:"query_contains,omitempty"`
	Base                             string   `yaml:"base" json:"base,omitempty"`
	Head                             string   `yaml:"head" json:"head,omitempty"`
	State                            string   `yaml:"state" json:"state,omitempty"`
	StateIn                          []string `yaml:"state_in" json:"state_in,omitempty"`
	LabelIn                          []string `yaml:"label_in" json:"label_in,omitempty"`
	AssigneeIn                       []string `yaml:"assignee_in" json:"assignee_in,omitempty"`
	TitleContains                    string   `yaml:"title_contains" json:"title_contains,omitempty"`
	BodyContains                     string   `yaml:"body_contains" json:"body_contains,omitempty"`
	Draft                            *bool    `yaml:"draft" json:"draft,omitempty"`
	Prerelease                       *bool    `yaml:"prerelease" json:"prerelease,omitempty"`
	Latest                           *bool    `yaml:"latest" json:"latest,omitempty"`
	Fill                             *bool    `yaml:"fill" json:"fill,omitempty"`
	Admin                            *bool    `yaml:"admin" json:"admin,omitempty"`
	Auto                             *bool    `yaml:"auto" json:"auto,omitempty"`
	DeleteBranch                     *bool    `yaml:"delete_branch" json:"delete_branch,omitempty"`
	MergeStrategy                    string   `yaml:"merge_strategy" json:"merge_strategy,omitempty"`
	MergeStrategyIn                  []string `yaml:"merge_strategy_in" json:"merge_strategy_in,omitempty"`
	RunID                            string   `yaml:"run_id" json:"run_id,omitempty"`
	Failed                           *bool    `yaml:"failed" json:"failed,omitempty"`
	Job                              string   `yaml:"job" json:"job,omitempty"`
	Debug                            *bool    `yaml:"debug" json:"debug,omitempty"`
	ExitStatus                       *bool    `yaml:"exit_status" json:"exit_status,omitempty"`
	Subcommand                       string   `yaml:"subcommand" json:"subcommand,omitempty"`
	SubcommandIn                     []string `yaml:"subcommand_in" json:"subcommand_in,omitempty"`
	GlobalChdir                      string   `yaml:"global_chdir" json:"global_chdir,omitempty"`
	WorkspaceSubcommand              string   `yaml:"workspace_subcommand" json:"workspace_subcommand,omitempty"`
	WorkspaceSubcommandIn            []string `yaml:"workspace_subcommand_in" json:"workspace_subcommand_in,omitempty"`
	StateSubcommand                  string   `yaml:"state_subcommand" json:"state_subcommand,omitempty"`
	StateSubcommandIn                []string `yaml:"state_subcommand_in" json:"state_subcommand_in,omitempty"`
	Target                           *bool    `yaml:"target" json:"target,omitempty"`
	TargetsContains                  []string `yaml:"targets_contains" json:"targets_contains,omitempty"`
	Replace                          *bool    `yaml:"replace" json:"replace,omitempty"`
	ReplacesContains                 []string `yaml:"replaces_contains" json:"replaces_contains,omitempty"`
	Destroy                          *bool    `yaml:"destroy" json:"destroy,omitempty"`
	AutoApprove                      *bool    `yaml:"auto_approve" json:"auto_approve,omitempty"`
	Lock                             *bool    `yaml:"lock" json:"lock,omitempty"`
	Refresh                          *bool    `yaml:"refresh" json:"refresh,omitempty"`
	RefreshOnly                      *bool    `yaml:"refresh_only" json:"refresh_only,omitempty"`
	Out                              string   `yaml:"out" json:"out,omitempty"`
	PlanFile                         string   `yaml:"plan_file" json:"plan_file,omitempty"`
	VarFilesContains                 []string `yaml:"var_files_contains" json:"var_files_contains,omitempty"`
	Vars                             *bool    `yaml:"vars" json:"vars,omitempty"`
	Backend                          *bool    `yaml:"backend" json:"backend,omitempty"`
	Upgrade                          *bool    `yaml:"upgrade" json:"upgrade,omitempty"`
	Reconfigure                      *bool    `yaml:"reconfigure" json:"reconfigure,omitempty"`
	MigrateState                     *bool    `yaml:"migrate_state" json:"migrate_state,omitempty"`
	Check                            *bool    `yaml:"check" json:"check,omitempty"`
	JSON                             *bool    `yaml:"json" json:"json,omitempty"`
}

type GitSemanticSpec struct {
	Verb            string
	VerbIn          []string
	Remote          string
	RemoteIn        []string
	Branch          string
	BranchIn        []string
	Ref             string
	RefIn           []string
	Force           *bool
	ForceWithLease  *bool
	ForceIfIncludes *bool
	Hard            *bool
	Recursive       *bool
	IncludeIgnored  *bool
	Cached          *bool
	Staged          *bool
	FlagsContains   []string
	FlagsPrefixes   []string
}

type AWSSemanticSpec struct {
	Service           string
	ServiceIn         []string
	Operation         string
	OperationIn       []string
	Profile           string
	ProfileIn         []string
	Region            string
	RegionIn          []string
	EndpointURL       string
	EndpointURLPrefix string
	DryRun            *bool
	NoCLIPager        *bool
	FlagsContains     []string
	FlagsPrefixes     []string
}

type KubectlSemanticSpec struct {
	Verb             string
	VerbIn           []string
	Subverb          string
	SubverbIn        []string
	ResourceType     string
	ResourceTypeIn   []string
	ResourceName     string
	ResourceNameIn   []string
	Namespace        string
	NamespaceIn      []string
	NamespaceMissing *bool
	Context          string
	ContextIn        []string
	Kubeconfig       string
	AllNamespaces    *bool
	Filename         string
	FilenameIn       []string
	FilenamePrefix   string
	Selector         string
	SelectorIn       []string
	SelectorContains []string
	SelectorMissing  *bool
	Container        string
	DryRun           *bool
	Force            *bool
	Recursive        *bool
	FlagsContains    []string
	FlagsPrefixes    []string
}

type GHSemanticSpec struct {
	Area                 string
	AreaIn               []string
	Verb                 string
	VerbIn               []string
	Repo                 string
	RepoIn               []string
	Org                  string
	OrgIn                []string
	EnvName              string
	EnvNameIn            []string
	Hostname             string
	HostnameIn           []string
	Web                  *bool
	Method               string
	MethodIn             []string
	Endpoint             string
	EndpointPrefix       string
	EndpointContains     []string
	Paginate             *bool
	Input                *bool
	Silent               *bool
	IncludeHeaders       *bool
	FieldKeysContains    []string
	RawFieldKeysContains []string
	HeaderKeysContains   []string
	PRNumber             string
	IssueNumber          string
	SecretName           string
	SecretNameIn         []string
	Tag                  string
	WorkflowName         string
	WorkflowID           string
	SearchType           string
	SearchTypeIn         []string
	QueryContains        string
	Base                 string
	Head                 string
	Ref                  string
	RefIn                []string
	State                string
	StateIn              []string
	LabelIn              []string
	AssigneeIn           []string
	TitleContains        string
	BodyContains         string
	Draft                *bool
	Prerelease           *bool
	Latest               *bool
	Fill                 *bool
	Force                *bool
	Admin                *bool
	Auto                 *bool
	DeleteBranch         *bool
	MergeStrategy        string
	MergeStrategyIn      []string
	RunID                string
	Failed               *bool
	Job                  string
	Debug                *bool
	ExitStatus           *bool
	FlagsContains        []string
	FlagsPrefixes        []string
}

type GwsSemanticSpec struct {
	Service              string
	ServiceIn            []string
	ResourcePath         []string
	ResourcePathContains []string
	Method               string
	MethodIn             []string
	Helper               *bool
	Mutating             *bool
	Destructive          *bool
	ReadOnly             *bool
	DryRun               *bool
	PageAll              *bool
	Upload               *bool
	Sanitize             *bool
	Params               *bool
	JSONBody             *bool
	Unmasked             *bool
	Scopes               []string
	FlagsContains        []string
	FlagsPrefixes        []string
}

type HelmfileSemanticSpec struct {
	Verb                             string
	VerbIn                           []string
	Environment                      string
	EnvironmentIn                    []string
	EnvironmentMissing               *bool
	File                             string
	FileIn                           []string
	FilePrefix                       string
	FileMissing                      *bool
	Namespace                        string
	NamespaceIn                      []string
	NamespaceMissing                 *bool
	KubeContext                      string
	KubeContextIn                    []string
	KubeContextMissing               *bool
	Selector                         string
	SelectorIn                       []string
	SelectorContains                 []string
	SelectorMissing                  *bool
	Interactive                      *bool
	DryRun                           *bool
	Wait                             *bool
	WaitForJobs                      *bool
	SkipDiff                         *bool
	SkipNeeds                        *bool
	IncludeNeeds                     *bool
	IncludeTransitiveNeeds           *bool
	Purge                            *bool
	Cascade                          string
	CascadeIn                        []string
	DeleteWait                       *bool
	StateValuesFile                  string
	StateValuesFileIn                []string
	StateValuesSetKeysContains       []string
	StateValuesSetStringKeysContains []string
	FlagsContains                    []string
	FlagsPrefixes                    []string
}

type ArgoCDSemanticSpec struct {
	Verb          string
	VerbIn        []string
	AppName       string
	AppNameIn     []string
	Project       string
	ProjectIn     []string
	Revision      string
	FlagsContains []string
	FlagsPrefixes []string
}

type TerraformSemanticSpec struct {
	Subcommand            string
	SubcommandIn          []string
	GlobalChdir           string
	WorkspaceSubcommand   string
	WorkspaceSubcommandIn []string
	StateSubcommand       string
	StateSubcommandIn     []string
	Target                *bool
	TargetsContains       []string
	Replace               *bool
	ReplacesContains      []string
	Destroy               *bool
	AutoApprove           *bool
	Input                 *bool
	Lock                  *bool
	Refresh               *bool
	RefreshOnly           *bool
	Out                   string
	PlanFile              string
	VarFilesContains      []string
	Vars                  *bool
	Backend               *bool
	Upgrade               *bool
	Reconfigure           *bool
	MigrateState          *bool
	Recursive             *bool
	Check                 *bool
	JSON                  *bool
	Force                 *bool
	FlagsContains         []string
	FlagsPrefixes         []string
}

func (s SemanticMatchSpec) Git() GitSemanticSpec {
	return GitSemanticSpec{
		Verb: s.Verb, VerbIn: s.VerbIn, Remote: s.Remote, RemoteIn: s.RemoteIn,
		Branch: s.Branch, BranchIn: s.BranchIn, Ref: s.Ref, RefIn: s.RefIn,
		Force: s.Force, ForceWithLease: s.ForceWithLease, ForceIfIncludes: s.ForceIfIncludes,
		Hard: s.Hard, Recursive: s.Recursive, IncludeIgnored: s.IncludeIgnored,
		Cached: s.Cached, Staged: s.Staged, FlagsContains: s.FlagsContains, FlagsPrefixes: s.FlagsPrefixes,
	}
}

func (s SemanticMatchSpec) AWS() AWSSemanticSpec {
	return AWSSemanticSpec{
		Service: s.Service, ServiceIn: s.ServiceIn, Operation: s.Operation, OperationIn: s.OperationIn,
		Profile: s.Profile, ProfileIn: s.ProfileIn, Region: s.Region, RegionIn: s.RegionIn,
		EndpointURL: s.EndpointURL, EndpointURLPrefix: s.EndpointURLPrefix, DryRun: s.DryRun,
		NoCLIPager: s.NoCLIPager, FlagsContains: s.FlagsContains, FlagsPrefixes: s.FlagsPrefixes,
	}
}

func (s SemanticMatchSpec) Kubectl() KubectlSemanticSpec {
	return KubectlSemanticSpec{
		Verb: s.Verb, VerbIn: s.VerbIn, Subverb: s.Subverb, SubverbIn: s.SubverbIn,
		ResourceType: s.ResourceType, ResourceTypeIn: s.ResourceTypeIn, ResourceName: s.ResourceName,
		ResourceNameIn: s.ResourceNameIn, Namespace: s.Namespace, NamespaceIn: s.NamespaceIn,
		NamespaceMissing: s.NamespaceMissing, Context: s.Context, ContextIn: s.ContextIn,
		Kubeconfig: s.Kubeconfig, AllNamespaces: s.AllNamespaces, Filename: s.Filename,
		FilenameIn: s.FilenameIn, FilenamePrefix: s.FilenamePrefix, Selector: s.Selector,
		SelectorIn: s.SelectorIn, SelectorContains: s.SelectorContains, SelectorMissing: s.SelectorMissing,
		Container: s.Container, DryRun: s.DryRun, Force: s.Force, Recursive: s.Recursive,
		FlagsContains: s.FlagsContains, FlagsPrefixes: s.FlagsPrefixes,
	}
}

func (s SemanticMatchSpec) GH() GHSemanticSpec {
	return GHSemanticSpec{
		Area: s.Area, AreaIn: s.AreaIn, Verb: s.Verb, VerbIn: s.VerbIn, Repo: s.Repo, RepoIn: s.RepoIn,
		Org: s.Org, OrgIn: s.OrgIn, EnvName: s.EnvName, EnvNameIn: s.EnvNameIn,
		Hostname: s.Hostname, HostnameIn: s.HostnameIn, Web: s.Web, Method: s.Method,
		MethodIn: s.MethodIn, Endpoint: s.Endpoint, EndpointPrefix: s.EndpointPrefix,
		EndpointContains: s.EndpointContains, Paginate: s.Paginate, Input: s.Input, Silent: s.Silent,
		IncludeHeaders: s.IncludeHeaders, FieldKeysContains: s.FieldKeysContains,
		RawFieldKeysContains: s.RawFieldKeysContains, HeaderKeysContains: s.HeaderKeysContains,
		PRNumber: s.PRNumber, IssueNumber: s.IssueNumber, SecretName: s.SecretName,
		SecretNameIn: s.SecretNameIn, Tag: s.Tag, WorkflowName: s.WorkflowName, WorkflowID: s.WorkflowID,
		SearchType: s.SearchType, SearchTypeIn: s.SearchTypeIn, QueryContains: s.QueryContains,
		Base: s.Base, Head: s.Head, Ref: s.Ref, RefIn: s.RefIn, State: s.State, StateIn: s.StateIn,
		LabelIn: s.LabelIn, AssigneeIn: s.AssigneeIn, TitleContains: s.TitleContains,
		BodyContains: s.BodyContains, Draft: s.Draft, Prerelease: s.Prerelease, Latest: s.Latest,
		Fill: s.Fill, Force: s.Force, Admin: s.Admin, Auto: s.Auto, DeleteBranch: s.DeleteBranch,
		MergeStrategy: s.MergeStrategy, MergeStrategyIn: s.MergeStrategyIn, RunID: s.RunID,
		Failed: s.Failed, Job: s.Job, Debug: s.Debug, ExitStatus: s.ExitStatus,
		FlagsContains: s.FlagsContains, FlagsPrefixes: s.FlagsPrefixes,
	}
}

func (s SemanticMatchSpec) Gws() GwsSemanticSpec {
	return GwsSemanticSpec{
		Service: s.Service, ServiceIn: s.ServiceIn, ResourcePath: s.ResourcePath,
		ResourcePathContains: s.ResourcePathContains, Method: s.Method, MethodIn: s.MethodIn,
		Helper: s.Helper, Mutating: s.Mutating, Destructive: s.Destructive, ReadOnly: s.ReadOnly,
		DryRun: s.DryRun, PageAll: s.PageAll, Upload: s.Upload, Sanitize: s.Sanitize,
		Params: s.Params, JSONBody: s.JSONBody, Unmasked: s.Unmasked, Scopes: s.Scopes,
		FlagsContains: s.FlagsContains, FlagsPrefixes: s.FlagsPrefixes,
	}
}

func (s SemanticMatchSpec) Helmfile() HelmfileSemanticSpec {
	return HelmfileSemanticSpec{
		Verb: s.Verb, VerbIn: s.VerbIn, Environment: s.Environment, EnvironmentIn: s.EnvironmentIn,
		EnvironmentMissing: s.EnvironmentMissing, File: s.File, FileIn: s.FileIn, FilePrefix: s.FilePrefix,
		FileMissing: s.FileMissing, Namespace: s.Namespace, NamespaceIn: s.NamespaceIn,
		NamespaceMissing: s.NamespaceMissing, KubeContext: s.KubeContext, KubeContextIn: s.KubeContextIn,
		KubeContextMissing: s.KubeContextMissing, Selector: s.Selector, SelectorIn: s.SelectorIn,
		SelectorContains: s.SelectorContains, SelectorMissing: s.SelectorMissing, Interactive: s.Interactive,
		DryRun: s.DryRun, Wait: s.Wait, WaitForJobs: s.WaitForJobs, SkipDiff: s.SkipDiff,
		SkipNeeds: s.SkipNeeds, IncludeNeeds: s.IncludeNeeds, IncludeTransitiveNeeds: s.IncludeTransitiveNeeds,
		Purge: s.Purge, Cascade: s.Cascade, CascadeIn: s.CascadeIn, DeleteWait: s.DeleteWait,
		StateValuesFile: s.StateValuesFile, StateValuesFileIn: s.StateValuesFileIn,
		StateValuesSetKeysContains:       s.StateValuesSetKeysContains,
		StateValuesSetStringKeysContains: s.StateValuesSetStringKeysContains,
		FlagsContains:                    s.FlagsContains, FlagsPrefixes: s.FlagsPrefixes,
	}
}

func (s SemanticMatchSpec) ArgoCD() ArgoCDSemanticSpec {
	return ArgoCDSemanticSpec{
		Verb: s.Verb, VerbIn: s.VerbIn, AppName: s.AppName, AppNameIn: s.AppNameIn,
		Project: s.Project, ProjectIn: s.ProjectIn, Revision: s.Revision,
		FlagsContains: s.FlagsContains, FlagsPrefixes: s.FlagsPrefixes,
	}
}

func (s SemanticMatchSpec) Terraform() TerraformSemanticSpec {
	return TerraformSemanticSpec{
		Subcommand: s.Subcommand, SubcommandIn: s.SubcommandIn, GlobalChdir: s.GlobalChdir,
		WorkspaceSubcommand: s.WorkspaceSubcommand, WorkspaceSubcommandIn: s.WorkspaceSubcommandIn,
		StateSubcommand: s.StateSubcommand, StateSubcommandIn: s.StateSubcommandIn,
		Target: s.Target, TargetsContains: s.TargetsContains, Replace: s.Replace,
		ReplacesContains: s.ReplacesContains, Destroy: s.Destroy, AutoApprove: s.AutoApprove,
		Input: s.Input, Lock: s.Lock, Refresh: s.Refresh, RefreshOnly: s.RefreshOnly,
		Out: s.Out, PlanFile: s.PlanFile, VarFilesContains: s.VarFilesContains, Vars: s.Vars,
		Backend: s.Backend, Upgrade: s.Upgrade, Reconfigure: s.Reconfigure,
		MigrateState: s.MigrateState, Recursive: s.Recursive, Check: s.Check, JSON: s.JSON,
		Force: s.Force, FlagsContains: s.FlagsContains, FlagsPrefixes: s.FlagsPrefixes,
	}
}
