package docstart

import (
	"iter"

	"github.com/mridang/yamllint-go/internal/types"
)

type Config struct {
	Present bool `yaml:"present"`
}

type Rule struct{}

func (r *Rule) ID() string {
	return "document-start"
}

func (r *Rule) Type() string {
	return "token"
}

func (r *Rule) DefaultConfig() Config {
	return Config{
		Present: true,
	}
}

func (r *Rule) Check(config Config, token *types.Token, prev, next, nextNext *types.Token, context map[string]any) iter.Seq[types.Problem] {
	return func(yield func(types.Problem) bool) {
		if config.Present {
			prevValid := false
			if prev != nil {
				prevValid = prev.Type == types.StreamStartToken ||
					prev.Type == types.DocumentEndToken ||
					prev.Type == types.DirectiveToken
			}

			tokenValid := token.Type == types.DocumentStartToken ||
				token.Type == types.DirectiveToken ||
				token.Type == types.StreamEndToken

			if prevValid && !tokenValid {
				if !yield(types.Problem{
					Line:   token.StartMark.Line,
					Column: 0,
					Desc:   "missing document start \"---\"",
				}) {
					return
				}
			}

		} else {
			if token.Type == types.DocumentStartToken {
				if !yield(types.Problem{
					Line:   token.StartMark.Line,
					Column: token.StartMark.Column,
					Desc:   "found forbidden document start \"---\"",
				}) {
					return
				}
			}
		}
	}
}
