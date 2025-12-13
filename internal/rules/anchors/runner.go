package anchors

import (
	"github.com/mridang/yamllint-go/internal/types"
)

// Runner implements linter.RuleRunner for the anchors rule
type Runner struct{}

func (r *Runner) Run(tokens []*types.Token, lines []*types.Line, comments []*types.Comment, cfg map[string]any) []types.Problem {
	var problems []types.Problem

	// Create rule config from map
	ruleConfig := Config{}
	if val, ok := cfg["forbid-undeclared-aliases"].(bool); ok {
		ruleConfig.ForbidUndeclaredAliases = val
	}
	if val, ok := cfg["forbid-duplicated-anchors"].(bool); ok {
		ruleConfig.ForbidDuplicatedAnchors = val
	}
	if val, ok := cfg["forbid-unused-anchors"].(bool); ok {
		ruleConfig.ForbidUnusedAnchors = val
	}

	// Create rule instance
	rule := &Rule{}

	// Run rule on all tokens
	context := make(map[string]any)
	for i, token := range tokens {
		var prev, next, nextNext *types.Token
		if i > 0 {
			prev = tokens[i-1]
		}
		if i < len(tokens)-1 {
			next = tokens[i+1]
		}
		if i < len(tokens)-2 {
			nextNext = tokens[i+2]
		}

		// Collect problems from iterator
		for problem := range rule.Check(ruleConfig, token, prev, next, nextNext, context) {
			problems = append(problems, problem)
		}
	}

	return problems
}
