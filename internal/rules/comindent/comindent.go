package commentsindentation

import (
	"iter"
	"strings"

	"github.com/mridang/yamllint-go/internal/types"
)

type Config struct{}

type Rule struct{}

func (r *Rule) ID() string {
	return "comments-indentation"
}

func (r *Rule) Type() string {
	return "comment"
}

func (r *Rule) DefaultConfig() Config {
	return Config{}
}

func getLineIndent(token *types.Token) int {
	if token == nil {
		return 0
	}
	start := token.StartMark.Index
	buffer := token.StartMark.Buffer

	lineStart := strings.LastIndexByte(buffer[:start], '\n')
	if lineStart == -1 {
		lineStart = 0
	} else {
		lineStart++
	}

	indent := 0
	for i := lineStart; i < len(buffer); i++ {
		if buffer[i] == ' ' {
			indent++
		} else {
			break
		}
	}
	return indent
}

func (r *Rule) Check(config Config, comment *types.Comment) iter.Seq[types.Problem] {
	return func(yield func(types.Problem) bool) {
		if comment.TokenBefore != nil && comment.TokenBefore.Type != types.StreamStartToken {
			if comment.TokenBefore.EndMark.Line == comment.LineNo-1 {
				return
			}
		}

		nextLineIndent := 0
		if comment.TokenAfter != nil {
			if comment.TokenAfter.Type == types.StreamEndToken {
				nextLineIndent = 0
			} else {
				nextLineIndent = comment.TokenAfter.StartMark.Column
			}
		}

		prevLineIndent := 0
		if comment.TokenBefore != nil {
			if comment.TokenBefore.Type == types.StreamStartToken {
				prevLineIndent = 0
			} else {
				prevLineIndent = getLineIndent(comment.TokenBefore)
			}
		}

		if nextLineIndent > prevLineIndent {
			prevLineIndent = nextLineIndent
		}

		if comment.CommentBefore != nil && !comment.CommentBefore.IsInline() {
			prevLineIndent = comment.CommentBefore.ColumnNo - 1
		}

		currentIndent := comment.ColumnNo - 1
		if currentIndent != prevLineIndent && currentIndent != nextLineIndent {
			if !yield(types.Problem{
				Line:   comment.LineNo,
				Column: comment.ColumnNo,
				Desc:   "comment not indented like content",
			}) {
				return
			}
		}
	}
}
