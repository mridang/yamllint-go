package newlines

import (
	"fmt"
	"iter"
	"runtime"
	"strings"

	"github.com/mridang/yamllint-go/internal/types"
)

type Config struct {
	Type string `yaml:"type"`
}

type Rule struct{}

func (r *Rule) ID() string {
	return "new-lines"
}

func (r *Rule) Type() string {
	return "line"
}

func (r *Rule) DefaultConfig() Config {
	return Config{
		Type: "unix",
	}
}

func (r *Rule) Check(config Config, line *types.Line) iter.Seq[types.Problem] {
	return func(yield func(types.Problem) bool) {
		var newlineChar string

		switch config.Type {
		case "unix":
			newlineChar = "\n"
		case "dos":
			newlineChar = "\r\n"
		case "platform":
			if runtime.GOOS == "windows" {
				newlineChar = "\r\n"
			} else {
				newlineChar = "\n"
			}
		default:
			newlineChar = "\n"
		}

		if line.Start == 0 && len(line.Buffer) > line.End {
			if line.End+len(newlineChar) <= len(line.Buffer) {
				actual := line.Buffer[line.End : line.End+len(newlineChar)]

				if actual != newlineChar {
					c := strings.Trim(fmt.Sprintf("%q", newlineChar), "\"")

					if !yield(types.Problem{
						Line:   0,
						Column: line.End - line.Start,
						Desc:   fmt.Sprintf("wrong new line character: expected %s", c),
					}) {
						return
					}
				}
			}
		}
	}
}
