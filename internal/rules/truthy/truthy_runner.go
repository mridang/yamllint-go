package truthy

import "github.com/mridang/yamllint-go/internal/types"

type Runner struct{}

func (r *Runner) Run(
	tokens []*types.Token,
	lines []*types.Line,
	comments []*types.Comment,
	cfg map[string]any,
) []types.Problem {

	rule := &Rule{}
	conf := rule.DefaultConfig()

	if v, ok := cfg["allowed-values"].([]any); ok {
		conf.AllowedValues = nil
		for _, x := range v {
			conf.AllowedValues = append(conf.AllowedValues, x.(string))
		}
	}

	if v, ok := cfg["check-keys"].(bool); ok {
		conf.CheckKeys = v
	}

	ctx := map[string]any{}
	var problems []types.Problem

	for i, tok := range tokens {
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

		for p := range rule.Check(conf, tok, prev, next, nextNext, ctx) {
			problems = append(problems, p)
		}
	}

	return problems
}
