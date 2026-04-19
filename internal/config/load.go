package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tasuku43/cmdproxy/internal/domain/policy"
	"gopkg.in/yaml.v3"
)

const LayerUser = "user"

type File struct {
	Version int               `yaml:"version"`
	Rules   []policy.RuleSpec `yaml:"rules"`
}

type Source = policy.Source

type Loaded struct {
	Rules  []policy.Rule
	Files  []Source
	Errors []error
}

type evalFile struct {
	Version int
	Rules   []evalRuleSpec
}

type evalRuleSpec struct {
	ID              string
	Pattern         string
	Match           policy.MatchSpec
	Message         string
	Reject          policy.RejectSpec
	Rewrite         policy.RewriteSpec
	BlockExampleLen int
	AllowExampleLen int
}

type evalCacheFile struct {
	Version       int              `json:"version"`
	SourcePath    string           `json:"source_path"`
	SourceSize    int64            `json:"source_size"`
	SourceModTime int64            `json:"source_mod_time"`
	CompiledRules []evalCachedRule `json:"compiled_rules"`
}

type evalCachedRule struct {
	ID      string             `json:"id"`
	Pattern string             `json:"pattern"`
	Match   policy.MatchSpec   `json:"match,omitempty"`
	Message string             `json:"message"`
	Reject  policy.RejectSpec  `json:"reject,omitempty"`
	Rewrite policy.RewriteSpec `json:"rewrite,omitempty"`
}

func ConfigPaths(home string, xdgConfigHome string) []Source {
	userConfigBase := xdgConfigHome
	if userConfigBase == "" {
		userConfigBase = filepath.Join(home, ".config")
	}
	return []Source{{
		Layer: LayerUser,
		Path:  filepath.Join(userConfigBase, "cmdproxy", "cmdproxy.yml"),
	}}
}

func CachePath(home string, xdgCacheHome string) string {
	cacheBase := xdgCacheHome
	if cacheBase == "" {
		cacheBase = filepath.Join(home, ".cache")
	}
	return filepath.Join(cacheBase, "cmdproxy", "eval-cache-v1.json")
}

func LoadEffective(home string, xdgConfigHome string) Loaded {
	return loadEffectiveWithLoader(home, xdgConfigHome, LoadFileIfPresent)
}

func LoadEffectiveForEval(home string, xdgConfigHome string, xdgCacheHome string) Loaded {
	loader := func(src Source) ([]policy.Rule, error) {
		return LoadFileForEvalIfPresent(src, CachePath(home, xdgCacheHome))
	}
	return loadEffectiveWithLoader(home, xdgConfigHome, loader)
}

func loadEffectiveWithLoader(home string, xdgConfigHome string, loader func(Source) ([]policy.Rule, error)) Loaded {
	var loaded Loaded
	for _, src := range ConfigPaths(home, xdgConfigHome) {
		rules, err := loader(src)
		if err != nil {
			loaded.Errors = append(loaded.Errors, err)
			continue
		}
		if len(rules) == 0 {
			continue
		}
		loaded.Files = append(loaded.Files, src)
		loaded.Rules = append(loaded.Rules, rules...)
	}
	loaded.Errors = append(loaded.Errors, policy.ValidateDuplicateIDs(loaded.Rules)...)
	return loaded
}

func LoadFileIfPresent(src Source) ([]policy.Rule, error) {
	data, err := readConfigFile(src)
	if err != nil {
		return nil, err
	}
	if data == "" {
		return nil, nil
	}
	file, err := decodeFullFile(src, data)
	if err != nil {
		return nil, err
	}
	issues := validateFile(file)
	if len(issues) > 0 {
		for i := range issues {
			issues[i] = fmt.Sprintf("%s config %s: %s", src.Layer, src.Path, issues[i])
		}
		return nil, &policy.ValidationError{Issues: issues}
	}
	rules := make([]policy.Rule, 0, len(file.Rules))
	for _, spec := range file.Rules {
		rules = append(rules, policy.NewRule(spec, src))
	}
	return rules, nil
}

