package keyduplicates

import (
	"fmt"
	"iter"

	"github.com/mridang/yamllint-go/internal/types"
)

const (
	MAP = iota
	SEQ
)

type Parent struct {
	Type int
	Keys []string
}

type Config struct {
	ForbidDuplicatedMergeKeys bool `yaml:"forbid-duplicated-merge-keys"`
}

type Rule struct{}

func (r *Rule) ID() string {
	return "key-duplicates"
}

func (r *Rule) Type() string {
	return "token"
}

func (r *Rule) DefaultConfig() Config {
	return Config{
		ForbidDuplicatedMergeKeys: false,
	}
}

func (r *Rule) Check(config Config, token *types.Token, prev, next, nextNext *types.Token, context map[string]any) iter.Seq[types.Problem] {
	return func(yield func(types.Problem) bool) {
		if _, ok := context["stack"]; !ok {
			context["stack"] = []*Parent{}
		}
		stack := context["stack"].([]*Parent)

		if token.Type == types.BlockMappingStartToken || token.Type == types.FlowMappingStartToken {
			context["stack"] = append(stack, &Parent{Type: MAP, Keys: []string{}})
		} else if token.Type == types.BlockSequenceStartToken || token.Type == types.FlowSequenceStartToken {
			context["stack"] = append(stack, &Parent{Type: SEQ, Keys: []string{}})
		} else if token.Type == types.BlockEndToken || token.Type == types.FlowMappingEndToken || token.Type == types.FlowSequenceEndToken {
			if len(stack) > 0 {
				context["stack"] = stack[:len(stack)-1]
			}
		} else if token.Type == types.KeyToken && next != nil && next.Type == types.ScalarToken {
			if len(stack) > 0 {
				top := stack[len(stack)-1]
				if top.Type == MAP {
					found := false
					for _, k := range top.Keys {
						if k == next.Value {
							found = true
							break
						}
					}

					if found && (next.Value != "<<" || config.ForbidDuplicatedMergeKeys) {
						if !yield(types.Problem{
							Line:   next.StartMark.Line,
							Column: next.StartMark.Column,
							Desc:   fmt.Sprintf("duplication of key \"%s\" in mapping", next.Value),
						}) {
							return
						}
					} else {
						top.Keys = append(top.Keys, next.Value)
					}
				}
			}
		}
	}
}
