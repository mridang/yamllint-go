package octalvalues

import "github.com/mridang/yamllint-go/internal/types"

// Runner implements linter.RuleRunner for the octal-values rule
type Runner struct{}

func (r *Runner) Run(
	tokens []*types.Token,
	lines []*types.Line,
	comments []*types.Comment,
	cfg map[string]any,
) []types.Problem {
	var problems []types.Problem

	// Default config
	rule := &Rule{}
	ruleConfig := rule.DefaultConfig()

	// Apply config
	if v, ok := cfg["forbid-implicit-octal"].(bool); ok {
		ruleConfig.ForbidImplicitOctal = v
	}
	if v, ok := cfg["forbid-explicit-octal"].(bool); ok {
		ruleConfig.ForbidExplicitOctal = v
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
