package command

import (
	"strings"

	"github.com/tasuku43/cc-bash-proxy/internal/domain/invocation"
	"mvdan.cc/sh/v3/syntax"
)

type Command struct {
	Raw              string
	Program          string
	ProgramToken     string
	Env              map[string]string
	GlobalOptions    []string
	ActionPath       []string
	Options          []string
	Args             []string
	WorkingDirectory string
	Namespace        string
	ResourceType     string
	ResourceName     string
	Parser           string
}

type CommandPlan struct {
	Raw                    string
	Commands               []Command
	Shape                  ShellShape
	SafeForStructuredAllow bool
	Diagnostics            []Diagnostic
}

type ShellShape struct {
	Kind                   ShellShapeKind
	HasPipeline            bool
	HasConditional         bool
	HasSequence            bool
	HasBackground          bool
	HasRedirection         bool
	HasCommandSubstitution bool
}

type ShellShapeKind string

const (
	ShellShapeSimple     ShellShapeKind = "simple"
	ShellShapeAndList    ShellShapeKind = "and_list"
	ShellShapeOrList     ShellShapeKind = "or_list"
	ShellShapeSequence   ShellShapeKind = "sequence"
	ShellShapePipeline   ShellShapeKind = "pipeline"
	ShellShapeBackground ShellShapeKind = "background"
	ShellShapeRedirect   ShellShapeKind = "redirect"
	ShellShapeSubshell   ShellShapeKind = "subshell"
	ShellShapeUnknown    ShellShapeKind = "unknown"
)

type Diagnostic struct {
	Severity string
	Message  string
}

func Parse(raw string) CommandPlan {
	plan := CommandPlan{
		Raw:   raw,
		Shape: ShellShape{Kind: ShellShapeUnknown},
	}

	if strings.TrimSpace(raw) == "" {
		plan.Diagnostics = append(plan.Diagnostics, Diagnostic{Severity: "error", Message: "empty command"})
		return plan
	}

	parser := syntax.NewParser()
	file, err := parser.Parse(strings.NewReader(raw), "")
	if err != nil {
		plan.Diagnostics = append(plan.Diagnostics, Diagnostic{Severity: "error", Message: err.Error()})
		return plan
	}

	walker := planWalker{raw: raw}
	if len(file.Stmts) > 1 || len(file.Last) > 0 {
		walker.shape.HasSequence = true
	}
	for _, stmt := range file.Stmts {
		walker.visitStmt(stmt)
	}

	plan.Commands = walker.commands
	plan.Shape = walker.shape.finalize()
	plan.SafeForStructuredAllow = plan.Shape.Kind == ShellShapeSimple &&
		len(plan.Commands) == 1 &&
		len(plan.Diagnostics) == 0 &&
		invocation.IsStructuredSafeForAllow(raw)
	return plan
}

type planWalker struct {
	raw      string
	shape    ShellShape
	commands []Command
}

func (w *planWalker) visitStmt(stmt *syntax.Stmt) {
	if stmt == nil {
		return
	}
	if stmt.Background || stmt.Coprocess || stmt.Disown {
		w.shape.HasBackground = true
	}
	if len(stmt.Redirs) > 0 {
		w.shape.HasRedirection = true
	}
	if len(stmt.Comments) > 0 || stmt.Negated {
		w.shape.Kind = ShellShapeUnknown
	}
	w.visitCommand(stmt.Cmd)
	for _, redir := range stmt.Redirs {
		w.visitWord(redir.Word)
		w.visitWord(redir.Hdoc)
	}
}

func (w *planWalker) visitCommand(cmd syntax.Command) {
	switch x := cmd.(type) {
	case nil:
		return
	case *syntax.CallExpr:
		w.visitCall(x)
	case *syntax.BinaryCmd:
		w.visitBinary(x)
	case *syntax.Subshell:
		w.shape.Kind = ShellShapeSubshell
		for _, stmt := range x.Stmts {
			w.visitStmt(stmt)
		}
	case *syntax.Block:
		w.shape.Kind = ShellShapeUnknown
		for _, stmt := range x.Stmts {
			w.visitStmt(stmt)
		}
	default:
		w.shape.Kind = ShellShapeUnknown
	}
}

func (w *planWalker) visitBinary(cmd *syntax.BinaryCmd) {
	switch cmd.Op {
	case syntax.AndStmt:
		w.shape.HasConditional = true
		if w.shape.Kind == "" {
			w.shape.Kind = ShellShapeAndList
		}
	case syntax.OrStmt:
		w.shape.HasConditional = true
		w.shape.Kind = ShellShapeOrList
	case syntax.Pipe, syntax.PipeAll:
		w.shape.HasPipeline = true
		w.shape.Kind = ShellShapePipeline
	default:
		w.shape.Kind = ShellShapeUnknown
	}
	w.visitStmt(cmd.X)
	w.visitStmt(cmd.Y)
}

