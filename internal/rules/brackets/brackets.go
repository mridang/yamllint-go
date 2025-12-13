package brackets

import (
	"iter"

	"github.com/mridang/yamllint-go/internal/types"
)

type Config struct {
	Forbid               any `yaml:"forbid"`
	MinSpacesInside      int `yaml:"min-spaces-inside"`
	MaxSpacesInside      int `yaml:"max-spaces-inside"`
	MinSpacesInsideEmpty int `yaml:"min-spaces-inside-empty"`
	MaxSpacesInsideEmpty int `yaml:"max-spaces-inside-empty"`
}

type Rule struct{}

func (r *Rule) ID() string   { return "brackets" }
func (r *Rule) Type() string { return "token" }

func (r *Rule) DefaultConfig() Config {
	return Config{
		Forbid:               false,
		MinSpacesInside:      0,
		MaxSpacesInside:      0,
		MinSpacesInsideEmpty: -1,
		MaxSpacesInsideEmpty: -1,
	}
}

func spacesAfter(t, next *types.Token) int {
	if t == nil || next == nil {
		return 0
	}
	if t.EndMark.Line != next.StartMark.Line {
		return 0
	}
	return next.StartMark.Index - t.EndMark.Index
}

func spacesBefore(prev, t *types.Token) int {
	if prev == nil || t == nil {
		return 0
	}
	if prev.EndMark.Line != t.StartMark.Line {
		return 0
	}
	return t.StartMark.Index - prev.EndMark.Index
}

func (r *Rule) Check(
	conf Config,
	token *types.Token,
	prev, next, nextNext *types.Token,
	context map[string]any,
) iter.Seq[types.Problem] {

	return func(yield func(types.Problem) bool) {

		// ---------- '[' ----------
		if token.Type == types.FlowSequenceStartToken {

			isEmpty :=
				next != nil &&
					(next.Type == types.FlowSequenceEndToken ||
						(nextNext != nil && nextNext.Type == types.FlowSequenceEndToken))

			// forbid
			if conf.Forbid == true ||
				(conf.Forbid == "non-empty" && !isEmpty) {

				context["skip-flow-seq"] = true
				yield(types.Problem{
					Line:   token.StartMark.Line,
					Column: token.EndMark.Column,
					Desc:   "forbidden flow sequence",
				})
				return
			}

			min := conf.MinSpacesInside
			max := conf.MaxSpacesInside

			if isEmpty {
				if conf.MinSpacesInsideEmpty != -1 {
					min = conf.MinSpacesInsideEmpty
				}
				if conf.MaxSpacesInsideEmpty != -1 {
					max = conf.MaxSpacesInsideEmpty
				}
				context["skip-flow-seq"] = true
			}

			spaces := spacesAfter(token, next)

			if spaces < min {
				yield(types.Problem{
					Line:   token.StartMark.Line,
					Column: token.EndMark.Column,
					Desc:   "too few spaces inside brackets",
				})
			} else if max >= 0 && spaces > max {
				yield(types.Problem{
					Line:   token.StartMark.Line,
					Column: token.EndMark.Column,
					Desc:   "too many spaces inside brackets",
				})
			}
			return
		}

		// ---------- ']' ----------
		if token.Type == types.FlowSequenceEndToken {

			if context["skip-flow-seq"] == true {
				delete(context, "skip-flow-seq")
				return
			}

			spaces := spacesBefore(prev, token)

			if spaces < conf.MinSpacesInside {
				yield(types.Problem{
					Line:   token.StartMark.Line,
					Column: token.StartMark.Column,
					Desc:   "too few spaces inside brackets",
				})
			} else if conf.MaxSpacesInside >= 0 && spaces > conf.MaxSpacesInside {
				yield(types.Problem{
					Line:   token.StartMark.Line,
					Column: token.StartMark.Column,
					Desc:   "too many spaces inside brackets",
				})
			}
		}
	}
}
