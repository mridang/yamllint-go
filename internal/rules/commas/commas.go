package commas

import (
	"iter"

	"github.com/mridang/yamllint-go/internal/common"
	"github.com/mridang/yamllint-go/internal/types"
)

type Config struct {
	MaxSpacesBefore int `yaml:"max-spaces-before"`
	MinSpacesAfter  int `yaml:"min-spaces-after"`
	MaxSpacesAfter  int `yaml:"max-spaces-after"`
}

type Rule struct{}

func (r *Rule) ID() string {
	return "commas"
}

func (r *Rule) Type() string {
	return "token"
}

func (r *Rule) DefaultConfig() Config {
	return Config{
		MaxSpacesBefore: 0,
		MinSpacesAfter:  1,
		MaxSpacesAfter:  1,
	}
}

func (r *Rule) Check(config Config, token *types.Token, prev, next, nextNext *types.Token, context map[string]any) iter.Seq[types.Problem] {
	return func(yield func(types.Problem) bool) {
		if token.Type == types.FlowEntryToken {
			if prev != nil && config.MaxSpacesBefore != -1 && prev.EndMark.Line < token.StartMark.Line {
				col := token.StartMark.Column
				if col < 1 {
					col = 1
				}
				if !yield(types.Problem{
					Line:   token.StartMark.Line,
					Column: col,
					Desc:   "too many spaces before comma",
				}) {
					return
				}
			} else {
				problem := common.SpacesBefore(token, prev, next,
					config.MaxSpacesBefore,
					"too many spaces before comma")
				if problem != nil {
					if !yield(*problem) {
						return
					}
				}
			}

			problem := common.SpacesAfter(token, prev, next,
				config.MaxSpacesAfter,
				"too many spaces after comma")
			if problem != nil {
				if !yield(*problem) {
					return
				}
			}

			if config.MinSpacesAfter != -1 && next != nil && token.EndMark.Line == next.StartMark.Line {
				if (next.StartMark.Column - token.EndMark.Column) < config.MinSpacesAfter {
					if !yield(types.Problem{
						Line:   token.StartMark.Line,
						Column: token.StartMark.Column,
						Desc:   "too few spaces after comma",
					}) {
						return
					}
				}
			}
		}
	}
}
