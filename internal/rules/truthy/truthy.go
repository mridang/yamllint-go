package truthy

import (
	"fmt"
	"iter"
	"sort"

	"github.com/mridang/yamllint-go/internal/types"
)

var truthy11 = []string{
	"YES", "Yes", "yes",
	"NO", "No", "no",
	"TRUE", "True", "true",
	"FALSE", "False", "false",
	"ON", "On", "on",
	"OFF", "Off", "off",
}

var truthy12 = []string{
	"TRUE", "True", "true",
	"FALSE", "False", "false",
}

type Config struct {
	AllowedValues []string `yaml:"allowed-values"`
	CheckKeys     bool     `yaml:"check-keys"`
}

type Rule struct{}

func (r *Rule) ID() string   { return "truthy" }
func (r *Rule) Type() string { return "token" }

func (r *Rule) DefaultConfig() Config {
	return Config{
		AllowedValues: []string{"true", "false"},
		CheckKeys:     true,
	}
}

func (r *Rule) Check(
	conf Config,
	token *types.Token,
	prev, next, nextNext *types.Token,
	ctx map[string]any,
) iter.Seq[types.Problem] {
	return func(yield func(types.Problem) bool) {

		// YAML directive handling
		if token.Type == types.DirectiveToken && token.Value == "YAML" {
			if next != nil {
				ctx["yaml_version"] = next.Value
			}
			return
		}

		if token.Type == types.DocumentEndToken {
			delete(ctx, "yaml_version")
			delete(ctx, "bad_truthy")
			return
		}

		// Ignore explicitly tagged scalars
		if prev != nil && prev.Type == types.TagToken {
			return
		}

		// Ignore keys when check-keys is false
		// (Go token stream uses ValueToken instead of KeyToken)
		if !conf.CheckKeys &&
			token.Type == types.ScalarToken &&
			token.Plain &&
			(next != nil && next.Type == types.ValueToken) {
			return
		}

		if token.Type != types.ScalarToken || !token.Plain {
			return
		}

		// Build bad truthy set once per document
		if _, ok := ctx["bad_truthy"]; !ok {
			var base []string
			if ctx["yaml_version"] == "1.2" {
				base = truthy12
			} else {
				base = truthy11
			}

			bad := map[string]struct{}{}
			for _, v := range base {
				bad[v] = struct{}{}
			}
			for _, v := range conf.AllowedValues {
				delete(bad, v)
			}
			ctx["bad_truthy"] = bad
		}

		bad := ctx["bad_truthy"].(map[string]struct{})
		if _, ok := bad[token.Value]; ok {
			allowed := append([]string{}, conf.AllowedValues...)
			sort.Strings(allowed)

			yield(types.Problem{
				Line:   token.StartMark.Line + 1,
				Column: token.StartMark.Column + 1,
				Desc: fmt.Sprintf(
					"truthy value should be one of [%s]",
					join(allowed),
				),
			})
		}
	}
}

func join(v []string) string {
	out := ""
	for i, s := range v {
		if i > 0 {
			out += ", "
		}
		out += s
	}
	return out
}