func (w *planWalker) visitCall(call *syntax.CallExpr) {
	for _, assign := range call.Assigns {
		if assign != nil {
			w.visitWord(assign.Value)
		}
	}
	for _, arg := range call.Args {
		w.visitWord(arg)
	}
	if len(call.Args) == 0 {
		return
	}

	raw := w.nodeRaw(call)
	parsed := invocation.Parse(raw)
	cmd := Command{
		Raw:          raw,
		Program:      parsed.Command,
		ProgramToken: parsed.CommandToken,
		Env:          parsed.EnvAssignments,
		Args:         parsed.Args,
		Parser:       parserName(parsed.Command),
	}
	cmd.GlobalOptions, cmd.ActionPath, cmd.Options = splitSemanticArgs(parsed.Command, parsed.Args)
	cmd.WorkingDirectory = workingDirectory(parsed.Command, parsed.Args)
	cmd.Namespace, cmd.ResourceType, cmd.ResourceName = kubectlResource(parsed.Command, parsed.Args)
	w.commands = append(w.commands, cmd)
}

func (w *planWalker) visitWord(word *syntax.Word) {
	if word == nil {
		return
	}
	for _, part := range word.Parts {
		w.visitWordPart(part)
	}
}

func (w *planWalker) visitWordPart(part syntax.WordPart) {
	switch x := part.(type) {
	case *syntax.CmdSubst:
		w.shape.HasCommandSubstitution = true
		for _, stmt := range x.Stmts {
			w.visitStmt(stmt)
		}
	case *syntax.DblQuoted:
		for _, nested := range x.Parts {
			w.visitWordPart(nested)
		}
	case *syntax.ParamExp, *syntax.ArithmExp, *syntax.ProcSubst, *syntax.ExtGlob, *syntax.BraceExp:
		w.shape.Kind = ShellShapeUnknown
	}
}

func (w *planWalker) nodeRaw(node syntax.Node) string {
	if node == nil || !node.Pos().IsValid() || !node.End().IsValid() {
		return ""
	}
	start := int(node.Pos().Offset())
	end := int(node.End().Offset())
	if start < 0 || end < start || end > len(w.raw) {
		return ""
	}
	return strings.TrimSpace(w.raw[start:end])
}

func (s ShellShape) finalize() ShellShape {
	if s.HasPipeline {
		s.Kind = ShellShapePipeline
		return s
	}
	if s.HasConditional {
		if s.Kind != ShellShapeOrList {
			s.Kind = ShellShapeAndList
		}
		return s
	}
	if s.HasSequence {
		s.Kind = ShellShapeSequence
		return s
	}
	if s.HasBackground {
		s.Kind = ShellShapeBackground
		return s
	}
	if s.HasRedirection {
		s.Kind = ShellShapeRedirect
		return s
	}
	if s.HasCommandSubstitution {
		s.Kind = ShellShapeUnknown
		return s
	}
	if s.Kind == "" || s.Kind == ShellShapeUnknown {
		s.Kind = ShellShapeSimple
	}
	return s
}

func parserName(program string) string {
	switch program {
	case "aws", "git", "kubectl":
		return program
	default:
		return "generic"
	}
}

func splitSemanticArgs(program string, args []string) ([]string, []string, []string) {
	globalOptions := []string{}
	actionPath := []string{}
	options := []string{}
	seenAction := false
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") && arg != "-" {
			if seenAction {
				options = append(options, arg)
			} else {
				globalOptions = append(globalOptions, arg)
			}
			if optionConsumesValue(program, arg) && i+1 < len(args) {
				if seenAction {
					options = append(options, args[i+1])
				} else {
					globalOptions = append(globalOptions, args[i+1])
				}
				i++
			}
			continue
		}
		seenAction = true
		actionPath = append(actionPath, arg)
	}
	return globalOptions, actionPath, options
}

func optionConsumesValue(program string, option string) bool {
	if strings.Contains(option, "=") {
		return false
	}
	switch program {
	case "git":
		switch option {
		case "-C", "-c", "--git-dir", "--work-tree", "--namespace":
			return true
		}
	case "aws":
		switch option {
		case "--profile", "--region", "--endpoint-url", "--output":
			return true
		}
	case "kubectl":
		switch option {
		case "-n", "--namespace", "--context", "--kubeconfig":
			return true
		}
	}
	return false
}

func workingDirectory(program string, args []string) string {
	if program != "git" {
		return ""
	}
	for i, arg := range args {
		if arg == "-C" && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func kubectlResource(program string, args []string) (string, string, string) {
	if program != "kubectl" {
		return "", "", ""
	}
	namespace := ""
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-n", "--namespace":
			if i+1 < len(args) {
				namespace = args[i+1]
				i++
			}
		default:
			if strings.HasPrefix(args[i], "--namespace=") {
				namespace = strings.TrimPrefix(args[i], "--namespace=")
			}
		}
	}
	for i, arg := range args {
		if arg == "get" || arg == "describe" || arg == "delete" || arg == "logs" {
			resourceType := ""
			resourceName := ""
			if i+1 < len(args) {
				resourceType = args[i+1]
			}
			if i+2 < len(args) {
				resourceName = args[i+2]
			}
			return namespace, resourceType, resourceName
		}
	}
	return namespace, "", ""
}
