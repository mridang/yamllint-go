package octalvalues

import (
	"fmt"
	"iter"
	"regexp"

	"github.com/mridang/yamllint-go/internal/types"
)

var (
	implicitOctalPattern = regexp.MustCompile(`^0[0-7]+$`)
	explicitOctalPattern = regexp.MustCompile(`^0o[0-7]+$`)
)

type Config struct {
	ForbidImplicitOctal bool `yaml:"forbid-implicit-octal"`
	ForbidExplicitOctal bool `yaml:"forbid-explicit-octal"`
}

type Rule struct{}

func (r *Rule) ID() string   { return "octal-values" }
func (r *Rule) Type() string { return "token" }

func (r *Rule) DefaultConfig() Config {
	return Config{
		ForbidImplicitOctal: false,
		ForbidExplicitOctal: false,
	}
}

func (r *Rule) Check(
	config Config,
	token *types.Token,
	prev, next, nextNext *types.Token,
	context map[string]any,
) iter.Seq[types.Problem] {
	return func(yield func(types.Problem) bool) {
		if prev != nil && prev.Type == types.TagToken {
			return
		}
		if token.Type != types.ScalarToken {
			return
		}
		if token.Style != "" {
			return
		}

		val := token.Value

		if config.ForbidImplicitOctal {
			if implicitOctalPattern.MatchString(val) {
				yield(types.Problem{
					Line:   token.StartMark.Line + 1,
					Column: token.StartMark.Column + 1 + len(val),
					Desc:   fmt.Sprintf("forbidden implicit octal value \"%s\"", val),
				})
			}
		}

		if config.ForbidExplicitOctal {
			if explicitOctalPattern.MatchString(val) {
				yield(types.Problem{
					Line:   token.StartMark.Line + 1,
					Column: token.StartMark.Column + 1 + len(val),
					Desc:   fmt.Sprintf("forbidden explicit octal value \"%s\"", val),
				})
			}
		}
	}
}
