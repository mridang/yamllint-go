package braces

import (
	"fmt"

	"github.com/mridang/yamllint-go/internal/types"
)

// Runner implements linter.RuleRunner for the braces rule
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

	// DEBUG: Print config
	fmt.Printf("\n=== BRACES CONFIG ===\n")
	fmt.Printf("Forbid: %v\n", ruleConfig.Forbid)
	fmt.Printf("MinSpacesInside: %d\n", ruleConfig.MinSpacesInside)
	fmt.Printf("MaxSpacesInside: %d\n", ruleConfig.MaxSpacesInside)
	fmt.Printf("MinSpacesInsideEmpty: %d\n", ruleConfig.MinSpacesInsideEmpty)
	fmt.Printf("MaxSpacesInsideEmpty: %d\n", ruleConfig.MaxSpacesInsideEmpty)
	fmt.Println()

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

		// DEBUG: Print flow mapping tokens
		if token.Type == types.FlowMappingStartToken || token.Type == types.FlowMappingEndToken {
			fmt.Printf("Token[%d]: Type=%d Line=%d Col=%d\n", i, token.Type, token.StartMark.Line, token.StartMark.Column)
			if next != nil {
				fmt.Printf("  Next: Type=%d Line=%d Col=%d StartIdx=%d\n", next.Type, next.StartMark.Line, next.StartMark.Column, next.StartMark.Index)
				spaces := next.StartMark.Index - token.EndMark.Index
				fmt.Printf("  Spaces=%d (next.Start=%d - token.End=%d)\n", spaces, next.StartMark.Index, token.EndMark.Index)
			}
		}

		// Collect problems from iterator
		for problem := range rule.Check(ruleConfig, token, prev, next, nextNext, context) {
			fmt.Printf("PROBLEM: Line=%d Col=%d Desc=%q\n", problem.Line, problem.Column, problem.Desc)
			problems = append(problems, problem)
		}
	}

	return problems
}
