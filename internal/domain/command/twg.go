package command

import "strings"

type TWGParser struct{}

type twgActionClass struct {
	verb     string
	readOnly bool
	mutating bool
	inferred bool
}

type twgActionGroup struct {
	prefix string
	read   string
	write  string
}

func init() {
	RegisterDefaultParser(TWGParser{})
}

func (TWGParser) Program() string {
	return "twg"
}

func (TWGParser) Parse(base Command) (Command, bool) {
	if base.Program != "twg" {
		return Command{}, false
	}

	cmd := base
	cmd.Parser = TWGParser{}.Program()
	cmd.SemanticParser = TWGParser{}.Program()

	words, globalOptions, complete := twgPathWords(base.RawWords)
	cmd.GlobalOptions = globalOptions
	semantic := &TWGSemantic{}
	if len(words) > 0 {
		semantic.Namespace = canonicalTWGNamespace(words[0])
	}

	if twgHasHelpFlag(base.RawWords) {
		semantic.Verb = "help"
		semantic.ReadOnly = true
		cmd.ActionPath = twgHelpActionPath(words)
		cmd.TWG = semantic
		return cmd, true
	}
	if twgHasVersionFlag(base.RawWords) {
		semantic.Verb = "version"
		semantic.ReadOnly = true
		cmd.ActionPath = []string{"version"}
		cmd.TWG = semantic
		return cmd, true
	}
	if len(words) == 0 {
		cmd.TWG = semantic
		return cmd, true
	}

	words[0] = semantic.Namespace
	if semantic.Namespace == "help" {
		semantic.Verb = "help"
		semantic.ReadOnly = true
		cmd.ActionPath = []string{"help"}
		cmd.TWG = semantic
		return cmd, true
	}
	if twgReadOnlyNamespaces[semantic.Namespace] {
		semantic.Verb = twgReadOnlyVerb(words)
		semantic.ReadOnly = true
		cmd.ActionPath = twgReadOnlyActionPath(words)
		cmd.TWG = semantic
		return cmd, true
	}

	path, class := longestTWGAction(words)
	if !complete && class.inferred {
		path = nil
		class = twgActionClass{}
	}
	cmd.ActionPath = path
	semantic.Verb = class.verb
	semantic.ReadOnly = class.readOnly
	semantic.Mutating = class.mutating
	cmd.TWG = semantic
	return cmd, true
}

func twgPathWords(words []string) ([]string, []Option, bool) {
	path := make([]string, 0, len(words))
	var options []Option
	for i := 0; i < len(words); i++ {
		word := words[i]
		if word == "--" {
			continue
		}
		if !strings.HasPrefix(word, "-") || word == "-" {
			path = append(path, word)
			continue
		}
		name, value, hasInlineValue := strings.Cut(word, "=")
		if twgGlobalFlagOptions[name] {
			options = append(options, Option{Name: name, Position: i})
			continue
		}
		if name == "--output-summary" {
			option := Option{Name: name, Position: i}
			if hasInlineValue {
				option.Value = value
				option.HasValue = true
			} else if i+1 < len(words) && twgOutputSummaryValues[words[i+1]] {
				i++
				option.Value = words[i]
				option.HasValue = true
			}
			options = append(options, option)
			continue
		}
		if !twgGlobalValueOptions[name] {
			return path, options, false
		}
		option := Option{Name: name, Position: i}
		if hasInlineValue {
			option.Value = value
			option.HasValue = true
		} else if i+1 < len(words) {
			i++
			option.Value = words[i]
			option.HasValue = true
		} else {
			options = append(options, option)
			return path, options, false
		}
		options = append(options, option)
	}
	return path, options, true
}

func longestTWGAction(words []string) ([]string, twgActionClass) {
	var path []string
	var class twgActionClass
	for i := 1; i <= len(words); i++ {
		candidate := strings.Join(words[:i], " ")
		if found, ok := twgActionClasses[candidate]; ok {
			path = append([]string(nil), words[:i]...)
			class = found
		}
	}
	return path, class
}

