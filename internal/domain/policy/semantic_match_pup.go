package policy

import commandpkg "github.com/tasuku43/cc-bash-guard/internal/domain/command"

func init() {
	registerSemanticHandler(semanticHandler{
		command:  "pup",
		match:    func(s SemanticMatchSpec, cmd commandpkg.Command) bool { return s.Pup().matches(cmd) },
		validate: ValidatePupSemanticMatchSpec,
	})
}

func (s PupSemanticSpec) matches(cmd commandpkg.Command) bool {
	if cmd.SemanticParser != "pup" || cmd.Pup == nil {
		return false
	}
	pup := cmd.Pup
	if s.Area != "" && pup.Area != s.Area {
		return false
	}
	if len(s.AreaIn) > 0 && !containsString(s.AreaIn, pup.Area) {
		return false
	}
	if s.SubArea != "" && pup.SubArea != s.SubArea {
		return false
	}
	if len(s.SubAreaIn) > 0 && !containsString(s.SubAreaIn, pup.SubArea) {
		return false
	}
	if s.Verb != "" && pup.Verb != s.Verb {
		return false
	}
	if len(s.VerbIn) > 0 && !containsString(s.VerbIn, pup.Verb) {
		return false
	}
	if s.Org != "" && pup.Org != s.Org {
		return false
	}
	if len(s.OrgIn) > 0 && !containsString(s.OrgIn, pup.Org) {
		return false
	}
	if s.Output != "" && pup.Output != s.Output {
		return false
	}
	if len(s.OutputIn) > 0 && !containsString(s.OutputIn, pup.Output) {
		return false
	}
	if s.Yes != nil && pup.Yes != *s.Yes {
		return false
	}
	if s.Agent != nil && pup.Agent != *s.Agent {
		return false
	}
	if s.NoAgent != nil && pup.NoAgent != *s.NoAgent {
		return false
	}
	for _, flag := range s.FlagsContains {
		if !containsString(pup.Flags, flag) {
			return false
		}
	}
	for _, prefix := range s.FlagsPrefixes {
		if !containsPrefix(pup.Flags, prefix) {
			return false
		}
	}
	return true
}
