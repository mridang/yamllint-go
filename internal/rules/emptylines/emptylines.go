package emptylines

import (
	"fmt"
	"iter"

	"github.com/mridang/yamllint-go/internal/types"
)

type Config struct {
	Max      int `yaml:"max"`
	MaxStart int `yaml:"max-start"`
	MaxEnd   int `yaml:"max-end"`
}

type Rule struct{}

func (r *Rule) ID() string {
	return "empty-lines"
}

func (r *Rule) Type() string {
	return "line"
}

func (r *Rule) DefaultConfig() Config {
	return Config{
		Max:      2,
		MaxStart: 0,
		MaxEnd:   0,
	}
}

func (r *Rule) Check(config Config, line *types.Line) iter.Seq[types.Problem] {
	return func(yield func(types.Problem) bool) {
		bufferLen := len(line.Buffer)

		if line.Start == line.End && line.End < bufferLen {
			if line.End+2 <= bufferLen && line.Buffer[line.End:line.End+2] == "\n\n" {
				return
			}
			if line.End+4 <= bufferLen && line.Buffer[line.End:line.End+4] == "\r\n\r\n" {
				return
			}

			blankLines := 0
			start := line.Start

			for start >= 2 && line.Buffer[start-2:start] == "\r\n" {
				blankLines++
				start -= 2
			}
			for start >= 1 && line.Buffer[start-1] == '\n' {
				blankLines++
				start -= 1
			}

			maxAllowed := config.Max

			if start == 0 {
				blankLines++
				maxAllowed = config.MaxStart
			}

			isLastCharNewline := line.End == bufferLen-1 && line.Buffer[line.End] == '\n'
			isLastCharsCRLF := line.End == bufferLen-2 && line.Buffer[line.End:line.End+2] == "\r\n"

			if isLastCharNewline || isLastCharsCRLF {
				if line.End == 0 {
					return
				}
				maxAllowed = config.MaxEnd
			}

			if blankLines > maxAllowed {
				if !yield(types.Problem{
					Line:   line.LineNo,
					Column: 0,
					Desc:   fmt.Sprintf("too many blank lines (%d > %d)", blankLines, maxAllowed),
				}) {
					return
				}
			}
		}
	}
}
