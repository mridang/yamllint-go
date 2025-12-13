package docend

import (
	"iter"

	"github.com/mridang/yamllint-go/internal/types"
)

type Config struct {
	Present bool `yaml:"present"`
}

type Rule struct{}

func (r *Rule) ID() string {
	return "document-end"
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
			isStreamEnd := token.Type == types.StreamEndToken
			isStart := token.Type == types.DocumentStartToken

			prevIsEndOrStreamStart := false
			if prev != nil {
				prevIsEndOrStreamStart = prev.Type == types.DocumentEndToken || prev.Type == types.StreamStartToken
			}

			prevIsDirective := false
			if prev != nil {
				prevIsDirective = prev.Type == types.DirectiveToken
			}

			if isStreamEnd && !prevIsEndOrStreamStart {
				if !yield(types.Problem{
					Line:   token.StartMark.Line,
					Column: 0,
					Desc:   "missing document end \"...\"",
				}) {
					return
				}
			} else if isStart && !(prevIsEndOrStreamStart || prevIsDirective) {
				if !yield(types.Problem{
					Line:   token.StartMark.Line,
					Column: 0,
					Desc:   "missing document end \"...\"",
				}) {
					return
				}
			}
		} else {
			if token.Type == types.DocumentEndToken {
				if !yield(types.Problem{
					Line:   token.StartMark.Line,
					Column: token.StartMark.Column,
					Desc:   "found forbidden document end \"...\"",
				}) {
					return
				}
			}
		}
	}
}
