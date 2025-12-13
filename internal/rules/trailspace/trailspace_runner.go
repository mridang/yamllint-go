// file: internal/rules/trailspace/trailspace_runner.go
package trailingspaces

import "github.com/mridang/yamllint-go/internal/types"

type Runner struct{}

func (r *Runner) Run(
	tokens []*types.Token,
	lines []*types.Line,
	comments []*types.Comment,
	cfg map[string]any,
) []types.Problem {
	rule := &Rule{}
	var problems []types.Problem

	for _, line := range lines {
		for p := range rule.Check(rule.DefaultConfig(), line) {
			problems = append(problems, p)
		}
	}

	return problems
}
