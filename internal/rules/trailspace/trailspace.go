// file: internal/rules/trailspace/trailspace.go
package trailingspaces

import (
	"iter"

	"github.com/mridang/yamllint-go/internal/types"
)

type Config struct{}

type Rule struct{}

func (r *Rule) ID() string   { return "trailing-spaces" }
func (r *Rule) Type() string { return "line" }

func (r *Rule) DefaultConfig() Config {
	return Config{}
}

func (r *Rule) Check(config Config, line *types.Line) iter.Seq[types.Problem] {
	return func(yield func(types.Problem) bool) {
		if line.End == 0 {
			return
		}

		pos := line.End
		for pos > line.Start {
			c := line.Buffer[pos-1]
			if c == '\n' || c == '\r' {
				pos--
				continue
			}
			if c == ' ' || c == '\t' {
				pos--
				continue
			}
			break
		}

		if pos != line.End {
			c := line.Buffer[pos]
			if c == ' ' || c == '\t' {
				yield(types.Problem{
					Line:   line.LineNo,
					Column: pos - line.Start + 1,
					Desc:   "trailing spaces",
				})
			}
		}
	}
}
