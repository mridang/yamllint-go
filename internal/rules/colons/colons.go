package colons

import (
	"iter"

	"github.com/mridang/yamllint-go/internal/common"
	"github.com/mridang/yamllint-go/internal/types"
)

type Config struct {
	MaxSpacesBefore int `yaml:"max-spaces-before"`
	MaxSpacesAfter  int `yaml:"max-spaces-after"`
}

type Rule struct{}

func (r *Rule) ID() string {
	return "colons"
}

func (r *Rule) Type() string {
	return "token"
}

func (r *Rule) DefaultConfig() Config {
	return Config{
		MaxSpacesBefore: 0,
		MaxSpacesAfter:  1,
	}
}

func (r *Rule) Check(config Config, token *types.Token, prev, next, nextNext *types.Token, context map[string]any) iter.Seq[types.Problem] {
	return func(yield func(types.Problem) bool) {
		if token.Type == types.ValueToken && !(prev != nil && prev.Type == types.AliasToken &&
			token.StartMark.Index-prev.EndMark.Index == 1) {

			problem := common.SpacesBefore(token, prev, next,
				config.MaxSpacesBefore,
				"too many spaces before colon")
			if problem != nil {
				if !yield(*problem) {
					return
				}
			}

			problem = common.SpacesAfter(token, prev, next,
				config.MaxSpacesAfter,
				"too many spaces after colon")
			if problem != nil {
				if !yield(*problem) {
					return
				}
			}
		}

		if token.Type == types.KeyToken && common.IsExplicitKey(token) {
			problem := common.SpacesAfter(token, prev, next,
				config.MaxSpacesAfter,
				"too many spaces after question mark")
			if problem != nil {
				if !yield(*problem) {
					return
				}
			}
		}
	}
}
