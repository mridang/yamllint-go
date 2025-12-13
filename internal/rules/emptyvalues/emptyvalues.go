package emptyvalues

import (
	"iter"

	"github.com/mridang/yamllint-go/internal/types"
)

type Config struct {
	ForbidInBlockMappings  bool `yaml:"forbid-in-block-mappings"`
	ForbidInFlowMappings   bool `yaml:"forbid-in-flow-mappings"`
	ForbidInBlockSequences bool `yaml:"forbid-in-block-sequences"`
}

type Rule struct{}

func (r *Rule) ID() string   { return "empty-values" }
func (r *Rule) Type() string { return "token" }

func (r *Rule) DefaultConfig() Config {
	return Config{
		ForbidInBlockMappings:  true,
		ForbidInFlowMappings:   true,
		ForbidInBlockSequences: true,
	}
}

func (r *Rule) Check(
	conf Config,
	token *types.Token,
	prev, next, nextNext *types.Token,
	ctx map[string]any,
) iter.Seq[types.Problem] {
	return func(yield func(types.Problem) bool) {

		// --- flow mapping depth tracking ---
		if token.Type == types.FlowMappingStartToken {
			ctx["flow_depth"] = ctx["flow_depth"].(int) + 1
		}
		if token.Type == types.FlowMappingEndToken {
			ctx["flow_depth"] = ctx["flow_depth"].(int) - 1
		}

		flowDepth, _ := ctx["flow_depth"].(int)

		// --- block mappings ---
		if conf.ForbidInBlockMappings &&
			flowDepth == 0 &&
			token.Type == types.ValueToken &&
			next != nil &&
			(next.Type == types.KeyToken ||
				next.Type == types.BlockEndToken) {

			yield(types.Problem{
				Line:   token.StartMark.Line + 1,
				Column: token.EndMark.Column + 1,
				Desc:   "empty value in block mapping",
			})
		}

		// --- flow mappings ---
		if conf.ForbidInFlowMappings &&
			flowDepth > 0 &&
			token.Type == types.ValueToken &&
			next != nil &&
			(next.Type == types.FlowEntryToken ||
				next.Type == types.FlowMappingEndToken) {

			yield(types.Problem{
				Line:   token.StartMark.Line + 1,
				Column: token.EndMark.Column + 1,
				Desc:   "empty value in flow mapping",
			})
		}

		// --- block sequences ---
		if conf.ForbidInBlockSequences &&
			token.Type == types.BlockEntryToken &&
			next != nil &&
			(next.Type == types.KeyToken ||
				next.Type == types.BlockEndToken ||
				next.Type == types.BlockEntryToken) {

			yield(types.Problem{
				Line:   token.StartMark.Line + 1,
				Column: token.EndMark.Column + 1,
				Desc:   "empty value in block sequence",
			})
		}
	}
}