func LoadFileForEvalIfPresent(src Source, cachePath string) ([]policy.Rule, error) {
	info, err := os.Stat(src.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s config read failed: %w", src.Layer, err)
	}
	if rules, ok := loadEvalCache(src, cachePath, info); ok {
		return rules, nil
	}
	data, err := readConfigFile(src)
	if err != nil {
		return nil, err
	}
	if data == "" {
		return nil, nil
	}
	file, err := decodeEvalFile(src, data)
	if err != nil {
		return nil, err
	}
	issues := validateEvalFile(file)
	if len(issues) > 0 {
		for i := range issues {
			issues[i] = fmt.Sprintf("%s config %s: %s", src.Layer, src.Path, issues[i])
		}
		return nil, &policy.ValidationError{Issues: issues}
	}
	rules := make([]policy.Rule, 0, len(file.Rules))
	cached := make([]evalCachedRule, 0, len(file.Rules))
	for _, spec := range file.Rules {
		ruleSpec := policy.RuleSpec{
			ID:      spec.ID,
			Pattern: spec.Pattern,
			Matcher: spec.Match,
			Message: spec.Message,
			Reject:  spec.Reject,
			Rewrite: spec.Rewrite,
		}
		rules = append(rules, policy.NewRule(ruleSpec, src))
		cached = append(cached, evalCachedRule{
			ID:      spec.ID,
			Pattern: spec.Pattern,
			Match:   spec.Match,
			Message: spec.Message,
			Reject:  spec.Reject,
			Rewrite: spec.Rewrite,
		})
	}
	writeEvalCache(cachePath, evalCacheFile{
		Version:       1,
		SourcePath:    src.Path,
		SourceSize:    info.Size(),
		SourceModTime: info.ModTime().UnixNano(),
		CompiledRules: cached,
	})
	return rules, nil
}

func readConfigFile(src Source) (string, error) {
	data, err := os.ReadFile(src.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", fmt.Errorf("%s config read failed: %w", src.Layer, err)
	}
	if strings.TrimSpace(string(data)) == "" {
		return "", fmt.Errorf("%s config %s is empty", src.Layer, src.Path)
	}
	return string(data), nil
}

func decodeFullFile(src Source, data string) (File, error) {
	dec := yaml.NewDecoder(strings.NewReader(data))
	dec.KnownFields(true)
	var file File
	if err := dec.Decode(&file); err != nil {
		return File{}, fmt.Errorf("%s config %s is invalid: %w", src.Layer, src.Path, err)
	}
	return file, nil
}

func decodeEvalFile(src Source, data string) (evalFile, error) {
	var root yaml.Node
	if err := yaml.Unmarshal([]byte(data), &root); err != nil {
		return evalFile{}, fmt.Errorf("%s config %s is invalid: %w", src.Layer, src.Path, err)
	}
	if len(root.Content) == 0 {
		return evalFile{}, fmt.Errorf("%s config %s is invalid: empty YAML document", src.Layer, src.Path)
	}
	doc := root.Content[0]
	if doc.Kind != yaml.MappingNode {
		return evalFile{}, fmt.Errorf("%s config %s is invalid: top-level must be a mapping", src.Layer, src.Path)
	}
	file := evalFile{}
	seenTopLevel := map[string]struct{}{}
	for i := 0; i < len(doc.Content); i += 2 {
		key := doc.Content[i]
		val := doc.Content[i+1]
		if _, ok := seenTopLevel[key.Value]; ok {
			continue
		}
		seenTopLevel[key.Value] = struct{}{}
		switch key.Value {
		case "version":
			if val.Kind != yaml.ScalarNode {
				return evalFile{}, fmt.Errorf("%s config %s is invalid: version must be a scalar", src.Layer, src.Path)
			}
			var version int
			if err := val.Decode(&version); err != nil {
				return evalFile{}, fmt.Errorf("%s config %s is invalid: version must be an integer", src.Layer, src.Path)
			}
			file.Version = version
		case "rules":
			rules, err := decodeEvalRules(src, val)
			if err != nil {
				return evalFile{}, err
			}
			file.Rules = rules
		default:
			return evalFile{}, fmt.Errorf("%s config %s is invalid: field %q not allowed", src.Layer, src.Path, key.Value)
		}
	}
	return file, nil
}

func decodeEvalRules(src Source, node *yaml.Node) ([]evalRuleSpec, error) {
	if node.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("%s config %s is invalid: rules must be a sequence", src.Layer, src.Path)
	}
	rules := make([]evalRuleSpec, 0, len(node.Content))
	for idx, item := range node.Content {
		if item.Kind != yaml.MappingNode {
			return nil, fmt.Errorf("%s config %s is invalid: rules[%d] must be a mapping", src.Layer, src.Path, idx)
		}
		ruleSpec, err := decodeEvalRule(src, idx, item)
		if err != nil {
			return nil, err
		}
		rules = append(rules, ruleSpec)
	}
	return rules, nil
}