func twgReadOnlyVerb(words []string) string {
	namespace := words[0]
	if len(words) > 1 && twgReadOnlyChildVerbs[namespace][words[1]] {
		return words[1]
	}
	return namespace
}

func twgReadOnlyActionPath(words []string) []string {
	verb := twgReadOnlyVerb(words)
	if verb == words[0] {
		return []string{words[0]}
	}
	return []string{words[0], verb}
}

func twgHelpActionPath(words []string) []string {
	if len(words) == 0 || canonicalTWGNamespace(words[0]) == "help" {
		return []string{"help"}
	}
	return []string{canonicalTWGNamespace(words[0]), "help"}
}

func canonicalTWGNamespace(namespace string) string {
	switch namespace {
	case "bb":
		return "bitbucket"
	case "prs":
		return "pull-requests"
	case "whoami":
		return "user"
	case "pull-requests-tree":
		return "pr-tree"
	case "issue-tree":
		return "workitem-tree"
	case "status-rollup":
		return "work-tree"
	case "ownership":
		return "responsibility"
	case "visualise":
		return "visualize"
	default:
		return namespace
	}
}

func twgHasHelpFlag(words []string) bool {
	for _, word := range words {
		if word == "--help" || word == "-h" {
			return true
		}
	}
	return false
}

func twgHasVersionFlag(words []string) bool {
	for _, word := range words {
		if word == "--version" || word == "-v" || word == "-V" {
			return true
		}
	}
	return false
}

var twgGlobalValueOptions = map[string]bool{
	"-u": true, "--user": true,
	"-s": true, "--site": true,
	"-o": true, "--output": true,
	"--api-version": true, "--timeout-ms": true,
	"--agent-fields": true,
	"--output-shape": true, "--output-file": true,
	"--mode": true,
}

var twgGlobalFlagOptions = map[string]bool{
	"--help": true, "-h": true, "--version": true, "-v": true, "-V": true,
}

var twgOutputSummaryValues = stringSet("stats auto inline")

var twgReadOnlyNamespaces = stringSet(`
access collaborators commits context csm deployments docs doctor focus-areas
focus-areas-tree meetings notifications org-tree people pr-tree pull-requests
recently-viewed resolve responsibility search search-code spaces talent user
user-search videos visualize work work-tree workitem-tree
`)

var twgReadOnlyChildVerbs = map[string]map[string]bool{
	"context":        stringSet("confluence get jira user"),
	"csm":            stringSet("channel context organization"),
	"docs":           stringSet("get query search"),
	"focus-areas":    stringSet("get query search"),
	"meetings":       stringSet("get query"),
	"people":         stringSet("bulk-lookup describe search"),
	"pull-requests":  stringSet("query search"),
	"responsibility": stringSet("get infer"),
	"search":         stringSet("list-apps"),
	"search-code":    stringSet("diff dir file overview scan search symbol"),
	"spaces":         stringSet("get query"),
	"talent":         stringSet("position"),
	"user":           stringSet("bulk-lookup direct-reports manager search"),
	"videos":         stringSet("get query"),
	"work":           stringSet("query search"),
}

