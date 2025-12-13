package trailinglines

import (
	"fmt"
	"iter"

	"github.com/mridang/yamllint-go/internal/types"
)

type Config struct{}

type Rule struct{}

func (r *Rule) ID() string {
	return "trailing-lines"
}

func (r *Rule) Type() string {
	return "line"
}

func (r *Rule) DefaultConfig() Config {
	return Config{}
}

func (r *Rule) Check(config Config, line *types.Line) iter.Seq[types.Problem] {
	return func(yield func(types.Problem) bool) {
		bufLen := len(line.Buffer)
		if line.Start > 0 && line.End == bufLen {
			pos := line.End - 1

			if pos >= 0 && line.Buffer[pos] == '\n' {
				pos--
			}

			if pos >= 0 && line.Buffer[pos] == '\r' {
				pos--
			}

			if pos >= 0 && (line.Buffer[pos] == '\n' || line.Buffer[pos] == '\r') {
				if !yield(types.Problem{
					Line:   line.LineNo,
					Column: 0,
					Desc:   fmt.Sprintf("too many blank lines at end of file (%d > 0)", 1),
				}) {
					return
				}
			}
		}
	}
}
