package command

import "strings"

type PupParser struct{}

func init() {
	RegisterDefaultParser(PupParser{})
}

func (PupParser) Program() string { return "pup" }

func (PupParser) Parse(base Command) (Command, bool) {
	if base.Program != "pup" {
		return Command{}, false
	}
	cmd := base
	cmd.Parser = "pup"
	cmd.SemanticParser = "pup"
	cmd.Args = []string{}

	var positionals []string
	for i := 0; i < len(base.RawWords); i++ {
		word := base.RawWords[i]
		switch {
		case pupOptionWithValue(word, "", "--org"):
			value, consumed := pupOptionValue(word, "--org", base.RawWords, i)
			cmd.Options = append(cmd.Options, Option{Name: "--org", Value: value, HasValue: consumed, Position: i})
			if consumed && !strings.Contains(word, "=") {
				i++
			}
		case pupOptionWithValue(word, "-o", "--output"):
			value, consumed := pupOptionValue(word, "--output", base.RawWords, i)
			name := "--output"
			if word == "-o" {
				name = "-o"
			}
			cmd.Options = append(cmd.Options, Option{Name: name, Value: value, HasValue: consumed, Position: i})
			if consumed && !strings.Contains(word, "=") {
				i++
			}
		case word == "--yes" || word == "--agent" || word == "--no-agent":
			cmd.Options = append(cmd.Options, Option{Name: word, Position: i})
		case strings.HasPrefix(word, "-") && word != "-":
			cmd.Options = append(cmd.Options, parseOptionWord(word, i))
		default:
			positionals = append(positionals, word)
		}
	}

	cmd.ActionPath, cmd.Args = pupActionPathAndArgs(positionals)
	cmd.Pup = buildPupSemantic(cmd.ActionPath, cmd.Options)
	return cmd, true
}

func pupOptionWithValue(word, short, long string) bool {
	if short != "" && word == short {
		return true
	}
	return word == long || strings.HasPrefix(word, long+"=")
}

func pupOptionValue(word, long string, words []string, i int) (string, bool) {
	if value, ok := strings.CutPrefix(word, long+"="); ok {
		return value, true
	}
	if i+1 >= len(words) {
		return "", false
	}
	return words[i+1], true
}

func pupActionPathAndArgs(positionals []string) ([]string, []string) {
	if len(positionals) == 0 {
		return []string{}, []string{}
	}
	if len(positionals) >= 3 {
		return append([]string(nil), positionals[:3]...), append([]string(nil), positionals[3:]...)
	}
	return append([]string(nil), positionals...), []string{}
}

func buildPupSemantic(actionPath []string, options []Option) *PupSemantic {
	semantic := &PupSemantic{Flags: normalizedPupFlags(options)}
	if len(actionPath) > 0 {
		semantic.Area = actionPath[0]
	}
	if len(actionPath) > 1 {
		semantic.Verb = actionPath[1]
	}
	if len(actionPath) > 2 {
		semantic.SubArea = actionPath[1]
		semantic.Verb = actionPath[2]
	}
	semantic.Org = pupLastOptionValue(options, "--org")
	semantic.Output = pupLastOptionValue(options, "-o", "--output")
	semantic.Yes = hasOption(options, "--yes")
	semantic.Agent = hasOption(options, "--agent")
	semantic.NoAgent = hasOption(options, "--no-agent")
	return semantic
}

func normalizedPupFlags(options []Option) []string {
	out := make([]string, 0, len(options))
	for _, opt := range options {
		if opt.Name == "" {
			continue
		}
		out = append(out, opt.Name)
	}
	return out
}

func pupLastOptionValue(options []Option, names ...string) string {
	for i := len(options) - 1; i >= 0; i-- {
		for _, n := range names {
			if options[i].Name == n && options[i].HasValue {
				return options[i].Value
			}
		}
	}
	return ""
}