var twgActionClasses = buildTWGActionClasses([]twgActionGroup{
	{prefix: "jira board", read: "get query", write: "create delete"},
	{prefix: "jira board backlog", read: "query"},
	{prefix: "jira board backlog-view", read: "query"},
	{prefix: "jira board cells", read: "query"},
	{prefix: "jira board projects", read: "query"},
	{prefix: "jira board quick-filter", read: "get query"},
	{prefix: "jira board scope", read: "query"},
	{prefix: "jira board sprints", read: "query"},
	{prefix: "jira board view-settings", read: "query"},
	{prefix: "jira dashboard", read: "get query", write: "bulk-edit copy create delete update"},
	{prefix: "jira dashboard gadget", read: "query", write: "add delete update"},
	{prefix: "jira dashboard gadget catalog", read: "query"},
	{prefix: "jira dashboard item-property", read: "get query", write: "delete set"},
	{prefix: "jira field", write: "cancel-delete create delete update"},
	{prefix: "jira filter", read: "get query", write: "add-favourite change-owner create delete remove-favourite reset-columns update"},
	{prefix: "jira filter columns", read: "get", write: "set"},
	{prefix: "jira filter share", read: "query", write: "add remove"},
	{prefix: "jira filter subscription", read: "query", write: "add remove"},
	{prefix: "jira space", read: "get issue-types"},
	{prefix: "jira space component", read: "counts get query", write: "create delete update"},
	{prefix: "jira space notification-scheme", read: "get"},
	{prefix: "jira space status", read: "query"},
	{prefix: "jira space types", read: "query"},
	{prefix: "jira space version", read: "counts get", write: "create delete merge move update"},
	{prefix: "jira space versions", read: "query"},
	{prefix: "jira sprint", read: "get snapshot", write: "complete create delete start update"},
	{prefix: "jira sprint workitems", read: "query"},
	{prefix: "jira workitem", read: "bulk-get get query search", write: "archive bulk-transition clone create create-bulk delete transition unarchive update"},
	{prefix: "jira workitem attachment", read: "download get query thumbnail", write: "delete upload"},
	{prefix: "jira workitem changelog", read: "query"},
	{prefix: "jira workitem comment", read: "query", write: "create delete update"},
	{prefix: "jira workitem field", read: "create-metadata update-metadata"},
	{prefix: "jira workitem link", read: "query", write: "artifact branch build commit deployment goal loom meeting page pr project repo weblink workitem"},
	{prefix: "jira workitem link-types", read: "query"},
	{prefix: "jira workitem priorities", read: "query"},
	{prefix: "jira workitem project-link-candidates", read: "query"},
	{prefix: "jira workitem property", read: "get query", write: "delete set"},
	{prefix: "jira workitem statuses", read: "query"},
	{prefix: "jira workitem transitions", read: "query"},
	{prefix: "jira workitem types", read: "get query"},
	{prefix: "jira workitem unlink", write: "artifact branch build commit deployment goal loom meeting page pr project repo weblink workitem"},
	{prefix: "jira workitem vote", read: "query", write: "add remove"},
	{prefix: "jira workitem watcher", read: "query", write: "add remove"},
	{prefix: "jira workitem worklog", read: "changed deleted get query", write: "add delete update"},

	{prefix: "confluence", read: "search tree"},
	{prefix: "confluence content", read: "body-formats export export-status get get-public-link list", write: "archive copy create delete delete-draft disable-public-link enable-public-link move publish set-status unarchive update"},
	{prefix: "confluence content attachments", read: "download get list", write: "delete upload"},
	{prefix: "confluence content comments", read: "get list", write: "create delete reopen reply resolve update"},
	{prefix: "confluence content history", read: "diff get list"},
	{prefix: "confluence content labels", read: "list", write: "add remove"},
	{prefix: "confluence content permissions", read: "list", write: "add clear remove replace"},
	{prefix: "confluence content reactions", read: "list", write: "add remove"},
	{prefix: "confluence content restriction-state", read: "get", write: "set"},
	{prefix: "confluence content tasks", read: "get list", write: "complete reopen"},
	{prefix: "confluence content versions", read: "diff get list", write: "restore"},
	{prefix: "confluence search", read: "query"},
	{prefix: "confluence space", read: "get list me", write: "archive create delete unarchive update"},
	{prefix: "confluence space instructions", read: "get", write: "set"},
	{prefix: "confluence templates", read: "get list"},

	{prefix: "jsm alert", read: "get query", write: "create delete update"},
	{prefix: "jsm approval", read: "get list", write: "approve reject"},
	{prefix: "jsm automation resolution-plan", read: "get"},
	{prefix: "jsm conversation", read: "query"},
	{prefix: "jsm conversation claim", write: "create"},
	{prefix: "jsm conversation close", write: "create"},
	{prefix: "jsm conversation message", read: "query"},
	{prefix: "jsm conversation settings", read: "query", write: "update"},
	{prefix: "jsm conversation workspace", read: "query", write: "create"},
	{prefix: "jsm help-article", read: "query"},
	{prefix: "jsm help-center", read: "query", write: "create"},
	{prefix: "jsm help-object-store", read: "query", write: "create"},
	{prefix: "jsm incident", read: "get query", write: "create delete transition update"},
	{prefix: "jsm incident affected-service", read: "query"},
	{prefix: "jsm incident alert", read: "query"},
	{prefix: "jsm incident link", write: "affected-service alert"},
	{prefix: "jsm incident responder", read: "query", write: "add remove"},
	{prefix: "jsm incident unlink", write: "affected-service alert"},
	{prefix: "jsm knowledge-app-link", read: "query"},
	{prefix: "jsm knowledge-article", read: "query"},
	{prefix: "jsm knowledge-base", read: "query", write: "create"},
	{prefix: "jsm knowledge-base search", read: "query"},
	{prefix: "jsm knowledge-capability", read: "query"},
	{prefix: "jsm knowledge-discovery", read: "query", write: "create"},
	{prefix: "jsm knowledge-permission bulk", read: "query"},
	{prefix: "jsm knowledge-source-type", read: "get"},
	{prefix: "jsm linked-source", read: "query", write: "link unlink"},
	{prefix: "jsm linked-source permission", write: "update"},
	{prefix: "jsm linked-source suggestion", read: "query"},
	{prefix: "jsm linked-source view", write: "update"},
	{prefix: "jsm portal", read: "query"},
	{prefix: "jsm post-incident-review", read: "get query", write: "create delete update"},
	{prefix: "jsm post-incident-review incident", read: "get"},
	{prefix: "jsm post-incident-review link", write: "incident"},
	{prefix: "jsm post-incident-review unlink", write: "incident"},
	{prefix: "jsm request", write: "create"},
	{prefix: "jsm request-type", read: "query"},
	{prefix: "jsm request-type field", read: "get"},
	{prefix: "jsm resolution-state", read: "get"},
	{prefix: "jsm service", read: "get query search"},
	{prefix: "jsm service-tier", read: "query"},
	{prefix: "jsm sla", read: "get metrics workitems"},
	{prefix: "jsm support-site-article", read: "query"},

	{prefix: "assets", read: "object objects query search"},
	{prefix: "assets object", read: "get query", write: "create delete update"},
	{prefix: "assets object-attribute-value", write: "update"},
	{prefix: "assets objects", read: "query"},
	{prefix: "assets objectschema", read: "get list", write: "create delete update"},
	{prefix: "assets objectschema attributes", read: "query"},
	{prefix: "assets objectschema settings", write: "update"},
	{prefix: "assets reference-type", read: "query", write: "create delete update"},
	{prefix: "assets service-object", read: "query"},
	{prefix: "assets type", read: "get query", write: "create delete update"},
	{prefix: "assets type attribute", write: "create delete update"},
	{prefix: "assets type list-attr", read: "query"},

	{prefix: "bitbucket", read: "inbox search"},
	{prefix: "bitbucket branch", read: "query", write: "create delete"},
	{prefix: "bitbucket commit", read: "query"},
	{prefix: "bitbucket deployment", read: "query"},
	{prefix: "bitbucket deployment environment variable", read: "query", write: "create delete update"},
	{prefix: "bitbucket pipeline", read: "get grep latest-failure query tail wait", write: "rerun-failed run trigger"},
	{prefix: "bitbucket pull-requests", read: "activity commits diff diffstat for-commit get merge-status patch query", write: "approve create decline merge remove-request-changes request-changes unapprove update"},
	{prefix: "bitbucket pull-requests comment", read: "get query", write: "create delete reopen resolve update"},
	{prefix: "bitbucket pull-requests default-reviewer", read: "get query", write: "add remove"},
	{prefix: "bitbucket pull-requests effective-default-reviewer", read: "query"},
	{prefix: "bitbucket pull-requests task", read: "get query", write: "create delete reopen resolve update"},
	{prefix: "bitbucket repo", read: "contributors file get query url"},
	{prefix: "bitbucket search", read: "prs"},
	{prefix: "bitbucket workspace", read: "get"},
	{prefix: "bitbucket workspace member", read: "query"},

	{prefix: "loom", read: "get"},
	{prefix: "loom invite", write: "accept"},
	{prefix: "loom space", read: "query", write: "create"},
	{prefix: "loom video", read: "comments get", write: "delete"},
	{prefix: "loom workspace", write: "join"},

	{prefix: "trello", read: "search"},
	{prefix: "trello board", read: "get", write: "close reopen update"},
	{prefix: "trello board list", read: "query"},
	{prefix: "trello board member", read: "query"},
	{prefix: "trello card", read: "get", write: "add-label add-member archive create mark-complete remove-label remove-member unarchive update"},
	{prefix: "trello card label", read: "query"},
	{prefix: "trello card member", read: "query"},
	{prefix: "trello list", read: "get"},
	{prefix: "trello list card", read: "query", write: "sort"},
	{prefix: "trello member", read: "get me"},
	{prefix: "trello member workspace", read: "query"},
	{prefix: "trello workspace", read: "get"},
	{prefix: "trello workspace member", read: "query"},

	{prefix: "teams", read: "get query", write: "archive create update"},
	{prefix: "teams members", read: "list", write: "add remove"},
	{prefix: "goals", read: "get query types", write: "archive create update"},
	{prefix: "goals status-update", write: "create update"},
	{prefix: "projects", read: "get query", write: "archive create update"},
	{prefix: "projects status-update", write: "create update"},

	{prefix: "admin auth", read: "status"},
	{prefix: "admin directory", read: "list"},
	{prefix: "admin group", read: "count get list role-assignments stats", write: "create"},
	{prefix: "admin group access", write: "grant revoke"},
	{prefix: "admin group member", write: "add remove"},
	{prefix: "admin org", read: "get list"},
	{prefix: "admin user", read: "capabilities count get last-active list stats", write: "cancel-delete delete invite restore suspend"},

	{prefix: "rovo", read: "list-apps search"},
	{prefix: "upkeep", read: "status", write: "disable enable run"},
})

