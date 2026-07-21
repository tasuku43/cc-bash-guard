package command

import "strings"

type PupParser struct{}

//go:generate go run ../../devtools/gen-pup-actions

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
	endOfOptions := false
	for i := 0; i < len(base.RawWords); i++ {
		word := base.RawWords[i]
		switch {
		case endOfOptions:
			positionals = append(positionals, word)
		case word == "--":
			endOfOptions = true
		case pupOptionWithValue(word, "", "--org"):
			value, consumed := pupOptionValue(word, "", "--org", base.RawWords, i)
			cmd.Options = append(cmd.Options, Option{Name: "--org", Value: value, HasValue: consumed, Position: i})
			if consumed && !pupOptionHasInlineValue(word, "", "--org") {
				i++
			}
		case pupOptionWithValue(word, "-o", "--output"):
			value, consumed := pupOptionValue(word, "-o", "--output", base.RawWords, i)
			name := "--output"
			if strings.HasPrefix(word, "-o") && !strings.HasPrefix(word, "--") {
				name = "-o"
			}
			cmd.Options = append(cmd.Options, Option{Name: name, Value: value, HasValue: consumed, Position: i})
			if consumed && !pupOptionHasInlineValue(word, "-o", "--output") {
				i++
			}
		case pupOptionWithValue(word, "", "--jq"):
			value, consumed := pupOptionValue(word, "", "--jq", base.RawWords, i)
			cmd.Options = append(cmd.Options, Option{Name: "--jq", Value: value, HasValue: consumed, Position: i})
			if consumed && !pupOptionHasInlineValue(word, "", "--jq") {
				i++
			}
		case word == "-y" || word == "--yes" || word == "--agent" || word == "--no-agent" ||
			word == "--read-only" || word == "--trust-site":
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
	if short != "" && (word == short || strings.HasPrefix(word, short+"=") ||
		(strings.HasPrefix(word, short) && !strings.HasPrefix(word, "--") && len(word) > len(short))) {
		return true
	}
	return word == long || strings.HasPrefix(word, long+"=")
}

func pupOptionValue(word, short, long string, words []string, i int) (string, bool) {
	if value, ok := strings.CutPrefix(word, long+"="); ok {
		return value, true
	}
	if short != "" {
		if value, ok := strings.CutPrefix(word, short+"="); ok {
			return value, true
		}
		if strings.HasPrefix(word, short) && !strings.HasPrefix(word, "--") && len(word) > len(short) {
			return strings.TrimPrefix(word, short), true
		}
	}
	if i+1 >= len(words) {
		return "", false
	}
	return words[i+1], true
}

func pupOptionHasInlineValue(word, short, long string) bool {
	if strings.HasPrefix(word, long+"=") {
		return true
	}
	return short != "" && strings.HasPrefix(word, short) && word != short && !strings.HasPrefix(word, "--")
}

func pupActionPathAndArgs(positionals []string) ([]string, []string) {
	if len(positionals) == 0 {
		return []string{}, []string{}
	}
	end := len(positionals)
	if end > pupMaxActionPathWords {
		end = pupMaxActionPathWords
	}
	for ; end > 0; end-- {
		if pupKnownActionPaths[strings.Join(positionals[:end], " ")] {
			return append([]string(nil), positionals[:end]...), append([]string(nil), positionals[end:]...)
		}
	}
	// Unknown command shapes retain only the top-level area. In particular, do
	// not guess that an argument or a future nested area is a leaf verb: doing
	// so could widen verb-based allow rules after pup adds commands.
	return append([]string(nil), positionals[:1]...), append([]string(nil), positionals[1:]...)
}

func pupActionPathSet(spec string) map[string]bool {
	set := make(map[string]bool)
	for _, path := range strings.Split(strings.TrimSpace(spec), "\n") {
		if path = strings.TrimSpace(path); path != "" {
			set[path] = true
		}
	}
	return set
}

func buildPupSemantic(actionPath []string, options []Option) *PupSemantic {
	semantic := &PupSemantic{Flags: normalizedPupFlags(options)}
	if len(actionPath) > 0 {
		semantic.Area = actionPath[0]
	}
	if len(actionPath) > 2 {
		semantic.SubArea = actionPath[1]
	}
	if len(actionPath) > 1 {
		semantic.Verb = actionPath[len(actionPath)-1]
	}
	semantic.Org = pupLastOptionValue(options, "--org")
	semantic.Output = pupLastOptionValue(options, "-o", "--output")
	semantic.Yes = hasOption(options, "--yes") || hasOption(options, "-y")
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
