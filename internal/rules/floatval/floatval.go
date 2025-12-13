package floatvalues

import (
	"fmt"
	"iter"
	"regexp"

	"github.com/mridang/yamllint-go/internal/types"
)

var (
	isNumeralBeforeDecimal = regexp.MustCompile(`^[-+]?(\.[0-9]+)([eE][-+]?[0-9]+)?$`)
	isScientificNotation   = regexp.MustCompile(`^[-+]?(\.[0-9]+|[0-9]+(\.[0-9]*)?)([eE][-+]?[0-9]+)$`)
	isInf                  = regexp.MustCompile(`^[-+]?(\.inf|\.Inf|\.INF)$`)
	isNaN                  = regexp.MustCompile(`^(\.nan|\.NaN|\.NAN)$`)
)

type Config struct {
	RequireNumeralBeforeDecimal bool `yaml:"require-numeral-before-decimal"`
	ForbidScientificNotation    bool `yaml:"forbid-scientific-notation"`
	ForbidNan                   bool `yaml:"forbid-nan"`
	ForbidInf                   bool `yaml:"forbid-inf"`
}

type Rule struct{}

func (r *Rule) ID() string   { return "float-values" }
func (r *Rule) Type() string { return "token" }

func (r *Rule) DefaultConfig() Config {
	return Config{
		RequireNumeralBeforeDecimal: false,
		ForbidScientificNotation:    false,
		ForbidNan:                   false,
		ForbidInf:                   false,
	}
}

func (r *Rule) Check(
	conf Config,
	token *types.Token,
	prev, next, nextNext *types.Token,
	ctx map[string]any,
) iter.Seq[types.Problem] {
	return func(yield func(types.Problem) bool) {

		if prev != nil && prev.Type == types.TagToken {
			return
		}

		if token.Type != types.ScalarToken || token.Style != "" {
			return
		}

		val := token.Value

		if conf.ForbidNan && isNaN.MatchString(val) {
			yield(types.Problem{
				Line:   token.StartMark.Line + 1,
				Column: token.StartMark.Column + 1,
				Desc:   fmt.Sprintf(`forbidden not a number value "%s"`, val),
			})
		}

		if conf.ForbidInf && isInf.MatchString(val) {
			yield(types.Problem{
				Line:   token.StartMark.Line + 1,
				Column: token.StartMark.Column + 1,
				Desc:   fmt.Sprintf(`forbidden infinite value "%s"`, val),
			})
		}

		if conf.ForbidScientificNotation && isScientificNotation.MatchString(val) {
			yield(types.Problem{
				Line:   token.StartMark.Line + 1,
				Column: token.StartMark.Column + 1,
				Desc:   fmt.Sprintf(`forbidden scientific notation "%s"`, val),
			})
		}

		if conf.RequireNumeralBeforeDecimal && isNumeralBeforeDecimal.MatchString(val) {
			yield(types.Problem{
				Line:   token.StartMark.Line + 1,
				Column: token.StartMark.Column + 1,
				Desc:   fmt.Sprintf(`forbidden decimal missing 0 prefix "%s"`, val),
			})
		}
	}
}
