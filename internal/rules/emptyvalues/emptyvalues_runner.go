package emptyvalues

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

	if v, ok := cfg["forbid-in-block-mappings"].(bool); ok {
		conf.ForbidInBlockMappings = v
	}
	if v, ok := cfg["forbid-in-flow-mappings"].(bool); ok {
		conf.ForbidInFlowMappings = v
	}
	if v, ok := cfg["forbid-in-block-sequences"].(bool); ok {
		conf.ForbidInBlockSequences = v
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