func buildTWGActionClasses(groups []twgActionGroup) map[string]twgActionClass {
	classes := map[string]twgActionClass{
		"jira workitem":            {verb: "get", readOnly: true, inferred: true},
		"jsm incident":             {verb: "query", readOnly: true, inferred: true},
		"jsm post-incident-review": {verb: "query", readOnly: true, inferred: true},
		"teams":                    {verb: "query", readOnly: true, inferred: true},
		"goals":                    {verb: "query", readOnly: true, inferred: true},
		"projects":                 {verb: "query", readOnly: true, inferred: true},
	}
	for _, group := range groups {
		for _, verb := range strings.Fields(group.read) {
			classes[group.prefix+" "+verb] = twgActionClass{verb: verb, readOnly: true}
		}
		for _, verb := range strings.Fields(group.write) {
			effectiveVerb := verb
			if strings.HasSuffix(group.prefix, " unlink") {
				effectiveVerb = "unlink"
			} else if strings.HasSuffix(group.prefix, " link") {
				effectiveVerb = "link"
			}
			classes[group.prefix+" "+verb] = twgActionClass{verb: effectiveVerb, mutating: true}
		}
	}
	classes["assets object"] = twgActionClass{verb: "object", readOnly: true, inferred: true}
	return classes
}

func stringSet(values string) map[string]bool {
	set := map[string]bool{}
	for _, value := range strings.Fields(values) {
		set[value] = true
	}
	return set
}
