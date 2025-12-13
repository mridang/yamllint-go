package keyordering

import (
	"fmt"
	"iter"
	"regexp"

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
	IgnoredKeys []string `yaml:"ignored-keys"`
}

type Rule struct{}

func (r *Rule) ID() string {
	return "key-ordering"
}

func (r *Rule) Type() string {
	return "token"
}

func (r *Rule) DefaultConfig() Config {
	return Config{
		IgnoredKeys: []string{},
	}
}

func (r *Rule) Check(config Config, token *types.Token, prev, next, nextNext *types.Token, context map[string]any) iter.Seq[types.Problem] {
	var ignoredRegexes []*regexp.Regexp
	for _, pattern := range config.IgnoredKeys {
		if re, err := regexp.Compile(pattern); err == nil {
			ignoredRegexes = append(ignoredRegexes, re)
		}
	}

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
					isIgnored := false
					for _, re := range ignoredRegexes {
						if re.MatchString(next.Value) {
							isIgnored = true
							break
						}
					}

					if !isIgnored {
						isWrongOrder := false
						for _, key := range top.Keys {
							if next.Value < key {
								isWrongOrder = true
								break
							}
						}

						if isWrongOrder {
							if !yield(types.Problem{
								Line:   next.StartMark.Line,
								Column: next.StartMark.Column,
								Desc:   fmt.Sprintf("wrong ordering of key \"%s\" in mapping", next.Value),
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
}