func decodeEvalRule(src Source, idx int, node *yaml.Node) (evalRuleSpec, error) {
	var spec evalRuleSpec
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		val := node.Content[i+1]
		switch key.Value {
		case "id":
			if val.Kind != yaml.ScalarNode {
				return evalRuleSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].id must be a string", src.Layer, src.Path, idx)
			}
			spec.ID = val.Value
		case "pattern":
			if val.Kind != yaml.ScalarNode {
				return evalRuleSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].pattern must be a string", src.Layer, src.Path, idx)
			}
			spec.Pattern = val.Value
		case "match":
			match, err := decodeEvalMatch(src, idx, val)
			if err != nil {
				return evalRuleSpec{}, err
			}
			spec.Match = match
		case "message":
			if val.Kind != yaml.ScalarNode {
				return evalRuleSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].message must be a string", src.Layer, src.Path, idx)
			}
			spec.Message = val.Value
		case "reject":
			reject, err := decodeEvalReject(src, idx, val)
			if err != nil {
				return evalRuleSpec{}, err
			}
			spec.Reject = reject
		case "rewrite":
			rewrite, err := decodeEvalRewrite(src, idx, val)
			if err != nil {
				return evalRuleSpec{}, err
			}
			spec.Rewrite = rewrite
		case "block_examples":
			if val.Kind != yaml.SequenceNode {
				return evalRuleSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].block_examples must be a sequence", src.Layer, src.Path, idx)
			}
			spec.BlockExampleLen = len(val.Content)
		case "allow_examples":
			if val.Kind != yaml.SequenceNode {
				return evalRuleSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].allow_examples must be a sequence", src.Layer, src.Path, idx)
			}
			spec.AllowExampleLen = len(val.Content)
		default:
			return evalRuleSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].%s not allowed", src.Layer, src.Path, idx, key.Value)
		}
	}
	return spec, nil
}

func decodeEvalRewrite(src Source, idx int, node *yaml.Node) (policy.RewriteSpec, error) {
	if node.Kind != yaml.MappingNode {
		return policy.RewriteSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].rewrite must be a mapping", src.Layer, src.Path, idx)
	}
	var rewrite policy.RewriteSpec
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		val := node.Content[i+1]
		switch key.Value {
		case "unwrap_shell_dash_c":
			if val.Kind != yaml.ScalarNode {
				return policy.RewriteSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].rewrite.unwrap_shell_dash_c must be a boolean", src.Layer, src.Path, idx)
			}
			var enabled bool
			if err := val.Decode(&enabled); err != nil {
				return policy.RewriteSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].rewrite.unwrap_shell_dash_c must be a boolean", src.Layer, src.Path, idx)
			}
			rewrite.UnwrapShellDashC = enabled
		case "move_flag_to_env":
			spec, err := decodeEvalMoveFlagToEnv(src, idx, val)
			if err != nil {
				return policy.RewriteSpec{}, err
			}
			rewrite.MoveFlagToEnv = spec
		case "move_env_to_flag":
			spec, err := decodeEvalMoveEnvToFlag(src, idx, val)
			if err != nil {
				return policy.RewriteSpec{}, err
			}
			rewrite.MoveEnvToFlag = spec
		case "unwrap_wrapper":
			spec, err := decodeEvalUnwrapWrapper(src, idx, val)
			if err != nil {
				return policy.RewriteSpec{}, err
			}
			rewrite.UnwrapWrapper = spec
		default:
			return policy.RewriteSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].rewrite.%s not allowed", src.Layer, src.Path, idx, key.Value)
		}
	}
	return rewrite, nil
}

