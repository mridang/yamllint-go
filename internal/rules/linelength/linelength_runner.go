// file: internal/rules/linelength/linelength_runner.go
package linelength

import "github.com/mridang/yamllint-go/internal/types"

type Runner struct{}

func (r *Runner) Run(
	_ []*types.Token,
	lines []*types.Line,
	_ []*types.Comment,
	cfg map[string]any,
) []types.Problem {
	rule := &Rule{}
	conf := rule.DefaultConfig()

	if v, ok := cfg["max"].(int); ok {
		conf.Max = v
	}
	if v, ok := cfg["allow-non-breakable-words"].(bool); ok {
		conf.AllowNonBreakableWords = v
	}
	if v, ok := cfg["allow-non-breakable-inline-mappings"].(bool); ok {
		conf.AllowNonBreakableInlineMappings = v
	}

	var problems []types.Problem
	for _, line := range lines {
		for p := range rule.Check(conf, line) {
			problems = append(problems, p)
		}
	}
	return problems
}
