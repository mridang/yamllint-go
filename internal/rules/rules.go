package rules

import (
	"iter"

	"github.com/mridang/yamllint-go/internal/types"
)

type TokenRule[C any] interface {
	ID() string
	Type() string
	DefaultConfig() C
	Check(config C, token *types.Token, prev, next, nextNext *types.Token, context map[string]any) iter.Seq[types.Problem]
}

type LineRule[C any] interface {
	ID() string
	Type() string
	DefaultConfig() C
	Check(config C, line *types.Line) iter.Seq[types.Problem]
}

type CommentRule[C any] interface {
	ID() string
	Type() string
	DefaultConfig() C
	Check(config C, comment *types.Comment) iter.Seq[types.Problem]
}
