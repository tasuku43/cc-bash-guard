package policy

import commandpkg "github.com/tasuku43/cc-bash-guard/internal/domain/command"

func init() {
	registerSemanticHandler(semanticHandler{
		command:  "twg",
		match:    func(s SemanticMatchSpec, cmd commandpkg.Command) bool { return s.TWG().matches(cmd) },
		validate: ValidateTWGSemanticMatchSpec,
	})
}

func (s TWGSemanticSpec) matches(cmd commandpkg.Command) bool {
	if cmd.SemanticParser != "twg" || cmd.TWG == nil {
		return false
	}
	t := cmd.TWG
	if s.Namespace != "" && t.Namespace != s.Namespace {
		return false
	}
	if len(s.NamespaceIn) > 0 && !containsString(s.NamespaceIn, t.Namespace) {
		return false
	}
	if s.Verb != "" && t.Verb != s.Verb {
		return false
	}
	if len(s.VerbIn) > 0 && !containsString(s.VerbIn, t.Verb) {
		return false
	}
	if s.ReadOnly != nil && t.ReadOnly != *s.ReadOnly {
		return false
	}
	if s.Mutating != nil && t.Mutating != *s.Mutating {
		return false
	}
	return true
}
