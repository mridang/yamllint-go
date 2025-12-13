package linter

import (
	"github.com/mridang/yamllint-go/internal/config"
	"github.com/mridang/yamllint-go/internal/parser"
	"github.com/mridang/yamllint-go/internal/types"
)

type Linter struct {
	config *config.Config
	rules  map[string]RuleRunner
}

// RuleRunner is the interface for running any rule
type RuleRunner interface {
	Run(tokens []*types.Token, lines []*types.Line, comments []*types.Comment, cfg map[string]any) []types.Problem
}

func New(cfg *config.Config) *Linter {
	return &Linter{
		config: cfg,
		rules:  make(map[string]RuleRunner),
	}
}

// RegisterRule registers a rule runner
func (l *Linter) RegisterRule(name string, runner RuleRunner) {
	l.rules[name] = runner
}

func (l *Linter) Run(data []byte) ([]types.Problem, error) {
	var problems []types.Problem

	// Parse using the goccy parser (which handles invalid YAML)
	tokens, comments, lines, err := parser.Parse(string(data))
	if err != nil {
		return nil, err
	}

	// Run enabled rules
	for ruleID, ruleCfg := range l.config.Rules {
		if ruleCfg.Level == "disable" {
			continue
		}

		if runner, ok := l.rules[ruleID]; ok {
			ruleProblems := runner.Run(tokens, lines, comments, ruleCfg.Config)
			for i := range ruleProblems {
				ruleProblems[i].Rule = ruleID
			}
			problems = append(problems, ruleProblems...)
		}
	}

	return problems, nil
}
