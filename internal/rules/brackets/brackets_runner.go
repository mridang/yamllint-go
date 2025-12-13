package brackets

import (
	"github.com/mridang/yamllint-go/internal/types"
)

// Runner implements linter.RuleRunner for the brackets rule
type Runner struct{}

func (r *Runner) Run(tokens []*types.Token, lines []*types.Line, comments []*types.Comment, cfg map[string]any) []types.Problem {
	var problems []types.Problem

	// Start with default config
	rule := &Rule{}
	ruleConfig := rule.DefaultConfig()
	
	// Override with provided config
	if val, ok := cfg["forbid"].(string); ok {
		ruleConfig.Forbid = val
	} else if val, ok := cfg["forbid"].(bool); ok {
		ruleConfig.Forbid = val
	}
	
	if val, ok := cfg["min-spaces-inside"].(int); ok {
		ruleConfig.MinSpacesInside = val
	}
	if val, ok := cfg["max-spaces-inside"].(int); ok {
		ruleConfig.MaxSpacesInside = val
	}
	if val, ok := cfg["min-spaces-inside-empty"].(int); ok {
		ruleConfig.MinSpacesInsideEmpty = val
	}
	if val, ok := cfg["max-spaces-inside-empty"].(int); ok {
		ruleConfig.MaxSpacesInsideEmpty = val
	}

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
