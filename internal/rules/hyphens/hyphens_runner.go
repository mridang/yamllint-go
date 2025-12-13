package hyphens

import "github.com/mridang/yamllint-go/internal/types"

// Runner implements linter.RuleRunner for the hyphens rule
type Runner struct{}

func (r *Runner) Run(
	tokens []*types.Token,
	lines []*types.Line,
	comments []*types.Comment,
	cfg map[string]any,
) []types.Problem {
	var problems []types.Problem

	// Rule + default config
	rule := &Rule{}
	ruleConfig := rule.DefaultConfig()

	// Apply config from YAML
	if v, ok := cfg["max-spaces-after"].(int); ok {
		ruleConfig.MaxSpacesAfter = v
	}

	context := make(map[string]any)

	for i, token := range tokens {
		var prev, next, nextNext *types.Token

		if i > 0 {
			prev = tokens[i-1]
		}
		if i+1 < len(tokens) {
			next = tokens[i+1]
		}
		if i+2 < len(tokens) {
			nextNext = tokens[i+2]
		}

		for p := range rule.Check(ruleConfig, token, prev, next, nextNext, context) {
			problems = append(problems, p)
		}
	}

	return problems
}