func decodeEvalMoveFlagToEnv(src Source, idx int, node *yaml.Node) (policy.MoveFlagToEnvSpec, error) {
	if node.Kind != yaml.MappingNode {
		return policy.MoveFlagToEnvSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].rewrite.move_flag_to_env must be a mapping", src.Layer, src.Path, idx)
	}
	var spec policy.MoveFlagToEnvSpec
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		val := node.Content[i+1]
		switch key.Value {
		case "flag":
			if val.Kind != yaml.ScalarNode {
				return policy.MoveFlagToEnvSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].rewrite.move_flag_to_env.flag must be a string", src.Layer, src.Path, idx)
			}
			spec.Flag = val.Value
		case "env":
			if val.Kind != yaml.ScalarNode {
				return policy.MoveFlagToEnvSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].rewrite.move_flag_to_env.env must be a string", src.Layer, src.Path, idx)
			}
			spec.Env = val.Value
		default:
			return policy.MoveFlagToEnvSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].rewrite.move_flag_to_env.%s not allowed", src.Layer, src.Path, idx, key.Value)
		}
	}
	return spec, nil
}

func decodeEvalMoveEnvToFlag(src Source, idx int, node *yaml.Node) (policy.MoveEnvToFlagSpec, error) {
	if node.Kind != yaml.MappingNode {
		return policy.MoveEnvToFlagSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].rewrite.move_env_to_flag must be a mapping", src.Layer, src.Path, idx)
	}
	var spec policy.MoveEnvToFlagSpec
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		val := node.Content[i+1]
		switch key.Value {
		case "env":
			if val.Kind != yaml.ScalarNode {
				return policy.MoveEnvToFlagSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].rewrite.move_env_to_flag.env must be a string", src.Layer, src.Path, idx)
			}
			spec.Env = val.Value
		case "flag":
			if val.Kind != yaml.ScalarNode {
				return policy.MoveEnvToFlagSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].rewrite.move_env_to_flag.flag must be a string", src.Layer, src.Path, idx)
			}
			spec.Flag = val.Value
		default:
			return policy.MoveEnvToFlagSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].rewrite.move_env_to_flag.%s not allowed", src.Layer, src.Path, idx, key.Value)
		}
	}
	return spec, nil
}

func decodeEvalUnwrapWrapper(src Source, idx int, node *yaml.Node) (policy.UnwrapWrapperSpec, error) {
	if node.Kind != yaml.MappingNode {
		return policy.UnwrapWrapperSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].rewrite.unwrap_wrapper must be a mapping", src.Layer, src.Path, idx)
	}
	var spec policy.UnwrapWrapperSpec
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		val := node.Content[i+1]
		switch key.Value {
		case "wrappers":
			values, err := decodeStringSequence(src, idx, "rewrite.unwrap_wrapper.wrappers", val)
			if err != nil {
				return policy.UnwrapWrapperSpec{}, err
			}
			spec.Wrappers = values
		default:
			return policy.UnwrapWrapperSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].rewrite.unwrap_wrapper.%s not allowed", src.Layer, src.Path, idx, key.Value)
		}
	}
	return spec, nil
}

func decodeEvalReject(src Source, idx int, node *yaml.Node) (policy.RejectSpec, error) {
	if node.Kind != yaml.MappingNode {
		return policy.RejectSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].reject must be a mapping", src.Layer, src.Path, idx)
	}
	var reject policy.RejectSpec
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		val := node.Content[i+1]
		switch key.Value {
		case "message":
			if val.Kind != yaml.ScalarNode {
				return policy.RejectSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].reject.message must be a string", src.Layer, src.Path, idx)
			}
			reject.Message = val.Value
		default:
			return policy.RejectSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].reject.%s not allowed", src.Layer, src.Path, idx, key.Value)
		}
	}
	return reject, nil
}

func decodeEvalMatch(src Source, idx int, node *yaml.Node) (policy.MatchSpec, error) {
	if node.Kind != yaml.MappingNode {
		return policy.MatchSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].match must be a mapping", src.Layer, src.Path, idx)
	}
	var match policy.MatchSpec
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		val := node.Content[i+1]
		switch key.Value {
		case "command":
			if val.Kind != yaml.ScalarNode {
				return policy.MatchSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].match.command must be a string", src.Layer, src.Path, idx)
			}
			match.Command = val.Value
		case "command_in":
			values, err := decodeStringSequence(src, idx, "match.command_in", val)
			if err != nil {
				return policy.MatchSpec{}, err
			}
			match.CommandIn = values
		case "subcommand":
			if val.Kind != yaml.ScalarNode {
				return policy.MatchSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].match.subcommand must be a string", src.Layer, src.Path, idx)
			}
			match.Subcommand = val.Value
		case "args_contains":
			values, err := decodeStringSequence(src, idx, "match.args_contains", val)
			if err != nil {
				return policy.MatchSpec{}, err
			}
			match.ArgsContains = values
		case "args_prefixes":
			values, err := decodeStringSequence(src, idx, "match.args_prefixes", val)
			if err != nil {
				return policy.MatchSpec{}, err
			}
			match.ArgsPrefixes = values
		case "env_requires":
			values, err := decodeStringSequence(src, idx, "match.env_requires", val)
			if err != nil {
				return policy.MatchSpec{}, err
			}
			match.EnvRequires = values
		case "env_missing":
			values, err := decodeStringSequence(src, idx, "match.env_missing", val)
			if err != nil {
				return policy.MatchSpec{}, err
			}
			match.EnvMissing = values
		default:
			return policy.MatchSpec{}, fmt.Errorf("%s config %s is invalid: rules[%d].match.%s not allowed", src.Layer, src.Path, idx, key.Value)
		}
	}
	return match, nil
}

