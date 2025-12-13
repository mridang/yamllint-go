package hyphens

import (
	"iter"

	"github.com/mridang/yamllint-go/internal/types"
)

type Config struct {
	MaxSpacesAfter int `yaml:"max-spaces-after"`
}

type Rule struct{}

func (r *Rule) ID() string   { return "hyphens" }
func (r *Rule) Type() string { return "token" }

func (r *Rule) DefaultConfig() Config {
	return Config{MaxSpacesAfter: 1}
}

func (r *Rule) Check(
	config Config,
	token *types.Token,
	prev, next, nextNext *types.Token,
	context map[string]any,
) iter.Seq[types.Problem] {
	return func(yield func(types.Problem) bool) {

		// Disabled
		if config.MaxSpacesAfter < 0 {
			return
		}

		if token.Type != types.BlockEntryToken {
			return
		}

		// If no following token, nothing to check
		if next == nil {
			return
		}

		// Count spaces between '-' and next token
		spaces := next.StartMark.Column - token.StartMark.Column - 1

		if spaces > config.MaxSpacesAfter {
			// yamllint reports the FIRST illegal space
			illegalColumn := token.StartMark.Column + config.MaxSpacesAfter + 2

			yield(types.Problem{
				Line:   token.StartMark.Line + 1,
				Column: illegalColumn,
				Desc:   "too many spaces after hyphen",
			})
		}
	}
}
