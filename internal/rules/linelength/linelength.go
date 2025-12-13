// file: internal/rules/linelength/linelength.go
package linelength

import (
	"fmt"
	"iter"
	"strings"

	"github.com/mridang/yamllint-go/internal/types"
)

type Config struct {
	Max                             int  `yaml:"max"`
	AllowNonBreakableWords          bool `yaml:"allow-non-breakable-words"`
	AllowNonBreakableInlineMappings bool `yaml:"allow-non-breakable-inline-mappings"`
}

type Rule struct{}

func (r *Rule) ID() string   { return "line-length" }
func (r *Rule) Type() string { return "line" }

func (r *Rule) DefaultConfig() Config {
	return Config{
		Max:                             80,
		AllowNonBreakableWords:          true,
		AllowNonBreakableInlineMappings: false,
	}
}

// Exact port of yamllint inline-mapping logic:
// - key must be non-breakable (no spaces before ':')
// - value must be non-breakable (no spaces after ': ')
func checkInlineMapping(content string) bool {
	colon := strings.Index(content, ":")
	if colon == -1 {
		return false
	}

	key := content[:colon]
	if strings.Contains(key, " ") {
		return false
	}

	i := colon + 1
	if i < len(content) && content[i] == ' ' {
		i++
	}
	if i >= len(content) {
		return false
	}

	value := content[i:]
	return !strings.Contains(value, " ")
}

func (r *Rule) Check(config Config, line *types.Line) iter.Seq[types.Problem] {
	return func(yield func(types.Problem) bool) {
		maxLength := config.Max
		length := line.End - line.Start

		if length <= maxLength {
			return
		}

		allowNonBreakable := config.AllowNonBreakableWords ||
			config.AllowNonBreakableInlineMappings

		if allowNonBreakable {
			start := line.Start
			for start < line.End && line.Buffer[start] == ' ' {
				start++
			}

			if start < line.End {
				switch line.Buffer[start] {
				case '#':
					for start < line.End && line.Buffer[start] == '#' {
						start++
					}
					if start < line.End {
						start++
					}
				case '-':
					if start+1 < line.End {
						start += 2
					}
				}

				if start < line.End {
					content := line.Buffer[start:line.End]
					if !strings.Contains(content, " ") {
						return
					}

					if config.AllowNonBreakableInlineMappings &&
						checkInlineMapping(line.Content()) {
						return
					}
				}
			}
		}

		yield(types.Problem{
			Line:   line.LineNo,
			Column: maxLength + 1,
			Desc:   fmt.Sprintf("line too long (%d > %d characters)", length, maxLength),
		})
	}
}