func decodeStringSequence(src Source, idx int, field string, node *yaml.Node) ([]string, error) {
	if node.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("%s config %s is invalid: rules[%d].%s must be a sequence", src.Layer, src.Path, idx, field)
	}
	values := make([]string, 0, len(node.Content))
	for _, item := range node.Content {
		if item.Kind != yaml.ScalarNode {
			return nil, fmt.Errorf("%s config %s is invalid: rules[%d].%s must contain only strings", src.Layer, src.Path, idx, field)
		}
		values = append(values, item.Value)
	}
	return values, nil
}

func validateFile(file File) []string {
	var issues []string
	if file.Version != 1 {
		issues = append(issues, "version must be 1")
	}
	if len(file.Rules) == 0 {
		issues = append(issues, "rules must be non-empty")
	}
	issues = append(issues, policy.ValidateRules(file.Rules)...)
	return issues
}

func validateEvalFile(file evalFile) []string {
	var issues []string
	if file.Version != 1 {
		issues = append(issues, "version must be 1")
	}
	if len(file.Rules) == 0 {
		issues = append(issues, "rules must be non-empty")
	}
	seen := map[string]struct{}{}
	for i, r := range file.Rules {
		prefix := fmt.Sprintf("rules[%d]", i)
		if !regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`).MatchString(r.ID) {
			issues = append(issues, prefix+".id must match [a-z0-9][a-z0-9-]*")
		}
		if _, ok := seen[r.ID]; ok && r.ID != "" {
			issues = append(issues, prefix+".id duplicates another rule in the same file")
		}
		seen[r.ID] = struct{}{}
		issues = append(issues, policy.ValidateRuleMatcher(prefix, r.Pattern, r.Match)...)
		issues = append(issues, policy.ValidateDirective(prefix, r.Message, r.Reject, r.Rewrite)...)
		if r.BlockExampleLen == 0 {
			issues = append(issues, prefix+".block_examples must be non-empty")
		}
		if r.AllowExampleLen == 0 {
			issues = append(issues, prefix+".allow_examples must be non-empty")
		}
	}
	return issues
}

func loadEvalCache(src Source, cachePath string, info os.FileInfo) ([]policy.Rule, bool) {
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, false
	}
	var cache evalCacheFile
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, false
	}
	if cache.Version != 1 || cache.SourcePath != src.Path || cache.SourceSize != info.Size() || cache.SourceModTime != info.ModTime().UnixNano() {
		return nil, false
	}
	rules := make([]policy.Rule, 0, len(cache.CompiledRules))
	for _, spec := range cache.CompiledRules {
		if strings.TrimSpace(spec.Pattern) != "" {
			if _, err := regexp.Compile(spec.Pattern); err != nil {
				return nil, false
			}
		}
		rules = append(rules, policy.NewRule(policy.RuleSpec{
			ID:      spec.ID,
			Pattern: spec.Pattern,
			Matcher: spec.Match,
			Message: spec.Message,
			Reject:  spec.Reject,
			Rewrite: spec.Rewrite,
		}, src))
	}
	return rules, true
}

func writeEvalCache(cachePath string, cache evalCacheFile) {
	if err := os.MkdirAll(filepath.Dir(cachePath), 0o755); err != nil {
		return
	}
	data, err := json.Marshal(cache)
	if err != nil {
		return
	}
	_ = os.WriteFile(cachePath, data, 0o644)
}
